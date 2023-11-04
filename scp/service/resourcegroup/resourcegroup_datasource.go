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
	scp.RegisterDataSource("scp_resource_group", DatasourceResourceGroup())
	scp.RegisterDataSource("scp_resource_group_resources", DatasourceResourceGroupResources())
	scp.RegisterDataSource("scp_resource_groups", DatasourceResourceGroups())
}

func DatasourceResourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceGroupRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"resource_group_id":   {Type: schema.TypeString, Required: true, Description: "Resource group id"},
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

func dataSourceResourceGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	resourceGroupId := rd.Get("resource_group_id").(string)

	result, err := inst.Client.ResourceGroup.GetResourceGroup(ctx, resourceGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.ResourceGroupId)
	rd.Set("resource_group_name", result.ResourceGroupName)
	rd.Set("target_resource_types", result.TargetResourceTypes)
	rd.Set("resource_group_description", result.ResourceGroupDescription)
	rd.Set("created_by_id", result.CreatedById)
	rd.Set("created_by_name", result.CreatedByName)
	rd.Set("created_by_email", result.CreatedByEmail)
	rd.Set("created_dt", result.CreatedDt.String())
	rd.Set("modified_by_id", result.ModifiedById)
	rd.Set("modified_by_name", result.ModifiedByName)
	rd.Set("modified_by_email", result.ModifiedByEmail)
	rd.Set("modified_dt", result.ModifiedDt.String())

	var tags common.HclSetObject
	for _, tag := range result.TargetResourceTag {
		kv := common.HclKeyValueObject{
			"tag_key":   tag.TagKey,
			"tag_value": tag.TagValue,
		}
		tags = append(tags, kv)
	}
	rd.Set("target_resource_tag", tags)

	return nil
}

func DatasourceResourceGroupResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceGroupResourcesRead,
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
			"contents":          {Type: schema.TypeList, Computed: true, Description: "Resource list", Elem: datasourceResourceGroupResourcesElem()},
		},
	}
}

func datasourceResourceGroupResourcesElem() *schema.Resource {
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
			"created_by":        {Type: schema.TypeString, Computed: true, Description: "Creator ID"},
			"created_by_name":   {Type: schema.TypeString, Computed: true, Description: "Creator name"},
			"created_by_email":  {Type: schema.TypeString, Computed: true, Description: "Creator email"},
			"created_dt":        {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":       {Type: schema.TypeString, Computed: true, Description: "Modifier ID"},
			"modified_by_name":  {Type: schema.TypeString, Computed: true, Description: "Modifier name"},
			"modified_by_email": {Type: schema.TypeString, Computed: true, Description: "Modifier email"},
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
func dataSourceResourceGroupResourcesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	contents := common.ConvertStructToMaps(info.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", info.TotalCount)

	return nil
}

func datasourceResourceGroupItemElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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

func DatasourceResourceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceGroupsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"created_by_id":       {Type: schema.TypeString, Optional: true, Computed: true, Description: "The user id which created the resource group"},
			"modified_by_id":      {Type: schema.TypeString, Optional: true, Computed: true, Description: "The user id which modified the resource group"},
			"modified_by_email":   {Type: schema.TypeString, Optional: true, Description: "The user email which modified the resource group"},
			"resource_group_name": {Type: schema.TypeString, Optional: true, Description: "Resource group name"},
			"total_count":         {Type: schema.TypeInt, Computed: true, Description: "total count"},
			"contents":            {Type: schema.TypeList, Computed: true, Description: "Resource group list", Elem: datasourceResourceGroupItemElem()},
		},
	}
}

func dataSourceResourceGroupsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, err := inst.Client.ResourceGroup.GetResourceGroupList(ctx, resourcegroup.ListResourceGroupRequest{
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
