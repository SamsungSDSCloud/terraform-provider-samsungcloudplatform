package autoscaling

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/autoscaling/autoscaling_common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/autoscaling2"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/loadbalancer2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_auto_scaling_group_load_balancer", ResourceAutoScalingGroupLoadBalancer())
}

func ResourceAutoScalingGroupLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutoScalingGroupLoadBalancerCreate,
		ReadContext:   resourceAutoScalingGroupLoadBalancerRead,
		UpdateContext: resourceAutoScalingGroupLoadBalancerUpdate,
		DeleteContext: resourceAutoScalingGroupLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"asg_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auto-Scaling Group ID",
			},
			"lb_rule_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    10,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "LB rule ID list connected to Auto-Scaling Group",
			},
		},
		Description: "Provides LB Service resource connected to Auto-Scaling Group.",
	}
}

func resourceAutoScalingGroupLoadBalancerCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rd.SetId(rd.Get("asg_id").(string))

	diagnostics, response := getAutoScalingGroupLoadBalancer(ctx, rd, meta)
	if len(diagnostics) > 0 {
		return diagnostics
	}

	request := generateAsgLoadBalancersUpdateRequest(rd, response)
	diagnostics = updateAutoScalingGroupLoadBalancer(ctx, rd, meta, request)
	if len(diagnostics) > 0 {
		return diagnostics
	}

	return resourceAutoScalingGroupLoadBalancerRead(ctx, rd, meta)
}

func getAutoScalingGroupLoadBalancer(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diag.Diagnostics, loadbalancer2.ListResponseLbServiceForAsgResponse) {
	inst := meta.(*client.Instance)

	response, _, err := inst.Client.LoadBalancer.GetLoadBalancerServiceConnectedToAsgList(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil, loadbalancer2.ListResponseLbServiceForAsgResponse{}
		}
		return diag.FromErr(err), loadbalancer2.ListResponseLbServiceForAsgResponse{}
	}

	return diag.Diagnostics{}, response
}

func generateAsgLoadBalancersUpdateRequest(rd *schema.ResourceData, response loadbalancer2.ListResponseLbServiceForAsgResponse) autoscaling2.AsgLoadBalancersUpdateRequest {
	currentLbRuleIds := extractLbRuleIdsFromResponse(response)
	desiredLbRuleIds := rd.Get("lb_rule_ids").(*schema.Set)
	attachLbRuleIds := autoscaling_common.CalculateSetDifference(desiredLbRuleIds, currentLbRuleIds).List()
	detachLbRuleIds := autoscaling_common.CalculateSetDifference(currentLbRuleIds, desiredLbRuleIds).List()

	request := autoscaling2.AsgLoadBalancersUpdateRequest{
		AttachLbRuleIds: common.ToStringList(attachLbRuleIds),
		DetachLbRuleIds: common.ToStringList(detachLbRuleIds),
	}

	return request
}

func extractLbRuleIdsFromResponse(response loadbalancer2.ListResponseLbServiceForAsgResponse) *schema.Set {
	currentLbRuleIds := schema.NewSet(schema.HashString, nil)
	for _, lbService := range response.Contents {
		for _, lbRule := range lbService.LbRules {
			currentLbRuleIds.Add(lbRule.LbRuleId)
		}
	}
	return currentLbRuleIds
}

func updateAutoScalingGroupLoadBalancer(ctx context.Context, rd *schema.ResourceData, meta interface{}, request autoscaling2.AsgLoadBalancersUpdateRequest) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if request.AttachLbRuleIds == nil && request.DetachLbRuleIds == nil {
		return diag.Diagnostics{}
	}

	_, _, err := inst.Client.AutoScaling.UpdateAutoScalingGroupLoadBalancer(ctx, rd.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForAutoScalingGroupStatus(ctx, inst.Client, rd.Id(), []string{"Attach to LB", "Detach from LB"}, []string{"In Service"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceAutoScalingGroupLoadBalancerRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diagnostics, response := getAutoScalingGroupLoadBalancer(ctx, rd, meta)
	if len(diagnostics) > 0 {
		return diagnostics
	}
	currentLbRuleIds := extractLbRuleIdsFromResponse(response)

	rd.Set("lb_rule_ids", currentLbRuleIds)

	return nil
}

func resourceAutoScalingGroupLoadBalancerUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if rd.HasChanges("lb_rule_ids") {
		diagnostics, response := getAutoScalingGroupLoadBalancer(ctx, rd, meta)
		if len(diagnostics) > 0 {
			return diagnostics
		}

		request := generateAsgLoadBalancersUpdateRequest(rd, response)
		diagnostics = updateAutoScalingGroupLoadBalancer(ctx, rd, meta, request)
		if len(diagnostics) > 0 {
			return diagnostics
		}
	}

	return resourceAutoScalingGroupLoadBalancerRead(ctx, rd, meta)
}

func resourceAutoScalingGroupLoadBalancerDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diagnostics, response := getAutoScalingGroupLoadBalancer(ctx, rd, meta)
	if len(diagnostics) > 0 {
		return diagnostics
	}
	currentLbRuleIds := extractLbRuleIdsFromResponse(response)

	request := autoscaling2.AsgLoadBalancersUpdateRequest{
		AttachLbRuleIds: []string{},
		DetachLbRuleIds: common.ToStringList(currentLbRuleIds.List()),
	}

	diagnostics = updateAutoScalingGroupLoadBalancer(ctx, rd, meta, request)
	if len(diagnostics) > 0 {
		return diagnostics
	}

	return nil
}
