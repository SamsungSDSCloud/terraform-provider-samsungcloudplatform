package autoscaling

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/autoscaling/autoscaling_common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/autoscaling2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func init() {
	scp.RegisterResource("scp_auto_scaling_group_policy", ResourceAutoScalingGroupPolicy())
}

func ResourceAutoScalingGroupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceAutoScalingGroupPolicyCreate,
		ReadContext:   ResourceAutoScalingGroupPolicyRead,
		UpdateContext: ResourceAutoScalingGroupPolicyUpdate,
		DeleteContext: ResourceAutoScalingGroupPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"asg_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Auto-Scaling Group ID",
			},
			"comparison_operator": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Comparison operator",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`GreaterThanOrEqualTo|GreaterThan|LessThanOrEqualTo|LessThan`), "Must be one of \"GreaterThanOrEqualTo\", \"GreaterThan\", \"LessThanOrEqualTo\"  or \"LessThan\"."),
			},
			"cooldown_seconds": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Cooldown seconds",
			},
			"evaluation_minutes": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Evaluation minutes",
			},
			"metric_method": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Metric method",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`AVG|MIN|MAX`), "Must be one of \"AVG\", \"MIN\" or \"MAX\"."),
			},
			"metric_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Metric type",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`CPU|MEMORY|NETWORK_IN|NETWORK_OUT|DISK_READ|DISK_WRITE|DISK_READ_OPER|DISK_WRITE_OPER`), "Must be one of \"CPU\", \"MEMORY\", \"NETWORK_IN\", \"NETWORK_OUT\", \"DISK_READ\", \"DISK_WRITE\", \"DISK_READ_OPER\" or \"DISK_WRITE_OPER\"."),
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy name",
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[a-z][ㄱ-ㅎㅏ-ㅣ가-힣a-zA-Z0-9-]*$`), "Must be 3 to 20, start with a lowercase letter, and use English, Korean, numbers, and -."),
				),
			},
			"scale_method": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Scale method",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`AMOUNT|PERCENTAGE|FIXED`), "Must be one of \"AMOUNT\", \"PERCENTAGE\" or \"FIXED\"."),
			},
			"scale_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Scale type",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`SCALE_OUT|SCALE_IN`), "Must be one of \"SCALE_OUT\" or \"SCALE_IN\"."),
			},
			"scale_value": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Scale value",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"threshold": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Threshold",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project ID",
			},
			"block_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Block ID",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy ID",
			},
			"policy_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy state",
			},
			"service_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service ID",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service zone ID",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The person who created the resource",
			},
			"created_dt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date",
			},
			"modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The person who modified the resource",
			},
			"modified_dt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification date",
			},
		},
		Description: "Provides a Auto-Scaling Group policy resource.",
	}
}

func ResourceAutoScalingGroupPolicyRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	info, _, err := inst.Client.AutoScaling.GetAutoScalingGroupPolicyDetail(ctx, rd.Get("asg_id").(string), rd.Id())

	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	autoscaling_common.SetResponseToResourceData(info, rd)

	return nil
}

func ResourceAutoScalingGroupPolicyCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	cooldownSeconds := int32(rd.Get("cooldown_seconds").(int))
	evaluationMinutes := int32(rd.Get("evaluation_minutes").(int))
	scaleValue := int32(rd.Get("scale_value").(int))
	asgId := rd.Get("asg_id").(string)
	createRequest := autoscaling2.AsgPolicyCreateRequest{
		ComparisonOperator: rd.Get("comparison_operator").(string),
		CooldownSeconds:    &cooldownSeconds,
		EvaluationMinutes:  &evaluationMinutes,
		MetricMethod:       rd.Get("metric_method").(string),
		MetricType:         rd.Get("metric_type").(string),
		PolicyName:         rd.Get("policy_name").(string),
		ScaleMethod:        rd.Get("scale_method").(string),
		ScaleType:          rd.Get("scale_type").(string),
		ScaleValue:         &scaleValue,
		Threshold:          rd.Get("threshold").(string),
	}

	result, _, err := inst.Client.AutoScaling.CreateAutoScalingGroupPolicy(ctx, asgId, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.PolicyId)

	return ResourceAutoScalingGroupPolicyRead(ctx, rd, meta)
}

func ResourceAutoScalingGroupPolicyUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	cooldownSeconds := int32(rd.Get("cooldown_seconds").(int))
	evaluationMinutes := int32(rd.Get("evaluation_minutes").(int))
	scaleValue := int32(rd.Get("scale_value").(int))
	asgId := rd.Get("asg_id").(string)
	updateRequest := autoscaling2.AsgPolicyUpdateRequest{
		ComparisonOperator: rd.Get("comparison_operator").(string),
		CooldownSeconds:    &cooldownSeconds,
		EvaluationMinutes:  &evaluationMinutes,
		MetricMethod:       rd.Get("metric_method").(string),
		MetricType:         rd.Get("metric_type").(string),
		PolicyName:         rd.Get("policy_name").(string),
		ScaleMethod:        rd.Get("scale_method").(string),
		ScaleValue:         &scaleValue,
		Threshold:          rd.Get("threshold").(string),
	}
	_, _, err := inst.Client.AutoScaling.UpdateAutoScalingGroupPolicy(ctx, asgId, rd.Id(), updateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceAutoScalingGroupPolicyRead(ctx, rd, meta)
}

func ResourceAutoScalingGroupPolicyDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, err := inst.Client.AutoScaling.DeleteAutoScalingGroupPolicy(ctx, rd.Get("asg_id").(string), rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}
	return nil
}
