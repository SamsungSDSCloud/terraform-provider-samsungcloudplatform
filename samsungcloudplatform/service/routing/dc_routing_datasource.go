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
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_direct_connect_routing_tables", DataSourceDCRoutingTable())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_direct_connect_routing_routes", DataSourceDCRoutingRoute())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_direct_connect_routing_rules", DataSourceDCRoutingRule())
}

func DataSourceDCRoutingTable() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDCRoutingTableRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("DirectConnectConnectionId"): {Type: schema.TypeString, Optional: true, Description: "DirectConnect Connention Id"},
			common.ToSnakeCase("RoutingTableId"):            {Type: schema.TypeString, Optional: true, Description: "Routing Table Id"},
			common.ToSnakeCase("RoutingTableName"):          {Type: schema.TypeString, Optional: true, Description: "Routing Table Name"},
			common.ToSnakeCase("CreatedBy"):                 {Type: schema.TypeString, Optional: true, Description: "Created By"},
			"contents":                                      {Type: schema.TypeList, Optional: true, Description: "DirectConnect Routing Table list", Elem: DCRoutingTableElem()},
			"total_count":                                   {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides a VPC Routing Table resource.",
	}
}

func DCRoutingTableElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("DirectConnectConnectionId"): {Type: schema.TypeString, Computed: true, Description: "DirectConnect Connention Id"},
			common.ToSnakeCase("RoutingRuleCount"):          {Type: schema.TypeInt, Computed: true, Description: "Routing Rule Count"},
			common.ToSnakeCase("RoutingTableId"):            {Type: schema.TypeString, Computed: true, Description: "Routing Table ID"},
			common.ToSnakeCase("RoutingTableName"):          {Type: schema.TypeString, Computed: true, Description: "Routing Table name"},
			common.ToSnakeCase("RoutingTableType"):          {Type: schema.TypeString, Computed: true, Description: "Routing Table Type"},
			common.ToSnakeCase("ServiceZoneId"):             {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			common.ToSnakeCase("T1RouterId"):                {Type: schema.TypeString, Computed: true, Description: "t1 Router ID"},
		},
	}
}

func resourceDCRoutingTableRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	routingTableId := rd.Get(common.ToSnakeCase("RoutingTableId")).(string)
	routingTableName := rd.Get(common.ToSnakeCase("RoutingTableName")).(string)
	directConnectConnectionId := rd.Get(common.ToSnakeCase("DirectConnectConnectionId")).(string)
	createdBy := rd.Get(common.ToSnakeCase("CreatedBy")).(string)

	responses, err := inst.Client.Routing.GetDCRoutingTableList(ctx, routingTableId, routingTableName, directConnectConnectionId, createdBy)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func DataSourceDCRoutingRoute() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDCRoutingRouteRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("RoutingTableId"): {Type: schema.TypeString, Required: true, Description: "Routing Table Id"},
			"contents":                           {Type: schema.TypeList, Optional: true, Description: "DirectConnect Routing Route list", Elem: DCRoutingRouteElem()},
			"total_count":                        {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
			"filter":                             common.DatasourceFilter(),
		},
		Description: "Provides a Direct Connect Routing Route resource.",
	}
}

func DCRoutingRouteElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("SourceServiceInterfaceId"):   {Type: schema.TypeString, Computed: true, Description: "Source Interface Id"},
			common.ToSnakeCase("SourceServiceInterfaceName"): {Type: schema.TypeString, Computed: true, Description: "Source Interface Name"},
		},
	}
}

func resourceDCRoutingRouteRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	routingTableId := rd.Get(common.ToSnakeCase("RoutingTableId")).(string)

	responses, err := inst.Client.Routing.GetDCRoutingRulesRoute(ctx, routingTableId)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DataSourceDCRoutingRoute().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func DataSourceDCRoutingRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDCRoutingRuleRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("RoutingTableId"):           {Type: schema.TypeString, Required: true, Description: "Routing Table Id"},
			common.ToSnakeCase("DestinationNetworkCidr"):   {Type: schema.TypeString, Optional: true, Description: "Destination Network Cidr"},
			common.ToSnakeCase("Editable"):                 {Type: schema.TypeString, Optional: true, Description: "is Editable (true | false)"},
			common.ToSnakeCase("RoutingRuleId"):            {Type: schema.TypeString, Optional: true, Description: "Routing Rule Id"},
			common.ToSnakeCase("SourceServiceInterfaceId"): {Type: schema.TypeString, Optional: true, Description: "Source Interface Id"},
			"contents":    {Type: schema.TypeList, Optional: true, Description: "DirectConnect Routing Rule list", Elem: DCRoutingRuleElem()},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides a DirectConnect Routing Table Rule resource.",
	}
}

func DCRoutingRuleElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("DestinationNetworkCidr"):          {Type: schema.TypeString, Computed: true, Description: "Destination Network Cidr"},
			common.ToSnakeCase("Editable"):                        {Type: schema.TypeBool, Computed: true, Description: "is Editable"},
			common.ToSnakeCase("RoutingRuleId"):                   {Type: schema.TypeString, Computed: true, Description: "Routing Rule Id"},
			common.ToSnakeCase("RoutingRuleState"):                {Type: schema.TypeString, Computed: true, Description: "Routing Rule State"},
			common.ToSnakeCase("sourceDirectConnectConnectionId"): {Type: schema.TypeString, Computed: true, Description: "Source DirectConnect Connection Id"},
			common.ToSnakeCase("SourceServiceInterfaceId"):        {Type: schema.TypeString, Computed: true, Description: "Source Interface Id"},
			common.ToSnakeCase("SourceServiceInterfaceName"):      {Type: schema.TypeString, Computed: true, Description: "Source Interface Name"},
			common.ToSnakeCase("CreatedBy"):                       {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):                       {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			common.ToSnakeCase("ModifiedBy"):                      {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):                      {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
			common.ToSnakeCase("ProjectId"):                       {Type: schema.TypeString, Computed: true, Description: "Project Id"},
		},
	}
}

func resourceDCRoutingRuleRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	routingTableId := rd.Get(common.ToSnakeCase("RoutingTableId")).(string)
	request := routing.ListVpcRoutingRulesRequest{
		DestinationNetworkCidr:   rd.Get(common.ToSnakeCase("DestinationNetworkCidr")).(string),
		Editable:                 rd.Get(common.ToSnakeCase("Editable")).(string),
		RoutingRuleId:            rd.Get(common.ToSnakeCase("RoutingRuleId")).(string),
		SourceServiceInterfaceId: rd.Get(common.ToSnakeCase("SourceServiceInterfaceId")).(string),
	}

	responses, err := inst.Client.Routing.GetDCRoutingRulesList(ctx, routingTableId, request)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}
