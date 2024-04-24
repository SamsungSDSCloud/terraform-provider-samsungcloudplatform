package resourcegroup

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_resource_group_resource_types", DatasourceResourceGroupResourceTypes())
	scp.RegisterDataSource("scp_resource_group_service_types", DatasourceResourceGroupServiceTypes())
}

func datasourceResourceTypeElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_type": {Type: schema.TypeString, Computed: true, Description: "Resource type"},
			"service_type":  {Type: schema.TypeString, Computed: true, Description: "Service type"},
			"tag_policy":    {Type: schema.TypeBool, Computed: true, Description: "Tag policy exists"},
		},
	}
}
func DatasourceResourceGroupResourceTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceGroupResourceTypesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"resource_type":  {Type: schema.TypeString, Optional: true, Description: "Resource type"},
			"service_type":   {Type: schema.TypeString, Optional: true, Description: "Service type"},
			"resource_types": {Type: schema.TypeList, Computed: true, Elem: datasourceResourceTypeElem()},
		},
	}
}

func dataSourceResourceGroupResourceTypesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	resourceType := rd.Get("resource_type").(string)
	serviceType := rd.Get("service_type").(string)

	info, err := inst.Client.ResourceGroup.GetResourceTypes(ctx, resourceType, serviceType)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceTypes := common.ConvertStructToMaps(info)

	rd.SetId(uuid.NewV4().String())
	rd.Set("resource_types", resourceTypes)

	return nil
}

func DatasourceResourceGroupServiceTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceGroupServiceTypesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"service_type":  {Type: schema.TypeString, Optional: true, Description: "Service type"},
			"service_types": {Type: schema.TypeList, Computed: true, Description: "Service types", Elem: &schema.Schema{Type: schema.TypeString, Description: "service type"}},
		},
	}
}

func dataSourceResourceGroupServiceTypesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	serviceType := rd.Get("service_type").(string)

	info, err := inst.Client.ResourceGroup.GetServiceTypes(ctx, serviceType)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("service_types", info)

	return nil
}
