package image

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	image "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/image2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_custom_image", DatasourceCustomImage())
	scp.RegisterDataSource("scp_custom_images", DatasourceCustomImages())
}
func elemDisk() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"boot_enabled":    {Type: schema.TypeBool, Computed: true},
		"device_node":     {Type: schema.TypeString, Computed: true},
		"disk_size":       {Type: schema.TypeInt, Computed: true},
		"encrypt_enabled": {Type: schema.TypeBool, Computed: true},
		"image_id":        {Type: schema.TypeString, Computed: true},
		"product_id":      {Type: schema.TypeString, Computed: true},
		"seq":             {Type: schema.TypeInt, Computed: true},
		"created_by":      {Type: schema.TypeString, Computed: true},
		"created_dt":      {Type: schema.TypeString, Computed: true},
		"modified_by":     {Type: schema.TypeString, Computed: true},
		"modified_dt":     {Type: schema.TypeString, Computed: true},
	}
}

func elemProduct() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"image_id":      {Type: schema.TypeString, Computed: true},
		"product_id":    {Type: schema.TypeString, Computed: true},
		"product_name":  {Type: schema.TypeString, Computed: true},
		"product_type":  {Type: schema.TypeString, Computed: true},
		"product_value": {Type: schema.TypeString, Computed: true},
		"seq":           {Type: schema.TypeInt, Computed: true},
		"created_dt":    {Type: schema.TypeString, Computed: true},
	}
}

func elemCustomImage() map[string]*schema.Schema {
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
		"image_type":               {Type: schema.TypeString, Computed: true, Description: "Image type (STANDARD, CUSTOM, MIGRATION)"},
		"origin_image_id":          {Type: schema.TypeString, Computed: true},
		"origin_image_name":        {Type: schema.TypeString, Computed: true},
		"origin_virtual_server_id": {Type: schema.TypeString, Computed: true},
		"os_type":                  {Type: schema.TypeString, Computed: true, Description: "OS type (Windows, Ubuntu, ..)"},
		"product_group_id":         {Type: schema.TypeString, Computed: true},
		"products":                 {Type: schema.TypeList, Computed: true, Elem: &schema.Resource{Schema: elemProduct()}},
		"properties":               {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"service_zone_id":          {Type: schema.TypeString, Computed: true},
		"image_description":        {Type: schema.TypeString, Computed: true, Description: "Custom image description. (Up to 50 characters)"},
		"created_by":               {Type: schema.TypeString, Computed: true},
		"created_dt":               {Type: schema.TypeString, Computed: true},
		"modified_by":              {Type: schema.TypeString, Computed: true},
		"modified_dt":              {Type: schema.TypeString, Computed: true},
	}
}
func DatasourceCustomImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceCustomImageRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:      elemCustomImage(),
		Description: "Provides a Custom Image details.",
	}
}

func DatasourceCustomImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceCustomImagesList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"image_name":        {Type: schema.TypeString, Optional: true},
			"image_state":       {Type: schema.TypeString, Optional: true, Default: ActiveState},
			"origin_image_name": {Type: schema.TypeString, Optional: true},
			"service_group":     {Type: schema.TypeString, Required: true, Description: "Service group (COMPUTE, DATABASE, EXTENSION, ...)"},
			"service":           {Type: schema.TypeString, Required: true, Description: "Service (Baremetal Server, Virtual Server, ...)"},
			"created_by":        {Type: schema.TypeString, Optional: true},
			"page":              {Type: schema.TypeInt, Optional: true, Description: "Page start number from which to get the list"},
			"size":              {Type: schema.TypeInt, Optional: true, Description: "Size to get list"},
			"sort":              {Type: schema.TypeString, Optional: true, Description: "Sort rule to get list"},
			"contents":          {Type: schema.TypeList, Optional: true, Description: "Custom image list", Elem: &schema.Resource{Schema: elemCustomImage()}},
			"total_count":       {Type: schema.TypeInt, Optional: true, Description: "Custom images total_count"},
			"region":            {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Region name"},
		},
		Description: "Provides list of custom images",
	}
}

func datasourceCustomImageRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	rd.SetId(rd.Get("image_id").(string))

	// 커스텀이미지 리스트 조회
	responseCustomImage, _, err := inst.Client.CustomImage.GetCustomImage(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	// response를 map형태로 변환 후 스키마에 반영
	mapCustomImageDetail := common.ToMap(responseCustomImage)
	for k, v := range mapCustomImageDetail {
		if rd.Set(k, v) != nil {
			return nil
		}
	}

	// 리스트형 속성일 경우 예외처리 필요
	if nil != rd.Set("disks", common.ConvertStructToMaps(responseCustomImage.Disks)) {
		return nil
	}
	if nil != rd.Set("products", common.ConvertStructToMaps(responseCustomImage.Products)) {
		return nil
	}

	return nil
}

func datasourceCustomImagesList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// 커스텀이미지 리스트 조회
	responseCustomImages, err := inst.Client.CustomImage.GetCustomImageList(ctx, image.CustomImageV2ApiListCustomImagesOpts{
		ImageName:        optional.NewString(rd.Get("image_name").(string)),
		ImageState:       optional.NewString(rd.Get("image_state").(string)),
		OriginImageName:  optional.NewString(rd.Get("origin_image_name").(string)),
		ServicedFor:      optional.NewString(rd.Get("service").(string)),
		ServicedGroupFor: optional.NewString(rd.Get("service_group").(string)),
		CreatedBy:        optional.NewString(rd.Get("created_by").(string)),
		ServiceZoneId:    optional.NewString(servicedZoneId),
		Page:             optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:             optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:             optional.NewInterface([]string{rd.Get("sort").(string)}),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(responseCustomImages.Contents) == 0 {
		return diag.Errorf("no matching custom image found")
	}

	// Key-Value 오브젝트화
	var setCustomImages = convertCustomImageListToHclSet(responseCustomImages.Contents)

	// 조회결과를 스키마에 저장
	rd.SetId(uuid.NewV4().String())
	if nil != rd.Set("contents", setCustomImages) {
		return nil
	}
	if nil != rd.Set("total_count", responseCustomImages.TotalCount) {
		return nil
	}

	return nil
}

func convertCustomImageListToHclSet(contents []image.CustomImageResponse) common.HclSetObject {
	var setContents common.HclSetObject

	for _, content := range contents {

		// response 데이터 중에서 elemCustomImage 포맷에 있는 값만 추출하여 새로운 map을 만듬
		elemContent := make(common.HclKeyValueObject)
		mapContent := common.ToMap(content)
		for k, _ := range elemCustomImage() {
			elemContent[k] = mapContent[k]
		}

		// 리스트형 속성일 경우 예외처리 필요
		elemContent["disks"] = common.ConvertStructToMaps(content.Disks)
		elemContent["products"] = common.ConvertStructToMaps(content.Products)

		// 추출한 map들을 반환할 리스트에 append
		setContents = append(setContents, elemContent)
	}

	return setContents
}
