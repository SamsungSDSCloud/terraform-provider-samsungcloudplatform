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
	scp.RegisterDataSource("scp_dns_records", DatasourceDnsRecords())
}

func DatasourceDnsRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDnsRecordList, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("DnsDomainId"):       {Type: schema.TypeString, Optional: true, Description: "DNS Domain Id"},
			common.ToSnakeCase("DnsRecordName"):     {Type: schema.TypeString, Optional: true, Description: "DNS Record Name"},
			common.ToSnakeCase("DnsRecordType"):     {Type: schema.TypeString, Optional: true, Description: "DNS Record Type"},
			common.ToSnakeCase("RecordDestination"): {Type: schema.TypeString, Optional: true, Description: "DNS Record Destination"},
			common.ToSnakeCase("CreatedBy"):         {Type: schema.TypeString, Optional: true, Description: "User ID who create the resources"},
			common.ToSnakeCase("Page"):              {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			common.ToSnakeCase("Size"):              {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			common.ToSnakeCase("Sort"):              {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":                              {Type: schema.TypeList, Computed: true, Description: "List of DNS Domain Services", Elem: dataSourceDnsRecordElem()},
			"total_count":                           {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides List of DNS Record.",
	}
}

func dataSourceDnsRecordList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	dnsDomainId := rd.Get("dns_domain_id").(string)
	if len(dnsDomainId) == 0 {
		return diag.Errorf("DNS Domain Id not found")
	}

	responses, _, err := inst.Client.Dns.GetDnsRecordList(ctx, dnsDomainId, &dns2.DnsOpenApiV2ControllerApiListDnsRecordOpts{
		DnsRecordName:     common.GetKeyString(rd, "dns_record_name"),
		DnsRecordType:     common.GetKeyString(rd, "dns_record_type"),
		RecordDestination: common.GetKeyString(rd, "record_destination"),
		CreatedBy:         common.GetKeyString(rd, "created_by"),
		Page:              optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:              optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:              optional.Interface{},
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := convertDnsRecordListToHclSet(responses)
	//contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func dataSourceDnsRecordElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("DnsRecordId"):        {Type: schema.TypeString, Computed: true, Description: "DNS Record Id"},
			common.ToSnakeCase("DnsRecordName"):      {Type: schema.TypeString, Computed: true, Description: "DNS Record Name"},
			common.ToSnakeCase("DnsRecordType"):      {Type: schema.TypeString, Computed: true, Description: "DNS Record Type"},
			common.ToSnakeCase("Ttl"):                {Type: schema.TypeInt, Computed: true, Description: "TTL"},
			common.ToSnakeCase("DnsState"):           {Type: schema.TypeString, Computed: true, Description: "DNS status"},
			common.ToSnakeCase("RecordDestinations"): {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "DNS Record Destinations"},
		},
	}
}

func convertDnsRecordListToHclSet(recordResponse dns2.ListResponseDnsDomainRecordListItemResponse) common.HclSetObject {
	var recordList common.HclSetObject
	for _, record := range recordResponse.Contents {
		if len(record.DnsRecordId) == 0 {
			continue
		}

		kv := common.HclKeyValueObject{
			"dns_record_id":       record.DnsRecordId,
			"dns_record_type":     record.DnsRecordType,
			"dns_record_name":     record.DnsRecordName,
			"ttl":                 record.Ttl,
			"dns_state":           record.DnsState,
			"record_destinations": record.RecordDestinations,
		}
		recordList = append(recordList, kv)
	}
	return recordList
}
