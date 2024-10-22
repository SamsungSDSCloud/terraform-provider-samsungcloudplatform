package routing

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client/routing"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_transit_gateway_routing_tables", DataSourceTGWRoutingTable())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_transit_gateway_routing_routes", DataSourceTGWRoutingRoute())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_transit_gateway_routing_rules", DataSourceTGWRoutingRule())
}

func DataSourceTGWRoutingTable() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceTgwRoutingTableRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("TransitGatewayConnectionId"): {Type: schema.TypeString, Optional: true, Description: "Transit Gateway Connection ID"},
			common.ToSnakeCase("RoutingTableId"):             {Type: schema.TypeString, Optional: true, Description: "Routing Table ID"},
			common.ToSnakeCase("RoutingTableName"):           {Type: schema.TypeString, Optional: true, Description: "Routing Table Name"},
			common.ToSnakeCase("CreatedBy"):                  {Type: schema.TypeString, Optional: true, Description: "Created By "},
			"contents":                                       {Type: schema.TypeList, Optional: true, Description: "Transit Gateway Connection's Routing Table List", Elem: TGWRoutingTableElem()},
			"total_counts":                                   {Type: schema.TypeInt, Optional: true, Description: "Total List size"},
		},
		Description: "Provides a TGW Routing Table List Resources.",
	}
}

func TGWRoutingTableElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("TransitGatewayConnectionId"): {Type: schema.TypeString, Optional: true, Description: "Transit Gateway Connection ID"},
			common.ToSnakeCase("RoutingRuleCount"):           {Type: schema.TypeInt, Computed: true, Description: "Routing Rule Count"},
			common.ToSnakeCase("RoutingTableId"):             {Type: schema.TypeString, Computed: true, Description: "Routing Table ID"},
			common.ToSnakeCase("RoutingTableName"):           {Type: schema.TypeString, Computed: true, Description: "Routing Table name"},
			common.ToSnakeCase("RoutingTableType"):           {Type: schema.TypeString, Computed: true, Description: "Routing Table Type"},
			common.ToSnakeCase("ServiceZoneId"):              {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			common.ToSnakeCase("T1RouterId"):                 {Type: schema.TypeString, Computed: true, Description: "t1 Router ID"},
		},
	}
}

func resourceTgwRoutingTableRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	request := routing.ListTgwRoutingTableRequest{
		RoutingTableId:             rd.Get(common.ToSnakeCase("RoutingTableId")).(string),
		RoutingTableName:           rd.Get(common.ToSnakeCase("RoutingTableName")).(string),
		TransitGatewayConnectionId: rd.Get(common.ToSnakeCase("TransitGatewayConnectionId")).(string),
		CreatedBy:                  rd.Get(common.ToSnakeCase("CreatedBy")).(string),
	}

	tableList, err := inst.Client.Routing.GetTgwRoutingTableList(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	TGWRoutingTableElem()

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", common.ConvertStructToMaps(tableList.Contents))
	rd.Set("total_counts", tableList.TotalCount)

	return nil

}
func DataSourceTGWRoutingRoute() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceTgwRoutingRouteRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("RoutingTableId"): {Type: schema.TypeString, Optional: true, Description: "Routing Table ID"},
			"contents":                           {Type: schema.TypeList, Optional: true, Description: "Transit Gateway Connection's Route List", Elem: TGWRoutingRouteElem()},
			"total_counts":                       {Type: schema.TypeInt, Optional: true, Description: "Total List size"},
		},
		Description: "Provides a TGW Routing Route List Resources.",
	}
}

func TGWRoutingRouteElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("SourceServiceInterfaceId"):   {Type: schema.TypeString, Computed: true, Description: "Source Interface Id"},
			common.ToSnakeCase("SourceServiceInterfaceName"): {Type: schema.TypeString, Computed: true, Description: "Source Interface Name"},
		},
	}
}

func resourceTgwRoutingRouteRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	routes, err := inst.Client.Routing.GetTgwRoutingRoutes(ctx, rd.Get(common.ToSnakeCase("RoutingTableId")).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", common.ConvertStructToMaps(routes.Contents))
	rd.Set("total_counts", routes.TotalCount)

	return nil
}
func DataSourceTGWRoutingRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceTgwRoutingRuleRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("RoutingTableId"):           {Type: schema.TypeString, Optional: true, Description: "Routing Table ID"},
			common.ToSnakeCase("DestinationNetworkCidr"):   {Type: schema.TypeString, Optional: true, Description: "Destination Network Cidr"},
			common.ToSnakeCase("Editable"):                 {Type: schema.TypeString, Optional: true, Description: "is Editable (true | false)"},
			common.ToSnakeCase("RoutingRuleId"):            {Type: schema.TypeString, Optional: true, Description: "Routing Rule Id"},
			common.ToSnakeCase("SourceServiceInterfaceId"): {Type: schema.TypeString, Optional: true, Description: "Source Interface Id"},
			"contents":     {Type: schema.TypeList, Optional: true, Description: "Transit Gateway Connection's Route List", Elem: TGWRoutingListElem()},
			"total_counts": {Type: schema.TypeInt, Optional: true, Description: "Total List size"},
		},
	}
}

func TGWRoutingListElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("DestinationNetworkCidr"):           {Type: schema.TypeString, Computed: true, Description: "Destination Network Cidr"},
			common.ToSnakeCase("Editable"):                         {Type: schema.TypeBool, Computed: true, Description: "is Editable"},
			common.ToSnakeCase("RoutingRuleId"):                    {Type: schema.TypeString, Computed: true, Description: "Routing Rule Id"},
			common.ToSnakeCase("RoutingRuleState"):                 {Type: schema.TypeString, Computed: true, Description: "Routing Rule State"},
			common.ToSnakeCase("SourceTransitGatewayConnectionId"): {Type: schema.TypeString, Computed: true, Description: "Source Transit Gateway Connection Id"},
			common.ToSnakeCase("SourceServiceInterfaceId"):         {Type: schema.TypeString, Computed: true, Description: "Source Interface Id"},
			common.ToSnakeCase("SourceServiceInterfaceName"):       {Type: schema.TypeString, Computed: true, Description: "Source Interface Name"},
			common.ToSnakeCase("CreatedBy"):                        {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):                        {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			common.ToSnakeCase("ModifiedBy"):                       {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):                       {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
			common.ToSnakeCase("ProjectId"):                        {Type: schema.TypeString, Computed: true, Description: "Project Id"},
		},
	}
}

func resourceTgwRoutingRuleRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	request := routing.ListTgwRoutingRuleRequest{
		DestinationNetworkCidr:   rd.Get(common.ToSnakeCase("DestinationNetworkCidr")).(string),
		Editable:                 rd.Get(common.ToSnakeCase("Editable")).(string),
		RoutingRuleId:            rd.Get(common.ToSnakeCase("RoutingRuleId")).(string),
		SourceServiceInterfaceId: rd.Get(common.ToSnakeCase("SourceServiceInterfaceId")).(string),
	}

	ruleList, err := inst.Client.Routing.GetTgwRoutingRuleList(ctx, rd.Get(common.ToSnakeCase("RoutingTableId")).(string), request)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", common.ConvertStructToMaps(ruleList.Contents))
	rd.Set("total_counts", ruleList.TotalCount)

	return nil
}
