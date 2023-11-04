package kubernetes

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	kubernetesapps "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/kubernetes-apps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterDataSource("scp_kubernetes_apps_image", DatasourceKubernetesAppsImage())
}

func DatasourceKubernetesAppsImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceKubernetesImageRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"base_image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contents": {Type: schema.TypeList, Computed: true, Description: "K8s app list", Elem: datasourceAppImageElem()},
		},
		Description: "Provides K8s app image details",
	}
}

func convertStandImageListToHclSet(images []kubernetesapps.ImagesResponse) (common.HclSetObject, []string) {
	var imageSet common.HclSetObject
	var ids []string

	for _, image := range images {
		if len(image.ImageId) == 0 {
			continue
		}

		ids = append(ids, image.ImageId)

		kv := common.HclKeyValueObject{
			"id":         image.ImageId,
			"base_image": image.BaseImage,
			"category":   image.Category,
			"image_name": image.ImageName,
			"version":    image.ImageAttr["display.version"],
		}

		imageSet = append(imageSet, kv)
	}

	return imageSet, ids
}

func datasourceKubernetesImageRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	images, _, err := inst.Client.KubernetesApps.ListImages(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	imageSet, _ := convertStandImageListToHclSet(images)

	if f, ok := rd.GetOk("filter"); ok {
		imageSet = common.ApplyFilter(DatasourceKubernetesAppsImage().Schema, f.(*schema.Set), imageSet)
	}

	if len(imageSet) == 0 {
		return diag.Errorf("no matching kubernetes apps image found")
	}

	for k, v := range imageSet[0] {
		if k == "id" {
			rd.SetId(v.(string))
			continue
		}
		rd.Set(k, v)
	}

	//rd.SetId(uuid.NewV4().String())
	//rd.Set("contents", imageSet[0])
	//rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceAppImageElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"base_image":     {Type: schema.TypeString, Computed: true, Description: "Project Id"},
			"category":       {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			"icon_file":      {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			"icon_file_name": {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			"image_id":       {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			"image_name":     {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			//"image_attr":         {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			"pool_id":            {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			"product_group_id":   {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			"product_group_name": {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			//"product_group_attr": {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			"zone_id":     {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			"description": {Type: schema.TypeString, Computed: true, Description: "K8s Version"},
			"created_by":  {Type: schema.TypeString, Computed: true, Description: "Created By"},
			"created_dt":  {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			"modified_by": {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			"modified_dt": {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
	}
}
