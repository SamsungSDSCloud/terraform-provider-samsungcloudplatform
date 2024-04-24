package dns

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/dns2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_dns_domains", DatasourceDnsDomains())
}

func DatasourceDnsDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDnsDomainList, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("DnsEnvUsage"):   {Type: schema.TypeString, Optional: true, Description: "DNS Domain Environment Usage"},
			common.ToSnakeCase("DnsDomainName"): {Type: schema.TypeString, Optional: true, Description: "DNS Domain Name"},
			common.ToSnakeCase("CreatedBy"):     {Type: schema.TypeString, Optional: true, Description: "User ID who create the resources"},
			common.ToSnakeCase("Page"):          {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			common.ToSnakeCase("Size"):          {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			common.ToSnakeCase("Sort"):          {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":                          {Type: schema.TypeList, Computed: true, Description: "List of DNS Domain Services", Elem: dataSourceDnsDomainElem()},
			"total_count":                       {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides List of DNS Domain Services.",
	}
}

func dataSourceDnsDomainList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.Dns.GetDnsDomainList(ctx, &dns2.DnsOpenApiV2ControllerApiListDnsDomainOpts{
		DnsDomainName: common.GetKeyString(rd, "dns_domain_name"),
		DnsEnvUsage:   common.GetKeyString(rd, "dns_env_usage"),
		CreatedBy:     common.GetKeyString(rd, "created_by"),
		Page:          optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:          optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:          optional.Interface{},
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

func dataSourceDnsDomainElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("DnsDomainId"):       {Type: schema.TypeString, Computed: true, Description: "DNS Domain Id"},
			common.ToSnakeCase("DnsDomainName"):     {Type: schema.TypeString, Computed: true, Description: "DNS Domain Name"},
			common.ToSnakeCase("DnsEnvUsage"):       {Type: schema.TypeString, Computed: true, Description: "DNS Environment Usage"},
			common.ToSnakeCase("DnsDomainType"):     {Type: schema.TypeString, Computed: true, Description: "DNS Domain Type"},
			common.ToSnakeCase("AutoExtension"):     {Type: schema.TypeBool, Computed: true, Description: "DNS Environment Auto Extension Option"},
			common.ToSnakeCase("DnsState"):          {Type: schema.TypeString, Computed: true, Description: "DNS Domain status"},
			common.ToSnakeCase("Region"):            {Type: schema.TypeString, Computed: true, Description: "DNS Domain Usage Region"},
			common.ToSnakeCase("LinkedRecordCount"): {Type: schema.TypeInt, Computed: true, Description: "Linked Record Count"},
			common.ToSnakeCase("StartDate"):         {Type: schema.TypeString, Computed: true, Description: "DNS Domain Start Date"},
			common.ToSnakeCase("ExpiredDate"):       {Type: schema.TypeString, Computed: true, Description: "DNS Domain Expired Date"},
		},
	}
}
