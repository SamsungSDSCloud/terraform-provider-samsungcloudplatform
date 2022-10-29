package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func ResourceSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupRuleCreate,
		ReadContext:   resourceSecurityGroupRuleRead,
		UpdateContext: resourceSecurityGroupRuleUpdate,
		DeleteContext: resourceSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target SecurityGroup id",
			},
			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SecurityGroup Rule direction (Can be 'in' or 'out')",
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					// Get attribute key
					attr := path[len(path)-1].(cty.GetAttrStep)
					attrKey := attr.Name

					// Get value
					value := strings.ToUpper(v.(string))

					if value == "IN" {
						return diag.Diagnostics{}
					}
					if value == "OUT" {
						return diag.Diagnostics{}
					}
					var diags diag.Diagnostics
					diags = append(diags, diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       fmt.Sprintf("Attribute %q only allows 'in' or 'out'. Please check : %s", attrKey, value),
						AttributePath: path,
					})
					return diags
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == strings.ToUpper(new) {
						return true
					}
					return false
				},
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "SecurityGroup Rule description. (Up to 50 characters)",
				ValidateDiagFunc: common.ValidateDescriptionMaxlength50,
			},
			"addresses_ipv4": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "SecurityGroup Rule target cidr addresses",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "SecurityGroup Rule service",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol type. (TCP, UDP, ICMP, ALL)",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Port value",
						},
					},
				},
			},
		},
		Description: "Provides a Security Group Rule resource.",
	}
}

func expandAddressesIpv4(rd *schema.ResourceData) []string {
	addressesIpv4List := rd.Get("addresses_ipv4").([]interface{})
	addressesIpv4 := make([]string, len(addressesIpv4List))
	for i, valueIpv4 := range addressesIpv4List {
		addressesIpv4[i] = valueIpv4.(string)
	}
	return addressesIpv4
}

func expandServices(rd *schema.ResourceData) ([]securitygroup.SecurityGroupServiceRule, error) {
	servicesSet := rd.Get("service").(*schema.Set).List()
	// Services
	services := make([]securitygroup.SecurityGroupServiceRule, len(servicesSet))
	for i, valueService := range servicesSet {
		s := valueService.(map[string]interface{})
		if t, ok := s["type"]; ok {
			services[i].ServiceType = strings.ToUpper(t.(string))
		} else {
			return nil, fmt.Errorf("Invalid input type found : " + t.(string))
		}
		if v, ok := s["value"]; ok {
			services[i].ServiceValue = v.(string)
		}
	}
	return services, nil
}

func resourceSecurityGroupRuleCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	sgId := rd.Get("security_group_id").(string)
	direction := rd.Get("direction").(string)
	description := rd.Get("description").(string)
	fmt.Println("OOO" + direction)
	// Addresses
	addressesIpv4 := expandAddressesIpv4(rd)

	// Services
	services, err := expandServices(rd)
	if err != nil {
		diag.FromErr(err)
	}

	inst := meta.(*client.Instance)

	response, err := inst.Client.SecurityGroup.CreateSecurityGroupRule(ctx, sgId, direction, addressesIpv4, description, services)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForSecurityGroupRuleStatus(ctx, inst.Client, response.ResourceId, sgId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}
	rd.SetId(response.ResourceId)
	return resourceSecurityGroupRuleRead(ctx, rd, meta)
}

func resourceSecurityGroupRuleRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	info, _, err := inst.Client.SecurityGroup.GetSecurityGroupRule(ctx, rd.Id(), rd.Get("security_group_id").(string))
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.Set("direction", info.RuleDirection)
	rd.Set("description", info.RuleDescription)
	rd.Set("addresses_ipv4", info.TargetNetworks)

	if info.IsAllService {
		s := common.HclSetObject{}
		s = append(s, common.HclKeyValueObject{
			"type": "all",
		})
		rd.Set("service", s)
	} else {
		s := common.HclSetObject{}
		for _, svc := range info.TcpServices {
			s = append(s, common.HclKeyValueObject{
				"type":  "tcp",
				"value": svc,
			})
		}
		for _, svc := range info.UdpServices {
			s = append(s, common.HclKeyValueObject{
				"type":  "udp",
				"value": svc,
			})
		}
		for _, svc := range info.IcmpServices {
			s = append(s, common.HclKeyValueObject{
				"type":  "icmp",
				"value": svc,
			})
		}
		rd.Set("service", s)
	}

	return nil
}

func resourceSecurityGroupRuleUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	if rd.HasChanges("direction", "addresses_ipv4", "service", "description") {

		addressesIpv4 := expandAddressesIpv4(rd)
		services, err := expandServices(rd)
		if err != nil {
			return diag.FromErr(err)
		}

		inst.Client.SecurityGroup.UpdateSecurityGroupRule(
			ctx,
			rd.Id(),
			rd.Get("security_group_id").(string),
			rd.Get("direction").(string),
			addressesIpv4,
			rd.Get("description").(string),
			services)

	}

	return resourceSecurityGroupRuleRead(ctx, rd, meta)
}

func resourceSecurityGroupRuleDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	_, err := inst.Client.SecurityGroup.DeleteSecurityGroupRule(
		ctx,
		rd.Get("security_group_id").(string),
		[]string{rd.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForSecurityGroupRuleStatus(
		ctx,
		inst.Client,
		rd.Id(),
		rd.Get("security_group_id").(string),
		[]string{}, []string{"DELETED"}, false)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForSecurityGroupRuleStatus(ctx context.Context, scpClient *client.SCPClient, id string, securityGroupId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.SecurityGroup.GetSecurityGroupRule(ctx, id, securityGroupId)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		if len(info.RuleId) == 0 {
			return info, "DELETED", nil
		}
		return info, info.RuleState, nil
	})
}
