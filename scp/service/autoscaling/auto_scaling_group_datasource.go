package autoscaling

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/autoscaling/autoscaling_common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/autoscaling2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_auto_scaling_groups", DataSourceAutoScalingGroups())
	scp.RegisterDataSource("scp_auto_scaling_group", DataSourceAutoScalingGroup())
}

func DataSourceAutoScalingGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoScalingGroupList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"filter":            common.DatasourceFilter(),
			"asg_name":          {Type: schema.TypeString, Optional: true, Description: "Auto-Scaling Group name"},
			"asg_state":         {Type: schema.TypeString, Optional: true, Description: "Auto-Scaling Group state"},
			"lc_name":           {Type: schema.TypeString, Optional: true, Description: "Launch Configuration name"},
			"local_subnet_id":   {Type: schema.TypeString, Optional: true, Description: "Local subnet ID"},
			"security_group_id": {Type: schema.TypeString, Optional: true, Description: "Security Group ID"},
			"service_id":        {Type: schema.TypeString, Optional: true, Description: "Service ID"},
			"service_zone_id":   {Type: schema.TypeString, Optional: true, Description: "Service zone ID"},
			"vpc_id":            {Type: schema.TypeString, Optional: true, Description: "VPC ID"},
			"subnet_id":         {Type: schema.TypeString, Optional: true, Description: "Subnet ID"},
			"created_by":        {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"page":              {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":              {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":              {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":          {Type: schema.TypeList, Computed: true, Description: "Auto-Scaling Group list", Elem: dataSourceAutoScalingGroupElem()},
			"total_count":       {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of Auto-Scaling Groups",
	}
}

func dataSourceAutoScalingGroupList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	// Call the API to retrieve the ASG list
	responses, _, err := inst.Client.AutoScaling.GetAutoScalingGroupList(context.Background(), &autoscaling2.AutoScalingGroupV2ApiGetAsgListV2Opts{
		AsgName:         optional.NewString(rd.Get("asg_name").(string)),
		AsgState:        common.GetKeyString(rd, "asg_state"),
		LcName:          optional.NewString(rd.Get("lc_name").(string)),
		LocalSubnetId:   optional.NewString(rd.Get("local_subnet_id").(string)),
		SecurityGroupId: optional.NewString(rd.Get("security_group_id").(string)),
		ServiceId:       optional.NewString(rd.Get("service_id").(string)),
		ServiceZoneId:   optional.NewString(rd.Get("service_zone_id").(string)),
		VpcId:           optional.NewString(rd.Get("vpc_id").(string)),
		SubnetId:        optional.NewString(rd.Get("subnet_id").(string)),
		CreatedBy:       optional.NewString(rd.Get("created_by").(string)),
		Page:            optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:            optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:            optional.NewString(rd.Get("sort").(string)),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DataSourceAutoScalingGroups().Schema, f.(*schema.Set), contents)
	}

	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))
	rd.SetId(uuid.NewV4().String())

	return nil
}

func dataSourceAutoScalingGroupElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"asg_id":                          {Type: schema.TypeString, Computed: true, Description: "Auto-Scaling Group ID"},
			"asg_name":                        {Type: schema.TypeString, Computed: true, Description: "Auto-Scaling Group name"},
			"asg_state":                       {Type: schema.TypeString, Computed: true, Description: "Auto-Scaling Group state"},
			"project_id":                      {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"block_id":                        {Type: schema.TypeString, Computed: true, Description: "Block ID"},
			"service_id":                      {Type: schema.TypeString, Computed: true, Description: "Service ID"},
			"service_zone_id":                 {Type: schema.TypeString, Computed: true, Description: "Service zone ID"},
			"lc_id":                           {Type: schema.TypeString, Computed: true, Description: "Launch Configuration ID"},
			"lc_name":                         {Type: schema.TypeString, Computed: true, Description: "Launch Configuration name"},
			"availability_zone_name":          {Type: schema.TypeString, Computed: true, Description: "Availability zone name"},
			"desired_server_count":            {Type: schema.TypeInt, Computed: true, Description: "Desired server count"},
			"min_server_count":                {Type: schema.TypeInt, Computed: true, Description: "Min server count"},
			"max_server_count":                {Type: schema.TypeInt, Computed: true, Description: "Max server count"},
			"desired_server_count_editable":   {Type: schema.TypeBool, Computed: true, Description: "Desired server count editable"},
			"multi_availability_zone_enabled": {Type: schema.TypeBool, Computed: true, Description: "Multi availability zone enabled"},
			"is_terminating":                  {Type: schema.TypeBool, Computed: true, Description: "Is terminating"},
			"vpc_id":                          {Type: schema.TypeString, Computed: true, Description: "VPC ID"},
			"subnet_id":                       {Type: schema.TypeString, Computed: true, Description: "Subnet ID"},
			"local_subnet_id":                 {Type: schema.TypeString, Computed: true, Description: "Local subnet ID"},
			"server_name_prefix":              {Type: schema.TypeString, Computed: true, Description: "Server name prefix"},
			"created_by":                      {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":                      {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":                     {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":                     {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func DataSourceAutoScalingGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoScalingGroupDetail,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"asg_id":                          {Type: schema.TypeString, Required: true, Description: "Auto-Scaling Group ID"},
			"asg_name":                        {Type: schema.TypeString, Computed: true, Description: "Auto-Scaling Group name"},
			"asg_state":                       {Type: schema.TypeString, Computed: true, Description: "Auto-Scaling Group state"},
			"project_id":                      {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"block_id":                        {Type: schema.TypeString, Computed: true, Description: "Block ID"},
			"service_id":                      {Type: schema.TypeString, Computed: true, Description: "Service ID"},
			"service_zone_id":                 {Type: schema.TypeString, Computed: true, Description: "Service zone ID"},
			"lc_id":                           {Type: schema.TypeString, Computed: true, Description: "Launch Configuration ID"},
			"lc_name":                         {Type: schema.TypeString, Computed: true, Description: "Launch Configuration name"},
			"availability_zone_name":          {Type: schema.TypeString, Computed: true, Description: "Availability zone name"},
			"desired_server_count":            {Type: schema.TypeInt, Computed: true, Description: "Desired server count"},
			"min_server_count":                {Type: schema.TypeInt, Computed: true, Description: "Min server count"},
			"max_server_count":                {Type: schema.TypeInt, Computed: true, Description: "Max server count"},
			"desired_server_count_editable":   {Type: schema.TypeBool, Computed: true, Description: "Desired server count editable"},
			"multi_availability_zone_enabled": {Type: schema.TypeBool, Computed: true, Description: "Multi availability zone enabled"},
			"is_terminating":                  {Type: schema.TypeBool, Computed: true, Description: "Is terminating"},
			"server_name_prefix":              {Type: schema.TypeString, Computed: true, Description: "Server name prefix"},
			"vpc_info":                        {Type: schema.TypeList, Computed: true, Description: "VPC information", Elem: resourceVpcInfo()},
			"security_group_ids":              {Type: schema.TypeList, Computed: true, Description: "Security Group ID", Elem: &schema.Schema{Type: schema.TypeString}},
			"file_storage_id":                 {Type: schema.TypeString, Computed: true, Description: "File Storage ID"},
			"created_by":                      {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":                      {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":                     {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":                     {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			"dns_enabled":                     {Type: schema.TypeBool, Computed: true, Description: "DNS enabled"},
		},
		Description: "Provides details of Auto-Scaling Group",
	}
}

func dataSourceAutoScalingGroupDetail(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	rd.SetId(rd.Get("asg_id").(string))
	response, _, err := inst.Client.AutoScaling.GetAutoScalingGroupDetail(ctx, rd.Id())

	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	autoscaling_common.SetResponseToResourceData(response, rd, "SubnetId", "VpcId", "LocalSubnetId", "DeploymentEnvType", "ServiceLevelProductId")

	return nil
}
