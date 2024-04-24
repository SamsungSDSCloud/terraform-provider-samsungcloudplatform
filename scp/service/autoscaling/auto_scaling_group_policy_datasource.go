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
	scp.RegisterDataSource("scp_auto_scaling_group_policies", DataSourceAutoScalingGroupPolicies())
	scp.RegisterDataSource("scp_auto_scaling_group_policy", DataSourceAutoScalingGroupPolicy())
}

func DataSourceAutoScalingGroupPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoScalingGroupPolicyList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":        common.DatasourceFilter(),
			"asg_id":        {Type: schema.TypeString, Required: true, Description: "Auto-Scaling Group ID"},
			"metric_method": {Type: schema.TypeString, Optional: true, Description: "Metric method"},
			"metric_type":   {Type: schema.TypeString, Optional: true, Description: "Metric type"},
			"policy_name":   {Type: schema.TypeString, Optional: true, Description: "Policy name"},
			"scale_type":    {Type: schema.TypeString, Optional: true, Description: "Scale type"},
			"page":          {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":          {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":          {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":      {Type: schema.TypeList, Computed: true, Description: "Auto-Scaling Group policy list", Elem: dataSourceAutoScalingGroupPolicyElem()},
			"total_count":   {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of Auto-Scaling Group policies",
	}
}

func dataSourceAutoScalingGroupPolicyElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"asg_id":              {Type: schema.TypeString, Computed: true, Description: "Auto-Scaling Group ID"},
			"comparison_operator": {Type: schema.TypeString, Computed: true, Description: "Comparison operator"},
			"cooldown_seconds":    {Type: schema.TypeInt, Computed: true, Description: "Cooldown seconds"},
			"evaluation_minutes":  {Type: schema.TypeInt, Computed: true, Description: "Evaluation minutes"},
			"metric_method":       {Type: schema.TypeString, Computed: true, Description: "Metric method"},
			"metric_type":         {Type: schema.TypeString, Computed: true, Description: "Metric type"},
			"policy_id":           {Type: schema.TypeString, Computed: true, Description: "Policy ID"},
			"policy_name":         {Type: schema.TypeString, Computed: true, Description: "Policy name"},
			"policy_state":        {Type: schema.TypeString, Computed: true, Description: "Policy state"},
			"scale_method":        {Type: schema.TypeString, Computed: true, Description: "Scale method"},
			"scale_type":          {Type: schema.TypeString, Computed: true, Description: "Scale type"},
			"scale_value":         {Type: schema.TypeInt, Computed: true, Description: "Scale value"},
			"threshold":           {Type: schema.TypeString, Computed: true, Description: "Threshold"},
			"threshold_unit":      {Type: schema.TypeString, Computed: true, Description: "Threshold unit"},
			"created_by":          {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":          {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":         {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":         {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func dataSourceAutoScalingGroupPolicyList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	response, _, err := inst.Client.AutoScaling.GetAutoScalingGroupPolicyList(ctx, rd.Get("asg_id").(string), &autoscaling2.AsgPolicyV2ApiGetAsgPolicyListV2Opts{
		MetricMethod: common.GetKeyString(rd, "metric_method"),
		MetricType:   common.GetKeyString(rd, "metric_type"),
		PolicyName:   common.GetKeyString(rd, "policy_name"),
		ScaleType:    common.GetKeyString(rd, "scale_type"),
		Page:         optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:         optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:         optional.NewString(rd.Get("sort").(string)),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(response.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DataSourceAutoScalingGroupPolicies().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}

func DataSourceAutoScalingGroupPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoScalingGroupPolicyDetail,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"asg_id":              {Type: schema.TypeString, Required: true, Description: "Auto-Scaling Group ID"},
			"policy_id":           {Type: schema.TypeString, Required: true, Description: "Policy ID"},
			"project_id":          {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"block_id":            {Type: schema.TypeString, Computed: true, Description: "Block ID"},
			"comparison_operator": {Type: schema.TypeString, Computed: true, Description: "Comparison operator"},
			"cooldown_seconds":    {Type: schema.TypeInt, Computed: true, Description: "Cooldown seconds"},
			"evaluation_minutes":  {Type: schema.TypeInt, Computed: true, Description: "Evaluation minutes"},
			"metric_method":       {Type: schema.TypeString, Computed: true, Description: "Metric method"},
			"metric_type":         {Type: schema.TypeString, Computed: true, Description: "Metric type"},
			"policy_name":         {Type: schema.TypeString, Computed: true, Description: "Policy name"},
			"policy_state":        {Type: schema.TypeString, Computed: true, Description: "Policy state"},
			"scale_method":        {Type: schema.TypeString, Computed: true, Description: "Scale method"},
			"scale_type":          {Type: schema.TypeString, Computed: true, Description: "Scale type"},
			"scale_value":         {Type: schema.TypeInt, Computed: true, Description: "Scale value"},
			"service_zone_id":     {Type: schema.TypeString, Computed: true, Description: "Service zone ID"},
			"threshold":           {Type: schema.TypeString, Computed: true, Description: "Threshold"},
			"threshold_unit":      {Type: schema.TypeString, Computed: true, Description: "Threshold unit"},
			"created_by":          {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":          {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":         {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":         {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
		Description: "Provides details of Auto-Scaling Group policy",
	}
}

func dataSourceAutoScalingGroupPolicyDetail(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	rd.SetId(rd.Get("policy_id").(string))
	response, _, err := inst.Client.AutoScaling.GetAutoScalingGroupPolicyDetail(ctx, rd.Get("asg_id").(string), rd.Id())

	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	autoscaling_common.SetResponseToResourceData(response, rd)

	return nil
}
