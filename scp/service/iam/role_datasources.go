package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_iam_role", DatasourceRole())
	scp.RegisterDataSource("scp_iam_roles", DatasourceRoles())
	scp.RegisterDataSource("scp_iam_role_policies", DatasourceRolePolicies())
}

func DatasourceRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoleRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role_id":           {Type: schema.TypeString, Required: true, Description: "Role ID"},
			"role_name":         {Type: schema.TypeString, Computed: true, Description: "Role name"},
			"project_id":        {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"description":       {Type: schema.TypeString, Computed: true, Description: "Description"},
			"role_policy_count": {Type: schema.TypeInt, Computed: true, Description: "Description"},
			"role_srn":          {Type: schema.TypeString, Computed: true, Description: "Role SRN"},
			"session_time":      {Type: schema.TypeInt, Computed: true, Description: "Session time"},
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
			"trust_principals": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Performing subjects",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Project IDs",
						},
						"user_srns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "User SRNs",
						},
					},
				},
			},
		},
	}
}

func dataSourceRoleRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	info, err := inst.Client.Iam.DetailRole(ctx, rd.Get("role_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(info.RoleId)

	rd.Set("role_id", info.RoleId)
	rd.Set("role_name", info.RoleName)
	rd.Set("project_id", info.ProjectId)
	rd.Set("description", info.Description)
	rd.Set("role_policy_count", info.RolePolicyCount)
	rd.Set("role_srn", info.RoleSrn)
	rd.Set("session_time", info.SessionTime)
	rd.Set("created_by", info.CreatedBy)
	rd.Set("created_by_name", info.CreatedByName)
	rd.Set("created_by_email", info.CreatedByEmail)
	rd.Set("created_dt", info.CreatedDt.String())
	rd.Set("modified_by", info.ModifiedBy)
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

	var principals common.HclSetObject
	if info.TrustPrincipals != nil {
		kv := common.HclKeyValueObject{
			"project_ids": info.TrustPrincipals.ProjectIds,
			"user_srns":   info.TrustPrincipals.UserSrns,
		}
		principals = append(principals, kv)
	}
	rd.Set("trust_principals", principals)
	return nil
}

func DatasourceRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRolesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":         common.DatasourceFilter(),
			"modifier_email": {Type: schema.TypeString, Optional: true, Description: "Modifier's email"},
			"role_name":      {Type: schema.TypeString, Optional: true, Description: "Role name"},
			"contents":       {Type: schema.TypeList, Computed: true, Description: "Contents", Elem: datasourceRolesElem()},
			"total_count":    {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceRolesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	email := rd.Get("modifier_email").(string)
	roleName := rd.Get("role_name").(string)

	responses, err := inst.Client.Iam.ListRoles(ctx, email, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceRoles().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceRolesElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":        {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"role_id":           {Type: schema.TypeString, Computed: true, Description: "Role ID"},
			"role_name":         {Type: schema.TypeString, Computed: true, Description: "Role name"},
			"description":       {Type: schema.TypeString, Computed: true, Description: "Description"},
			"created_by":        {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_by_name":   {Type: schema.TypeString, Computed: true, Description: "Creator's name"},
			"created_by_email":  {Type: schema.TypeString, Computed: true, Description: "Creator's email"},
			"created_dt":        {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":       {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_by_name":  {Type: schema.TypeString, Computed: true, Description: "Modifier's name"},
			"modified_by_email": {Type: schema.TypeString, Computed: true, Description: "Modifier's email"},
			"modified_dt":       {Type: schema.TypeString, Computed: true, Description: "Modified date"},
		},
	}
}

func DatasourceRolePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRolePoliciesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":      common.DatasourceFilter(),
			"role_id":     {Type: schema.TypeString, Required: true, Description: "Role ID"},
			"policy_name": {Type: schema.TypeString, Optional: true, Description: "Modifier's email"},
			"policy_type": {Type: schema.TypeString, Optional: true, Description: "Role name"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Contents", Elem: datasourceRolePolicyElem()},
			"total_count": {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceRolePoliciesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	roleId := rd.Get("role_id").(string)
	policyName := rd.Get("policy_name").(string)
	policyType := rd.Get("policy_type").(string)

	responses, err := inst.Client.Iam.ListRolePolicies(ctx, roleId, policyName, policyType)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceRolePolicies().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceRolePolicyElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_id":         {Type: schema.TypeString, Computed: true, Description: "Policy ID"},
			"policy_name":       {Type: schema.TypeString, Computed: true, Description: "Policy name"},
			"policy_type":       {Type: schema.TypeString, Computed: true, Description: "Policy type"},
			"description":       {Type: schema.TypeString, Computed: true, Description: "Description"},
			"created_by":        {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_by_name":   {Type: schema.TypeString, Computed: true, Description: "Creator's name"},
			"created_by_email":  {Type: schema.TypeString, Computed: true, Description: "Creator's email"},
			"created_dt":        {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":       {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_by_name":  {Type: schema.TypeString, Computed: true, Description: "Modifier's name"},
			"modified_by_email": {Type: schema.TypeString, Computed: true, Description: "Modifier's email"},
			"modified_dt":       {Type: schema.TypeString, Computed: true, Description: "Modified date"},
		},
	}
}
