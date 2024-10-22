package autoscaling

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/autoscaling/autoscaling_common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/autoscaling2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_launch_configurations", DataSourceLaunchConfigurations())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_launch_configuration", DataSourceLaunchConfiguration())
}

func DataSourceLaunchConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLaunchConfigurationList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":          common.DatasourceFilter(),
			"image_id":        {Type: schema.TypeString, Optional: true, Description: "Image ID"},
			"lc_name":         {Type: schema.TypeString, Optional: true, Description: "Launch Configuration name"},
			"service_zone_id": {Type: schema.TypeString, Optional: true, Description: "Service zone ID"},
			"created_by":      {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"page":            {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":            {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":            {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":        {Type: schema.TypeList, Computed: true, Description: "Launch Configuration list", Elem: dataSourceLaunchConfigurationElem()},
			"total_count":     {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of Launch Configurations",
	}
}

func dataSourceLaunchConfigurationElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":      {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"asg_count":       {Type: schema.TypeInt, Computed: true, Description: "Number of Auto-Scaling Group using this Launch Configuration"},
			"block_id":        {Type: schema.TypeString, Computed: true, Description: "Block ID"},
			"image_id":        {Type: schema.TypeString, Computed: true, Description: "Image ID"},
			"key_pair_id":     {Type: schema.TypeString, Computed: true, Description: "Key pair ID"},
			"lc_id":           {Type: schema.TypeString, Computed: true, Description: "Launch Configuration ID"},
			"lc_name":         {Type: schema.TypeString, Computed: true, Description: "Launch Configuration name"},
			"service_zone_id": {Type: schema.TypeString, Computed: true, Description: "Service zone ID"},
			"created_by":      {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":      {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":     {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":     {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func dataSourceLaunchConfigurationList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.AutoScaling.GetLaunchConfigurationList(ctx, &autoscaling2.AsgLaunchConfigurationV2ApiGetLaunchConfigListV2Opts{
		ImageId:       optional.NewString(rd.Get("image_id").(string)),
		LcName:        optional.NewString(rd.Get("lc_name").(string)),
		ServiceZoneId: optional.NewString(rd.Get("service_zone_id").(string)),
		CreatedBy:     optional.NewString(rd.Get("created_by").(string)),
		Page:          optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:          optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:          optional.NewString(rd.Get("sort").(string)),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DataSourceLaunchConfigurations().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}

func DataSourceLaunchConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLaunchConfigurationDetail,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id":          {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"asg_ids":             {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Auto-Scaling Group ID list"},
			"block_id":            {Type: schema.TypeString, Computed: true, Description: "Block ID"},
			"block_storages":      {Type: schema.TypeList, Computed: true, Elem: dataSourceLaunchConfigurationBlockStorageElem(), Description: "Block Storage list"},
			"contract_product_id": {Type: schema.TypeString, Computed: true, Description: "Contract product ID"},
			"image_id":            {Type: schema.TypeString, Computed: true, Description: "Image ID"},
			"initial_script":      {Type: schema.TypeString, Computed: true, Description: "Virtual Server's initial script"},
			"key_pair_id":         {Type: schema.TypeString, Computed: true, Description: "Key pair ID"},
			"lc_id":               {Type: schema.TypeString, Required: true, Description: "Launch Configuration ID"},
			"lc_name":             {Type: schema.TypeString, Computed: true, Description: "Launch Configuration name"},
			"os_product_id":       {Type: schema.TypeString, Computed: true, Description: "OS product ID"},
			"os_type":             {Type: schema.TypeString, Computed: true, Description: "OS type"},
			"product_group_id":    {Type: schema.TypeString, Computed: true, Description: "Product group ID"},
			"scale_product_id":    {Type: schema.TypeString, Computed: true, Description: "Scale product ID"},
			"server_type":         {Type: schema.TypeString, Computed: true, Description: "Server type"},
			"service_zone_id":     {Type: schema.TypeString, Computed: true, Description: "Service zone ID"},
			"created_by":          {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":          {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":         {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":         {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
		Description: "Provides details of Launch Configuration",
	}
}

func dataSourceLaunchConfigurationBlockStorageElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_storage_size": {Type: schema.TypeInt, Computed: true, Description: "Block Storage size (GB)"},
			"disk_type":          {Type: schema.TypeString, Computed: true, Description: "Block storage product (default value : SSD)"},
			"encryption_enabled": {Type: schema.TypeBool, Computed: true, Description: "Encryption enabled"},
			"is_boot_disk":       {Type: schema.TypeBool, Computed: true, Description: "Is boot disk or not"},
			"product_id":         {Type: schema.TypeString, Computed: true, Description: "Product ID"},
		},
	}
}

func dataSourceLaunchConfigurationDetail(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	rd.SetId(rd.Get("lc_id").(string))
	response, _, err := inst.Client.AutoScaling.GetLaunchConfigurationDetail(ctx, rd.Id())

	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	autoscaling_common.SetResponseToResourceData(response, rd)

	return nil
}
