package firewall

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/firewall2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
	"sync"
)

func init() {
	scp.RegisterResource("scp_firewall_rule", ResourceFirewallRule())
	scp.RegisterResource("scp_firewall_bulk_rule", ResourceFirewallBulkRule())
}

func ResourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallRuleCreate,
		ReadContext:   resourceFirewallRuleRead,
		UpdateContext: resourceFirewallRuleUpdate,
		DeleteContext: resourceFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"firewall_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Firewall id",
			},
			"target_id": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "Target firewall resource id",
			},
			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule direction. (IN, OUT, IN_OUT)",
				ValidateFunc: validation.StringInSlice([]string{
					"IN",
					"OUT",
					"IN_OUT",
				}, false),
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule action. (ALLOW, DROP)",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"DROP",
				}, false),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Rule enabled state.",
			},
			"source_addresses_ipv4": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Source ip addresses list",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
				},
			},
			"destination_addresses_ipv4": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Destination ip addresses list",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
				},
			},
			"service": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Firewall Rule service",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol type. (TCP, UDP, ICMP, ALL)",
							ValidateFunc: validation.StringInSlice([]string{
								"TCP",
								"UDP",
								"ICMP",
								"ALL",
								"TCP_ALL",
								"UDP_ALL",
								"ICMP_ALL",
							}, false),
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Port value",
						},
					},
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Rule description. (0 to 100 characters)",
				ValidateFunc: validation.StringLenBetween(0, 100),
			},
			"location_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location Rule id",
			},
			"rule_location_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule location type. (FIRST, BEFORE, AFTER, LAST)",
				ValidateFunc: validation.StringInSlice([]string{
					"FIRST",
					"BEFORE",
					"AFTER",
					"LAST",
				}, false),
			},
		},
		Description: "Provides a Firewall Rule resource.",
	}
}

func ResourceFirewallBulkRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallBulkRuleCreate,
		ReadContext:   resourceFirewallBulkRuleRead,
		DeleteContext: resourceFirewallBulkRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"firewall_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Firewall id",
			},
			"bulk_rule_location_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bulk rule location type",
				ValidateFunc: validation.StringInSlice([]string{
					"FIRST",
					"BEFORE",
					"AFTER",
					"LAST",
				}, false),
			},
			"bulk_rule_location_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Bulk rule location id",
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
							Description: "Rule direction. (IN, OUT, IN_OUT)",
							ValidateFunc: validation.StringInSlice([]string{
								"IN",
								"OUT",
								"IN_OUT",
							}, false),
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule action. (ALLOW, DROP)",
							ValidateFunc: validation.StringInSlice([]string{
								"ALLOW",
								"DROP",
							}, false),
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Rule enabled state.",
						},
						"source_addresses_ipv4": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Source ip addresses cidr list",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"destination_addresses_ipv4": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Destination ip addresses cidr list",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"service": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Firewall Rule service",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Protocol type. (TCP, UDP, ICMP, ALL)",
										ValidateFunc: validation.StringInSlice([]string{
											"TCP",
											"UDP",
											"ICMP",
											"ALL",
											"TCP_ALL",
											"UDP_ALL",
											"ICMP_ALL",
										}, false),
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Port value",
									},
								},
							},
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Rule description. (0 to 100 characters)",
							ValidateFunc: validation.StringLenBetween(0, 100),
						},
					},
				},
			},
		},
		Description: "Provides a Firewall Rule resource.",
	}
}

