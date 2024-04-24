package gslb

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/gslb2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_gslbs", DatasourceGslbs())
}

func DatasourceGslbs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGslbList, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"gslb_env_usage": {Type: schema.TypeString, Optional: true, Description: "GSLB Environment Usage"},
			"gslb_name":      {Type: schema.TypeString, Optional: true, Description: "GSLB Name"},
			"created_by":     {Type: schema.TypeString, Optional: true, Description: "User ID who create the resources"},
			"page":           {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":           {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":           {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":       {Type: schema.TypeList, Computed: true, Description: "List of GSLB Services", Elem: DatasourceGslbElem()},
			"total_count":    {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides List of GSLB Services.",
	}
}

func DatasourceGslbElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gslb_id":               {Type: schema.TypeString, Computed: true, Description: "GSLB Id"},
			"gslb_name":             {Type: schema.TypeString, Computed: true, Description: "GSLB Name"},
			"gslb_env_usage":        {Type: schema.TypeString, Computed: true, Description: "GSLB Environment Usage"},
			"gslb_algorithm":        {Type: schema.TypeString, Computed: true, Description: "GSLB Algorithm"},
			"gslb_state":            {Type: schema.TypeString, Computed: true, Description: "GSLB status"},
			"linked_resource_count": {Type: schema.TypeInt, Computed: true, Description: "GSLB Resource Count"},
			"created_dt":            {Type: schema.TypeString, Computed: true, Description: "Creation date"},
		},
	}
}

func dataSourceGslbList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.Gslb.GetGslbList(ctx, &gslb2.GslbOpenApiV2ControllerApiListGslbsOpts{
		GslbName:     common.GetKeyString(rd, "gslb_name"),
		GslbEnvUsage: common.GetKeyString(rd, "gslb_env_usage"),
		CreatedBy:    common.GetKeyString(rd, "created_by"),
		Page:         optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:         optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:         optional.Interface{},
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
