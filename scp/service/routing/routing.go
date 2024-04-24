package routing

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/routing"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	scp.RegisterResource("scp_vpc_routing", ResourceVpcRouting())
}

func ResourceVpcRouting() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceVpcRoutingCreate,
		ReadContext:   resourceVpcRoutingRead,
		UpdateContext: resourceVpcRoutingUpdate,
		DeleteContext: resourceVpcRoutingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"routing_table_id": {
				Type:         schema.TypeString,
				Required:     true, //필수 작성
				Description:  "Routing Table id",
				ValidateFunc: validation.StringLenBetween(3, 100),
			},
			"destination_network_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Network CIDR",
				ValidateFunc: validation.IsCIDRNetwork(24, 27),
			},
			"source_service_interface_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Source Interface Id",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"source_service_interface_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Source Interface Name",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
		},
		Description: "Provides a VPC Routing Rule resource.",
	}
}

func resourceVpcRoutingCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	routingTableId := rd.Get("routing_table_id").(string)
	destinationNetworkCidr := rd.Get("destination_network_cidr").(string)
	sourceServiceInterfaceId := rd.Get("source_service_interface_id").(string)
	sourceServiceInterfaceName := rd.Get("source_service_interface_name").(string)

	request := routing.CreateRoutingRulesRequest{}
	request.RoutingRules = append(request.RoutingRules, routing.RoutingRule{
		DestinationNetworkCidr:     destinationNetworkCidr,
		SourceServiceInterfaceId:   sourceServiceInterfaceId,
		SourceServiceInterfaceName: sourceServiceInterfaceName,
	})

	inst := meta.(*client.Instance)

	// check if routingTableId is valid
	_, err := inst.Client.Routing.GetVpcRoutingTableDetail(ctx, routingTableId)
	if err != nil {
		return diag.FromErr(err)
	}

	// duplication check
	_, err = inst.Client.Routing.CheckDuplicationRoutingRule(ctx, routingTableId, destinationNetworkCidr)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Try create vpc dns zone : "+routingTableId+", "+destinationNetworkCidr)

	err = inst.Client.Routing.CreateRoutingRules(ctx, routingTableId, request)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitRoutingRuleCreating(ctx, inst.Client, routingTableId, destinationNetworkCidr)
	if err != nil {
		return diag.FromErr(err)
	}

	info, _, err := inst.Client.Routing.GetVpcRoutingRulesByCidr(ctx, routingTableId, destinationNetworkCidr)
	rd.SetId(inst.Client.Routing.MergeRoutingRuleId(routingTableId, info.RoutingRuleId))

	// Refresh
	return resourceVpcRoutingRead(ctx, rd, meta)
}

func resourceVpcRoutingRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	ruleInfo, _, err := inst.Client.Routing.GetVpcRoutingRulesById(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}
	routingTableId, _ := inst.Client.Routing.SplitRoutingRuleId(rd.Id())

	rd.Set("routing_table_id", routingTableId)
	rd.Set("destination_network_cidr", ruleInfo.DestinationNetworkCidr)
	rd.Set("source_service_interface_id", ruleInfo.SourceServiceInterfaceId)
	rd.Set("source_service_interface_name", ruleInfo.SourceServiceInterfaceName)

	return nil
}

func resourceVpcRoutingUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("Update function is not implemented")
}

func resourceVpcRoutingDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	routingTableId, routingRuleId := inst.Client.Routing.SplitRoutingRuleId(rd.Id())
	err := inst.Client.Routing.DeleteRoutingRules(ctx, routingTableId, routingRuleId)
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitRoutingRuleDeleting(ctx, inst.Client, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitRoutingRuleCreating(ctx context.Context, scpClient *client.SCPClient, routingTableId string, destinationNetworkCidr string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		return scpClient.Routing.GetVpcRoutingRulesByCidr(ctx, routingTableId, destinationNetworkCidr)
	})
}

func waitRoutingRuleDeleting(ctx context.Context, scpClient *client.SCPClient, ruleId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		return scpClient.Routing.GetVpcRoutingRulesById(ctx, ruleId)
	})
}
