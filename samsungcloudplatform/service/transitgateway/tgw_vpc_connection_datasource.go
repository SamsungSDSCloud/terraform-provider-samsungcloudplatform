package transitgateway

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	transitgateway2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/transit-gateway2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_transit_gateway_connections", DataSourceTransitGatewayConnections())
}

func DataSourceTransitGatewayConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: tgwConnDataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("TransitGatewayConnectionName"): {Type: schema.TypeString, Optional: true, Description: "TGW VPC Connection Name"},
			common.ToSnakeCase("RequesterTransitGatewayId"):    {Type: schema.TypeString, Optional: true, Description: "Requester TGW ID"},
			common.ToSnakeCase("ApproverVpcId"):                {Type: schema.TypeString, Optional: true, Description: "Approver VPC ID"},
			common.ToSnakeCase("CreatedBy"):                    {Type: schema.TypeString, Optional: true, Description: "User ID who create the resources"},
			common.ToSnakeCase("Page"):                         {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page number"},
			common.ToSnakeCase("Size"):                         {Type: schema.TypeInt, Optional: true, Default: 20, Description: "List size per a page"},
			"contents":                                         {Type: schema.TypeList, Optional: true, Description: "List of TGW VPC Connections", Elem: tgwConnDataSourceElem()},
			"total_count":                                      {Type: schema.TypeInt, Computed: true, Description: "Total Count of TGW VPC Connections"},
		},
		Description: "provides Lists of TGW-VPC Connection",
	}
}

func tgwConnDataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.TransitGateway.GetTransitGatewayConnectionList(ctx, &transitgateway2.TransitGatewayConnectionOpenApiControllerApiListTransitGatewayConnectionsOpts{
		TransitGatewayConnectionName: optional.NewString(rd.Get("transit_gateway_connection_name").(string)),
		RequesterTransitGatewayId:    optional.NewString(rd.Get("requester_transit_gateway_id").(string)),
		ApproverVpcId:                optional.NewString(rd.Get("approver_vpc_id").(string)),
		CreatedBy:                    optional.NewString(rd.Get("created_by").(string)),
		Page:                         optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:                         optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:                         optional.Interface{},
	})
	if err != nil {
		diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func tgwConnDataSourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ProjectId"):                           {Type: schema.TypeString, Computed: true, Description: "ProjectId"},
			common.ToSnakeCase("ServiceZoneId"):                       {Type: schema.TypeString, Computed: true, Description: "Service Zone Id"},
			common.ToSnakeCase("CreatedBy"):                           {Type: schema.TypeString, Computed: true, Description: "User ID who create the resources"},
			common.ToSnakeCase("CreatedDt"):                           {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			common.ToSnakeCase("ModifiedBy"):                          {Type: schema.TypeString, Computed: true, Description: "User ID who modified the resources last"},
			common.ToSnakeCase("ModifiedDt"):                          {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			common.ToSnakeCase("TransitGatewayConnectionId"):          {Type: schema.TypeString, Computed: true, Description: "TGW - VPC Connection ID"},
			common.ToSnakeCase("TransitGatewayConnectionName"):        {Type: schema.TypeString, Computed: true, Description: "TGW - VPC Connection Name"},
			common.ToSnakeCase("TransitGatewayConnectionState"):       {Type: schema.TypeString, Computed: true, Description: "TGW - VPC Connection State"},
			common.ToSnakeCase("TransitGatewayConnectionType"):        {Type: schema.TypeString, Computed: true, Description: "TGW - VPC Connection Type"},
			common.ToSnakeCase("TransitGatewayConnectionDescription"): {Type: schema.TypeString, Computed: true, Description: "TGW - VPC Connection Description"},
			common.ToSnakeCase("RequesterProjectId"):                  {Type: schema.TypeString, Computed: true, Description: "Requester TGW's ProjectId"},
			common.ToSnakeCase("RequesterTransitGatewayId"):           {Type: schema.TypeString, Computed: true, Description: "Requester TGW ID"},
			common.ToSnakeCase("ApproverProjectId"):                   {Type: schema.TypeString, Computed: true, Description: "Approver VPC's ProjectId"},
			common.ToSnakeCase("ApproverVpcId"):                       {Type: schema.TypeString, Computed: true, Description: "Approver ProjectId"},
		},
	}

}
