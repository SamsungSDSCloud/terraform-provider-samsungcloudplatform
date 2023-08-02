package loadbalancer

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/loadbalancer2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_lb_services", DatasourceLBServices())
	scp.RegisterDataSource("scp_lb_services_connectable_to_asg", DatasourceLBServicesConnectableToAsg())
	scp.RegisterDataSource("scp_lb_services_connected_to_asg", DatasourceLBServicesConnectedToAsg())
}

func DatasourceLBServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBServiceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id":   {Type: schema.TypeString, Required: true, Description: "Load balancer id"},
			"layer_type":         {Type: schema.TypeString, Optional: true, Description: "Protocol layer type (L4, L7)"},
			"lb_service_name":    {Type: schema.TypeString, Optional: true, Description: "Load balancer service Name"},
			"load_balancer_name": {Type: schema.TypeString, Optional: true, Description: "Load balancer name"},
			"protocol":           {Type: schema.TypeString, Optional: true, Description: "The file storage protocol type to create (NFS, CIFS)"},
			"service_ip_address": {Type: schema.TypeString, Optional: true, Description: "Service ip address"},
			"status_check":       {Type: schema.TypeBool, Optional: true, Description: "check status"},
			"page":               {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":               {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":               {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"created_by":         {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"contents":           {Type: schema.TypeList, Optional: true, Description: "Load balancer service list", Elem: datasourceLbServiceElem()},
			"total_count":        {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of Load Balancer services",
	}
}

func dataSourceLBServiceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	loadBalancerId := rd.Get("load_balancer_id").(string)
	if len(loadBalancerId) == 0 {
		return diag.Errorf("Load balancer id not found")
	}

	responses, _, err := inst.Client.LoadBalancer.GetLbServiceList(ctx, loadBalancerId, &loadbalancer2.LbServiceOpenApiControllerApiGetLoadBalancerServiceListOpts{
		LayerType:        optional.NewString(rd.Get("layer_type").(string)),
		LbServiceName:    optional.NewString(rd.Get("lb_service_name").(string)),
		LoadBalancerName: optional.NewString(rd.Get("load_balancer_name").(string)),
		Protocol:         optional.NewString(rd.Get("protocol").(string)),
		ServiceIpAddress: optional.NewString(rd.Get("service_ip_address").(string)),
		StatusCheck:      optional.NewBool(rd.Get("status_check").(bool)),
		CreatedBy:        optional.NewString(rd.Get("created_by").(string)),
		Page:             optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:             optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:             optional.Interface{},
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceLbServiceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":               {Type: schema.TypeString, Computed: true, Description: "Load balancer service ip id"},
			"block_id":                 {Type: schema.TypeString, Computed: true, Description: "Block id of this region"},
			"default_forwarding_ports": {Type: schema.TypeString, Computed: true, Description: "Default forwarding ports"},
			"layer_type":               {Type: schema.TypeString, Computed: true, Description: "Protocol layer type (L4, L7)"},
			"lb_service_id":            {Type: schema.TypeString, Computed: true, Description: "Load balancer service id"},
			"lb_service_ip_id":         {Type: schema.TypeString, Computed: true, Description: "Load balancer service ip id"},
			"lb_service_name":          {Type: schema.TypeString, Computed: true, Description: "Load balancer service name"},
			"lb_service_state":         {Type: schema.TypeString, Computed: true, Description: "Load balancer service status"},
			"nat_ip_address":           {Type: schema.TypeString, Computed: true, Description: "Nat ip address"},
			"protocol":                 {Type: schema.TypeString, Computed: true, Description: "Protocol"},
			"service_ip_address":       {Type: schema.TypeString, Computed: true, Description: "Service ip address"},
			"service_ports":            {Type: schema.TypeString, Computed: true, Description: "Service ports"},
			"service_zone_id":          {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			"created_by":               {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":               {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":              {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":              {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func DatasourceLBServicesConnectableToAsg() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBServiceConnectableToAsgList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id":      {Type: schema.TypeString, Required: true, Description: "VPC ID"},
			"contents":    {Type: schema.TypeList, Optional: true, Description: "Load balancer service connectable to asg list", Elem: datasourceLBServiceConnectableOrConnectedToAsgElem()},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of load balancer service Connectable to ASG",
	}
}

func DatasourceLBServicesConnectedToAsg() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBServiceConnectedToAsgList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_id": {Type: schema.TypeString, Required: true, Description: "ASG ID"},
			"contents":              {Type: schema.TypeList, Optional: true, Description: "Load balancer service connected to asg list", Elem: datasourceLBServiceConnectableOrConnectedToAsgElem()},
			"total_count":           {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of load balancer service Connected to ASG",
	}
}

func dataSourceLBServiceConnectableToAsgList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	VpcId := rd.Get("vpc_id").(string)

	responses, _, err := inst.Client.LoadBalancer.GetLoadBalancerServiceConnectableToAsgList(ctx, VpcId)

	if err != nil {
		return diag.FromErr(err)
	}

	contents := convertLbServiceConnectableOrConnectedToAsgListToHclSet(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func dataSourceLBServiceConnectedToAsgList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	AutoScalingGroupId := rd.Get("auto_scaling_group_id").(string)

	responses, _, err := inst.Client.LoadBalancer.GetLoadBalancerServiceConnectedToAsgList(ctx, AutoScalingGroupId)

	if err != nil {
		return diag.FromErr(err)
	}

	contents := convertLbServiceConnectableOrConnectedToAsgListToHclSet(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func convertLbServiceConnectableOrConnectedToAsgListToHclSet(lbServices []loadbalancer2.LbServiceForAsgResponse) common.HclSetObject {
	var lbServiceList common.HclSetObject
	for _, lbService := range lbServices {

		var lbRuleList common.HclListObject
		for _, lbRule := range lbService.LbRules {
			lbRuleKv := common.HclKeyValueObject{
				common.ToSnakeCase("AutoScaleGroupIds"): lbRule.AutoScaleGroupIds,
				common.ToSnakeCase("LbRuleId"):          lbRule.LbRuleId,
				common.ToSnakeCase("LbServerGroupId"):   lbRule.LbServerGroupId,
				common.ToSnakeCase("PatternUrl"):        lbRule.PatternUrl,
				common.ToSnakeCase("Seq"):               lbRule.Seq,
			}

			lbRuleList = append(lbRuleList, lbRuleKv)
		}

		kv := common.HclKeyValueObject{
			common.ToSnakeCase("AutoScaleGroupIds"):      lbService.AutoScaleGroupIds,
			common.ToSnakeCase("DefaultForwardingPorts"): lbService.DefaultForwardingPorts,
			common.ToSnakeCase("LayerType"):              lbService.LayerType,
			common.ToSnakeCase("LbRules"):                lbRuleList,
			common.ToSnakeCase("LbServiceId"):            lbService.LbServiceId,
			common.ToSnakeCase("LbServiceIpAddress"):     lbService.LbServiceIpAddress,
			common.ToSnakeCase("LbServiceName"):          lbService.LbServiceName,
			common.ToSnakeCase("LbServiceState"):         lbService.LbServiceState,
			common.ToSnakeCase("LoadBalancerId"):         lbService.LoadBalancerId,
			common.ToSnakeCase("LoadBalancerName"):       lbService.LoadBalancerName,
			common.ToSnakeCase("NatIpAddress"):           lbService.NatIpAddress,
			common.ToSnakeCase("Persistence"):            lbService.Persistence,
			common.ToSnakeCase("Protocol"):               lbService.Protocol,
			common.ToSnakeCase("ServicePorts"):           lbService.ServicePorts,
		}

		lbServiceList = append(lbServiceList, kv)
	}
	return lbServiceList
}

func datasourceLBServiceConnectableOrConnectedToAsgElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("AutoScaleGroupIds"): {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Auto Scaling Group ID",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			common.ToSnakeCase("DefaultForwardingPorts"): {Type: schema.TypeString, Computed: true, Description: "Default Forwarding Ports"},
			common.ToSnakeCase("LayerType"):              {Type: schema.TypeString, Computed: true, Description: "Layer Type"},
			common.ToSnakeCase("LbRules"): {Type: schema.TypeList, Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						common.ToSnakeCase("AutoScaleGroupIds"): {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Auto Scaling Group ID",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						common.ToSnakeCase("LbRuleId"):        {Type: schema.TypeString, Computed: true, Description: "LB Rule ID"},
						common.ToSnakeCase("LbServerGroupId"): {Type: schema.TypeString, Computed: true, Description: "LB Server Group ID"},
						common.ToSnakeCase("PatternUrl"):      {Type: schema.TypeString, Computed: true, Description: "LB Rule Pattern URL"},
						common.ToSnakeCase("Seq"):             {Type: schema.TypeInt, Computed: true, Description: "LB Rule Sequence"},
					},
					Description: "LB Rules",
				},
			},
			common.ToSnakeCase("LbServiceId"):        {Type: schema.TypeString, Computed: true, Description: "LB Service ID"},
			common.ToSnakeCase("LbServiceIpAddress"): {Type: schema.TypeString, Computed: true, Description: "LB Service IP Address"},
			common.ToSnakeCase("LbServiceName"):      {Type: schema.TypeString, Computed: true, Description: "LB Service Name"},
			common.ToSnakeCase("LbServiceState"):     {Type: schema.TypeString, Computed: true, Description: "LB Service State"},
			common.ToSnakeCase("LoadBalancerId"):     {Type: schema.TypeString, Computed: true, Description: "Load Balancer ID"},
			common.ToSnakeCase("LoadBalancerName"):   {Type: schema.TypeString, Computed: true, Description: "Load Balancer Name"},
			common.ToSnakeCase("NatIpAddress"):       {Type: schema.TypeString, Computed: true, Description: "NAT IP Address"},
			common.ToSnakeCase("Persistence"):        {Type: schema.TypeString, Computed: true, Description: "Persistence Type"},
			common.ToSnakeCase("Protocol"):           {Type: schema.TypeString, Computed: true, Description: "Protocol"},
			common.ToSnakeCase("ServicePorts"):       {Type: schema.TypeString, Computed: true, Description: "Service Ports"},
		},
		Description: "Provides list of Load Balancer services Connectable Or Connected to ASG",
	}
}
