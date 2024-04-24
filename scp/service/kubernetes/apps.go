package kubernetes

import (
	"context"
	"fmt"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_kubernetes_apps", ResourceKubernetesApps())
}

func ResourceKubernetesApps() *schema.Resource {
	return &schema.Resource{
		CreateContext: createApps,
		ReadContext:   readApps,
		UpdateContext: resourceKubernetesAppsUpdate,
		DeleteContext: deleteApps,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Kubernetes app name",
			},
			"engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of scp_kubernetes_engine resource",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Namespace name",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image ID (use scp_standard_image data source)",
			},
			"additional_params": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Additional Params",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a K8s Apps resource.",
	}
}

func createApps(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	name := data.Get("name").(string)
	engineId := data.Get("engine_id").(string)
	namespace := data.Get("namespace").(string)
	imageId := data.Get("image_id").(string)
	additionalParams := data.Get("additional_params").(map[string]interface{})
	tags := data.Get("tags").(map[string]interface{})

	image, _, err := inst.Client.KubernetesApps.ReadImage(ctx, imageId)
	if err != nil {
		return diag.FromErr(err)
	}

	apps, _, err := inst.Client.KubernetesApps.CreateApps(ctx, engineId, namespace, imageId, image.ProductGroupId, name, additionalParams, tags)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(apps.ReleaseId)

	err = client.WaitForStatus(ctx, inst.Client, []string{}, []string{"deployed"}, refreshApps(ctx, meta, data.Id(), true))
	if err != nil {
		return diag.FromErr(err)
	}

	return readApps(ctx, data, meta)
}

func readApps(ctx context.Context, data *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)
	apps, _, err := inst.Client.KubernetesApps.ReadApps(ctx, data.Id())

	if err != nil {
		data.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	data.Set("name", apps.ReleaseName)
	data.Set("engine_id", apps.ClusterId)
	data.Set("namespace", apps.NamespaceName)
	// TODO: Cannot retrieve image id
	tfTags.SetTags(ctx, data, meta, data.Id())

	return nil
}

func resourceKubernetesAppsUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	err = tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return
	}

	return readApps(ctx, rd, meta)
}

func deleteApps(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	engineId := data.Get("engine_id").(string)
	namespace := data.Get("namespace").(string)

	_, err := inst.Client.KubernetesApps.DeleteApps(ctx, engineId, namespace, data.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = client.WaitForStatus(ctx, inst.Client, []string{}, []string{"DELETED"}, refreshApps(ctx, meta, data.Id(), false))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func refreshApps(ctx context.Context, meta interface{}, id string, errorOnNotFound bool) func() (interface{}, string, error) {
	inst := meta.(*client.Instance)

	return func() (interface{}, string, error) {
		apps, httpStatus, err := inst.Client.KubernetesApps.ReadApps(ctx, id)

		if httpStatus == 200 {
			return apps, apps.ReleaseState, nil
		} else if httpStatus == 404 {
			if errorOnNotFound {
				return nil, "", fmt.Errorf("kubernetes apps with id=%s not found", id)
			}

			return apps, "DELETED", nil
		} else if err != nil {
			return nil, "", err
		}

		return nil, "", fmt.Errorf("failed to read kubernetes apps(%s) status:%d", id, httpStatus)
	}
}
