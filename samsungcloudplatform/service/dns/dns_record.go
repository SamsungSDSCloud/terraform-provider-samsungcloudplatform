package dns

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client/dns"
	common "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"time"
)

func init() {
	samsungcloudplatform.RegisterResource("samsungcloudplatform_dns_record", ResourceDnsRecord())
}

func ResourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDnsRecordCreate,
		ReadContext:   resourceDnsRecordRead,
		UpdateContext: resourceDnsRecordUpdate,
		DeleteContext: resourceDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dns_domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNS Domain Id",
			},
			"dns_record_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNS Record Name (0 to 63, lowercase, number and -_.@)",
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 63),
					validation.StringMatch(regexp.MustCompile(`^[a-z\d@*-_.]*$`), "must contain only lowercase, number and -_.@"),
				),
			},
			"dns_record_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "DNS Record Type. One of A, TXT, CNAME, MX, AAAA, SPF",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"A", "TXT", "CNAME", "MX", "AAAA", "SPF"}, false)),
			},
			"ttl": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "DNS TTL. (300 to 86400)",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(300, 86400)),
			},
			"dns_record_mapping": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "DNS Record Mappings. Record Type CNAME, SPF, TXT can have only one record mapping. Record Type A, AAAA, MX can have 1 or more record mappings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_destination": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "DnsDomain Resource Destination",
							ValidateDiagFunc: validateIpv4,
						},
						"preference": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "DnsDomain Resource Weight",
							ValidateFunc: validation.IntBetween(0, 65535),
						},
					},
				},
			},
		},
		Description: "Provides a DnsDomain resource.",
	}
}

