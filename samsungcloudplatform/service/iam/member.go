package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	samsungcloudplatform.RegisterResource("samsungcloudplatform_iam_member", ResourceMember())
}
func ResourceMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberCreate,
		ReadContext:   resourceMemberRead,
		UpdateContext: resourceMemberUpdate,
		DeleteContext: resourceMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Group ID list",
			},
			"user_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User email",
			},
			"tags": tfTags.TagsSchema(),

			"project_id":       {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"company_name":     {Type: schema.TypeString, Computed: true, Description: "Company name"},
			"email":            {Type: schema.TypeString, Computed: true, Description: "Email"},
			"last_access_date": {Type: schema.TypeString, Computed: true, Description: "Last access data"},
			"organization_id":  {Type: schema.TypeString, Computed: true, Description: "Organization ID"},
			"position_name":    {Type: schema.TypeString, Computed: true, Description: "Position within the company"},
			"registered_by":    {Type: schema.TypeString, Computed: true, Description: "Register's email"},
			"registered_dt":    {Type: schema.TypeString, Computed: true, Description: "Registered date"},
			"user_group_count": {Type: schema.TypeInt, Computed: true, Description: "Number of user's groups"},
			"user_id":          {Type: schema.TypeString, Computed: true, Description: "User ID"},
			"user_name":        {Type: schema.TypeString, Computed: true, Description: "User name"},
			"user_srn":         {Type: schema.TypeString, Computed: true, Description: "User SRN"},
		},
		Description: "Provides IAM member resource",
	}
}

func resourceMemberCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	groupIds := common.ToStringList(rd.Get("group_ids").(*schema.Set).List())
	email := rd.Get("user_email").(string)

	_, _, err := inst.Client.Iam.AddMember(ctx, groupIds, email, rd.Get("tags").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	memberResponse, err := inst.Client.Iam.ListMembers(ctx, "", email, "")
	if err != nil {
		return diag.FromErr(err)
	}
	if memberResponse.TotalCount != 1 {
		return diag.Errorf("fail to find user email : %s", email)
	}

	rd.SetId(memberResponse.Contents[0].UserId)

	return resourceMemberRead(ctx, rd, meta)
}

func resourceMemberRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	member, err := inst.Client.Iam.DetailMember(ctx, rd.Id())

	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	memberGroups, err := inst.Client.Iam.ListMemberGroups(ctx, rd.Id(), "")
	groupIds := make([]string, memberGroups.TotalCount)
	for i, group := range memberGroups.Contents {
		groupIds[i] = group.GroupId
	}

	rd.Set("group_ids", groupIds)
	rd.Set("project_id", member.ProjectId)
	rd.Set("company_name", member.CompanyName)
	rd.Set("email", member.Email)
	rd.Set("last_access_date", member.LastAccessDt)
	rd.Set("organization_id", member.OrganizationId)
	rd.Set("position_name", member.PositionName)
	rd.Set("registered_by", member.RegisteredBy)
	rd.Set("registered_dt", member.RegisteredDt)
	rd.Set("user_group_count", member.UserGroupCount)
	rd.Set("user_id", member.UserId)
	rd.Set("user_name", member.UserName)
	rd.Set("user_srn", member.UserSrn)

	tfTags.SetTags(ctx, rd, meta, rd.Id())
	return nil
}

func resourceMemberUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("group_ids") {
		// 추가는 group id, 삭제는 usergroup id 로 해야함
		// user group id는 목록 조회를 통해 얻어와야 함
		// 그럴수 있지
		addIds, removeIds := common.GetAddRemoveItemsStringListFromSet(rd, "group_ids")
		if len(addIds) > 0 {
			_, err := inst.Client.Iam.AddMemberGroups(ctx, rd.Id(), addIds)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if len(removeIds) > 0 {
			userGroupIds := make([]string, len(removeIds))
			for i, id := range removeIds {
				groupDetail, err := inst.Client.Iam.DetailGroup(ctx, id)
				if err != nil {
					return diag.FromErr(err)
				}
				listGroup, err := inst.Client.Iam.ListMemberGroups(ctx, rd.Id(), groupDetail.GroupName)
				if err != nil {
					return diag.FromErr(err)
				}
				if listGroup.TotalCount == 0 || listGroup.TotalCount > 1 {
					return diag.Errorf("fail to find proper group information")
				}
				userGroupIds[i] = listGroup.Contents[0].UserGroupId
			}
			_, err := inst.Client.Iam.RemoveMemberGroups(ctx, rd.Id(), userGroupIds)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMemberRead(ctx, rd, meta)
}

func resourceMemberDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.Iam.RemoveMember(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
