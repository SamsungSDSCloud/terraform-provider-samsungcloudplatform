package loadbalancer

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/loadbalancer2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_lb_server_groups", DatasourceLBServerGroups())
}

func DatasourceLBServerGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id":     {Type: schema.TypeString, Required: true, Description: "Load balancer id"},
			"lb_server_group_name": {Type: schema.TypeString, Optional: true, Description: "Load balancer server group name"},
			"lb_service_name":      {Type: schema.TypeString, Optional: true, Description: "Load balancer Service name"},
			"load_balancer_name":   {Type: schema.TypeString, Optional: true, Description: "Load balancer name"},
			"member_ip_address":    {Type: schema.TypeString, Optional: true, Description: "Member ip address"},
			"created_by":           {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"page":                 {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                 {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":                 {Type: schema.TypeString, Optional: true, Default: 20, Description: "Sort"},
			"contents":             {Type: schema.TypeList, Optional: true, Description: "Load balancer server group list", Elem: datasourceElem()},
			"total_count":          {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of Load Balancer server groups",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	loadBalancerId := rd.Get("load_balancer_id").(string)
	if len(loadBalancerId) == 0 {
		return diag.Errorf("Load balancer id not found")
	}

	responses, _, err := inst.Client.LoadBalancer.GetLbServerGroupList(ctx, loadBalancerId, &loadbalancer2.LbServerGroupOpenApiControllerApiGetLoadBalancerServerGroupListOpts{
		LbServerGroupName: optional.NewString(rd.Get("lb_server_group_name").(string)),
		LbServiceName:     optional.NewString(rd.Get("lb_service_name").(string)),
		LoadBalancerName:  optional.NewString(rd.Get("load_balancer_name").(string)),
		MemberIpAddress:   optional.NewString(rd.Get("member_ip_address").(string)),
		CreatedBy:         optional.NewString(rd.Get("created_by").(string)),
		Page:              optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:              optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:              optional.Interface{},
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

func datasourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":         {Type: schema.TypeString, Computed: true, Description: "Project Id"},
			"block_id":           {Type: schema.TypeString, Computed: true, Description: "Block Id"},
			"lb_server_group_id": {Type: schema.TypeString, Computed: true, Description: "Load balancer server group id"},
			//"lb_server_group_members": {Type: schema.TypeString, Computed: true, Description: "Load Balancer Server Group Name", Elem: datasourceServerGroupMembersElem()},

			"lb_server_group_name":  {Type: schema.TypeString, Computed: true, Description: "Load balancer server group name"},
			"lb_server_group_state": {Type: schema.TypeString, Computed: true, Description: "Load balancer server status"},
			"lb_server_group_type":  {Type: schema.TypeString, Computed: true, Description: "Load balancer server group type"},
			"load_balancer_id":      {Type: schema.TypeString, Computed: true, Description: "Load balancer id"},
			"load_balancer_name":    {Type: schema.TypeString, Computed: true, Description: "Load balancer name"},
			"load_balancer_state":   {Type: schema.TypeString, Computed: true, Description: "Load balancer status"},
			"persistence":           {Type: schema.TypeString, Computed: true, Description: "Persistence"},
			"service_zone_id":       {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			"created_by":            {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":            {Type: schema.TypeString, Computed: true, Description: "Creation Date"},
			"modified_by":           {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":           {Type: schema.TypeString, Computed: true, Description: "Modification Date"},
		},
	}
}

func datasourceServerGroupMembersElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_id":     {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"join_state":                {Type: schema.TypeString, Computed: true, Description: "Load balancer id"},
			"lb_server_group_member_id": {Type: schema.TypeString, Computed: true, Description: "Load balancer server group name"},
			"member_ip_address":         {Type: schema.TypeString, Computed: true, Description: "Load balancer server group name"},
			"member_port":               {Type: schema.TypeInt, Computed: true, Description: "Load balancer server group name"},
			"member_weight":             {Type: schema.TypeInt, Computed: true, Description: "Load balancer service name"},
			"object_id":                 {Type: schema.TypeString, Computed: true, Description: "Load balancer name"},
			"object_name":               {Type: schema.TypeString, Computed: true, Description: "Member ip address"},
			"object_type":               {Type: schema.TypeString, Computed: true, Description: "Object type"},
		},
	}
}
