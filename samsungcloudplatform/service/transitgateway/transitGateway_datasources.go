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
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_transit_gateways", DataSourceTransitGateways())
}

func DataSourceTransitGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("TransitGatewayId"):   {Type: schema.TypeString, Optional: true, Description: "Transit Gateway ID"},
			common.ToSnakeCase("TransitGatewayName"): {Type: schema.TypeString, Optional: true, Description: "Transit Gateway Name"},
			common.ToSnakeCase("CreatedBy"):          {Type: schema.TypeString, Optional: true, Description: "User ID who create the resources"},
			common.ToSnakeCase("Page"):               {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page number"},
			common.ToSnakeCase("Size"):               {Type: schema.TypeInt, Optional: true, Default: 20, Description: "List size per a page"},
			"contents":                               {Type: schema.TypeList, Optional: true, Description: "List of Transit Gateways", Elem: dataSourceElem()},
			"total_count":                            {Type: schema.TypeInt, Computed: true, Description: "Total Count of Transit Gateways"},
		},
		Description: "Provides List of Transit Gateway.",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.TransitGateway.GetTransitGatewayList(ctx, &transitgateway2.TransitGatewayOpenApiControllerApiListTransitGatewaysOpts{
		TransitGatewayId:   optional.NewString(rd.Get("transit_gateway_id").(string)),
		TransitGatewayName: optional.NewString(rd.Get("transit_gateway_name").(string)),
		CreatedBy:          optional.NewString(rd.Get("created_by").(string)),
		Page:               optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:               optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:               optional.Interface{},
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

func dataSourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ProjectId"):                 {Type: schema.TypeString, Computed: true, Description: "Project id"},
			common.ToSnakeCase("ServiceZoneId"):             {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			common.ToSnakeCase("CreatedBy"):                 {Type: schema.TypeString, Computed: true, Description: "User ID who create the resources"},
			common.ToSnakeCase("CreatedDt"):                 {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			common.ToSnakeCase("ModifiedBy"):                {Type: schema.TypeString, Computed: true, Description: "User ID who modified the resources last"},
			common.ToSnakeCase("ModifiedDt"):                {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			common.ToSnakeCase("TransitGatewayId"):          {Type: schema.TypeString, Computed: true, Description: "TransitGateway ID"},
			common.ToSnakeCase("TransitGatewayName"):        {Type: schema.TypeString, Computed: true, Description: "TransitGateway Name"},
			common.ToSnakeCase("TransitGatewayState"):       {Type: schema.TypeString, Computed: true, Description: "TransitGateway State"},
			common.ToSnakeCase("TransitGatewayDescription"): {Type: schema.TypeString, Computed: true, Description: "TransitGateway Description"},
			common.ToSnakeCase("UplinkEnabled"):             {Type: schema.TypeBool, Computed: true, Description: "UplinkEnabled"},
			common.ToSnakeCase("BandwidthGbps"):             {Type: schema.TypeInt, Computed: true, Description: "BandwidthGbps"},
			common.ToSnakeCase("vpcCount"):                  {Type: schema.TypeInt, Computed: true, Description: "vpcCount"},
		},
	}
}
