package vpc

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

func init() {
	scp.RegisterResource("scp_vpc_dns", ResourceVpcDns())
}

func ResourceVpcDns() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceVpcDnsCreate,
		ReadContext:   resourceVpcDnsRead,
		UpdateContext: resourceVpcDnsUpdate,
		DeleteContext: resourceVpcDnsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true, //필수 작성
				Description:  "VPC id",
				ValidateFunc: validation.StringLenBetween(3, 100),
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true, //필수 작성
				ForceNew:    true, //해당 필드 수정시 자원 삭제 후 다시 생성됨
				Description: "domain (ex: abc.com)",
				ValidateFunc: validation.All( //입력 값 Validation 체크
					validation.StringLenBetween(3, 100),
					validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9]+\\.[a-zA-Z0-9\\.]+$"), "dns name has to be valid"),
				),
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Subnet Id for Source IP",
				ValidateFunc: validation.StringLenBetween(1, 50),
			},
			"dns_ip": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "DNS ip address",
				ValidateDiagFunc: common.ValidateIpv4,
			},
			"source_ip": {
				Type:             schema.TypeString,
				Optional:         true, // optional
				Description:      "Source Ip address",
				ValidateDiagFunc: common.ValidateIpv4,
			},
		},
		Description: "Provides a VPC Dns resource.",
	}
}

func resourceVpcDnsCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	vpcId := rd.Get("vpc_id").(string)
	dnsUserZoneDomain := rd.Get("domain").(string)
	subnetId := rd.Get("subnet_id").(string)
	dnsUserZoneServerIp := rd.Get("dns_ip").(string)
	dnsUserZoneSourceIp := rd.Get("source_ip").(string)

	inst := meta.(*client.Instance)

	// check if vpcId is valid
	_, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)
	if err != nil {
		return diag.FromErr(err)
	}

	// check if subnetId is valid
	subnetInfo, _, err := inst.Client.Subnet.GetSubnet(ctx, subnetId)
	if err != nil {
		return diag.FromErr(err)
	}

	if dnsUserZoneSourceIp != "" && !common.IsSubnetContainsIp(subnetInfo.SubnetCidrBlock, dnsUserZoneSourceIp) {
		return diag.Errorf("Source IP is not in subnet")
	}

	tflog.Debug(ctx, "Try create vpc dns zone : "+vpcId+", "+dnsUserZoneDomain+", "+subnetId+", "+dnsUserZoneServerIp)

	_, err = inst.Client.Vpc.CreateVpcDns(ctx, vpcId, dnsUserZoneDomain, dnsUserZoneServerIp, dnsUserZoneSourceIp, subnetId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitVpcDnsCreating(ctx, inst.Client, vpcId, dnsUserZoneDomain)
	if err != nil {
		return diag.FromErr(err)
	}

	info, _, err := inst.Client.Vpc.GetVpcDnsInfoByDomain(ctx, vpcId, dnsUserZoneDomain)
	rd.SetId(inst.Client.Vpc.MergeVpcDnsId(vpcId, info.DnsUserZoneId))

	// Refresh
	return resourceVpcDnsRead(ctx, rd, meta)
}

func resourceVpcDnsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	dnsInfo, _, err := inst.Client.Vpc.GetVpcDnsInfoById(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	rd.Set("domain", dnsInfo.DnsUserZoneDomain)
	rd.Set("source_ip", dnsInfo.DnsUserZoneSourceIp)
	rd.Set("dns_ip", dnsInfo.DnsUserZoneServerIp)

	return nil
}

func resourceVpcDnsUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceVpcDnsDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	err := inst.Client.Vpc.DeleteVpcDns(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitVpcDnsDeleting(ctx, inst.Client, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitVpcDnsCreating(ctx context.Context, scpClient *client.SCPClient, vpcId string, dnsUserZoneDomain string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		return scpClient.Vpc.GetVpcDnsInfoByDomain(ctx, vpcId, dnsUserZoneDomain)
	})
}

func waitVpcDnsDeleting(ctx context.Context, scpClient *client.SCPClient, dnsId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		return scpClient.Vpc.GetVpcDnsInfoById(ctx, dnsId)
	})
}
