package dns

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/dns"
	common "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"strings"
	"time"
)

func init() {
	scp.RegisterResource("scp_dns_domain", ResourceDnsDomain())
}

func ResourceDnsDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDnsDomainCreate,
		ReadContext:   resourceDnsDomainRead,
		UpdateContext: resourceDnsDomainUpdate,
		DeleteContext: resourceDnsDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dns_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNS Name",
				ValidateFunc: validation.All(
					validation.StringLenBetween(2, 63),
					validation.StringMatch(regexp.MustCompile(`^[a-z\d]$|^[a-z\d][a-z\d-]*[a-z\d]$`), "Enter 2-63 char (lowercase, number, -)."),
				),
			},
			"dns_root_domain_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "DNS Root Domain Name",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{".com", ".net", ".org", ".kr", ".한국", ".pe.kr", ".biz", ".info", ".cn", ".tv", ".in", ".eu", ".ac", ".tw", ".mobi", ".name", ".cc", ".jp", ".asia", ".me", ".tel", ".pro", ".so", ".sx", ".co", ".xxx", ".pw", ".ru", ".ph", ".co.kr", ".app", ".io"}, false)),
			},
			"dns_description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "DNS Domain Description",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 200)),
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Dns Domain resource. (Only available for PRIVATE environment usage type)",
	}
}

func resourceDnsDomainCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	dnsDomainName := rd.Get("dns_domain_name").(string)
	dnsRootDomainName := rd.Get("dns_root_domain_name").(string)
	dnsDescription := rd.Get("dns_description").(string)

	createRequest := dns.CreateDnsDomainRequest{
		DnsDomainName:  dnsDomainName + dnsRootDomainName,
		DnsDescription: dnsDescription,
	}

	result, _, err := inst.Client.Dns.CreateDnsDomain(ctx, createRequest, rd.Get("tags").(map[string]interface{}))

	if err != nil {
		return
	}

	// wait for server state
	time.Sleep(10 * time.Second)

	err = waitForDnsDomainStatus(ctx, inst.Client, result.ResourceId, []string{"CREATING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return
	}

	rd.SetId(result.ResourceId)

	return resourceDnsDomainRead(ctx, rd, meta)
}

func resourceDnsDomainRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	dnsInfo, _, err := inst.Client.Dns.GetDnsDomainDetail(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("dns_domain_name", strings.Split(dnsInfo.DnsDomainName, ".")[0])
	rd.Set("dns_root_domain_name", strings.TrimLeft(dnsInfo.DnsDomainName, strings.Split(dnsInfo.DnsDomainName, ".")[0]))
	rd.Set("dns_description", dnsInfo.DnsDescription)
	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceDnsDomainUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	if rd.HasChanges("dns_description") {
		_, err := inst.Client.Dns.UpdateDnsDomainDescription(ctx, rd.Id(), rd.Get("dns_description").(string))
		if err != nil {
			return
		}
	}

	err = tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return
	}

	return resourceDnsDomainRead(ctx, rd, meta)
}

func resourceDnsDomainDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Dns.DeleteDnsDomain(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	time.Sleep(10 * time.Second)

	err = waitForDnsDomainStatus(ctx, inst.Client, rd.Id(), []string{"TERMINATING"}, []string{"DELETED"}, false)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForDnsDomainStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.Dns.GetDnsDomainDetail(ctx, id)

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
