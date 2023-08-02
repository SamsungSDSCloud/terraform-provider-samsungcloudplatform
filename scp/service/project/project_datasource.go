package project

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client/project"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_project", DatasourceProjects())
}

func DatasourceProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
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
			"contents":                {Type: schema.TypeList, Computed: true, Description: "Project info list", Elem: datasourceElem()},
			"total_count":             {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of my projects.",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func datasourceElem() *schema.Resource {
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
