package project

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client/project"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_projects", DatasourceProjects())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_project", DatasourceProject())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_project_zones", DatasourceProjectZones())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_project_products", DatasourceProjectProducts())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_project_product_resources", DatasourceProjectProductResources())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_project_user_products_resources", DatasourceProjectUserProductsResources())
}

func DatasourceProject() *schema.Resource {

	var projectResource schema.Resource
	projectResource.ReadContext = datasourceProjectDetailRead
	projectResource.Schema = datasourceProjectsElem().Schema
	projectResource.Schema["budget"] = &schema.Schema{
		Type:        schema.TypeSet,
		Computed:    true,
		Description: "Budget information",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"budget_amount":            {Type: schema.TypeFloat, Computed: true, Description: "Budget amount"},
				"create_prevent_threshold": {Type: schema.TypeInt, Computed: true, Description: "Creation prevention threshold"},
				"is_budget_used":           {Type: schema.TypeBool, Computed: true, Description: "Budget use"},
				"is_create_prevent":        {Type: schema.TypeBool, Computed: true, Description: "Prevent creation"},
				"request_guide":            {Type: schema.TypeString, Computed: true, Description: "Request guide"},
				"alarm_thresholds":         {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeInt}, Description: "List of alarm thresholds"},
			},
		},
	}
	projectResource.Description = "Provides detailed information about the project defined in configuration"

	return &projectResource
}

func datasourceProjectDetailRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	result, err := inst.Client.Project.GetProjectInfo(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.ProjectId)
	mm := common.ToMap(result)

	for v, m := range mm {
		if v != "budget" {
			rd.Set(v, m)
		}
	}

	var budget common.HclSetObject
	if result.Budget != nil {
		kv := common.HclKeyValueObject{
			"budget_amount":            result.Budget.BudgetAmount,
			"create_prevent_threshold": result.Budget.CreatePreventThreshold,
			"is_budget_used":           result.Budget.IsBudgetUsed,
			"is_create_prevent":        result.Budget.IsCreatePrevent,
			"request_guide":            result.Budget.RequestGuide,
			"alarm_thresholds":         result.Budget.AlarmThresholds,
		}
		budget = append(budget, kv)
		//rd.Set("budget", kv)
		rd.Set("budget", budget)
	}

	return nil
}

func DatasourceProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceProjectList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_name":            {Type: schema.TypeString, Optional: true, Description: "Account name"},
			"bill_year_month":         {Type: schema.TypeString, Optional: true, Description: "Billing year and month"},
			"is_billing_info_demand":  {Type: schema.TypeBool, Optional: true, Description: "Whether to provide billing information"},
			"is_resource_info_demand": {Type: schema.TypeBool, Optional: true, Description: "Whether to provide resource information"},
			"is_user_info_demand":     {Type: schema.TypeBool, Optional: true, Description: "Whether to provide user information"},
			"project_name":            {Type: schema.TypeString, Optional: true, Description: "Project name"},
			"created_by_email":        {Type: schema.TypeString, Optional: true, Description: "Creator's email"},
			"contents":                {Type: schema.TypeList, Computed: true, Description: "Project info list", Elem: datasourceProjectsElem()},
			"total_count":             {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of my projects.",
	}
}

