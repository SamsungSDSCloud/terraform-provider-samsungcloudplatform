package loadbalancer

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/loadbalancer2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func DatasourceLBServiceIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBServiceIpList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id":   {Type: schema.TypeString, Required: true, Description: "Load balancer id"},
			"lb_service_name":    {Type: schema.TypeString, Optional: true, Description: "Load balancer service name"},
			"nat_ip_address":     {Type: schema.TypeString, Optional: true, Description: "Nat ip address"},
			"service_ip_address": {Type: schema.TypeString, Optional: true, Description: "Service ip address"},
			"page":               {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":               {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":               {Type: schema.TypeString, Optional: true, Default: 20, Description: "Sort"},
			"contents":           {Type: schema.TypeList, Optional: true, Description: "Load balancer ip list", Elem: datasourceLbServiceIpElem()},
			"total_count":        {Type: schema.TypeInt, Computed: true, Description: "Total load balancer ip list"},
		},
		Description: "Provides list of Load Balancer service ips",
	}
}

func dataSourceLBServiceIpList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	loadBalancerId := rd.Get("load_balancer_id").(string)
	if len(loadBalancerId) == 0 {
		return diag.Errorf("Load balancer id not found")
	}

	responses, _, err := inst.Client.LoadBalancer.GetLbServiceIpList(ctx, loadBalancerId, &loadbalancer2.LbServiceOpenApiControllerApiGetLoadBalancerServiceIpListOpts{
		LbServiceName:    optional.NewString(rd.Get("lb_service_name").(string)),
		NatIpAddress:     optional.NewString(rd.Get("nat_ip_address").(string)),
		ServiceIpAddress: optional.NewString(rd.Get("service_ip_address").(string)),
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

func datasourceLbServiceIpElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"lb_service_ip_id":    {Type: schema.TypeString, Computed: true, Description: "Id of load balancer service ip"},
			"lb_service_ip_state": {Type: schema.TypeString, Computed: true, Description: "Status of load balancer service"},
			//lbServices
			"nat_ip_address":     {Type: schema.TypeString, Computed: true, Description: "Nat ip address"},
			"nat_ip_id":          {Type: schema.TypeString, Computed: true, Description: "Nat ip id"},
			"service_ip_address": {Type: schema.TypeString, Computed: true, Description: "Service ip address"},
			"service_ip_cidr":    {Type: schema.TypeString, Computed: true, Description: "Service ip cidr"},
			"service_ip_id":      {Type: schema.TypeString, Computed: true, Description: "Service ip id"},
			"service_ip_pool_id": {Type: schema.TypeString, Computed: true, Description: "pool id of service ip"},
		},
	}
}
