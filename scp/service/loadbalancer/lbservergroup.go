package loadbalancer

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/common"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func ResourceLbServerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLbServerGroupCreate,
		ReadContext:   resourceLbServerGroupRead,
		UpdateContext: resourceLbServerGroupUpdate,
		DeleteContext: resourceLbServerGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"lb_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: common.ValidateName3to20DashInMiddle,
				Description:      "Load-Balancer server group name. (3 to 20 characters with dash in middle)",
			},
			"algorithm": {
				Type:     schema.TypeString,
				Required: true,
				//ValidateDiagFunc: ValidateOneOfAlgorithm
				Description: "Balancing algorithm. (ROUND_ROBIN, WEIGHTED_ROUND_ROBIN, LEAST_CONNECTION, WEIGHTED_LEAST_CONNECTION, IP_HASH)",
			},
			"server_group_member": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Server-Group members",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"join_state": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Target service joining state. (ENABLED, DISABLED, GRACEFUL_DISABLED)",
						},
						"object_type": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: ValidateLbServerGroupObjectType,
							Description:      "Target object type. (INSTANCE, BAREMETAL, MANUAL)",
						},
						"object_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Target object id (VM server or BareMetal server). This can not be set with 'object_ipv4'. Input resource should be in the same VPC.",
						},
						"object_ipv4": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateDiagFunc: common.ValidateCidrIpv4,
							Description:      "Target object ipv4 for manual setting.",
						},
						"object_port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: common.ValidatePortRange,
							Description:      "Target object port for manual setting. (1 to 65535)",
						},
						"weight": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: ValidateLbServerGroupWeight,
							Description:      "Balancing weight. This is used with when weighted algorithm is set. (1 to 256)",
						},
					},
				},
			},
			"monitor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monitor_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"monitor_count": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"monitor_interval_sec": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"monitor_port": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: common.ValidatePortRange,
			},
			"monitor_timeout_sec": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: common.ValidatePositiveInt,
			},
			"monitor_http_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Monitor http method. (Only HTTP monitor_protocol. GET, POST)",
			},
			"monitor_http_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Monitor http version. (Only HTTP monitor_protocol. 1.0, 1.1)",
			},
			"monitor_http_url": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc : 0 <= str length <= 50
				Description: "Monitor http url path. (Only HTTP monitor_protocol. 0 to 50 alpha-numeric characters with period, dash, underscore)",
			},
			"monitor_http_request_body": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc : 0 <= str length <= 300
				Description: "Request body content. (Only POST monitor_http_method. 0 to 300 byte characters)",
			},
			"monitor_http_response_body": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateDiagFunc : 0 <= str length <= 300
				Description: "Response body content. (Only HTTP monitor_protocol. 0 to 300 byte characters)",
			},
		},
		Description: "Provides a Load Balancer Server Group resource.",
	}
}

func expandMembers(rd *schema.ResourceData) ([]loadbalancer.LbServerGroupMember, error) {
	memberList := rd.Get("server_group_member").([]interface{})
	// Services
	members := make([]loadbalancer.LbServerGroupMember, len(memberList))
	for i, member := range memberList {
		m := member.(map[string]interface{})
		if t, ok := m["join_state"]; ok {
			members[i].JoinState = strings.ToUpper(t.(string))
		} else {
			return nil, fmt.Errorf("wrong input param : " + t.(string))
		}

		if v, ok := m["weight"]; ok {
			members[i].Weight = int32(v.(int))
		} else {
			return nil, fmt.Errorf("weight value is not set")
		}

		if v, ok := m["object_type"]; ok {
			members[i].ObjectType = v.(string)
		} else {
			return nil, fmt.Errorf("object type must be set")
		}

		if v, ok := m["object_port"]; ok {
			members[i].ObjectPort = int32(v.(int))
		} else {
			return nil, fmt.Errorf("object_port must be set")
		}

		// VM or BM
		requiresObjectId := members[i].ObjectType == "INSTANCE" || members[i].ObjectType == "BAREMETAL"
		if t, ok := m["object_id"]; ok && requiresObjectId {
			members[i].ObjectId = t.(string)

			/*
				if v, ok := m["object_ipv4"]; ok {
					if len(v.(string)) > 0 {
						return nil, fmt.Errorf("following field must not be set with object_ipv4")
					}
				}
			*/
			//delete "break", because cannot include 2 virtual machines in lb group
			//break
		}

		// Manual
		if t, ok := m["object_ipv4"]; ok {
			members[i].ObjectIpAddress = t.(string)
		} else {
			return nil, fmt.Errorf("object_ipv4 must be set")
		}
	}
	return members, nil
}

func resourceLbServerGroupCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	loadBalancerId := rd.Get("lb_id").(string)
	name := rd.Get("name").(string)
	algorithm := rd.Get("algorithm").(string)

	members, err := expandMembers(rd)
	if err != nil {
		return diag.FromErr(err)
	}

	for i, m := range members {
		if m.ObjectType == "INSTANCE" {
			vsInfo, _, err := inst.Client.VirtualServer.GetVirtualServer(ctx, m.ObjectId)
			if err != nil {
				return diag.FromErr(err)
			}
			members[i].ObjectIpAddress = vsInfo.Ip
			if len(members[i].ObjectIpAddress) == 0 {
				return diag.Errorf("failed to retrieve ip information from VirtualServer.")
			}
		}
	}

	monitor := loadbalancer.LbServerGroupMonitor{
		HttpMethod:        rd.Get("monitor_http_method").(string),
		HttpVersion:       rd.Get("monitor_http_version").(string),
		LbMonitorCount:    int32(rd.Get("monitor_count").(int)),
		LbMonitorInterval: int32(rd.Get("monitor_interval_sec").(int)),
		LbMonitorPort:     int32(rd.Get("monitor_port").(int)),
		LbMonitorTimeout:  int32(rd.Get("monitor_timeout_sec").(int)),
		LbMonitorUrl:      rd.Get("monitor_http_url").(string),
		Protocol:          rd.Get("monitor_protocol").(string),
		RequestBody:       rd.Get("monitor_http_request_body").(string),
		ResponseBody:      rd.Get("monitor_http_response_body").(string),
	}

	if len(monitor.RequestBody) != 0 && monitor.HttpMethod == "GET" {
		return diag.Errorf("GET method can not have request body")
	}

	tcpMultiplexingEnabled := monitor.Protocol == "HTTP" && monitor.HttpVersion == "1.1"

	isNameInvalid, err := inst.Client.LoadBalancer.CheckServerGroupNameDuplicated(ctx, loadBalancerId, name)
	if err != nil {
		return diag.FromErr(err)
	}
	if isNameInvalid {
		return diag.Errorf("Input server group name is invalid (maybe duplicated) : " + name)
	}

	response, err := inst.Client.LoadBalancer.CreateLbServerGroup(ctx, loadBalancerId, algorithm, name, &monitor, members, tcpMultiplexingEnabled)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForLbServerGroupStatus(ctx, inst.Client, response.ResourceId, loadBalancerId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}
	rd.SetId(response.ResourceId)

	return resourceLbServerGroupRead(ctx, rd, meta)
}

func resourceLbServerGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.LoadBalancer.GetLbServerGroup(ctx, rd.Id(), rd.Get("lb_id").(string))
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.Set("name", info.LbServerGroupName)
	rd.Set("algorithm", info.LbServerGroupAlgorithm)
	rd.Set("server_group_member", info.LbServerGroupMembers)
	rd.Set("monitor_id", info.LbMonitor.LbMonitorId)
	rd.Set("monitor_http_method", info.LbMonitor.HttpMethod)
	rd.Set("monitor_count", info.LbMonitor.LbMonitorCount)
	rd.Set("monitor_interval", info.LbMonitor.LbMonitorInterval)
	rd.Set("monitor_port", info.LbMonitor.LbMonitorPort)
	rd.Set("monitor_timeout", info.LbMonitor.LbMonitorPort)
	rd.Set("monitor_http_url", info.LbMonitor.LbMonitorUrl)
	rd.Set("monitor_protocol", info.LbMonitor.Protocol)
	rd.Set("monitor_http_request_body", info.LbMonitor.RequestBody)
	rd.Set("monitor_http_response_body", info.LbMonitor.ResponseBody)
	//rd.Set("tcp_multiplexing_enabled", info.TcpMultiplexingEnabled)

	return nil
}

func resourceLbServerGroupUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	if rd.HasChanges(
		"algorithm",
		"server_group_member",
		"monitor_http_method",
		"monitor_http_version",
		"monitor_count",
		"monitor_interval_sec",
		"monitor_port",
		"monitor_timeout_sec",
		"monitor_http_url",
		"monitor_protocol",
		"monitor_http_request_body",
		"monitor_http_response_body") {

		members, err := expandMembers(rd)
		if err != nil {
			return diag.FromErr(err)
		}

		monitor := loadbalancer.LbServerGroupMonitor{
			HttpMethod:        rd.Get("monitor_http_method").(string),
			HttpVersion:       rd.Get("monitor_http_version").(string),
			LbMonitorCount:    int32(rd.Get("monitor_count").(int)),
			LbMonitorInterval: int32(rd.Get("monitor_interval_sec").(int)),
			LbMonitorPort:     int32(rd.Get("monitor_port").(int)),
			LbMonitorTimeout:  int32(rd.Get("monitor_timeout_sec").(int)),
			LbMonitorUrl:      rd.Get("monitor_http_url").(string),
			Protocol:          rd.Get("monitor_protocol").(string),
			RequestBody:       rd.Get("monitor_http_request_body").(string),
			ResponseBody:      rd.Get("monitor_http_response_body").(string),
		}

		tcpMultiplexingEnabled := true
		//if rd.Get("monitor_http_version").(string) == "1.1" {
		//	tcpMultiplexingEnabled = true
		//}

		inst.Client.LoadBalancer.UpdateLbServerGroup(
			ctx,
			tcpMultiplexingEnabled,
			rd.Get("algorithm").(string),
			rd.Id(),
			rd.Get("lb_id").(string),
			&monitor,
			members)
	}

	return resourceLbServerGroupRead(ctx, rd, meta)
}

func resourceLbServerGroupDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	_, err := inst.Client.LoadBalancer.DeleteLbServerGroup(ctx, rd.Id(), rd.Get("lb_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForLbServerGroupStatus(ctx, inst.Client, rd.Id(), rd.Get("lb_id").(string), []string{}, []string{"DELETED"}, false)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ValidateLbServerGroupWeight(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	value := int32(v.(int))

	err := common.CheckInt32Range(value, 1, 256)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateLbServerGroupObjectType(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	const (
		ObjectTypeInstance  string = "INSTANCE"
		ObjectTypeBareMetal string = "BAREMETAL"
		ObjectTypeManual    string = "MANUAL"
	)

	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	value := v.(string)

	if (value != ObjectTypeInstance) && (value != ObjectTypeBareMetal) && (value != ObjectTypeManual) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has invalid value : %s", attrKey, value),
			AttributePath: path,
		})
	}

	return diags
}

func waitForLbServerGroupStatus(ctx context.Context, scpClient *client.SCPClient, id string, loadBalancerId string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.LoadBalancer.GetLbServerGroup(ctx, id, loadBalancerId)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.LbServerGroupState, nil
	})
}
