package routing

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/routing"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_vpc_routing_rules", DataSourceVpcRoutingRule())
	scp.RegisterDataSource("scp_vpc_routing_routes", DataSourceVpcRoutingRoute())
	scp.RegisterDataSource("scp_vpc_routing_tables", DataSourceVpcRoutingTable())
}

func DataSourceVpcRoutingTable() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVpcRoutingTableRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("RoutingTableId"):   {Type: schema.TypeString, Optional: true, Description: "Routing Table Id"},
			common.ToSnakeCase("RoutingTableName"): {Type: schema.TypeString, Optional: true, Description: "Routing Table Name"},
			common.ToSnakeCase("VpcId"):            {Type: schema.TypeString, Optional: true, Description: "VPC Id"},
			common.ToSnakeCase("CreatedBy"):        {Type: schema.TypeString, Optional: true, Description: "Created By"},
			"contents":                             {Type: schema.TypeList, Optional: true, Description: "VPC Routing Table list", Elem: RoutingTableElem()},
			"total_count":                          {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides a VPC Routing Table resource.",
	}
}

func RoutingTableElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("RoutingRuleCount"): {Type: schema.TypeInt, Computed: true, Description: "Routing Rule Count"},
			common.ToSnakeCase("RoutingTableId"):   {Type: schema.TypeString, Computed: true, Description: "Routing Table ID"},
			common.ToSnakeCase("RoutingTableName"): {Type: schema.TypeString, Computed: true, Description: "Routing Table name"},
			common.ToSnakeCase("RoutingTableType"): {Type: schema.TypeString, Computed: true, Description: "Routing Table Type"},
			common.ToSnakeCase("ServiceZoneId"):    {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			common.ToSnakeCase("T1RouterId"):       {Type: schema.TypeString, Computed: true, Description: "t1 Router ID"},
			common.ToSnakeCase("VpcId"):            {Type: schema.TypeString, Computed: true, Description: "VPC ID"},
		},
	}
}

func resourceVpcRoutingTableRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	request := routing.ListVpcRoutingTableRequest{
		RoutingTableId:   rd.Get(common.ToSnakeCase("RoutingTableId")).(string),
		RoutingTableName: rd.Get(common.ToSnakeCase("RoutingTableName")).(string),
		VpcId:            rd.Get(common.ToSnakeCase("VpcId")).(string),
		CreatedBy:        rd.Get(common.ToSnakeCase("CreatedBy")).(string),
	}

	responses, err := inst.Client.Routing.GetVpcRoutingTableListV2(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func DataSourceVpcRoutingRoute() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVpcRoutingRouteRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("RoutingTableId"): {Type: schema.TypeString, Required: true, Description: "Routing Table Id"},
			"contents":                           {Type: schema.TypeList, Optional: true, Description: "VPC Routing Route list", Elem: RoutingRouteElem()},
			"total_count":                        {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides a VPC Routing Route resource.",
	}
}

func RoutingRouteElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("SourceServiceInterfaceId"):   {Type: schema.TypeString, Computed: true, Description: "Source Interface Id"},
			common.ToSnakeCase("SourceServiceInterfaceName"): {Type: schema.TypeString, Computed: true, Description: "Source Interface Name"},
		},
	}
}

func resourceVpcRoutingRouteRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	routingTableId := rd.Get(common.ToSnakeCase("RoutingTableId")).(string)

	responses, err := inst.Client.Routing.GetVpcRoutingRulesRoute(ctx, routingTableId)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func DataSourceVpcRoutingRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVpcRoutingRuleRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("RoutingTableId"):           {Type: schema.TypeString, Required: true, Description: "Routing Table Id"},
			common.ToSnakeCase("DestinationNetworkCidr"):   {Type: schema.TypeString, Optional: true, Description: "Destination Network Cidr"},
			common.ToSnakeCase("Editable"):                 {Type: schema.TypeString, Optional: true, Description: "is Editable (true | false)"},
			common.ToSnakeCase("RoutingRuleId"):            {Type: schema.TypeString, Optional: true, Description: "Routing Rule Id"},
			common.ToSnakeCase("SourceServiceInterfaceId"): {Type: schema.TypeString, Optional: true, Description: "Source Interface Id"},
			"contents":    {Type: schema.TypeList, Optional: true, Description: "VPC Routing Rule list", Elem: RoutingRuleElem()},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides a VPC Routing Table resource.",
	}
}

func RoutingRuleElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("BlockId"):                    {Type: schema.TypeString, Computed: true, Description: "Block Id"},
			common.ToSnakeCase("DestinationNetworkCidr"):     {Type: schema.TypeString, Computed: true, Description: "Destination Network Cidr"},
			common.ToSnakeCase("Editable"):                   {Type: schema.TypeBool, Computed: true, Description: "is Editable"},
			common.ToSnakeCase("RoutingRuleId"):              {Type: schema.TypeString, Computed: true, Description: "Routing Rule Id"},
			common.ToSnakeCase("RoutingRuleState"):           {Type: schema.TypeString, Computed: true, Description: "Routing Rule State"},
			common.ToSnakeCase("ServiceZoneId"):              {Type: schema.TypeString, Computed: true, Description: "Service Zone Id"},
			common.ToSnakeCase("SourceServiceInterfaceId"):   {Type: schema.TypeString, Computed: true, Description: "Source Interface Id"},
			common.ToSnakeCase("SourceServiceInterfaceName"): {Type: schema.TypeString, Computed: true, Description: "Source Interface Name"},
			common.ToSnakeCase("SourceVpcId"):                {Type: schema.TypeString, Computed: true, Description: "Source VPC Id"},
			common.ToSnakeCase("CreatedBy"):                  {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):                  {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			common.ToSnakeCase("ModifiedBy"):                 {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):                 {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
			common.ToSnakeCase("ProjectId"):                  {Type: schema.TypeString, Computed: true, Description: "Project Id"},
		},
	}
}

func resourceVpcRoutingRuleRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	routingTableId := rd.Get(common.ToSnakeCase("RoutingTableId")).(string)
	request := routing.ListVpcRoutingRulesRequest{
		DestinationNetworkCidr:   rd.Get(common.ToSnakeCase("DestinationNetworkCidr")).(string),
		Editable:                 rd.Get(common.ToSnakeCase("Editable")).(string),
		RoutingRuleId:            rd.Get(common.ToSnakeCase("RoutingRuleId")).(string),
		SourceServiceInterfaceId: rd.Get(common.ToSnakeCase("SourceServiceInterfaceId")).(string),
	}

	responses, err := inst.Client.Routing.GetVpcRoutingRulesList(ctx, routingTableId, request)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}
