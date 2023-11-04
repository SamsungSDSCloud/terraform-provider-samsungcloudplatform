package vpc

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_vpc_dnss", DatasourceVpcDns())
}

func DatasourceVpcDns() *schema.Resource {
	return &schema.Resource{
		ReadContext: DnsList, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{ //스키마 정의
			common.ToSnakeCase("VpcId"): {Type: schema.TypeString, Required: true, Description: "VPC id"},
			"contents":                  {Type: schema.TypeList, Optional: true, Description: "VPC DNS list", Elem: DnsElem()},
			"total_count":               {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of vpc DNS's.",
	}
}

func DnsList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	vpcId := rd.Get(common.ToSnakeCase("VpcId")).(string)

	responses, err := inst.Client.Vpc.GetVpcDnsList(ctx, vpcId)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func DnsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("DnsUserZoneDomain"):   {Type: schema.TypeString, Computed: true, Description: "Zone Domain"},
			common.ToSnakeCase("DnsUserZoneId"):       {Type: schema.TypeString, Computed: true, Description: "Zone Domain Id"},
			common.ToSnakeCase("DnsUserZoneName"):     {Type: schema.TypeString, Computed: true, Description: "Zone Name"},
			common.ToSnakeCase("DnsUserZoneServerIp"): {Type: schema.TypeString, Computed: true, Description: "Zone Dns IP"},
			common.ToSnakeCase("DnsUserZoneSourceIp"): {Type: schema.TypeString, Computed: true, Description: "Zone Source IP"},
			common.ToSnakeCase("DnsUserZoneState"):    {Type: schema.TypeString, Computed: true, Description: "Zone State"},
		},
	}
}
