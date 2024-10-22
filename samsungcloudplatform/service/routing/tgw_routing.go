package routing

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client/routing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	samsungcloudplatform.RegisterResource("samsungcloudplatform_transit_gateway_routing", ResourceTGWRouting())
}

func ResourceTGWRouting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTGWRoutingCreate,
		ReadContext:   resourceTGWRoutingRead,
		DeleteContext: resourceTGWRoutingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"routing_table_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Routing Table ID for Transit Gateway Connection",
				ValidateFunc: validation.StringLenBetween(3, 60),
			},
			"destination_network_cidr": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network CIDR",
			},
			"source_service_interface_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Source Interface ID",
				ValidateFunc: validation.StringLenBetween(3, 255),
			},
			"source_service_interface_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Source Interface Name",
				ValidateFunc: validation.StringLenBetween(3, 255),
			},
		},
		Description: "Provides a Transit Gateway Connection Routing Rule Resources",
	}
}

func resourceTGWRoutingCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

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

	// Routing Table 존재여부 확인
	_, err := inst.Client.Routing.GetTgwRoutingTableDetail(ctx, routingTableId)
	if err != nil {
		return diag.FromErr(err)
	}

	//Rule 생성
	err = inst.Client.Routing.CreateTgwRoutingRules(ctx, routingTableId, request)
	if err != nil {
		return diag.FromErr(err)
	}
	//Rule 생성 완료 체크
	err = waitTGWRoutingRuleCreating(ctx, inst.Client, routingTableId, destinationNetworkCidr)
	if err != nil {
		return diag.FromErr(err)
	}

	info, _, err := inst.Client.Routing.GetTgwRoutingRuleByCidr(ctx, routingTableId, destinationNetworkCidr)
	if err != nil {
		return diag.FromErr(err)
	}
	rd.SetId(inst.Client.Routing.MergeRoutingRuleId(routingTableId, info.RoutingRuleId))

	return resourceTGWRoutingRead(ctx, rd, meta)
}

func resourceTGWRoutingRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	routingTableId, routingRuleId := inst.Client.Routing.SplitRoutingRuleId(rd.Id())

	ruleInfo, _, err := inst.Client.Routing.GetTgwRoutingRuleById(ctx, routingTableId, routingRuleId)
	if err != nil {
		diag.FromErr(err)
	}

	rd.Set("routing_table_id", routingTableId)
	rd.Set("destination_network_cidr", ruleInfo.DestinationNetworkCidr)
	rd.Set("source_service_interface_id", ruleInfo.SourceServiceInterfaceId)
	rd.Set("source_service_interface_name", ruleInfo.SourceServiceInterfaceName)

	return nil
}

func resourceTGWRoutingDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	routingTableId, routingRuleId := inst.Client.Routing.SplitRoutingRuleId(rd.Id())

	//Rule 삭제
	err := inst.Client.Routing.DeleteTgwRoutingRules(ctx, routingTableId, routingRuleId)
	if err != nil {
		diag.FromErr(err)
	}
	err = waitTGWRoutingRuleDeleting(ctx, inst.Client, routingTableId, routingRuleId)

	if err != nil {
		diag.FromErr(err)
	}

	return nil
}

func waitTGWRoutingRuleCreating(ctx context.Context, scpClient *client.SCPClient, routingTableId string, destinationNetworkCidr string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		return scpClient.Routing.GetTgwRoutingRuleByCidr(ctx, routingTableId, destinationNetworkCidr)
	})
}

func waitTGWRoutingRuleDeleting(ctx context.Context, scpClient *client.SCPClient, routingTableId string, routingRuleId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		return scpClient.Routing.GetTgwRoutingRuleById(ctx, routingTableId, routingRuleId)
	})

}
