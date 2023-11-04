package routing

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/routing"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	scp.RegisterResource("scp_direct_connect_routing", ResourceDCRouting())
}

func ResourceDCRouting() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceDCRoutingCreate,
		ReadContext:   resourceDCRoutingRead,
		// UpdateContext: resourceDCRoutingUpdate,
		DeleteContext: resourceDCRoutingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"routing_table_id": {
				Type:         schema.TypeString,
				Required:     true, //필수 작성
				ForceNew:     true,
				Description:  "Routing Table id",
				ValidateFunc: validation.StringLenBetween(3, 100),
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
				Description:  "Source Interface Id",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"source_service_interface_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Source Interface Name",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
		},
		Description: "Provides a DirectConnect Routing Rule resource.",
	}
}

func resourceDCRoutingCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

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

	// // check if routingTableId is valid
	// _, err := inst.Client.Routing.GetDCRoutingTableDetail(ctx, routingTableId)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// duplication check
	if _, err := inst.Client.Routing.CheckDCDuplicationRoutingRule(ctx, routingTableId, destinationNetworkCidr); err != nil {
		return diag.FromErr(err)
	}

	if err := inst.Client.Routing.CreateDCRoutingRules(ctx, routingTableId, request); err != nil {
		return diag.FromErr(err)
	}

	if err := waitDCRoutingRuleCreating(ctx, inst.Client, routingTableId, destinationNetworkCidr); err != nil {
		return diag.FromErr(err)
	}

	info, _, err := inst.Client.Routing.GetDCRoutingRulesByCidr(ctx, routingTableId, destinationNetworkCidr)
	if err != nil {
		return diag.FromErr(err)
	}
	rd.SetId(inst.Client.Routing.MergeRoutingRuleId(routingTableId, info.RoutingRuleId))

	// Refresh
	return resourceDCRoutingRead(ctx, rd, meta)
}

func resourceDCRoutingRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	ruleInfo, _, err := inst.Client.Routing.GetDCRoutingRulesById(ctx, rd.Id())
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

func resourceDCRoutingDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	routingTableId, routingRuleId := inst.Client.Routing.SplitRoutingRuleId(rd.Id())
	err := inst.Client.Routing.DeleteDCRoutingRules(ctx, routingTableId, routingRuleId)
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitDCRoutingRuleDeleting(ctx, inst.Client, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitDCRoutingRuleCreating(ctx context.Context, scpClient *client.SCPClient, routingTableId string, destinationNetworkCidr string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		return scpClient.Routing.GetDCRoutingRulesByCidr(ctx, routingTableId, destinationNetworkCidr)
	})
}

func waitDCRoutingRuleDeleting(ctx context.Context, scpClient *client.SCPClient, ruleId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		return scpClient.Routing.GetDCRoutingRulesById(ctx, ruleId)
	})
}