func datasourceProjectList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	request := project.ListProjectRequest{
		AccountName:    rd.Get("account_name").(string),
		BillYearMonth:  rd.Get("bill_year_month").(string),
		ProjectName:    rd.Get("project_name").(string),
		CreatedByEmail: rd.Get("created_by_email").(string),
	}

	if b, ok := rd.GetOk("is_billing_info_demand"); ok {
		request.IsBillingInfoDemand = b.(bool)
	}
	if b, ok := rd.GetOk("is_resource_info_demand"); ok {
		request.IsBillingInfoDemand = b.(bool)
	}
	if b, ok := rd.GetOk("is_user_info_demand"); ok {
		request.IsBillingInfoDemand = b.(bool)
	}

	responses, err := inst.Client.Project.GetProjectList(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceProjectsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":                {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"account_admin_email":       {Type: schema.TypeString, Computed: true, Description: "Account administrator's email"},
			"account_admin_name":        {Type: schema.TypeString, Computed: true, Description: "Account administrator's name"},
			"account_id":                {Type: schema.TypeString, Computed: true, Description: "Account id"},
			"account_name":              {Type: schema.TypeString, Computed: true, Description: "Account name"},
			"account_type":              {Type: schema.TypeString, Computed: true, Description: "Account type"},
			"billing_organization_code": {Type: schema.TypeString, Computed: true, Description: "Billing organization code"},
			"business_category_id":      {Type: schema.TypeString, Computed: true, Description: "Business category id"},
			"business_category_name":    {Type: schema.TypeString, Computed: true, Description: "Business category name"},
			"business_group_id":         {Type: schema.TypeString, Computed: true, Description: "Business group id"},
			"company_id":                {Type: schema.TypeString, Computed: true, Description: "Company id"},
			"company_name":              {Type: schema.TypeString, Computed: true, Description: "Company name"},
			"current_month_bill_amount": {Type: schema.TypeInt, Computed: true, Description: "The amount of current month bill"},
			"default_zone_id":           {Type: schema.TypeString, Computed: true, Description: "Default zone id"},
			"estimated_used_amount":     {Type: schema.TypeFloat, Computed: true, Description: "Estimated used amount"},
			"free_trial_expired_date":   {Type: schema.TypeString, Computed: true, Description: "Free trial expiration date"},
			"free_trial_expired_dday":   {Type: schema.TypeString, Computed: true, Description: "Free trial expires D-day"},
			"free_trial_start_date":     {Type: schema.TypeString, Computed: true, Description: "Free trial start date"},
			"last_month_bill_amount":    {Type: schema.TypeInt, Computed: true, Description: "Previous month's bill"},
			"igw_create_yn":             {Type: schema.TypeString, Computed: true, Description: "Internet connection"},
			"network_type":              {Type: schema.TypeString, Computed: true, Description: "Network type"},
			"price_system_year":         {Type: schema.TypeString, Computed: true, Description: "Price system year"},
			"project_name":              {Type: schema.TypeString, Computed: true, Description: "Project name"},
			"project_state":             {Type: schema.TypeString, Computed: true, Description: "Project status"},
			"project_member_count":      {Type: schema.TypeInt, Computed: true, Description: "Project members count"},
			"project_resource_count":    {Type: schema.TypeInt, Computed: true, Description: "Project resources count"},
			"project_service_count":     {Type: schema.TypeInt, Computed: true, Description: "Project services count"},
			"project_description":       {Type: schema.TypeString, Computed: true, Description: "Project description"},
			"created_by":                {Type: schema.TypeString, Computed: true, Description: "The id of the person who created the resource"},
			"created_by_name":           {Type: schema.TypeString, Computed: true, Description: "The name of the person who created the resource"},
			"created_by_email":          {Type: schema.TypeString, Computed: true, Description: "The e-mail address of the person who created the resource"},
			"created_dt":                {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":               {Type: schema.TypeString, Computed: true, Description: "The id of the person who modified the resource"},
			"modified_by_name":          {Type: schema.TypeString, Computed: true, Description: "The name of the person who modified the resource"},
			"modified_by_email":         {Type: schema.TypeString, Computed: true, Description: "The e-mail of the person who modified the resource"},
			"modified_dt":               {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			"vpc_version":               {Type: schema.TypeString, Computed: true, Description: "VPC version"},
			"business_category_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"business_category_user_email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Business category user email",
						},
						"business_category_user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Business category user name",
						},
					},
				},
				Description: "List of business category users",
			},
			"service_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Block ID",
						},
						"is_multi_availability_zone": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it is an availability zone",
						},
						"service_zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service zone ID",
						},
						"service_zone_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service zone location",
						},
						"service_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service zone name",
						},
						"availability_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone_name": {Type: schema.TypeString, Computed: true, Description: "Availability zone name"},
								},
							},
							Description: "List of availability zones",
						},
					},
				},
				Description: "List of service zones",
			},
		},
	}
}

func DatasourceProjectZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceProjectZoneList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":      common.DatasourceFilter(),
			"project_id":  {Type: schema.TypeString, Required: true, Description: "Project ID"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Zones in project", Elem: datasourceProjectZonesElem()},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of service zones in project",
	}
}

func datasourceProjectZonesElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_id":                   {Type: schema.TypeString, Computed: true, Description: "Block ID"},
			"is_multi_availability_zone": {Type: schema.TypeBool, Computed: true, Description: "Multi availability zone"},
			"service_zone_id":            {Type: schema.TypeString, Computed: true, Description: "Service zone ID"},
			"service_zone_location":      {Type: schema.TypeString, Computed: true, Description: "Service zone location"},
			"service_zone_name":          {Type: schema.TypeString, Computed: true, Description: "Service zone name"},
			"availability_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone_name": {Type: schema.TypeString, Computed: true, Description: "Availability zone name"},
					},
				},
				Description: "List of availability zones",
			},
		},
	}
}