func expandServices(rd *schema.ResourceData) ([]firewall2.ServiceVo, error) {
	servicesSet := rd.Get("service").(*schema.Set).List()
	// Services
	services := make([]firewall2.ServiceVo, len(servicesSet))
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

func expandSourcesIpv4(rd *schema.ResourceData) []string {
	addressesIpv4List := rd.Get("source_addresses_ipv4").([]interface{})
	addressesIpv4 := make([]string, len(addressesIpv4List))
	for i, valueIpv4 := range addressesIpv4List {
		addressesIpv4[i] = valueIpv4.(string)
	}
	return addressesIpv4
}
func expandDestinationsIpv4(rd *schema.ResourceData) []string {
	addressesIpv4List := rd.Get("destination_addresses_ipv4").([]interface{})
	addressesIpv4 := make([]string, len(addressesIpv4List))
	for i, valueIpv4 := range addressesIpv4List {
		addressesIpv4[i] = valueIpv4.(string)
	}
	return addressesIpv4
}

func expandRules(rd *schema.ResourceData) ([]firewall2.FirewallCreateRuleRequest, error) {

	ruleSet := rd.Get("rule").(*schema.Set).List()
	rules := make([]firewall2.FirewallCreateRuleRequest, len(ruleSet))

	for i, rule := range ruleSet {
		fmt.Println(rule)

		itemObject := rule.(common.HclKeyValueObject)

		if direction, ok := itemObject["direction"]; ok {
			rules[i].RuleDirection = strings.ToUpper(direction.(string))
		}

		if action, ok := itemObject["action"]; ok {
			rules[i].RuleAction = strings.ToUpper(action.(string))
		}

		if enabled, ok := itemObject["enabled"]; ok {
			v := enabled.(bool)
			rules[i].IsRuleEnabled = &v

		}

		if sourceAddressesIpv4, ok := itemObject["source_addresses_ipv4"]; ok {
			s := make([]string, 0)
			for _, sourceAddressIpv4 := range sourceAddressesIpv4.([]interface{}) {
				s = append(s, sourceAddressIpv4.(string))
				fmt.Println(s)
			}
			rules[i].SourceIpAddresses = s
		}

		if destinationAddressesIpv4, ok := itemObject["destination_addresses_ipv4"]; ok {
			s := make([]string, 0)
			for _, destinationAddressIpv4 := range destinationAddressesIpv4.([]interface{}) {
				s = append(s, destinationAddressIpv4.(string))
				fmt.Println(s)
			}
			rules[i].DestinationIpAddresses = s
		}

		if services, ok := itemObject["service"]; ok {

			servicesSet := services.(*schema.Set).List()
			// Services
			services := make([]firewall2.ServiceVo, len(servicesSet))
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
			rules[i].Services = services
		}

		if description, ok := itemObject["description"]; ok {
			rules[i].RuleDescription = description.(string)
		}

		rules[i].RuleLocationType = "LAST"
		rules[i].RuleLocationId = ""
	}

	return rules, nil
}

func addSubnetMask(ipAddress []string) []string {
	var ipList []string
	for _, ip := range ipAddress {
		if !strings.Contains(ip, "/") {
			ipList = append(ipList, ip+"/32")
		} else {
			ipList = append(ipList, ip)
		}
	}
	return ipList
}

var mutex sync.Mutex

func resourceFirewallRuleCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO Implement locking for each firewall_id
	mutex.Lock()
	defer mutex.Unlock()

	firewallId := rd.Get("firewall_id").(string)
	action := rd.Get("action").(string)
	direction := rd.Get("direction").(string)

	isEnabled := rd.Get("enabled").(bool)

	description := rd.Get("description").(string)

	// Addresses
	sourceIpv4 := expandSourcesIpv4(rd)
	destinationIpv4 := expandDestinationsIpv4(rd)

	// Services
	services, err := expandServices(rd)
	if err != nil {
		return diag.FromErr(err)
	}

	inst := meta.(*client.Instance)

	response, _, err := inst.Client.Firewall.CreateFirewallRule(ctx, firewallId, firewall2.FirewallCreateRuleRequest{
		SourceIpAddresses:      sourceIpv4,
		DestinationIpAddresses: destinationIpv4,
		Services:               services,
		RuleDirection:          direction,
		RuleAction:             action,
		IsRuleEnabled:          &isEnabled,
		RuleLocationType:       "LAST",
		RuleLocationId:         "",
		RuleDescription:        description,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	targetState := common.ActiveState

	err = waitForFirewallRuleStatus(ctx, inst.Client, firewallId, response.ResourceId, []string{common.CreatingState}, []string{targetState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ResourceId)

	return resourceFirewallRuleRead(ctx, rd, meta)
}

func resourceFirewallRuleRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	firewallId := rd.Get("firewall_id").(string)

	info, _, err := inst.Client.Firewall.GetFirewallRule(ctx, firewallId, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	info.SourceIpAddresses = addSubnetMask(info.SourceIpAddresses)
	info.DestinationIpAddresses = addSubnetMask(info.DestinationIpAddresses)

	rd.Set("action", info.RuleAction)
	rd.Set("direction", info.RuleDirection)
	rd.Set("enabled", info.IsRuleEnabled)
	rd.Set("description", info.RuleDescription)

	rd.Set("source_addresses_ipv4", info.SourceIpAddresses)
	rd.Set("destination_addresses_ipv4", info.DestinationIpAddresses)

	if *info.IsAllService {
		s := common.HclSetObject{}
		s = append(s, common.HclKeyValueObject{
			"type": "ALL",
		})
		rd.Set("service", s)
	} else {
		s := common.HclSetObject{}
		for _, svc := range info.TcpServices {
			s = append(s, common.HclKeyValueObject{
				"type":  "TCP",
				"value": svc,
			})
		}
		for _, svc := range info.UdpServices {
			s = append(s, common.HclKeyValueObject{
				"type":  "UDP",
				"value": svc,
			})
		}
		for _, svc := range info.IcmpServices {
			s = append(s, common.HclKeyValueObject{
				"type":  "ICMP",
				"value": svc,
			})
		}
		rd.Set("service", s)
	}

	return nil
}

func resourceFirewallRuleUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO Implement locking for each firewall_id
	mutex.Lock()
	defer mutex.Unlock()

	firewallId := rd.Get("firewall_id").(string)
	inst := meta.(*client.Instance)

	isEnabled := rd.Get("enabled").(bool)

	if rd.HasChanges("enabled") {
		_, _, err := inst.Client.Firewall.UpdateFirewallRuleEnable(ctx, firewallId, rd.Id(), isEnabled)
		if err != nil {
			return diag.FromErr(err)
		}
		targetState := common.ActiveState
		err = waitForFirewallRuleStatus(ctx, inst.Client, firewallId, rd.Id(), []string{common.DeployingState}, []string{targetState}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("action", "direction", "description", "source_addresses_ipv4", "destination_addresses_ipv4", "service") {

		action := rd.Get("action").(string)
		direction := rd.Get("direction").(string)
		description := rd.Get("description").(string)

		// Addresses
		sourceIpv4 := expandSourcesIpv4(rd)
		destinationIpv4 := expandDestinationsIpv4(rd)

		// Services
		services, err := expandServices(rd)
		if err != nil {
			return diag.FromErr(err)
		}

		_, _, err = inst.Client.Firewall.UpdateFirewallRule(ctx, firewallId, rd.Id(), firewall2.FirewallRuleUpdateRequest{
			SourceIpAddresses:      sourceIpv4,
			DestinationIpAddresses: destinationIpv4,
			Services:               services,
			RuleDirection:          direction,
			RuleAction:             action,
			RuleDescription:        description,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		err = waitForFirewallRuleStatus(ctx, inst.Client, firewallId, rd.Id(), []string{common.DeployingState}, []string{common.ActiveState}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("location_rule_id", "rule_location_type") {

		locationRuleId := rd.Get("location_rule_id").(string)
		ruleLocationType := rd.Get("rule_location_type").(string)

		_, _, err := inst.Client.Firewall.UpdateFirewallRuleLocation(ctx, firewallId, rd.Id(), firewall2.FirewallRuleChangeLocationRequest{
			LocationRuleId:   locationRuleId,
			RuleLocationType: ruleLocationType,
		})

		if err != nil {
			return diag.FromErr(err)
		}
		err = waitForFirewallRuleStatus(ctx, inst.Client, firewallId, rd.Id(), []string{common.DeployingState}, []string{common.ActiveState}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceFirewallRuleRead(ctx, rd, meta)
}

func resourceFirewallRuleDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO Implement locking for each firewall_id
	mutex.Lock()
	defer mutex.Unlock()

	inst := meta.(*client.Instance)

	firewallId := rd.Get("firewall_id").(string)

	_, _, err := inst.Client.Firewall.DeleteFirewallRule(ctx, firewallId, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitForFirewallRuleStatus(ctx, inst.Client, firewallId, rd.Id(), []string{common.TerminatingState}, []string{common.DeletedState}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceFirewallBulkRuleCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO Implement locking for each firewall_id
	mutex.Lock()
	defer mutex.Unlock()

	firewallId := rd.Get("firewall_id").(string)

	bulkRuleLocationType := rd.Get("bulk_rule_location_type").(string)
	bulkRuleLocationId := rd.Get("bulk_rule_location_id").(string)
	bulkRules, err := expandRules(rd)

	if err != nil {
		return diag.FromErr(err)
	}

	inst := meta.(*client.Instance)

	response, _, err := inst.Client.Firewall.CreateFirewallBulkRule(ctx, firewallId, firewall2.FirewallRuleCreateBulkRequest{
		BulkRuleLocationType:   bulkRuleLocationType,
		BulkRuleLocationId:     bulkRuleLocationId,
		BulkRules:              bulkRules,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	targetState := common.ActiveState

	err = WaitForFirewallStatus(ctx, inst.Client, firewallId, FirewallPendingStates(), []string{targetState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ResourceId)

	return nil
}

func resourceFirewallBulkRuleRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceFirewallBulkRuleDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func waitForFirewallRuleStatus(ctx context.Context, scpClient *client.SCPClient, firewallId string, ruleId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.Firewall.GetFirewallRule(ctx, firewallId, ruleId)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			if c == 403 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			return nil, "", err
		}
		return info, info.RuleState, nil
	})
}
