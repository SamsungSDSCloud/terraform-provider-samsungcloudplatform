package image

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/library/image2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ActiveState string = "ACTIVE"
)

func DatasourceStandardImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceStandardImageRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"service_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service group (COMPUTE, CONTAINER, ...)",
			},
			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service (Virtual Server, Kubernetes Engine VM, ...)",
			},
			"base_image": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Base image for service",
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image category",
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image name",
			},
			"image_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image type (STANDARD)",
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS type (Windows, Ubuntu, ..)",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region name",
			},
		},
		Description: "Provides standard image details",
	}
}

func convertStandImageListToHclSet(standardImages []image2.StandardImageResponse, serviceGroup string, service string) (common.HclSetObject, []string) {
	var setStandardImages common.HclSetObject
	var ids []string
	// Convert to HclSet
	for _, si := range standardImages {
		if len(si.ImageId) == 0 {
			continue
		}
		ids = append(ids, si.ImageId)
		kv := common.HclKeyValueObject{
			"id":            si.ImageId,
			"base_image":    si.BaseImage,
			"category":      si.Category,
			"image_name":    si.ImageName,
			"image_type":    si.ImageType,
			"os_type":       si.OsType,
			"description":   si.ImageDescription,
			"service":       service,
			"service_group": serviceGroup,
		}
		setStandardImages = append(setStandardImages, kv)

	}
	return setStandardImages, ids
}

func datasourceStandardImageRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	projectInfo, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	servicedZoneId := projectInfo.DefaultZoneId
	if len(servicedZoneId) == 0 {
		vpcLocation := rd.Get("region").(string)
		servicedZoneId, _, err = client.FindServiceZoneIdAndProductGroupId(ctx, inst.Client, vpcLocation, common.NetworkProductGroup, common.VpcProductName)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	serviceGroup := rd.Get("service_group").(string)
	service := rd.Get("service").(string)

	responseStandardImages, err := inst.Client.Image.GetStandardImageList(ctx, servicedZoneId, ActiveState, serviceGroup, service)
	if err != nil {
		return diag.FromErr(err)
	}

	setStandardImages, _ := convertStandImageListToHclSet(responseStandardImages.Contents, serviceGroup, service)

	if f, ok := rd.GetOk("filter"); ok {
		setStandardImages = common.ApplyFilter(DatasourceStandardImage().Schema, f.(*schema.Set), setStandardImages)
	}

	if len(setStandardImages) == 0 {
		return diag.Errorf("no matching standard image found")
	}

	for k, v := range setStandardImages[0] {
		if k == "id" {
			rd.SetId(v.(string))
			continue
		}
		rd.Set(k, v)
	}

	return nil
}
