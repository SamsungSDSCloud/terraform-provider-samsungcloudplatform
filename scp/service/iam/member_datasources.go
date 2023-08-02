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
	scp.RegisterDataSource("scp_iam_member", DatasourceMember())
	scp.RegisterDataSource("scp_iam_members", DatasourceMembers())
	scp.RegisterDataSource("scp_iam_member_groups", DatasourceMemberGroups())
	scp.RegisterDataSource("scp_iam_member_systemgroups", DatasourceMemberSystemGroups())
}

func DatasourceMember() *schema.Resource {
	var memberResource schema.Resource
	memberResource.ReadContext = dataSourceMemberRead
	memberResource.Schema = datasourceMemberElem().Schema

	delete(memberResource.Schema, "user_id")

	memberResource.Schema["user_id"] = &schema.Schema{Type: schema.TypeString, Required: true, Description: "User ID"}
	memberResource.Schema["user_srn"] = &schema.Schema{Type: schema.TypeString, Computed: true, Description: "User SRN"}
	memberResource.Schema["registered_by"] = &schema.Schema{Type: schema.TypeString, Computed: true, Description: "Registrant's Email"}
	memberResource.Schema["registered_dt"] = &schema.Schema{Type: schema.TypeString, Computed: true, Description: "Registration date"}
	memberResource.Schema["tags"] = &schema.Schema{
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
		Description: "Tag list",
	}

	return &memberResource
}

func dataSourceMemberRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, err := inst.Client.Iam.DetailMember(ctx, rd.Get("user_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(info.UserId)
	rd.Set("company_name", info.CompanyName)
	rd.Set("email", info.Email)
	rd.Set("last_access_dt", info.LastAccessDt.String())
	rd.Set("organization_id", info.OrganizationId)
	rd.Set("position_name", info.PositionName)
	rd.Set("user_group_count", info.UserGroupCount)
	rd.Set("user_id", info.UserId)
	rd.Set("user_name", info.UserName)
	rd.Set("created_by", info.CreatedBy)
	rd.Set("created_by_name", info.CreatedByName)
	rd.Set("created_by_email", info.CreatedByEmail)
	rd.Set("created_dt", info.CreatedDt.String())
	rd.Set("modified_by", info.ModifiedBy)
	rd.Set("modified_by_name", info.ModifiedByName)
	rd.Set("modified_by_email", info.ModifiedByEmail)
	rd.Set("modified_dt", info.ModifiedDt.String())
	rd.Set("user_srn", info.UserSrn)
	rd.Set("registered_by", info.RegisteredBy)
	rd.Set("registered_dt", info.RegisteredDt.String())

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

func DatasourceMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMembersRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":       common.DatasourceFilter(),
			"company_name": {Type: schema.TypeString, Optional: true, Description: "Company name"},
			"email":        {Type: schema.TypeString, Optional: true, Description: "Member's email"},
			"user_name":    {Type: schema.TypeString, Optional: true, Description: "Member's name"},
			"contents":     {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: datasourceMemberElem()},
			"total_count":  {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceMembersRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	companyName := rd.Get("company_name").(string)
	email := rd.Get("email").(string)
	userName := rd.Get("user_name").(string)

	responses, err := inst.Client.Iam.ListMembers(ctx, companyName, email, userName)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)
	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceMembers().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceMemberElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"company_name":      {Type: schema.TypeString, Computed: true, Description: "Company name"},
			"email":             {Type: schema.TypeString, Computed: true, Description: "User email"},
			"last_access_dt":    {Type: schema.TypeString, Computed: true, Description: "Last access date"},
			"organization_id":   {Type: schema.TypeString, Computed: true, Description: "Organization ID"},
			"position_name":     {Type: schema.TypeString, Computed: true, Description: "Position name"},
			"user_group_count":  {Type: schema.TypeInt, Computed: true, Description: "User group count"},
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
func DatasourceMemberGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMemberGroupsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":      common.DatasourceFilter(),
			"member_id":   {Type: schema.TypeString, Required: true, Description: "Company name"},
			"group_name":  {Type: schema.TypeString, Optional: true, Description: "Member's email"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: datasourceMemberGroupsElem()},
			"total_count": {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceMemberGroupsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	memberId := rd.Get("member_id").(string)
	groupName := rd.Get("group_name").(string)

	responses, err := inst.Client.Iam.ListMemberGroups(ctx, memberId, groupName)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)
	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceMemberGroups().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceMemberGroupsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_id":          {Type: schema.TypeString, Computed: true, Description: "Group ID"},
			"group_name":        {Type: schema.TypeString, Computed: true, Description: "Group name"},
			"group_type":        {Type: schema.TypeString, Computed: true, Description: "Group type"},
			"user_group_id":     {Type: schema.TypeString, Computed: true, Description: "User group ID"},
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

func DatasourceMemberSystemGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMemberSystemGroupsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":      common.DatasourceFilter(),
			"member_id":   {Type: schema.TypeString, Required: true, Description: "Company name"},
			"contents":    {Type: schema.TypeList, Computed: true, Description: "Contents list", Elem: datasourceMemberSystemGroupsElem()},
			"total_count": {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceMemberSystemGroupsRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	memberId := rd.Get("member_id").(string)

	responses, err := inst.Client.Iam.ListMemberSystemGroups(ctx, memberId)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)
	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceMemberSystemGroups().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceMemberSystemGroupsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_id":      {Type: schema.TypeString, Computed: true, Description: "Group ID"},
			"group_name":    {Type: schema.TypeString, Computed: true, Description: "Group name"},
			"group_type":    {Type: schema.TypeString, Computed: true, Description: "Group type"},
			"user_group_id": {Type: schema.TypeString, Computed: true, Description: "User group ID"},
		},
	}
}
