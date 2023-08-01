package image

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/image2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_migration_image", DatasourceMigrationImage())
	scp.RegisterDataSource("scp_migration_images", DatasourceMigrationImages())
}

func DatasourceMigrationImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceMigrationImageRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:      elemMigrationImage(),
		Description: "Provides a Migration Image details.",
	}
}
func DatasourceMigrationImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceMigrationImagesList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"image_name":        {Type: schema.TypeString, Optional: true},
			"image_state":       {Type: schema.TypeString, Optional: true, Default: "ACTIVE"},
			"origin_image_name": {Type: schema.TypeString, Optional: true},
			"service_group":     {Type: schema.TypeString, Required: true, Description: "Service group (COMPUTE, DATABASE, EXTENSION, ...)"},
			"service":           {Type: schema.TypeString, Required: true, Description: "Service (Baremetal Server, Virtual Server, ...)"},
			"created_by":        {Type: schema.TypeString, Optional: true},
			"page":              {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":              {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":              {Type: schema.TypeString, Optional: true, Default: "imageName:asc", Description: "Sort rule to get list"},
			"contents":          {Type: schema.TypeList, Optional: true, Description: "Migration image list", Elem: &schema.Resource{Schema: elemMigrationImage()}},
			"total_count":       {Type: schema.TypeInt, Optional: true, Description: "Migration images total_count"},
			"region":            {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Region name"},
		},
		Description: "Provides list of migration images",
	}
}

func elemMigrationImage() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id":               {Type: schema.TypeString, Computed: true},
		"availability_zone_name":   {Type: schema.TypeString, Computed: true},
		"base_image":               {Type: schema.TypeString, Computed: true},
		"block_id":                 {Type: schema.TypeString, Computed: true},
		"category":                 {Type: schema.TypeString, Computed: true},
		"default_disk_size":        {Type: schema.TypeInt, Computed: true},
		"disk_size":                {Type: schema.TypeInt, Computed: true, Description: "Extra disk size."},
		"disks":                    {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: elemDisk()}},
		"icon":                     {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"image_id":                 {Type: schema.TypeString, Required: true, ForceNew: true},
		"image_name":               {Type: schema.TypeString, Computed: true},
		"image_state":              {Type: schema.TypeString, Computed: true, Description: "Image state (ACTIVE)"},
		"image_type":               {Type: schema.TypeString, Computed: true},
		"origin_image_id":          {Type: schema.TypeString, Computed: true},
		"origin_image_name":        {Type: schema.TypeString, Computed: true},
		"origin_virtual_server_id": {Type: schema.TypeString, Computed: true},
		"os_type":                  {Type: schema.TypeString, Computed: true, Description: "OS type (Windows, Ubuntu, ..)"},
		"product_group_id":         {Type: schema.TypeString, Computed: true},
		"products":                 {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: elemProduct()}},
		"properties":               {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"service_zone_id":          {Type: schema.TypeString, Required: true},
		"image_description":        {Type: schema.TypeString, Computed: true, Description: "Migration image description. (Up to 50 characters)"},
		"created_by":               {Type: schema.TypeString, Computed: true},
		"created_dt":               {Type: schema.TypeString, Computed: true},
		"modified_by":              {Type: schema.TypeString, Computed: true},
		"modified_dt":              {Type: schema.TypeString, Computed: true},
	}
}
func datasourceMigrationImageRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	rd.SetId(rd.Get("image_id").(string))

	responseMigrationImage, _, err := inst.Client.MigrationImage.GetMigrationImageInfo(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	// response를 map형태로 변환 후 스키마에 반영
	mapMigrationImageDetail := common.ToMap(responseMigrationImage)
	for k, v := range mapMigrationImageDetail {
		if rd.Set(k, v) != nil {
			return nil
		}
	}

	// 리스트형 속성일 경우 예외처리 필요
	if nil != rd.Set("disks", common.ConvertStructToMaps(responseMigrationImage.Disks)) {
		return nil
	}
	if nil != rd.Set("products", common.ConvertStructToMaps(responseMigrationImage.Products)) {
		return nil
	}

	return nil
}

func datasourceMigrationImagesList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	projectInfo, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var servicedZoneId string
	vpcLocation := rd.Get("region").(string)

	if len(vpcLocation) == 0 {
		servicedZoneId = projectInfo.DefaultZoneId
	} else {
		servicedZoneId, _, err = client.FindServiceZoneIdAndProductGroupId(ctx, inst.Client, vpcLocation, common.NetworkProductGroup, common.VpcProductName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	responseMigrationImages, err := inst.Client.MigrationImage.GetMigrationImageList(ctx, image2.MigrationImageV2ApiListMigrationImagesOpts{
		ServiceZoneId:    optional.NewString(servicedZoneId),
		ImageState:       optional.NewString(rd.Get("image_state").(string)),
		ServicedFor:      optional.NewString(rd.Get("service").(string)),
		ServicedGroupFor: optional.NewString(rd.Get("service_group").(string)),
		Page:             optional.NewInt32(0),
		Size:             optional.NewInt32(10000),
		Sort:             optional.NewInterface([]string{"imageName:asc"}),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(responseMigrationImages.Contents) == 0 {
		return diag.Errorf("no matching Migration image found")
	}

	// Key-Value 오브젝트화
	var setMigrationImages = convertMigrationImageListToHclSet(responseMigrationImages.Contents)

	// 조회결과를 스키마에 저장
	rd.SetId(uuid.NewV4().String())
	if nil != rd.Set("contents", setMigrationImages) {
		return nil
	}
	if nil != rd.Set("total_count", responseMigrationImages.TotalCount) {
		return nil
	}

	return nil
}

func convertMigrationImageListToHclSet(contents []image2.MigrationImageResponse) common.HclSetObject {
	var setContents common.HclSetObject

	for _, content := range contents {

		elemContent := make(common.HclKeyValueObject)
		mapContent := common.ToMap(content)
		for k, _ := range elemMigrationImage() {
			elemContent[k] = mapContent[k]
		}

		elemContent["disks"] = common.ConvertStructToMaps(content.Disks)
		elemContent["products"] = common.ConvertStructToMaps(content.Products)

		setContents = append(setContents, elemContent)
	}

	return setContents
}
