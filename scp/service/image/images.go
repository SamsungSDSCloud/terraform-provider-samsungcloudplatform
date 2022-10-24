package image

import (
	"context"
	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/ScpDevTerra/trf-provider/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DatasourceStandardImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceStandardImagesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter": common.DatasourceFilter(),
			"service_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service group (COMPUTE, CONTAINER, DATABASE, ...)",
			},
			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service (Virtual Server, Kubernetes Engine VM, Postgresql, ...)",
			},
			"standard_images": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Standard image list",
				Elem:        common.GetDatasourceItemsSchema(DatasourceStandardImage()),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region name",
			},
		},
		Description: "Provides list of standard images",
	}
}

func datasourceStandardImagesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	setStandardImages, ids := convertStandImageListToHclSet(responseStandardImages.Contents, serviceGroup, service)

	if f, ok := rd.GetOk("filter"); ok {
		setStandardImages = common.ApplyFilter(DatasourceStandardImages().Schema, f.(*schema.Set), setStandardImages)
	}

	rd.SetId(common.GenerateHash(ids))
	rd.Set("ids", ids)
	rd.Set("standard_images", setStandardImages)

	return nil
}
