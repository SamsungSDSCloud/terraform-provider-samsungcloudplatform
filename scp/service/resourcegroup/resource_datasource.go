package resourcegroup

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/resourcegroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_resources", DatasourceResources())
	scp.RegisterDataSource("scp_resource", DatasourceResource())
	//scp.RegisterDataSource("scp_resource_group_resource_srn", DatasourceResourceSrn())
	scp.RegisterDataSource("scp_resources_in_my_project", DatasourceMyProjectResources())
	scp.RegisterDataSource("scp_resource_in_my_project", DatasourceMyProjectResource())
}

func DatasourceResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":                common.DatasourceFilter(),
			"created_by_id":         {Type: schema.TypeString, Optional: true, Description: "Creator's ID"},
			"display_service_names": {Type: schema.TypeList, Optional: true, Description: "Display service names", Elem: &schema.Schema{Type: schema.TypeString}},
			"include_deleted":       {Type: schema.TypeString, Optional: true, Computed: true, Description: "Includes deleted resources"},
			"location":              {Type: schema.TypeString, Optional: true, Computed: true, Description: "Location"},
			"modified_by_id":        {Type: schema.TypeString, Optional: true, Description: "Modifier's ID"},
			"my_create":             {Type: schema.TypeString, Optional: true, Description: "Whether I created it or not"},
			"partitions":            {Type: schema.TypeList, Optional: true, Description: "Partition list", Elem: &schema.Schema{Type: schema.TypeString}},
			"regions":               {Type: schema.TypeList, Optional: true, Description: "Region list", Elem: &schema.Schema{Type: schema.TypeString}},
			"resource_id":           {Type: schema.TypeString, Optional: true, Description: "Resource ID"},
			"resource_name":         {Type: schema.TypeString, Optional: true, Description: "Resource name"},
			"resource_types":        {Type: schema.TypeList, Optional: true, Description: "Resource type list", Elem: &schema.Schema{Type: schema.TypeString}},
			"service_types":         {Type: schema.TypeList, Optional: true, Description: "Service type list", Elem: &schema.Schema{Type: schema.TypeString}},
			"service_zones":         {Type: schema.TypeList, Optional: true, Description: "Service zone list", Elem: &schema.Schema{Type: schema.TypeString}},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value",
						},
					},
				},
				Description: "Tag list",
			},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "total count"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Resource list", Elem: datasourceResourcesElem()},
		},
	}
}

func datasourceResourcesElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"event_state":       {Type: schema.TypeString, Computed: true, Description: "Event state"},
			"partition":         {Type: schema.TypeString, Computed: true, Description: "Partition"},
			"region":            {Type: schema.TypeString, Computed: true, Description: "Region"},
			"resource_id":       {Type: schema.TypeString, Computed: true, Description: "Resource ID"},
			"resource_name":     {Type: schema.TypeString, Computed: true, Description: "Resource name"},
			"resource_srn":      {Type: schema.TypeString, Computed: true, Description: "Resource SRN"},
			"resource_state":    {Type: schema.TypeString, Computed: true, Description: "Resource state"},
			"resource_type":     {Type: schema.TypeString, Computed: true, Description: "Resource type"},
			"service_type":      {Type: schema.TypeString, Computed: true, Description: "Service type"},
			"zone":              {Type: schema.TypeString, Computed: true, Description: "Service zone"},
			"created_by":        {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_by_name":   {Type: schema.TypeString, Computed: true, Description: "Creator's name"},
			"created_by_email":  {Type: schema.TypeString, Computed: true, Description: "Creator's email"},
			"created_dt":        {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":       {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_by_name":  {Type: schema.TypeString, Computed: true, Description: "Modifier's name"},
			"modified_by_email": {Type: schema.TypeString, Computed: true, Description: "Modifier's email"},
			"modified_dt":       {Type: schema.TypeString, Computed: true, Description: "Modified date"},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag key",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag value",
						},
					},
				},
				Description: "Tag list"},
		},
	}
}

func dataSourceResourcesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, err := inst.Client.ResourceGroup.GetResources(ctx, resourcegroup.ListResourceRequest{
		CreatedById:         rd.Get("created_by_id").(string),
		DisplayServiceNames: rd.Get("display_service_names").([]interface{}),
		IncludeDeleted:      rd.Get("include_deleted").(string),
		Location:            rd.Get("location").(string),
		ModifiedById:        rd.Get("modified_by_id").(string),
		MyCreate:            rd.Get("my_create").(string),
		Partitions:          rd.Get("partitions").([]interface{}),
		Regions:             rd.Get("regions").([]interface{}),
		ResourceId:          rd.Get("resource_id").(string),
		ResourceName:        rd.Get("resource_name").(string),
		ResourceTypes:       rd.Get("resource_types").([]interface{}),
		ServiceTypes:        rd.Get("service_types").([]interface{}),
		ServiceZones:        rd.Get("service_zones").([]interface{}),
		Tags:                rd.Get("tags").([]interface{}),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(info.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceResources().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", info.TotalCount)

	return nil
}

func DatasourceResource() *schema.Resource {
	var resourceResource schema.Resource
	resourceResource.ReadContext = dataSourceResourceRead
	resourceResource.Schema = datasourceResourcesElem().Schema
	resourceResource.Schema["resource_id"] = &schema.Schema{Type: schema.TypeString, Required: true, Description: "Resource ID"}
	resourceResource.Schema["include_deleted"] = &schema.Schema{Type: schema.TypeString, Optional: true, Description: "Include deleted"}

	return &resourceResource
}

func dataSourceResourceRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	resourceId := rd.Get("resource_id").(string)
	includeDeleted := common.GetKeyString(rd, "include_deleted")

	info, err := inst.Client.ResourceGroup.GetResource(ctx, resourceId, includeDeleted)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(info.ResourceId)
	rd.Set("event_state", info.EventState)
	rd.Set("partition", info.Partition)
	rd.Set("region", info.Region)
	rd.Set("resource_name", info.ResourceName)
	rd.Set("resource_srn", info.ResourceSrn)
	rd.Set("resource_state", info.ResourceState)
	rd.Set("resource_type", info.ResourceType)
	rd.Set("service_type", info.ServiceType)
	rd.Set("zone", info.Zone)
	rd.Set("created_by_id", info.CreatedBy)
	rd.Set("created_by_name", info.CreatedByName)
	rd.Set("created_by_email", info.CreatedByEmail)
	rd.Set("created_dt", info.CreatedDt.String())
	rd.Set("modified_by_id", info.ModifiedBy)
	rd.Set("modified_by_name", info.ModifiedByName)
	rd.Set("modified_by_email", info.ModifiedByEmail)
	rd.Set("modified_dt", info.ModifiedDt.String())

	var tags common.HclSetObject
	for _, tag := range info.Tags {
		kv := common.HclKeyValueObject{
			"tag_key":   tag.TagKey,
			"tag_value": tag.TagValue,
		}
		tags = append(tags, kv)
	}
	rd.Set("tags", tags)

	return nil
}

func DatasourceResourceSrn() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceSrnRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"resource_id":     {Type: schema.TypeString, Required: true, Description: "Resource id"},
			"include_deleted": {Type: schema.TypeString, Optional: true, Description: "Includes deleted Resource (ex. Y or N)"},
			"resource_srn":    {Type: schema.TypeString, Computed: true, Description: "Resource srn"},
		},
	}
}

func dataSourceResourceSrnRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	resourceId := rd.Get("resource_id").(string)
	includeDeleted := rd.Get("include_deleted").(string)

	srn, err := inst.Client.ResourceGroup.GetResourceSrn(ctx, resourceId, includeDeleted)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(resourceId)
	rd.Set("resource_srn", srn)

	return nil
}

func DatasourceMyProjectResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMyProjectResourcesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":                common.DatasourceFilter(),
			"created_by_id":         {Type: schema.TypeString, Optional: true, Description: "Creator's ID"},
			"display_service_names": {Type: schema.TypeList, Optional: true, Description: "Display service names", Elem: &schema.Schema{Type: schema.TypeString}},
			"include_deleted":       {Type: schema.TypeString, Optional: true, Computed: true, Description: "Includes deleted resources"},
			"location":              {Type: schema.TypeString, Optional: true, Computed: true, Description: "Location"},
			"modified_by_id":        {Type: schema.TypeString, Optional: true, Description: "Modifier's ID"},
			"my_create":             {Type: schema.TypeString, Optional: true, Description: "Whether I created it or not"},
			"partitions":            {Type: schema.TypeList, Optional: true, Description: "Partition list", Elem: &schema.Schema{Type: schema.TypeString}},
			"project_ids":           {Type: schema.TypeList, Optional: true, Description: "Project ID list", Elem: &schema.Schema{Type: schema.TypeString}},
			"regions":               {Type: schema.TypeList, Optional: true, Description: "Region list", Elem: &schema.Schema{Type: schema.TypeString}},
			"resource_id":           {Type: schema.TypeString, Optional: true, Description: "Resource ID"},
			"resource_name":         {Type: schema.TypeString, Optional: true, Description: "Resource name"},
			"resource_types":        {Type: schema.TypeList, Optional: true, Description: "Resource type list", Elem: &schema.Schema{Type: schema.TypeString}},
			"service_types":         {Type: schema.TypeList, Optional: true, Description: "Service type list", Elem: &schema.Schema{Type: schema.TypeString}},
			"service_zones":         {Type: schema.TypeList, Optional: true, Description: "Service zone list", Elem: &schema.Schema{Type: schema.TypeString}},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value",
						},
					},
				},
				Description: "Tag list",
			},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "total count"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Resource list", Elem: datasourceMyProjectResourcesElem()},
		},
	}
}

func datasourceMyProjectResourcesElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_name":      {Type: schema.TypeString, Computed: true, Description: "Account name"},
			"project_id":        {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"project_name":      {Type: schema.TypeString, Computed: true, Description: "Project name"},
			"event_state":       {Type: schema.TypeString, Computed: true, Description: "Event state"},
			"partition":         {Type: schema.TypeString, Computed: true, Description: "Partition"},
			"region":            {Type: schema.TypeString, Computed: true, Description: "Region"},
			"resource_id":       {Type: schema.TypeString, Computed: true, Description: "Resource ID"},
			"resource_name":     {Type: schema.TypeString, Computed: true, Description: "Resource name"},
			"resource_srn":      {Type: schema.TypeString, Computed: true, Description: "Resource SRN"},
			"resource_state":    {Type: schema.TypeString, Computed: true, Description: "Resource state"},
			"resource_type":     {Type: schema.TypeString, Computed: true, Description: "Resource type"},
			"service_type":      {Type: schema.TypeString, Computed: true, Description: "Service type"},
			"zone":              {Type: schema.TypeString, Computed: true, Description: "Service zone"},
			"created_by":        {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_by_name":   {Type: schema.TypeString, Computed: true, Description: "Creator's name"},
			"created_by_email":  {Type: schema.TypeString, Computed: true, Description: "Creator's email"},
			"created_dt":        {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":       {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_by_name":  {Type: schema.TypeString, Computed: true, Description: "Modifier's name"},
			"modified_by_email": {Type: schema.TypeString, Computed: true, Description: "Modifier's email"},
			"modified_dt":       {Type: schema.TypeString, Computed: true, Description: "Modified date"},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag key",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag value",
						},
					},
				},
				Description: "Tag list"},
		},
	}
}

func dataSourceMyProjectResourcesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	projectIds := common.ToStringList(rd.Get("project_ids").([]interface{}))
	info, err := inst.Client.ResourceGroup.GetMyProjectResources(ctx, projectIds, resourcegroup.ListResourceRequest{
		CreatedById:         rd.Get("created_by_id").(string),
		DisplayServiceNames: rd.Get("display_service_names").([]interface{}),
		IncludeDeleted:      rd.Get("include_deleted").(string),
		Location:            rd.Get("location").(string),
		ModifiedById:        rd.Get("modified_by_id").(string),
		MyCreate:            rd.Get("my_create").(string),
		Partitions:          rd.Get("partitions").([]interface{}),
		Regions:             rd.Get("regions").([]interface{}),
		ResourceId:          rd.Get("resource_id").(string),
		ResourceName:        rd.Get("resource_name").(string),
		ResourceTypes:       rd.Get("resource_types").([]interface{}),
		ServiceTypes:        rd.Get("service_types").([]interface{}),
		ServiceZones:        rd.Get("service_zones").([]interface{}),
		Tags:                rd.Get("tags").([]interface{}),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(info.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceResources().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", info.TotalCount)

	return nil
}

func DatasourceMyProjectResource() *schema.Resource {
	var resourceResource schema.Resource
	resourceResource.ReadContext = dataSourceMyProjectResourceRead
	resourceResource.Schema = datasourceMyProjectResourcesElem().Schema
	resourceResource.Schema["resource_id"] = &schema.Schema{Type: schema.TypeString, Required: true, Description: "Resource ID"}
	resourceResource.Schema["project_id"] = &schema.Schema{Type: schema.TypeString, Required: true, Description: "Project ID"}
	resourceResource.Schema["include_deleted"] = &schema.Schema{Type: schema.TypeString, Optional: true, Description: "Include deleted"}

	return &resourceResource
}

func dataSourceMyProjectResourceRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	resourceId := rd.Get("resource_id").(string)
	projectId := rd.Get("project_id").(string)
	includeDeleted := common.GetKeyString(rd, "include_deleted")

	info, err := inst.Client.ResourceGroup.GetMyProjectResource(ctx, resourceId, projectId, includeDeleted)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(info.ResourceId)
	rd.Set("account_name", info.AccountName)
	rd.Set("project_id", info.ProjectId)
	rd.Set("project_name", info.ProjectName)
	rd.Set("event_state", info.EventState)
	rd.Set("partition", info.Partition)
	rd.Set("region", info.Region)
	rd.Set("resource_name", info.ResourceName)
	rd.Set("resource_srn", info.ResourceSrn)
	rd.Set("resource_state", info.ResourceState)
	rd.Set("resource_type", info.ResourceType)
	rd.Set("service_type", info.ServiceType)
	rd.Set("zone", info.Zone)
	rd.Set("created_by_id", info.CreatedBy)
	rd.Set("created_by_name", info.CreatedByName)
	rd.Set("created_by_email", info.CreatedByEmail)
	rd.Set("created_dt", info.CreatedDt.String())
	rd.Set("modified_by_id", info.ModifiedBy)
	rd.Set("modified_by_name", info.ModifiedByName)
	rd.Set("modified_by_email", info.ModifiedByEmail)
	rd.Set("modified_dt", info.ModifiedDt.String())

	var tags common.HclSetObject
	for _, tag := range info.Tags {
		kv := common.HclKeyValueObject{
			"tag_key":   tag.TagKey,
			"tag_value": tag.TagValue,
		}
		tags = append(tags, kv)
	}
	rd.Set("tags", tags)

	return nil
}
