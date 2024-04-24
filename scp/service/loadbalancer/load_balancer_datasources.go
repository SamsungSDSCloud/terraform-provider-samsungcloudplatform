package loadbalancer

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/loadbalancer2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_load_balancers", DatasourceLoadBalancers())
}

func DatasourceLoadBalancers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoadBalancerList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_name":  {Type: schema.TypeString, Optional: true, Description: "Load balancer name"},
			"load_balancer_size":  {Type: schema.TypeString, Optional: true, Description: "Size of load balancer to be created (SMALL,MEDIUM,LARGE)"},
			"load_balancer_state": {Type: schema.TypeString, Optional: true, Description: "Load balancer status"},
			"vpc_name":            {Type: schema.TypeString, Optional: true, Description: "Vpc name"},
			"created_by":          {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"page":                {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":                {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":            {Type: schema.TypeList, Optional: true, Description: "Load balancer list", Elem: datasourceLoadBalancerElem()},
			"total_count":         {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of load balancers",
	}
}

func dataSourceLoadBalancerList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.LoadBalancer.GetLoadBalancerList(ctx, &loadbalancer2.LoadBalancerOpenApiControllerApiGetLoadBalancerListOpts{
		LoadBalancerName:  optional.NewString(rd.Get("load_balancer_name").(string)),
		LoadBalancerSize:  optional.NewString(rd.Get("load_balancer_size").(string)),
		LoadBalancerState: optional.NewString(rd.Get("load_balancer_state").(string)),
		VpcName:           optional.NewString(rd.Get("vpc_name").(string)),
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

func datasourceLoadBalancerElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":          {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"block_id":            {Type: schema.TypeString, Computed: true, Description: "Block id"},
			"load_balancer_id":    {Type: schema.TypeString, Computed: true, Description: "Load balancer id"},
			"load_balancer_name":  {Type: schema.TypeString, Computed: true, Description: "Load balancer name"},
			"load_balancer_size":  {Type: schema.TypeString, Computed: true, Description: "Size of load balancer to be created (SMALL,MEDIUM,LARGE)"},
			"load_balancer_state": {Type: schema.TypeString, Computed: true, Description: "Load balancer status"},
			"service_zone_id":     {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			"vpc_id":              {Type: schema.TypeString, Computed: true, Description: "Vpc id"},
			"created_by":          {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":          {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":         {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":         {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}
