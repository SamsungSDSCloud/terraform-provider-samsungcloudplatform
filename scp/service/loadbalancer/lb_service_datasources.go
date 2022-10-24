package loadbalancer

import (
	"context"
	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/ScpDevTerra/trf-provider/scp/common"
	"github.com/ScpDevTerra/trf-sdk/library/loadbalancer2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

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
