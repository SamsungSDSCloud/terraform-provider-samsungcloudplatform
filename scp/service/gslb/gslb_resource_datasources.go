package gslb

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_gslb_resources", DatasourceGslbResources())
}

func DatasourceGslbResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGslbResourceList, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("GslbId"):    {Type: schema.TypeString, Optional: true, Description: "GSLB Id"},
			common.ToSnakeCase("CreatedBy"): {Type: schema.TypeString, Optional: true, Description: "User ID who create the resources"},
			common.ToSnakeCase("Page"):      {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			common.ToSnakeCase("Size"):      {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			common.ToSnakeCase("Sort"):      {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":                      {Type: schema.TypeList, Computed: true, Description: "List of GSLB Resource", Elem: dataSourceGslbResourceElem()},
			"total_count":                   {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides List of DNS Record.",
	}
}

func dataSourceGslbResourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	gslbId := rd.Get("gslb_id").(string)
	if len(gslbId) == 0 {
		return diag.Errorf("DNS Domain Id not found")
	}

	responses, err := inst.Client.Gslb.GetGslbResource(ctx, gslbId)

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func dataSourceGslbResourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("GslbDestination"):         {Type: schema.TypeString, Computed: true, Description: "GSLB Destination"},
			common.ToSnakeCase("GslbRegion"):              {Type: schema.TypeString, Computed: true, Description: "GSLB Region"},
			common.ToSnakeCase("GslbResourceId"):          {Type: schema.TypeString, Computed: true, Description: "GSLB Resource Id"},
			common.ToSnakeCase("GslbResourceName"):        {Type: schema.TypeString, Computed: true, Description: "GSLB Resource Name"},
			common.ToSnakeCase("GslbResourceWeight"):      {Type: schema.TypeInt, Computed: true, Description: "GSLB Resource Weight"},
			common.ToSnakeCase("GslbResourceDescription"): {Type: schema.TypeString, Computed: true, Description: "GSLB Resource Description"},
		},
	}
}
