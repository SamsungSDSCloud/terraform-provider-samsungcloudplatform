package internetgateway

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	internetgateway2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/internet-gateway2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_internet_gateways", DatasourceInternetGateways())
}

func DatasourceInternetGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{ //스키마 정의
			common.ToSnakeCase("ServiceZoneId"):              {Type: schema.TypeString, Optional: true, Description: "Service zone id"},
			common.ToSnakeCase("VpcId"):                      {Type: schema.TypeString, Optional: true, Description: "VPC id"},
			common.ToSnakeCase("internetGatewayId"):          {Type: schema.TypeString, Optional: true, Description: "Internet Gateway name"},
			common.ToSnakeCase("internetGatewayName"):        {Type: schema.TypeString, Optional: true, Description: "Internet Gateway name"},
			common.ToSnakeCase("internetGatewayDescription"): {Type: schema.TypeString, Optional: true, Description: "Description of internet gateway"},
			common.ToSnakeCase("internetGatewayState"):       {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Internet Gateway status"},
			common.ToSnakeCase("CreatedBy"):                  {Type: schema.TypeString, Optional: true, Description: "Person who created the resource"},
			common.ToSnakeCase("Page"):                       {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			common.ToSnakeCase("Size"):                       {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":                                       {Type: schema.TypeList, Optional: true, Description: "Internet Gateway list", Elem: datasourceElem()},
			"total_count":                                    {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of internet gateways.",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.InternetGateway.GetInternetGatewayList(ctx, &internetgateway2.InternetGatewayV2ControllerV2ApiListInternetGatewaysOpts{
		VpcId:               optional.NewString(rd.Get("vpc_id").(string)),
		InternetGatewayId:   optional.NewString(rd.Get("internet_gateway_id").(string)),
		InternetGatewayName: optional.NewString(rd.Get("internet_gateway_name").(string)),
		CreatedBy:           optional.NewString(rd.Get("created_by").(string)),
		Page:                optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:                optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:                optional.Interface{},
	})

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
			common.ToSnakeCase("ProjectId"):                  {Type: schema.TypeString, Computed: true, Description: "Project id"},
			common.ToSnakeCase("ServiceZoneId"):              {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			common.ToSnakeCase("CreatedBy"):                  {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			common.ToSnakeCase("CreatedDt"):                  {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			common.ToSnakeCase("ModifiedBy"):                 {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			common.ToSnakeCase("ModifiedDt"):                 {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			common.ToSnakeCase("InternetGatewayId"):          {Type: schema.TypeString, Computed: true, Description: "InternetGateway id"},
			common.ToSnakeCase("InternetGatewayName"):        {Type: schema.TypeString, Computed: true, Description: "InternetGateway name"},
			common.ToSnakeCase("InternetGatewayDescription"): {Type: schema.TypeString, Computed: true, Description: "InternetGateway description"},
			common.ToSnakeCase("InternetGatewayState"):       {Type: schema.TypeString, Computed: true, Description: "InternetGateway status"},
			common.ToSnakeCase("VpcId"):                      {Type: schema.TypeString, Computed: true, Description: "VPC id"},
		},
	}
}
