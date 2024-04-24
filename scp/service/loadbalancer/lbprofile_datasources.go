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
	scp.RegisterDataSource("scp_lb_profiles", DatasourceLbProfiles())
}

func DatasourceLbProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLbProfileList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id":    {Type: schema.TypeString, Required: true, Description: "Load balancer id"},
			"lb_profile_category": {Type: schema.TypeString, Optional: true, Description: "Load balancer profile category"},
			"lb_profile_name":     {Type: schema.TypeString, Optional: true, Description: "Load balancer profile name"},
			"lb_service_name":     {Type: schema.TypeString, Optional: true, Description: "Load balancer service name"},
			"load_balancer_name":  {Type: schema.TypeString, Optional: true, Description: "Load balancer name"},
			"created_by":          {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"page":                {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":                {Type: schema.TypeString, Optional: true, Description: "Sort"},
			"contents":            {Type: schema.TypeList, Computed: true, Description: "Load balancer profile list", Elem: datasourceLbProfileElem()},
			"total_count":         {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of load balancer profiles",
	}
}

func dataSourceLbProfileList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	loadBalancerId := rd.Get("load_balancer_id").(string)
	if len(loadBalancerId) == 0 {
		return diag.Errorf("Load balancer id not found")
	}

	responses, _, err := inst.Client.LoadBalancer.GetLbProfileList(ctx, loadBalancerId, &loadbalancer2.LbProfileOpenApiControllerApiGetLoadBalancerProfileListOpts{
		LbProfileCategory: optional.NewString(rd.Get("lb_profile_category").(string)),
		LbProfileName:     optional.NewString(rd.Get("lb_profile_name").(string)),
		LbServiceName:     optional.NewString(rd.Get("lb_service_name").(string)),
		LoadBalancerName:  optional.NewString(rd.Get("load_balancer_name").(string)),
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

func datasourceLbProfileElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":          {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"block_id":            {Type: schema.TypeString, Computed: true, Description: "Block id"},
			"lb_profile_category": {Type: schema.TypeString, Computed: true, Description: "Load balancer profile category"},
			"lb_profile_name":     {Type: schema.TypeString, Computed: true, Description: "Load balancer profile name"},
			"lb_profile_state":    {Type: schema.TypeString, Computed: true, Description: "Load balancer profile state"},
			"lb_profile_type":     {Type: schema.TypeString, Computed: true, Description: "Load balancer profile type"},
			"lb_service_names": {Type: schema.TypeList, Computed: true, Description: "Load balancer services' names",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				}},
			"service_zone_id": {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			"created_by":      {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":      {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":     {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":     {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}
