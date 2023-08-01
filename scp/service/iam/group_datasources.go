package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_iam_group", DatasourceGroup())
	scp.RegisterDataSource("scp_iam_groups", DatasourceGroups())
	scp.RegisterDataSource("scp_iam_group_members", DatasourceGroupMembers())
	scp.RegisterDataSource("scp_iam_group_policies", DatasourceGroupPolicies())
}

func DatasourceGroup() *schema.Resource {
	var groupResource schema.Resource
	groupResource.ReadContext = datasourceGroupRead
	groupResource.Schema = datasourceGroupsElem().Schema

	delete(groupResource.Schema, "group_id")
	delete(groupResource.Schema, "group_name")

	groupResource.Schema["group_id"] = &schema.Schema{Type: schema.TypeString, Optional: true, Description: "Group ID"}
	groupResource.Schema["group_name"] = &schema.Schema{Type: schema.TypeString, Optional: true, Description: "Group name"}
	groupResource.Schema["group_policy_count"] = &schema.Schema{Type: schema.TypeInt, Computed: true, Description: "Group policy count"}
	groupResource.Schema["group_srn"] = &schema.Schema{Type: schema.TypeString, Computed: true, Description: "Group SRN"}

	groupResource.Description = "Provides detailed information for a given group name or ID"

	return &groupResource
}

func datasourceGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	var groupId string

	if groupNameRaw, ok := rd.GetOk("group_name"); ok {
		groupName := groupNameRaw.(string)

		listGroup, err := inst.Client.Iam.ListGroups(ctx, groupName, "")
		if err != nil {
			return diag.FromErr(err)
		}

		if len(listGroup.Contents) == 0 {
			return diag.Errorf("Group name \"%s\" not found...", groupName)
		}

		groupId = listGroup.Contents[0].GroupId
	} else if groupIdRaw, ok := rd.GetOk("group_id"); ok {
		groupId = groupIdRaw.(string)
	} else {
		listGroup, err := inst.Client.Iam.ListGroups(ctx, "", "")
		if err != nil {
			return diag.FromErr(err)
		}

		if len(listGroup.Contents) == 0 {
			return diag.Errorf("No group exists")
		}

		groupId = listGroup.Contents[0].GroupId
	}

	info, err := inst.Client.Iam.DetailGroup(ctx, groupId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(groupId)
	rd.Set("project_id", info.ProjectId)
	rd.Set("group_id", info.GroupId)
	rd.Set("group_name", info.GroupName)
	rd.Set("group_type", info.GroupType)
	rd.Set("group_user_count", info.GroupUserCount)
	rd.Set("description", info.Description)
	rd.Set("created_by", info.CreatedBy)
	rd.Set("created_by_name", info.CreatedByName)
	rd.Set("created_by_email", info.CreatedByEmail)
	rd.Set("created_dt", info.CreatedDt.String())
	rd.Set("modified_by", info.ModifiedBy)
	rd.Set("modified_by_name", info.ModifiedByName)
	rd.Set("modified_by_email", info.ModifiedByEmail)
	rd.Set("modified_dt", info.ModifiedDt.String())
	rd.Set("group_policy_count", info.GroupPolicyCount)
	rd.Set("group_srn", info.GroupSrn)

	return nil
}

func DatasourceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceGroupsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":            common.DatasourceFilter(),
			"group_name":        {Type: schema.TypeString, Optional: true, Description: "Group name"},
			"modified_by_email": {Type: schema.TypeString, Optional: true, Description: "Modifier's email"},
			"contents":          {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: datasourceGroupsElem()},
			"total_count":       {Type: schema.TypeInt, Computed: true, Description: "Total count"},
		},
	}
}

func datasourceGroupsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	groupName := rd.Get("group_name").(string)
	email := rd.Get("modified_by_email").(string)

	responses, err := inst.Client.Iam.ListGroups(ctx, groupName, email)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceGroups().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceGroupsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":        {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"group_id":          {Type: schema.TypeString, Computed: true, Description: "Group ID"},
			"group_name":        {Type: schema.TypeString, Computed: true, Description: "Group name"},
			"group_type":        {Type: schema.TypeString, Computed: true, Description: "Group type"},
			"group_user_count":  {Type: schema.TypeInt, Computed: true, Description: "Number of group users"},
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

func DatasourceGroupMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceGroupMembersRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":       common.DatasourceFilter(),
			"group_id":     {Type: schema.TypeString, Required: true, Description: "Group ID"},
			"company_name": {Type: schema.TypeString, Optional: true, Description: "Company name"},
			"email":        {Type: schema.TypeString, Optional: true, Description: "Email"},
			"user_name":    {Type: schema.TypeString, Optional: true, Description: "User name"},
			"contents":     {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: datasourceGroupMembersElem()},
			"total_count":  {Type: schema.TypeInt, Computed: true, Description: "Total count"},
		},
	}
}

func datasourceGroupMembersRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	groupId := rd.Get("group_id").(string)
	company := rd.Get("company_name").(string)
	email := rd.Get("email").(string)
	userName := rd.Get("user_name").(string)

	result, err := inst.Client.Iam.ListGroupMembers(ctx, groupId, iam.ListMemberRequest{
		CompanyName: company,
		Email:       email,
		UserName:    userName,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(result.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceGroupMembers().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", result.TotalCount)

	return nil
}

func datasourceGroupMembersElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"company_name":      {Type: schema.TypeString, Computed: true, Description: "Company name"},
			"email":             {Type: schema.TypeString, Computed: true, Description: "Email"},
			"last_access_dt":    {Type: schema.TypeString, Computed: true, Description: "Last access date"},
			"user_group_id":     {Type: schema.TypeString, Computed: true, Description: "User group ID"},
			"user_id":           {Type: schema.TypeString, Computed: true, Description: "User ID"},
			"user_name":         {Type: schema.TypeString, Computed: true, Description: "User name"},
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

func DatasourceGroupPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceGroupPoliciesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":      common.DatasourceFilter(),
			"group_id":    {Type: schema.TypeString, Required: true, Description: "Group ID"},
			"policy_name": {Type: schema.TypeString, Optional: true, Description: "Policy name"},
			"policy_type": {Type: schema.TypeString, Optional: true, Description: "Policy type"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: datasourceGroupPoliciesElem()},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total count"},
		},
	}
}

func datasourceGroupPoliciesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	groupId := rd.Get("group_id").(string)
	policyName := rd.Get("policy_name").(string)
	policyType := rd.Get("policy_type").(string)

	result, err := inst.Client.Iam.ListGroupPolicies(ctx, groupId, iam.ListPolicyRequest{
		PolicyName: policyName,
		PolicyType: policyType,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(result.Contents)

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceGroupPolicies().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", result.TotalCount)

	return nil
}

func datasourceGroupPoliciesElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_id":           {Type: schema.TypeString, Computed: true, Description: "Policy ID"},
			"policy_name":         {Type: schema.TypeString, Computed: true, Description: "Policy name"},
			"policy_type":         {Type: schema.TypeString, Computed: true, Description: "Policy type"},
			"principal_policy_id": {Type: schema.TypeString, Computed: true, Description: "Principal policy ID"},
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
