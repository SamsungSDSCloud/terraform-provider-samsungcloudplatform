package directconnect

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	directconnect2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/direct-connect2"
	"github.com/antihax/optional"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_direct_connects", DatasourceDirectConnects())
}

func DatasourceDirectConnects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{ //스키마 정의
			common.ToSnakeCase("DirectConnectId"):   {Type: schema.TypeString, Optional: true, Description: "Direct connect id"},
			common.ToSnakeCase("DirectConnectName"): {Type: schema.TypeString, Optional: true, Description: "Direct connect name"},
			common.ToSnakeCase("CreatedBy"):         {Type: schema.TypeString, Optional: true, Description: "Person who created the resource"},
			common.ToSnakeCase("Page"):              {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			common.ToSnakeCase("Size"):              {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":                              {Type: schema.TypeList, Optional: true, Description: "Direct Connect list", Elem: datasourceElem()},
			"total_count":                           {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of direct connects.",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.DirectConnect.GetDirectConnectList(ctx, &directconnect2.DirectConnectOpenApiControllerApiListDirectConnectsOpts{
		DirectConnectId:   optional.NewString(rd.Get("direct_connect_id").(string)),
		DirectConnectName: optional.NewString(rd.Get("direct_connect_name").(string)),
		CreatedBy:         optional.NewString(rd.Get("created_by").(string)),
		Page:              optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:              optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:              optional.Interface{},
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
			common.ToSnakeCase("ProjectId"):                {Type: schema.TypeString, Computed: true, Description: "Project id"},
			common.ToSnakeCase("ServiceZoneId"):            {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			common.ToSnakeCase("CreatedBy"):                {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			common.ToSnakeCase("CreatedDt"):                {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			common.ToSnakeCase("ModifiedBy"):               {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			common.ToSnakeCase("ModifiedDt"):               {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			common.ToSnakeCase("DirectConnectId"):          {Type: schema.TypeString, Computed: true, Description: "DirectConnect id"},
			common.ToSnakeCase("DirectConnectName"):        {Type: schema.TypeString, Computed: true, Description: "DirectConnect name"},
			common.ToSnakeCase("DirectConnectDescription"): {Type: schema.TypeString, Computed: true, Description: "DirectConnect description"},
			common.ToSnakeCase("DirectConnectState"):       {Type: schema.TypeString, Computed: true, Description: "DirectConnect status"},
			common.ToSnakeCase("ProductGroupId"):           {Type: schema.TypeString, Computed: true, Description: "ProductGroup id"},
			common.ToSnakeCase("BandwidthGbps"):            {Type: schema.TypeInt, Computed: true, Description: "Bandwidth gbps"},
			common.ToSnakeCase("UplinkEnabled"):            {Type: schema.TypeBool, Computed: true, Description: "Uplink enabled"},
		},
	}
}
