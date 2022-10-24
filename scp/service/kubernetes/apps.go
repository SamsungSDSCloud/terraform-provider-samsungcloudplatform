package kubernetes

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceKubernetesApps() *schema.Resource {
	return &schema.Resource{
		CreateContext: createApps,
		ReadContext:   readApps,
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

	image, _, err := inst.Client.KubernetesApps.ReadImage(ctx, imageId)
	if err != nil {
		return diag.FromErr(err)
	}

	apps, _, err := inst.Client.KubernetesApps.CreateApps(ctx, engineId, namespace, imageId, image.ProductGroupId, name)
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
		return diag.FromErr(err)
	}

	data.Set("name", apps.ReleaseName)
	data.Set("engine_id", apps.ClusterId)
	data.Set("namespace", apps.NamespaceName)
	// TODO: Cannot retrieve image id

	return nil
}

func deleteApps(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	engineId := data.Get("engine_id").(string)
	namespace := data.Get("namespace").(string)

	_, err := inst.Client.KubernetesApps.DeleteApps(ctx, engineId, namespace, data.Id())
	if err != nil {
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
