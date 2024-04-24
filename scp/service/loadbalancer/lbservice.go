package loadbalancer

import (
	"context"
	"fmt"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	scp.RegisterResource("scp_lb_service", ResourceLbService())
}

func ResourceLbService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLbServiceCreate,
		ReadContext:   resourceLbServiceRead,
		UpdateContext: resourceLbServiceUpdate,
		DeleteContext: resourceLbServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"lb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target Load-Balancer id.",
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: common.ValidateName3to20DashInMiddle,
				Description:      "Name of Load-Balancer Service. (3 to 20 characters with dash in middle)",
			},
			"app_profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Application Profile ID",
			},
			"forwarding_ports": {
				Type:     schema.TypeString,
				Optional: true,
				//ValidateDiagFunc: common.ValidateStringPortRange,
				Description: "Forwarding port numbers. Multiple ports can be inserted using comma and dash. (e.g. 8000-8100,8200)",
			},
			"service_ports": {
				Type:     schema.TypeString,
				Required: true,
				//ValidateDiagFunc: common.ValidateStringPortRange,
				Description: "Servicing port numbers. Multiple ports can be inserted using comma and dash. (e.g. 8000-8100,8200)",
			},
			"layer_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ValidateLbServiceLayerType,
				Description:      "Servicing protocol layer. (L4 for TCP, L7 for HTTP or HTTPS)",
			},
			"protocol": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: common.ValidateProtocol,
				Description:      "Servicing protocol. (TCP, HTTP, HTTPS)",
			},
			"lb_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Server-Group rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lb_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lb_rule_seq": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"lb_server_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Target server-group id.",
						},
						"pattern_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Pattern URL.",
						},
					},
				},
			},
			"persistence": {
				Type:     schema.TypeString,
				Required: true,
				//ValidateDiagFunc: persistence validation.. do i need it?
				Description: "Persistence option. (DISABLED, SOURCE_IP, COOKIE)",
			},
			"persistence_profile_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Persistence target profile id.",
			},
			"service_ipv4": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "Servicing IP address",
			},
			"lb_service_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// NOTE: NOT YET
			"client_certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSL client certification id.",
			},
			"client_ssl_security_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSL client security level.",
			},
			"server_certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSL server certification id.",
			},
			"server_ssl_security_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSL server security level.",
			},
			"use_access_log": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nat_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Wheter to use NAT IP (public IP) or not.",
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "NAT IP attached to LB service IP.",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Load Balancer Service resource.",
	}
}

func expandRules(rd *schema.ResourceData) ([]loadbalancer.LbServiceRule, error) {
	ruleList := rd.Get("lb_rules").([]interface{})

	rules := make([]loadbalancer.LbServiceRule, len(ruleList))
	for i, rule := range ruleList {
		r := rule.(map[string]interface{})
		if t, ok := r["lb_server_group_id"]; ok {
			rules[i].LbServerGroupId = t.(string)
		}
		if t, ok := r["pattern_url"]; ok {
			rules[i].PatternUrl = t.(string)
		}
		if v, ok := r["lb_rule_seq"]; ok {
			rules[i].LbRuleSeq = int32(v.(int))
		}
		if t, ok := r["lb_rule_id"]; ok {
			rules[i].LbRuleId = t.(string)
		}
	}

	return rules, nil
}

func resourceLbServiceCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	loadBalancerId := rd.Get("lb_id").(string)
	lbServiceName := rd.Get("name").(string)
	natActive := rd.Get("nat_active").(bool)

	layerType := rd.Get("layer_type").(string)

	persistence := rd.Get("persistence").(string)
	persistenceProfileId := rd.Get("persistence_profile_id").(string)

	defaultForwardingPorts := rd.Get("forwarding_ports").(string)

	protocol := rd.Get("protocol").(string)

	serviceIpAddr := rd.Get("service_ipv4").(string)
	servicePorts := rd.Get("service_ports").(string)
	serviceIpId := rd.Get("lb_service_ip_id").(string)

	serverCertificateId := rd.Get("server_certificate_id").(string)
	serverSslSecurityLevel := rd.Get("server_ssl_security_level").(string)
	clientCertificateId := rd.Get("client_certificate_id").(string)
	clientSslSecurityLevel := rd.Get("client_ssl_security_level").(string)

	useAccessLog := rd.Get("use_access_log").(bool)

	inst := meta.(*client.Instance)

	// check if name is duplicated
	isNameInvalid, err := inst.Client.LoadBalancer.CheckLbServiceNameDuplicated(ctx, loadBalancerId, lbServiceName)
	if err != nil {
		return diag.FromErr(err)
	}
	if isNameInvalid {
		return diag.Errorf("Input service name is invalid (maybe duplicated) : " + lbServiceName)
	}

	// check if ip-port pair is duplicated
	responses, _, err := inst.Client.LoadBalancer.CheckLbServiceIpPortDuplicated(ctx, loadBalancerId, serviceIpAddr, servicePorts)
	if err != nil {
		return diag.FromErr(err)
	}
	if responses.TotalCount != 0 {
		return diag.Errorf("Input service ip and ports pair is invalid (maybe duplicated) : " + serviceIpAddr + ", " + servicePorts)
	}

	// if layer type is L4, make 1 empty rule and set rule's server group id
	// else if layer type is L7, get rule from input data
	rules, err := expandRules(rd)
	if err != nil {
		return diag.FromErr(err)
	}

	// Setup rule sequence index
	for i, _ := range rules {
		rules[i].LbRuleSeq = int32(i)
	}

	// persistence : none -> no persistence_profile_id
	// persistence : source_it | cookie -> get profile per each option
	//if persistence == "DISABLED" {
	//	persistenceProfileId = ""
	//}

	// app_profile_id : from input, nullable
	//appProfileId := "DEFAULT" // rd.Get("app_profile_id").(string)
	appProfileId := rd.Get("app_profile_id").(string)

	//if len(serviceIpAddr) == 0 {
	//serviceIpAddr = "192.168.1.1"
	//}

	if len(rules) == 0 {
		if layerType == "L4" {
			rules = append(rules, loadbalancer.LbServiceRule{
				LbRuleId:        "",
				LbRuleSeq:       0,
				LbServerGroupId: "",
				PatternUrl:      "",
			})
		} else if layerType == "L7" {
			rules = append(rules, loadbalancer.LbServiceRule{
				LbRuleId:        "",
				LbRuleSeq:       0,
				LbServerGroupId: "",
				PatternUrl:      "/(default)",
			})
		}
	} else {
		if layerType == "L4" {
			if len(rules) > 1 {
				return diag.Errorf("Only one 'lb_rules' is allowed for L4 layer.")
			}
			if len(rules[0].PatternUrl) != 0 {
				return diag.Errorf("No 'pattern_url' is allowed for L4 layer.")
			}
		} else if layerType == "L7" {
			for _, rule := range rules {
				if len(rule.PatternUrl) == 0 {
					return diag.Errorf("Empty 'pattern_url' is not allowed. Use '/(default)' for default url.")
				}
			}
		}
	}

	tags := rd.Get("tags").(map[string]interface{})
	response, err := inst.Client.LoadBalancer.CreateLbService(ctx, loadBalancerId, appProfileId, defaultForwardingPorts, layerType,
		lbServiceName, natActive, persistence, persistenceProfileId, protocol, rules, serviceIpAddr, servicePorts, serviceIpId,
		serverCertificateId, serverSslSecurityLevel, clientCertificateId, clientSslSecurityLevel, useAccessLog, tags)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForLbServiceStatus(ctx, inst.Client, response.ResourceId, loadBalancerId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}
	rd.SetId(response.ResourceId)
	return resourceLbServiceRead(ctx, rd, meta)
}

func resourceLbServiceRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.LoadBalancer.GetLbService(ctx, rd.Id(), rd.Get("lb_id").(string))
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("name", info.LbServiceName)
	rd.Set("app_profile_id", info.ApplicationProfileId)
	rd.Set("client_certificate_id", info.ClientCertificateId)
	rd.Set("client_ssl_security_level", info.ClientSslSecurityLevel)
	rd.Set("server_certificate_id", info.ServerCertificateId)
	rd.Set("server_ssl_security_level", info.ServerSslSecurityLevel)
	rd.Set("forwarding_ports", info.DefaultForwardingPorts)
	rd.Set("layer_type", info.LayerType)
	rd.Set("lb_rules", common.ConvertStructToMaps(info.LbRules))
	rd.Set("lb_service_ip_id", info.LbServiceIpId)
	rd.Set("persistence", info.Persistence)
	rd.Set("persistence_profile_id", info.PersistenceProfileId)
	rd.Set("protocol", info.Protocol)
	rd.Set("service_ipv4", info.ServiceIpAddress)
	rd.Set("service_ports", info.ServicePorts)
	rd.Set("use_access_log", info.UseAccessLog)
	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceLbServiceUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	// lb rules cannot be changed with other fields in one api call
	if rd.HasChanges("lb_rules") {
		lbRules, err := expandRules(rd)
		if err != nil {
			return diag.FromErr(err)
		}

		if len(lbRules) < 1 {
			rule := loadbalancer.LbServiceRule{
				LbRuleSeq:  0,
				PatternUrl: "",
				LbRuleId:   "",
			}
			lbRules = append(lbRules, rule)
		}

		_, err = inst.Client.LoadBalancer.UpdateLbRules(ctx, rd.Id(), rd.Get("lb_id").(string), lbRules)
		if err != nil {
			return diag.FromErr(err)
		}
		err = waitForLbServiceStatus(ctx, inst.Client, rd.Id(), rd.Get("lb_id").(string), common.NetworkProcessingStates(), []string{common.ActiveState}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if rd.HasChanges("app_profile_id", "client_certificate_id", "client_ssl_security_level", "forwarding_ports",
		"persistence", "persistence_profile_id", "server_certificate_id", "server_ssl_security_level", "service_ports", "use_access_log") {

		applicationProfileId := rd.Get("app_profile_id").(string)
		clientCertificateId := rd.Get("client_certificate_id").(string)
		client_ssl_security_level := rd.Get("client_ssl_security_level").(string)
		defaultForwardingPorts := rd.Get("forwarding_ports").(string)
		persistence := rd.Get("persistence").(string)
		persistenceProfileId := rd.Get("persistence_profile_id").(string)
		serverCertificateId := rd.Get("server_certificate_id").(string)
		server_ssl_security_level := rd.Get("server_ssl_security_level").(string)
		servicePorts := rd.Get("service_ports").(string)
		useAccessLog := rd.Get("use_access_log").(bool)

		_, err := inst.Client.LoadBalancer.UpdateLbService(
			ctx,
			rd.Id(), rd.Get("lb_id").(string),
			applicationProfileId,
			clientCertificateId,
			client_ssl_security_level,
			defaultForwardingPorts,
			persistence,
			persistenceProfileId,
			serverCertificateId,
			server_ssl_security_level,
			servicePorts,
			useAccessLog)
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForLbServiceStatus(ctx, inst.Client, rd.Id(), rd.Get("lb_id").(string), common.NetworkProcessingStates(), []string{common.ActiveState}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if rd.HasChanges("nat_active", "public_ip_id") {
		lbServiceIpId := rd.Get("lb_service_ip_id").(string)
		natActive := rd.Get("nat_active").(bool)
		publicIpId := rd.Get("public_ip_id").(string)

		_, _, err := inst.Client.LoadBalancer.AttachNatIpToLoadBalancerServiceIp(ctx, rd.Get("lb_id").(string), lbServiceIpId, natActive, publicIpId)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLbServiceRead(ctx, rd, meta)
}

func resourceLbServiceDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	// First Delete all rules
	info, _, err := inst.Client.LoadBalancer.GetLbService(ctx, rd.Id(), rd.Get("lb_id").(string))
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	// Disconnect rules
	var rules []loadbalancer.LbServiceRule // make([]loadbalancer.LbServiceRule, len(info.LbRules))
	for _, r := range info.LbRules {
		if len(r.LbServerGroupId) == 0 {
			rules = append(rules, loadbalancer.LbServiceRule{
				LbRuleId:        r.LbRuleId,
				LbRuleSeq:       r.LbRuleSeq,
				LbServerGroupId: r.LbServerGroupId,
				PatternUrl:      r.PatternUrl,
			})
		}
	}

	inst.Client.LoadBalancer.UpdateLbRules(ctx, rd.Id(), rd.Get("lb_id").(string), rules)

	inst.Client.LoadBalancer.UpdateLbService(
		ctx,
		rd.Id(), rd.Get("lb_id").(string),
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		*info.UseAccessLog)
	err = waitForLbServiceStatus(ctx, inst.Client, rd.Id(), rd.Get("lb_id").(string), []string{}, []string{"ACTIVE"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = inst.Client.LoadBalancer.DeleteLbService(ctx, rd.Id(), rd.Get("lb_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForLbServiceStatus(ctx, inst.Client, rd.Id(), rd.Get("lb_id").(string), []string{}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ValidateLbServiceLayerType(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	const (
		LbL4 string = "L4"
		LbL7 string = "L7"
	)

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get Value
	value := v.(string)

	// Check size string
	if (value != LbL4) && (value != LbL7) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has invalid value : %s", attrKey, value),
			AttributePath: path,
		})
	}

	return diags
}

func waitForLbServiceStatus(ctx context.Context, scpClient *client.SCPClient, id string, loadBalancerId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.LoadBalancer.GetLbService(ctx, id, loadBalancerId)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.LbServiceState, nil
	})
}
