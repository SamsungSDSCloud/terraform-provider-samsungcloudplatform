package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	samsungcloudplatform.RegisterResource("Kubernetes", "samsungcloudplatform_kubernetes_namespace", ResourceKubernetesNamespace())
}

func ResourceKubernetesNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: createNamespace,
		ReadContext:   readNamespace,
		DeleteContext: deleteNamespace,

		Importer: &schema.ResourceImporter{
			// Composite import id: "<engineId>/<namespaceName>". Read uses
			// data.Get("engine_id") so passthrough alone leaves it empty,
			// causing a 400 from the API on the import refresh.
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.SplitN(d.Id(), "/", 2)
				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("import id must be in format <engineId>/<namespaceName>, got %q", d.Id())
				}
				if err := d.Set("engine_id", parts[0]); err != nil {
					return nil, err
				}
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Namespace name",
				ValidateFunc: validation.StringLenBetween(1, 63),
			},
			"engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of scp_kubernetes_engine resource",
			},
		},
		Description: "Provides a K8s Namespace resource.",
	}
}

func createNamespace(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	engineId := data.Get("engine_id").(string)
	name := data.Get("name").(string)

	_, _, err := inst.Client.Kubernetes.CreateNamespace(ctx, engineId, name)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(name)
	return readNamespace(ctx, data, meta)
}

func readNamespace(ctx context.Context, data *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil

	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)
	engineId := data.Get("engine_id").(string)
	name := data.Id()
	ns, _, err := inst.Client.Kubernetes.ReadNamespace(ctx, engineId, name)

	if err != nil {
		data.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	data.Set("name", ns.NamespaceName)
	return nil
}

func deleteNamespace(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	engineId := data.Get("engine_id").(string)
	name := data.Id()

	_, err := inst.Client.Kubernetes.DeleteNamespace(ctx, engineId, name)
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	return nil
}
