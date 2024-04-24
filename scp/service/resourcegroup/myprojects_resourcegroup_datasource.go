package resourcegroup

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/resourcegroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_resource_group_in_my_projects", DatasourceResourceGroupInMyProjects())
	scp.RegisterDataSource("scp_resource_group_resources_in_my_projects", DatasourceResourceGroupResourcesInMyProjects())
	scp.RegisterDataSource("scp_resource_groups_in_my_projects", DatasourceResourceGroupsInMyProjects())
}

func datasourceMyProjectsResourceGroupItemElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_name":               {Type: schema.TypeString, Computed: true, Description: "Account name"},
			"project_id":                 {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"project_name":               {Type: schema.TypeString, Computed: true, Description: "Project name"},
			"resource_group_id":          {Type: schema.TypeString, Computed: true, Description: "Resource group id"},
			"resource_group_name":        {Type: schema.TypeString, Computed: true, Description: "Resource group name"},
			"resource_group_description": {Type: schema.TypeString, Computed: true, Description: "Resource group description"},
			"created_by_id":              {Type: schema.TypeString, Computed: true, Description: "The user id which created the resource group"},
			"created_by_name":            {Type: schema.TypeString, Computed: true, Description: "The user name which created the resource group"},
			"created_by_email":           {Type: schema.TypeString, Computed: true, Description: "The user email which created the resource group"},
			"created_dt":                 {Type: schema.TypeString, Computed: true, Description: "The created date of the resource group"},
			"modified_by_id":             {Type: schema.TypeString, Computed: true, Description: "The user id which modified the resource group"},
			"modified_by_name":           {Type: schema.TypeString, Computed: true, Description: "The user name which modified the resource group"},
			"modified_by_email":          {Type: schema.TypeString, Computed: true, Description: "The user email which modified the resource group"},
			"modified_dt":                {Type: schema.TypeString, Computed: true, Description: "The modified date of the resource group"},
		},
	}
}

func DatasourceResourceGroupInMyProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMyProjectsResourceGroupRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"resource_group_id":   {Type: schema.TypeString, Required: true, Description: "Resource group id"},
			"account_name":        {Type: schema.TypeString, Computed: true, Description: "Account name"},
			"project_id":          {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"project_name":        {Type: schema.TypeString, Computed: true, Description: "Project name"},
			"resource_group_name": {Type: schema.TypeString, Computed: true, Description: "Resource group name"},
			"target_resource_tag": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateName1to256DotDashUnderscore,
							Description:      "Tag key",
						},
						"tag_value": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: common.ValidateName1to256DotDashUnderscore,
							Description:      "Tag value",
						},
					},
				},
				Description: "Tag list",
			},
			"target_resource_types":      {Type: schema.TypeList, Optional: true, Computed: true, Description: "Resource group types", Elem: &schema.Schema{Type: schema.TypeString, Description: "type"}},
			"resource_group_description": {Type: schema.TypeString, Optional: true, Computed: true, Description: "Resource group description"},
			"created_by_id":              {Type: schema.TypeString, Optional: true, Description: "The user id which created the resource group"},
			"created_by_name":            {Type: schema.TypeString, Optional: true, Description: "The user name which created the resource group"},
			"created_by_email":           {Type: schema.TypeString, Computed: true, Description: "The user email which created the resource group"},
			"created_dt":                 {Type: schema.TypeString, Computed: true, Description: "The created date of the resource group"},
			"modified_by_id":             {Type: schema.TypeString, Optional: true, Description: "The user id which modified the resource group"},
			"modified_by_name":           {Type: schema.TypeString, Optional: true, Description: "The user name which modified the resource group"},
			"modified_by_email":          {Type: schema.TypeString, Computed: true, Description: "The user email which modified the resource group"},
			"modified_dt":                {Type: schema.TypeString, Computed: true, Description: "The modified date of the resource group"},
		},
	}
}

func dataSourceMyProjectsResourceGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	resourceGroupId := rd.Get("resource_group_id").(string)

	result, err := inst.Client.ResourceGroup.GetResourceGroupInMyProjects(ctx, resourceGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.ResourceGroupId)
	rd.Set("account_name", result.AccountName)
	rd.Set("project_id", result.ProjectId)
	rd.Set("project_name", result.ProjectName)
	rd.Set("resource_group_name", result.ResourceGroupName)
	rd.Set("target_resource_types", result.TargetResourceTypes)
	rd.Set("resource_group_description", result.ResourceGroupDescription)
	rd.Set("created_by_id", result.CreatedById)
	rd.Set("created_by_name", result.CreatedByName)
	rd.Set("created_by_email", result.CreatedByEmail)
	rd.Set("created_dt", result.CreatedDt)
	rd.Set("modified_by_id", result.ModifiedById)
	rd.Set("modified_by_name", result.ModifiedByName)
	rd.Set("modified_by_email", result.ModifiedByEmail)
	rd.Set("modified_dt", result.ModifiedDt)

	var tags common.HclSetObject
	for _, tag := range result.TargetResourceTags {
		kv := common.HclKeyValueObject{
			"tag_key":   tag.TagKey,
			"tag_value": tag.TagValue,
		}
		tags = append(tags, kv)
	}
	rd.Set("target_resource_tag", tags)

	return nil
}

func DatasourceResourceGroupResourcesInMyProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMyProjectsResourceGroupResourcesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"resource_group_id": {Type: schema.TypeString, Required: true, Description: "Resource group id"},
			"created_by_id":     {Type: schema.TypeString, Optional: true, Computed: true, Description: "The user id which created the resource"},
			"modified_by_id":    {Type: schema.TypeString, Optional: true, Computed: true, Description: "The user id which modified the resource"},
			"resource_id":       {Type: schema.TypeString, Optional: true, Description: "Resource id"},
			"resource_name":     {Type: schema.TypeString, Optional: true, Description: "Resource name"},
			"total_count":       {Type: schema.TypeInt, Computed: true, Description: "total count"},
			"contents":          {Type: schema.TypeList, Computed: true, Description: "Resource list", Elem: &schema.Schema{Type: schema.TypeMap, Description: "resource"}},
		},
	}
}

func dataSourceMyProjectsResourceGroupResourcesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	resourceGroupId := rd.Get("resource_group_id").(string)

	info, err := inst.Client.ResourceGroup.GetResourceGroupResourcesList(ctx, resourceGroupId, resourcegroup.ListResourceGroupResourcesRequest{
		CreatedById:  rd.Get("created_by_id").(string),
		ModifiedById: rd.Get("modified_by_id").(string),
		ResourceId:   rd.Get("resource_id").(string),
		ResourceName: rd.Get("resource_name").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	resources := common.HclListObject{}
	for _, resource := range info.Contents {
		kv := common.HclKeyValueObject{}

		kv["event_state"] = resource.EventState
		kv["partition"] = resource.Partition
		kv["region"] = resource.Region
		kv["resource_id"] = resource.ResourceId
		kv["resource_name"] = resource.ResourceName
		kv["resource_srn"] = resource.ResourceSrn
		kv["resource_state"] = resource.ResourceState
		kv["resource_type"] = resource.ResourceType
		kv["service_type"] = resource.ServiceType

		for _, tag := range resource.Tags {
			var tags common.HclSetObject
			kvo := common.HclKeyValueObject{
				"tag_key":   tag.TagKey,
				"tag_value": tag.TagValue,
			}
			tags = append(tags, kvo)
			kv["tags"] = tags
		}
		kv["zone"] = resource.Zone
		kv["created_by_id"] = resource.CreatedBy
		kv["created_by_name"] = resource.CreatedByName
		kv["created_by_email"] = resource.CreatedByEmail
		kv["created_dt"] = resource.CreatedDt.String()
		kv["modified_by_id"] = resource.ModifiedBy
		kv["modified_by_name"] = resource.ModifiedByName
		kv["modified_by_email"] = resource.ModifiedByEmail
		kv["modified_dt"] = resource.ModifiedDt.String()
		resources = append(resources, kv)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", resources)
	rd.Set("total_count", info.TotalCount)
	rd.Set("created_by_id", info.CreatedById)
	rd.Set("modified_by_id", info.ModifiedById)

	return nil
}

func DatasourceResourceGroupsInMyProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMyProjectsResourceGroupsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"created_by_id":       {Type: schema.TypeString, Optional: true, Computed: true, Description: "The user id which created the resource group"},
			"modified_by_id":      {Type: schema.TypeString, Optional: true, Computed: true, Description: "The user id which modified the resource group"},
			"modified_by_email":   {Type: schema.TypeString, Optional: true, Description: "The user email which modified the resource group"},
			"project_ids":         {Type: schema.TypeList, Optional: true, Description: "Project id list", Elem: &schema.Schema{Type: schema.TypeString, Description: "Project id"}},
			"resource_group_name": {Type: schema.TypeString, Optional: true, Description: "Resource group name"},
			"total_count":         {Type: schema.TypeInt, Computed: true, Description: "Total count"},
			"contents":            {Type: schema.TypeList, Computed: true, Description: "Resource group list", Elem: datasourceMyProjectsResourceGroupItemElem()},
		},
	}
}

func dataSourceMyProjectsResourceGroupsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	projectIds := rd.Get("project_ids").([]interface{})

	info, err := inst.Client.ResourceGroup.GetResourceGroupListInMyProjects(ctx, projectIds, resourcegroup.ListResourceGroupRequest{
		CreatedById:       rd.Get("created_by_id").(string),
		ModifiedById:      rd.Get("modified_by_id").(string),
		ModifiedByEmail:   rd.Get("modified_by_email").(string),
		ResourceGroupName: rd.Get("resource_group_name").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(info.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", info.TotalCount)
	rd.Set("created_by_id", info.CreatedById)
	rd.Set("modified_by_id", info.ModifiedById)

	return nil
}
