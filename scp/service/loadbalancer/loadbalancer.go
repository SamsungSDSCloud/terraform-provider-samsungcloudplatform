package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"strings"
)

func init() {
	scp.RegisterResource("scp_load_balancer", ResourceLoadBalancer())
}

func ResourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerCreate,
		ReadContext:   resourceLoadBalancerRead,
		UpdateContext: resourceLoadBalancerUpdate,
		DeleteContext: resourceLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target VPC id",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Load balancer name. (3 to 20 without specials)",
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[0-9A-Za-z_-]+$`), "must contain only alphanumeric,-,_ characters"),
				),
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Load balancer description. (0 to 100 characters)",
				ValidateDiagFunc: common.ValidateDescriptionMaxlength100,
			},
			"cidr_ipv4": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Load balancer cidr ipv4",
				ValidateDiagFunc: common.ValidateCidrIpv4,
			},
			"link_ip_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Load balancer link IP band",
				ValidateDiagFunc: common.ValidateCidrIpv4,
			},
			"size": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Load balancer size. (SMALL, MEDIUM, LARGE)",
				ValidateDiagFunc: ValidateLoadBalancerSize, // valid func needed
			},
			"link_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Link ip address",
			},
			"tags": tfTags.TagsSchema(),
			//"firewall_enabled": {
			//	Type:        schema.TypeBool,
			//	Required:    true,
			//	ForceNew:    true,
			//	Description: "set filewall enable state",
			//},
		},
		Description: "Provides a Load Balancer resource.",
	}
}

func resourceLoadBalancerCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get values from schema
	vpcId := rd.Get("vpc_id").(string)
	name := rd.Get("name").(string)
	description := rd.Get("description").(string)
	cidrIpv4 := rd.Get("cidr_ipv4").(string)
	size := strings.ToUpper(rd.Get("size").(string))
	firewallEnabled := false // rd.Get("firewall_enabled").(bool)
	linkIpCidr := rd.Get("link_ip_cidr").(string)
	isFirewallLoggable := false

	inst := meta.(*client.Instance)

	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)
	if err != nil {
		return diag.FromErr(err)
	}

	isNameInvalid, err := inst.Client.LoadBalancer.CheckLoadBalancerName(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}
	if isNameInvalid {
		return diag.Errorf("Input load balancer name is invalid (maybe duplicated) : " + name)
	}

	isCidrInvalid, err := inst.Client.Subnet.CheckSubnetCidrIpv4(ctx, cidrIpv4, vpcId) // console : /27~/22, API : /24, need to fix backend
	if err != nil {
		return diag.FromErr(err)
	}
	if isCidrInvalid {
		return diag.Errorf("Input cidr is invalid (maybe duplicated) : " + cidrIpv4)
	}

	// check limit value?
	isSizeValid, err := inst.Client.LoadBalancer.CheckLoadBalancerLimitValue(ctx, size, vpcId)
	if err != nil {
		return diag.FromErr(err)
	}
	if !isSizeValid {
		return diag.Errorf("Cannot create " + size + "-size load balancer because the capability of selected VPC is exceeded.")
	}

	projectInfo, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	var targetBlockId string
	for _, serviceZone := range projectInfo.ServiceZones {
		if serviceZone.ServiceZoneId == vpcInfo.ServiceZoneId {
			targetBlockId = serviceZone.BlockId
		}
	}
	if len(targetBlockId) == 0 {
		return diag.Errorf("Failed to find target block")
	}

	tags := rd.Get("tags").(map[string]interface{})
	result, err := inst.Client.LoadBalancer.CreateLoadBalancer(ctx, targetBlockId, firewallEnabled, isFirewallLoggable, size, name, cidrIpv4, linkIpCidr, vpcInfo.ServiceZoneId, vpcId, description, tags)
	if err != nil {
		return diag.FromErr(err)
	}
	err = waitForLoadBalancerStatus(ctx, inst.Client, result.ResourceId, []string{}, []string{"ACTIVE"}, true)

	// Get linkIpAddress
	info, _, err := inst.Client.LoadBalancer.GetLoadBalancer(ctx, result.ResourceId)
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.SetId(result.ResourceId)
	rd.Set("link_ip", info.LinkIpAddress)

	return resourceLoadBalancerRead(ctx, rd, meta)
}

func resourceLoadBalancerRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.LoadBalancer.GetLoadBalancer(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("vpc_id", info.VpcId)
	rd.Set("cidr_ipv4", info.ServiceIpCidr)
	rd.Set("name", info.LoadBalancerName)
	rd.Set("description", info.LoadBalancerDescription)
	rd.Set("size", info.LoadBalancerSize)
	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceLoadBalancerUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("description") {
		_, err := inst.Client.LoadBalancer.UpdateLoadBalancerDescription(ctx, rd.Id(), rd.Get("description").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLoadBalancerRead(ctx, rd, meta)
}

func resourceLoadBalancerDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, err := inst.Client.LoadBalancer.DeleteLoadBalancer(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}
	err = waitForLoadBalancerStatus(ctx, inst.Client, rd.Id(), []string{"TERMINATING"}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ValidateLoadBalancerSize(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	const (
		LbSmall  string = "SMALL"
		LbMedium string = "MEDIUM"
		LbLarge  string = "LARGE"
	)

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get Value
	value := strings.ToUpper(v.(string))

	// Check size string
	if (value != LbSmall) && (value != LbMedium) && (value != LbLarge) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has invalid value : %s", attrKey, value),
			AttributePath: path,
		})
	}

	return diags
}

func waitForLoadBalancerStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.LoadBalancer.GetLoadBalancer(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.LoadBalancerState, nil
	})
}
