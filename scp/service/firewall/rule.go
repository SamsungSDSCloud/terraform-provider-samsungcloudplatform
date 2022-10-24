package firewall

import (
	"context"
	"fmt"
	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/ScpDevTerra/trf-provider/scp/common"
	"github.com/ScpDevTerra/trf-sdk/library/firewall2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
	"sync"
)

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
				Description: "Source ip addresses cidr list",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},
			"destination_addresses_ipv4": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Destination ip addresses cidr list",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
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
							ValidateFunc: validation.StringInSlice([]string{
								"TCP",
								"UDP",
								"ICMP",
								"ALL",
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
		IsRuleEnabled:          isEnabled,
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
		return diag.FromErr(err)
	}

	rd.Set("action", info.RuleAction)
	rd.Set("direction", info.RuleDirection)
	rd.Set("enabled", info.IsRuleEnabled)
	rd.Set("description", info.RuleDescription)

	rd.Set("source_addresses_ipv4", info.SourceIpAddresses)
	rd.Set("destination_addresses_ipv4", info.DestinationIpAddresses)

	if info.IsAllService {
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

	return resourceFirewallRuleRead(ctx, rd, meta)
}

func resourceFirewallRuleDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// TODO Implement locking for each firewall_id
	mutex.Lock()
	defer mutex.Unlock()

	inst := meta.(*client.Instance)

	firewallId := rd.Get("firewall_id").(string)

	_, _, err := inst.Client.Firewall.DeleteFirewallRule(ctx, firewallId, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForFirewallRuleStatus(ctx, inst.Client, firewallId, rd.Id(), []string{common.TerminatingState}, []string{common.DeletedState}, false)
	if err != nil {
		return diag.FromErr(err)
	}

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
