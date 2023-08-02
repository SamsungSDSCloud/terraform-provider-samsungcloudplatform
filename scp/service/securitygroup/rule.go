package securitygroup

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func init() {
	scp.RegisterResource("scp_security_group_rule", ResourceSecurityGroupRule())
	scp.RegisterResource("scp_security_group_bulk_rule", ResourceSecurityGroupBulkRule())
}

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

func ResourceSecurityGroupBulkRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupBulkRuleCreate,
		ReadContext:   resourceSecurityGroupBulkRuleRead,
		DeleteContext: resourceSecurityGroupRuleAllDelete,
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
			"rule": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
				Description: "Security Group Rule List",
			},
		},
		Description: "Provides a Security Group Bulk Rule resource.",
	}
}

func expandRules(rd *schema.ResourceData) ([]securitygroup.SecurityGroupRule, error) {
	ruleSet := rd.Get("rule").(*schema.Set).List()
	// Rules
	rules := make([]securitygroup.SecurityGroupRule, len(ruleSet))

	for i, rule := range ruleSet {
		fmt.Println(rule)

		itemObject := rule.(common.HclKeyValueObject)

		if direction, ok := itemObject["direction"]; ok {
			rules[i].RuleDirection = strings.ToUpper(direction.(string))
		}

		if description, ok := itemObject["description"]; ok {
			rules[i].RuleDescription = description.(string)
		}

		if addressesIpv4, ok := itemObject["addresses_ipv4"]; ok {
			//s := rules[i].DestinationIpAddresses
			s := make([]string, 0)
			for _, addressIpv4 := range addressesIpv4.([]interface{}) {
				s = append(s, addressIpv4.(string))
				fmt.Println(s)
			}

			if strings.ToUpper(rules[i].RuleDirection) == "IN" {
				rules[i].SourceIpAddresses = s
			} else if strings.ToUpper(rules[i].RuleDirection) == "OUT" {
				rules[i].DestinationIpAddresses = s
			}
		}

		if services, ok := itemObject["service"]; ok {

			servicesSet := services.(*schema.Set).List()
			//Services
			services := make([]securitygroup.SecurityGroupServiceRule, len(servicesSet))
			for j, valueService := range servicesSet {
				s := valueService.(map[string]interface{})
				if t, ok := s["type"]; ok {
					services[j].ServiceType = strings.ToUpper(t.(string))
				} else {
					return nil, fmt.Errorf("Invalid input type found : " + t.(string))
				}
				if v, ok := s["value"]; ok {
					services[j].ServiceValue = v.(string)
				}
			}

			rules[i].Services = services
		}

	}

	return rules, nil
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

func resourceSecurityGroupBulkRuleCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	sgId := rd.Get("security_group_id").(string)
	rules, err := expandRules(rd)
	if err != nil {
		diag.FromErr(err)
	}

	inst := meta.(*client.Instance)

	response, err := inst.Client.SecurityGroup.CreateSecurityGroupBulkRule(ctx, sgId, rules)
	if err != nil {
		return diag.FromErr(err)
	}
	rd.SetId(response.ResourceId)

	return nil
}

func resourceSecurityGroupRuleRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	info, _, err := inst.Client.SecurityGroup.GetSecurityGroupRule(ctx, rd.Id(), rd.Get("security_group_id").(string))
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("direction", info.RuleDirection)
	rd.Set("description", info.RuleDescription)
	rd.Set("addresses_ipv4", info.TargetNetworks)

	if *info.IsAllService {
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

	err := waitForSecurityGroupRuleStatus(ctx, inst.Client, rd.Id(), rd.Get("security_group_id").(string), []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSecurityGroupRuleRead(ctx, rd, meta)
}

func resourceSecurityGroupRuleDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	_, err := inst.Client.SecurityGroup.DeleteSecurityGroupRule(
		ctx,
		rd.Get("security_group_id").(string),
		[]string{rd.Id()})
	if err != nil && !common.IsDeleted(err) {
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

func resourceSecurityGroupRuleAllDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.SecurityGroup.DeleteSecurityGroupRuleAll(
		ctx,
		rd.Get("security_group_id").(string))
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSecurityGroupBulkRuleRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
