package project

import (
	"context"

	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client/project"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func DatasourceProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"access_level":          {Type: schema.TypeString, Optional: true, Description: "Access level"},
			"action_name":           {Type: schema.TypeString, Optional: true, Description: "Action name"},
			"cmp_service_name":      {Type: schema.TypeString, Optional: true, Description: "CMP service name"},
			"is_user_authorization": {Type: schema.TypeBool, Optional: true, Description: "Check whether to have user authorization or not"},
			"contents":              {Type: schema.TypeList, Optional: true, Description: "Project info list", Elem: datasourceElem()},
			"total_count":           {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of my projects.",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	requestParam := project.ListProjectRequest{
		AccessLevel:         rd.Get("access_level").(string),
		ActionName:          rd.Get("action_name").(string),
		CmpServiceName:      rd.Get("cmp_service_name").(string),
		IsUserAuthorization: rd.Get("is_user_authorization").(bool),
	}

	responses, err := inst.Client.Project.GetProjectList(ctx, requestParam)
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
			"project_id":                {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"account_code":              {Type: schema.TypeString, Computed: true, Description: "Account code"},
			"account_id":                {Type: schema.TypeString, Computed: true, Description: "Account id"},
			"account_name":              {Type: schema.TypeString, Computed: true, Description: "Account name"},
			"account_type":              {Type: schema.TypeString, Computed: true, Description: "Account type"},
			"billing_organization_code": {Type: schema.TypeString, Computed: true, Description: "Billing organization code"},
			"billing_unit":              {Type: schema.TypeString, Computed: true, Description: "Billing unit"},
			"business_category_id":      {Type: schema.TypeString, Computed: true, Description: "Business category id"},
			"business_group_id":         {Type: schema.TypeString, Computed: true, Description: "Business group id"},
			"company_id":                {Type: schema.TypeString, Computed: true, Description: "Company id"},
			"created_dt_str":            {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"fixed_cost_amount":         {Type: schema.TypeString, Computed: true, Description: "Fixed cost amount"},
			"fixed_exchange_rate":       {Type: schema.TypeString, Computed: true, Description: "Fixed exchange rate"},
			"is_fixed_exchange_rate":    {Type: schema.TypeString, Computed: true, Description: "Check whether to be fixed exchange rate or not"},
			"modified_dt_str":           {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			"organization_id":           {Type: schema.TypeString, Computed: true, Description: "Organization id"},
			"owner_id":                  {Type: schema.TypeString, Computed: true, Description: "Owner id"},
			"owner_name":                {Type: schema.TypeString, Computed: true, Description: "Owner name"},
			"project_attributes":        {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Project attributes"},
			"project_budget":            {Type: schema.TypeInt, Computed: true, Description: "Project budget"},
			"project_name":              {Type: schema.TypeString, Computed: true, Description: "Project name"},
			"project_state":             {Type: schema.TypeString, Computed: true, Description: "Project status"},
			"start_dt":                  {Type: schema.TypeString, Computed: true, Description: "Project started date"},
			"start_dt_str":              {Type: schema.TypeString, Computed: true, Description: "Project started date"},
			"project_description":       {Type: schema.TypeString, Computed: true, Description: "Project description"},
			"created_by":                {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":                {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":               {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":               {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}

func DatasourceAccountsByMyProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountsByMyProject,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"access_level":          {Type: schema.TypeString, Optional: true, Description: "Access Level"},
			"action_name":           {Type: schema.TypeString, Optional: true, Description: "Action Name"},
			"cmp_service_name":      {Type: schema.TypeString, Optional: true, Description: "CMP Service Name"},
			"is_user_authorization": {Type: schema.TypeBool, Optional: true, Description: "Is User Authorization"},
			"my_project":            {Type: schema.TypeBool, Optional: true, Description: "My Project"},
			"contents":              {Type: schema.TypeList, Computed: true, Description: "Block Storage list", Elem: datasourceAccountsByMyProjectElem()},
			"total_count":           {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceAccountsByMyProject(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	requestParam := project.ListAccountRequest{
		AccessLevel:         rd.Get("access_level").(string),
		ActionName:          rd.Get("action_name").(string),
		CmpServiceName:      rd.Get("cmp_service_name").(string),
		IsUserAuthorization: rd.Get("is_user_authorization").(bool),
		MyProject:           rd.Get("my_project").(bool),
	}

	responses, err := inst.Client.Project.GetAccountList(ctx, requestParam)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceAccountsByMyProjectElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":           {Type: schema.TypeString, Computed: true},
			"account_id":           {Type: schema.TypeString, Computed: true},
			"account_name":         {Type: schema.TypeString, Computed: true},
			"account_type":         {Type: schema.TypeString, Computed: true},
			"account_user_id":      {Type: schema.TypeString, Computed: true},
			"is_project_creatable": {Type: schema.TypeBool, Computed: true},
			"network_type":         {Type: schema.TypeString, Computed: true},
			"created_by":           {Type: schema.TypeString, Computed: true},
			"created_dt":           {Type: schema.TypeString, Computed: true},
			"modified_by":          {Type: schema.TypeString, Computed: true},
			"modified_dt":          {Type: schema.TypeString, Computed: true},
		},
	}
}

func DatasourceProjectDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectDetail,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ProjectId"): {Type: schema.TypeString, Optional: true, Description: "Project ID"},
			"content":                       {Type: schema.TypeString, Computed: true, Description: "Content"},
		},
	}
}

func dataSourceProjectDetail(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	response, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", common.ToMap(response))

	return nil
}