func datasourceProjectZoneList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	id := rd.Get("project_id").(string)

	result, err := inst.Client.Project.GetProjectZoneList(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(result.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceProjectZones().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}

func DatasourceProjectProducts() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceProjectProductsList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":        common.DatasourceFilter(),
			"project_id":    {Type: schema.TypeString, Required: true, Description: "Project ID"},
			"language_code": {Type: schema.TypeString, Optional: true, Description: "Language code"},
			"contents":      {Type: schema.TypeList, Computed: true, Description: "List of products  in project", Elem: datasourceProductCategoryElem()},
			"total_count":   {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of products in given project",
	}
}

func datasourceProductCategoryElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_category_id":          {Type: schema.TypeString, Computed: true, Description: "Product category id"},
			"product_category_name":        {Type: schema.TypeString, Computed: true, Description: "Product category name"},
			"product_category_state":       {Type: schema.TypeString, Computed: true, Description: "Product category state"},
			"product_category_description": {Type: schema.TypeString, Computed: true, Description: "Product category description"},
			"product_set":                  {Type: schema.TypeString, Computed: true, Description: "Product category set"},
			"products": {Type: schema.TypeList, Computed: true, Description: "List of product resources",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_product_creatable":         {Type: schema.TypeString, Computed: true, Description: "Product creation availability"},
						"product_offering_detail_info": {Type: schema.TypeString, Computed: true, Description: "Product offering details"},
						"product_offering_id":          {Type: schema.TypeString, Computed: true, Description: "Product offering ID"},
						"product_offering_name":        {Type: schema.TypeString, Computed: true, Description: "Product offering name"},
						"product_offering_state":       {Type: schema.TypeString, Computed: true, Description: "Product offering state"},
						"product_offering_description": {Type: schema.TypeString, Computed: true, Description: "Product offering description"},
					},
				},
			},
		},
	}
}

func datasourceProjectProductsList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	id := rd.Get("project_id").(string)
	code := common.GetKeyString(rd, "language_code")

	result, err := inst.Client.Project.GetProjectProductsList(ctx, id, code)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(result.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceProjectProducts().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}

func DatasourceProjectUserProductsResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceProjectUserProductsResourcesList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":              common.DatasourceFilter(),
			"product_category_id": {Type: schema.TypeString, Optional: true, Description: "Product category ID"},
			"contents":            {Type: schema.TypeList, Computed: true, Description: "List of products resources in project", Elem: datasourceProductCategoryResourceElem()},
			"total_count":         {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides product resources for projects the user belongs to",
	}
}

func datasourceProductCategoryResourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_category_name": {Type: schema.TypeString, Computed: true, Description: "Product category name"},
			"product_resources": {Type: schema.TypeList, Computed: true, Description: "List of product resources",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_offering_name":           {Type: schema.TypeString, Computed: true, Description: "Product offering names"},
						"product_offering_resource_count": {Type: schema.TypeInt, Computed: true, Description: "Number of resources provided by product"},
					},
				},
			},
		},
	}
}

func datasourceProjectUserProductsResourcesList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	result, err := inst.Client.Project.GetProductResourceList(ctx, common.GetKeyString(rd, "product_category_id"))
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(result.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceProjectUserProductsResources().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}

func DatasourceProjectProductResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceProjectProductResourcesList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":              common.DatasourceFilter(),
			"project_id":          {Type: schema.TypeString, Required: true, Description: "Project ID"},
			"product_category_id": {Type: schema.TypeString, Optional: true, Description: "Product category ID"},
			"contents":            {Type: schema.TypeList, Computed: true, Description: "List of product resources in project", Elem: datasourceProjectProductResourceElem()},
			"total_count":         {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides product resources for given project",
	}
}

func datasourceProjectProductResourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_category_name": {Type: schema.TypeString, Computed: true, Description: "Product category name"},
			"product_resources": {Type: schema.TypeList, Computed: true, Description: "List of product resources",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_offering_name":           {Type: schema.TypeString, Computed: true, Description: "Product offering names"},
						"product_offering_resource_count": {Type: schema.TypeInt, Computed: true, Description: "Number of resources provided by product"},
					},
				},
			},
		},
	}
}

func datasourceProjectProductResourcesList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	id := rd.Get("project_id").(string)
	prodCatId := common.GetKeyString(rd, "product_category_id")

	result, err := inst.Client.Project.GetProjectProductResourcesList(ctx, id, prodCatId)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(result.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceProjectProductResources().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}
