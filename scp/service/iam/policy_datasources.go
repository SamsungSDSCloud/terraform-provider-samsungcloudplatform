package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/iam"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_iam_policy", DatasourcePolicy())
	scp.RegisterDataSource("scp_iam_policies", DatasourcePolicies())
	scp.RegisterDataSource("scp_iam_policy_groups", DatasourcePolicyGroups())
	scp.RegisterDataSource("scp_iam_policy_roles", DatasourcePolicyRoles())
}

func DatasourcePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"policy_id":              {Type: schema.TypeString, Required: true, Description: "Policy ID"},
			"policy_name":            {Type: schema.TypeString, Computed: true, Description: "Policy name"},
			"policy_json":            {Type: schema.TypeString, Computed: true, Description: "Policy JSON"},
			"description":            {Type: schema.TypeString, Computed: true, Description: "Description"},
			"project_id":             {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"policy_principal_count": {Type: schema.TypeInt, Computed: true, Description: "Policy principal count"},
			"policy_srn":             {Type: schema.TypeString, Computed: true, Description: "Policy SRN"},
			"policy_type":            {Type: schema.TypeString, Computed: true, Description: "Policy type"},
			"policy_version":         {Type: schema.TypeString, Computed: true, Description: "Policy version"},
			"created_by":             {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_by_name":        {Type: schema.TypeString, Computed: true, Description: "Creator's name"},
			"created_by_email":       {Type: schema.TypeString, Computed: true, Description: "Creator's email"},
			"created_dt":             {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":            {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_by_name":       {Type: schema.TypeString, Computed: true, Description: "Modifier's name"},
			"modified_by_email":      {Type: schema.TypeString, Computed: true, Description: "Modifier's email"},
			"modified_dt":            {Type: schema.TypeString, Computed: true, Description: "Modified date"},
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
		Description: "Provides detailed information of a policy",
	}
}

func dataSourcePolicyRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, err := inst.Client.Iam.DetailPolicy(ctx, rd.Get("policy_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(info.PolicyId)

	rd.Set("policy_name", info.PolicyName)
	rd.Set("policy_json", info.PolicyJson)
	rd.Set("description", info.Description)
	rd.Set("project_id", info.ProjectId)
	rd.Set("policy_id", info.PolicyId)
	rd.Set("policy_principal_count", info.PolicyPrincipalCount)
	rd.Set("policy_srn", info.PolicySrn)
	rd.Set("policy_type", info.PolicyType)
	rd.Set("policy_version", info.PolicyVersion)
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

	return nil
}

func DatasourcePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePoliciesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":            common.DatasourceFilter(),
			"modified_by_email": {Type: schema.TypeString, Optional: true, Description: "Modifier's email"},
			"policy_name":       {Type: schema.TypeString, Optional: true, Description: "Policy name"},
			"contents":          {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: dataSourcePoliciesElem()},
			"total_count":       {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of policies",
	}
}

func dataSourcePoliciesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	requestParam := iam.ListMemberRequest{
		CompanyName: rd.Get(common.ToSnakeCase("modifiedByEmail")).(string),
		Email:       rd.Get(common.ToSnakeCase("policyName")).(string),
	}

	responses, err := inst.Client.Iam.ListPolicies(ctx, requestParam)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceGroupMembers().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func dataSourcePoliciesElem() *schema.Resource {
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

func DatasourcePolicyGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyGroupsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":      common.DatasourceFilter(),
			"policy_id":   {Type: schema.TypeString, Required: true, Description: "Policy ID"},
			"group_name":  {Type: schema.TypeString, Optional: true, Description: "Group name"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: dataSourcePolicyGroupElem()},
			"total_count": {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provide a list of groups that contain this policy",
	}
}

func dataSourcePolicyGroupsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	policyId := rd.Get("policy_id").(string)
	groupName := rd.Get("group_name").(string)

	responses, _, err := inst.Client.Iam.ListPolicyGroups(ctx, policyId, groupName)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceGroupMembers().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func dataSourcePolicyGroupElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_id":            {Type: schema.TypeString, Computed: true, Description: "Policy ID"},
			"group_name":          {Type: schema.TypeString, Computed: true, Description: "Policy name"},
			"group_type":          {Type: schema.TypeString, Computed: true, Description: "Policy type"},
			"principal_policy_id": {Type: schema.TypeString, Computed: true, Description: "Principal policy ID"},
			"created_by":          {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_by_name":     {Type: schema.TypeString, Computed: true, Description: "Creator's name"},
			"created_by_email":    {Type: schema.TypeString, Computed: true, Description: "Creator's email"},
			"created_dt":          {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":         {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_by_name":    {Type: schema.TypeString, Computed: true, Description: "Modifier's name"},
			"modified_by_email":   {Type: schema.TypeString, Computed: true, Description: "Modifier's email"},
			"modified_dt":         {Type: schema.TypeString, Computed: true, Description: "Modified date"},
		},
	}
}

func DatasourcePolicyRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyRolesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":      common.DatasourceFilter(),
			"policy_id":   {Type: schema.TypeString, Required: true, Description: "Policy ID"},
			"role_name":   {Type: schema.TypeString, Optional: true, Description: "Role name"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: dataSourcePolicyRoleElem()},
			"total_count": {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provide a list of roles that contain this policy",
	}
}

func dataSourcePolicyRolesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	policyId := rd.Get("policy_id").(string)
	roleName := rd.Get("role_name").(string)

	responses, _, err := inst.Client.Iam.ListPolicyRoles(ctx, policyId, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceGroupMembers().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func dataSourcePolicyRoleElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"principal_policy_id": {Type: schema.TypeString, Computed: true, Description: "Principal policy ID"},
			"role_id":             {Type: schema.TypeString, Computed: true, Description: "Role ID"},
			"role_name":           {Type: schema.TypeString, Computed: true, Description: "Role name"},
			"description":         {Type: schema.TypeString, Computed: true, Description: "Description"},
			"created_by":          {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_by_name":     {Type: schema.TypeString, Computed: true, Description: "Creator's name"},
			"created_by_email":    {Type: schema.TypeString, Computed: true, Description: "Creator's email"},
			"created_dt":          {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":         {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_by_name":    {Type: schema.TypeString, Computed: true, Description: "Modifier's name"},
			"modified_by_email":   {Type: schema.TypeString, Computed: true, Description: "Modifier's email"},
			"modified_dt":         {Type: schema.TypeString, Computed: true, Description: "Modified date"},
		},
	}
}