func validateIpv4(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check characters
	if !regexp.MustCompile("^(25[0-5]|2[0-4]\\d|1\\d{2}|[1-9]\\d|[1-9])(\\.(25[0-5]|2[0-4]\\d|1\\d{2}|[1-9]\\d|\\d)){3}$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q: Supports only IPv4 format.", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func validateDnsRecordName(rd *schema.ResourceData) diag.Diagnostics {
	//dnsRecordType := rd.Get("dns_record_type").(string)
	//dnsRecordName := rd.Get("dns_record_name").(string)

	// Check characters
	//switch dnsRecordType {
	//case "A":
	//case "AAAA":
	//case "MX":
	//	if !regexp.MustCompile(`^[a-z\\d]$|^([a-z\\d]|(_(?!_))|(-(?!-)))([a-z\\d]|(\\.(?!\\.))|(_(?!_))|(-(?!-)))*[a-z\\d]$|^(\\*\\.[a-z\\d])([a-z\\d]|(\\.(?!\\.))|(_(?!_))|(-(?!-)))*[a-z\\d]$|^(\\*\\.[a-z\\d])$`).MatchString(dnsRecordName) {
	//		return diag.Errorf("Invalid DNS Record Name")
	//	}
	//default:
	//	if !regexp.MustCompile(`^[a-z\\d@]$|^([a-z\\d]|(_(?!_))|(-(?!-)))([a-z\\d]|(\\.(?!\\.))|(_(?!_))|(-(?!-)))*[a-z\\d]$|^(\\*\\.[a-z\\d])([a-z\\d]|(\\.(?!\\.))|(_(?!_))|(-(?!-)))*[a-z\\d]$|^(\\*\\.[a-z\\d])$`).MatchString(dnsRecordName) {
	//		return diag.Errorf("Invalid DNS Record Name")
	//	}
	//}

	return nil
}

func convertDnsDomainResources(list common.HclListObject) ([]dns.DnsRecordMappingRequest, error) {
	var result []dns.DnsRecordMappingRequest
	for _, l := range list {
		itemObject := l.(common.HclKeyValueObject)
		info := dns.DnsRecordMappingRequest{}
		if recordDestination, ok := itemObject["record_destination"]; ok {
			info.RecordDestination = recordDestination.(string)
		}
		if preference, ok := itemObject["preference"]; ok {
			info.Preference = int32(preference.(int))
		}

		result = append(result, info)
	}

	return result, nil
}

func resourceDnsRecordCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dnsDomainId := rd.Get("dns_domain_id").(string)
	dnsRecordType := rd.Get("dns_record_type").(string)
	dnsRecordName := rd.Get("dns_record_name").(string)
	ttl := int32(rd.Get("ttl").(int))

	dnsRecordMapping, err := convertDnsDomainResources(rd.Get("dns_record_mapping").(*schema.Set).List())

	createRequest := dns.CreateDnsRecordRequest{
		DnsRecordType:    dnsRecordType,
		DnsRecordName:    dnsRecordName,
		Ttl:              ttl,
		DnsRecordMapping: dnsRecordMapping,
	}

	validateDnsRecordName(rd)

	result, _, err := inst.Client.Dns.CreateDnsRecord(ctx, dnsDomainId, createRequest)
	if err != nil {
		return
	}

	// wait for server state
	time.Sleep(10 * time.Second)

	err = waitForDnsRecordStatus(ctx, inst.Client, dnsDomainId, result.ResourceId, []string{"CREATING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return
	}

	rd.SetId(result.ResourceId)

	return resourceDnsRecordRead(ctx, rd, meta)
}

func resourceDnsRecordRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	dnsDomainId := rd.Get("dns_domain_id").(string)

	info, _, err := inst.Client.Dns.GetDnsRecordDetail(ctx, dnsDomainId, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("dns_domain_id", dnsDomainId)
	rd.Set("dns_record_name", info.DnsRecordName)
	rd.Set("dns_record_type", info.DnsRecordType)
	rd.Set("ttl", info.Ttl)

	var dnsRecordMapping []common.HclKeyValueObject

	for _, recordMapInfo := range info.DnsRecordMapping {
		dnsRecordMapping = append(dnsRecordMapping, common.HclKeyValueObject{
			"record_destination": recordMapInfo.RecordDestination,
			"preference":         recordMapInfo.Preference,
		})
	}

	rd.Set("dns_record_mapping", dnsRecordMapping)

	return nil
}

func resourceDnsRecordUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	dnsDomainId := rd.Get("dns_domain_id").(string)
	dnsRecordType := rd.Get("dns_record_type").(string)

	if rd.HasChanges("ttl", "dns_record_mapping") {
		ttl := int32(rd.Get("ttl").(int))

		dnsRecordMappingList := rd.Get("dns_record_mapping").(*schema.Set).List()
		dnsRecordMapping := make([]dns.DnsRecordMappingRequest, len(dnsRecordMappingList))
		for i, recordMapInfo := range dnsRecordMappingList {
			r := recordMapInfo.(map[string]interface{})
			if t, ok := r["record_destination"]; ok {
				dnsRecordMapping[i].RecordDestination = t.(string)
			}
			if v, ok := r["preference"]; ok {
				dnsRecordMapping[i].Preference = int32(v.(int))
			}
		}

		updateDnsRecordRequest := dns.ChangeDnsRecordRequest{
			DnsRecordType:    dnsRecordType,
			Ttl:              ttl,
			DnsRecordMapping: dnsRecordMapping,
		}

		_, _, err := inst.Client.Dns.UpdateDnsRecord(ctx, dnsDomainId, rd.Id(), updateDnsRecordRequest)
		if err != nil {
			return diag.FromErr(err)
		}

		// wait for server state
		time.Sleep(10 * time.Second)

		err = waitForDnsRecordStatus(ctx, inst.Client, dnsDomainId, rd.Id(), []string{"EDITING"}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	return resourceDnsRecordRead(ctx, rd, meta)
}

func resourceDnsRecordDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Dns.DeleteDnsRecord(ctx, rd.Get("dns_domain_id").(string), rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	time.Sleep(10 * time.Second)

	err = waitForDnsRecordStatus(ctx, inst.Client, rd.Get("dns_domain_id").(string), rd.Id(), []string{"TERMINATING"}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForDnsRecordStatus(ctx context.Context, scpClient *client.SCPClient, dnsDomainId string, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.Dns.GetDnsRecordDetail(ctx, dnsDomainId, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.DnsState, nil
	})
}
