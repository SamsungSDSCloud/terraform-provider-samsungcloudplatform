package endpoint

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/endpoint2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_endpoint", DatasourceEndpoint())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_endpoints", DatasourceEndpoints())
}
func DatasourceEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			common.ToSnakeCase("EndpointIpAddress"):   {Type: schema.TypeString, Optional: true, Description: "Endpoint Ip Address"},
			common.ToSnakeCase("EndpointId"):          {Type: schema.TypeString, Optional: true, Description: "Endpoint Id"},
			common.ToSnakeCase("EndpointState"):       {Type: schema.TypeString, Optional: true, Description: "Endpoint status"},
			common.ToSnakeCase("EndpointType"):        {Type: schema.TypeString, Optional: true, Description: "Endpoint type"},
			common.ToSnakeCase("ObjectId"):            {Type: schema.TypeString, Optional: true, Description: "Object Id"},
			common.ToSnakeCase("EndpointDescription"): {Type: schema.TypeString, Optional: true, Description: "Endpoint Description"},
			common.ToSnakeCase("VpcId"):               {Type: schema.TypeString, Optional: true, Description: "Vpc Id"},
			common.ToSnakeCase("ServiceZoneId"):       {Type: schema.TypeString, Optional: true, Description: "Service zone id"},
			common.ToSnakeCase("CreatedBy"):           {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
		},
		Description: "Provide Detail of public ip",
	}
}

func dataSourceRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	endpointId := rd.Get(common.ToSnakeCase("endpointId")).(string)
	info, _, err := inst.Client.Endpoint.GetEndpoint(ctx, endpointId)

	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set(common.ToSnakeCase("EndpointIpAddress"), info.EndpointIpAddress)
	rd.Set(common.ToSnakeCase("EndpointName"), info.EndpointName)
	rd.Set(common.ToSnakeCase("EndpointType"), info.EndpointType)
	rd.Set(common.ToSnakeCase("EndpointState"), info.EndpointState)
	rd.Set(common.ToSnakeCase("ObjectId"), info.ObjectId)
	rd.Set(common.ToSnakeCase("EndpointDescription"), info.EndpointDescription)
	rd.Set(common.ToSnakeCase("ServiceZoneId"), info.ServiceZoneId)
	rd.Set(common.ToSnakeCase("VpcId"), info.VpcId)
	rd.Set(common.ToSnakeCase("CreatedBy"), info.CreatedBy)

	return nil
}

func DatasourceEndpoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{ //스키마 정의
			common.ToSnakeCase("objectId"):          {Type: schema.TypeString, Optional: true, Description: "Endpoint id"},
			common.ToSnakeCase("EndpointId"):        {Type: schema.TypeString, Optional: true, Description: "Endpoint id"},
			common.ToSnakeCase("EndpointName"):      {Type: schema.TypeString, Optional: true, Description: "Endpoint name"},
			common.ToSnakeCase("endpointType"):      {Type: schema.TypeString, Optional: true, Description: "Endpoint type"},
			common.ToSnakeCase("EndpointStates"):    {Type: schema.TypeString, Optional: true, Description: "Endpoint status"},
			common.ToSnakeCase("EndpointIpAddress"): {Type: schema.TypeString, Optional: true, Description: "Endpoint Ip Address"},
			common.ToSnakeCase("VpcId"):             {Type: schema.TypeString, Optional: true, Description: "Vpc id"},
			common.ToSnakeCase("ServiceZoneId"):     {Type: schema.TypeString, Optional: true, Description: "Service zone id"},
			common.ToSnakeCase("CreatedBy"):         {Type: schema.TypeString, Optional: true, Description: "Person who created the resource"},
			common.ToSnakeCase("Page"):              {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			common.ToSnakeCase("Size"):              {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":                              {Type: schema.TypeList, Optional: true, Description: "Endpoint list", Elem: datasourceElem()},
			"total_count":                           {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of endpoints.",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	requestParam := &endpoint2.EndpointOpenApiControllerApiListEndpointOpts{
		EndpointIpAddress: common.GetKeyString(rd, common.ToSnakeCase("EndpointIpAddress")),
		EndpointId:        common.GetKeyString(rd, common.ToSnakeCase("EndpointId")),
		EndpointName:      common.GetKeyString(rd, common.ToSnakeCase("EndpointName")),
		EndpointType:      common.GetKeyString(rd, common.ToSnakeCase("EndpointType")),
		EndpointStates:    optional.NewInterface(rd.Get("endpoint_states").(string)),
		ObjectId:          common.GetKeyString(rd, common.ToSnakeCase("ObjectId")),
		VpcId:             common.GetKeyString(rd, common.ToSnakeCase("VpcId")),
		ServiceZoneId:     common.GetKeyString(rd, common.ToSnakeCase("ServiceZoneId")),
		CreatedBy:         common.GetKeyString(rd, common.ToSnakeCase("CreatedBy")),
		Page:              optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:              optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:              optional.NewInterface([]string{"createdDt:desc"}),
	}

	responses, _, err := inst.Client.Endpoint.GetEndpointList(ctx, requestParam)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("CreatedBy"):         {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			common.ToSnakeCase("CreatedDt"):         {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			common.ToSnakeCase("ModifiedBy"):        {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			common.ToSnakeCase("ModifiedDt"):        {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			common.ToSnakeCase("ProjectId"):         {Type: schema.TypeString, Computed: true, Description: "Project id"},
			common.ToSnakeCase("ServiceZoneId"):     {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			common.ToSnakeCase("VpcId"):             {Type: schema.TypeString, Computed: true, Description: "VPC id"},
			common.ToSnakeCase("EndpointId"):        {Type: schema.TypeString, Computed: true, Description: "Endpoint id"},
			common.ToSnakeCase("EndpointName"):      {Type: schema.TypeString, Computed: true, Description: "Endpoint name"},
			common.ToSnakeCase("EndpointType"):      {Type: schema.TypeString, Computed: true, Description: "Endpoint type"},
			common.ToSnakeCase("EndpointState"):     {Type: schema.TypeString, Computed: true, Description: "Endpoint status"},
			common.ToSnakeCase("EndpointIpAddress"): {Type: schema.TypeString, Computed: true, Description: "Endpoint ip address"},
			common.ToSnakeCase("ObjectId"):          {Type: schema.TypeString, Computed: true, Description: "Object id"},
		},
	}
}
