package iam

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"unicode/utf8"
)

func init() {
	scp.RegisterResource("scp_iam_group", ResourceGroup())
}

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of the authority-group (3 to 24 characters using Korean, English, numbers, +=,.@-_)",
				ValidateDiagFunc: validateGroupName,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Description (1000 characters or less)",
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"user_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of user IDs",
			},
			"policy_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of policy IDs",
			},
			"group_policy_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Group policies' count",
			},
			"group_srn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Group SRN",
			},
			"group_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Group type",
			},
			"group_user_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Group users' count",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator's ID",
			},
			"created_by_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator's name",
			},
			"created_by_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator's email",
			},
			"created_dt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created date",
			},
			"modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modifier's ID",
			},
			"modified_by_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modifier's name",
			},
			"modified_by_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modifier's email",
			},
			"modified_dt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modified date",
			},
		},
		Description: "Provides IAM group resource.",
	}
}

func resourceGroupCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	groupName := rd.Get("group_name").(string)
	desc := rd.Get("description").(string)

	result, _, err := inst.Client.Iam.CreateGroup(ctx, groupName, desc)

	if err != nil {
		return diag.FromErr(err)
	}

	userIdsRaw := rd.Get("user_ids").(*schema.Set).List()
	userIds := common.ToStringList(userIdsRaw)
	if len(userIds) > 0 {
		_, err = inst.Client.Iam.AddGroupMembers(ctx, result.GroupId, userIds)
		if err != nil {
			rd.SetId("")
			return diag.FromErr(err)
		}
	}

	policyIdsRaw := rd.Get("policy_ids").(*schema.Set).List()
	policyIds := common.ToStringList(policyIdsRaw)
	if len(policyIds) > 0 {
		_, err = inst.Client.Iam.AddGroupPolicies(ctx, result.GroupId, policyIds)
		if err != nil {
			rd.SetId("")
			return diag.FromErr(err)
		}
	}

	rd.SetId(result.GroupId)
	return resourceGroupRead(ctx, rd, meta)
}

func resourceGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	result, err := inst.Client.Iam.DetailGroup(ctx, rd.Id())
	if err != nil {
		return
	}

	groupMembers, err := inst.Client.Iam.ListGroupMembers(ctx, rd.Id(), iam.ListMemberRequest{})
	if err != nil {
		return
	}

	userIds := make([]string, groupMembers.TotalCount)
	for i, user := range groupMembers.Contents {
		userIds[i] = user.UserId
	}

	groupPolicies, err := inst.Client.Iam.ListGroupPolicies(ctx, rd.Id(), iam.ListPolicyRequest{})
	if err != nil {
		return
	}
	policyIds := make([]string, groupPolicies.TotalCount)
	for i, policy := range groupPolicies.Contents {
		policyIds[i] = policy.PolicyId
	}

	rd.Set("user_ids", userIds)
	rd.Set("policy_ids", policyIds)
	rd.Set("group_policy_count", result.GroupPolicyCount)
	rd.Set("group_srn", result.GroupSrn)
	rd.Set("group_type", result.GroupType)
	rd.Set("group_user_count", result.GroupUserCount)
	rd.Set("created_by", result.CreatedBy)
	rd.Set("created_by_name", result.CreatedByName)
	rd.Set("created_by_email", result.CreatedByEmail)
	rd.Set("created_dt", result.CreatedDt.String())
	rd.Set("modified_by", result.ModifiedBy)
	rd.Set("modified_by_name", result.ModifiedByName)
	rd.Set("modified_by_email", result.ModifiedByEmail)
	rd.Set("modified_dt", result.ModifiedDt.String())

	return nil
}

func resourceGroupUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	groupId := rd.Id()

	if rd.HasChanges("group_name", "description") {
		groupName := rd.Get("group_name").(string)
		desc := rd.Get("description").(string)

		_, err = inst.Client.Iam.UpdateGroup(ctx, groupId, groupName, desc)
		if err != nil {
			return
		}
	}

	if rd.HasChanges("user_ids") {
		addList, removeList := common.GetAddRemoveItemsStringListFromSet(rd, "user_ids")

		if len(addList) > 0 {
			_, err = inst.Client.Iam.AddGroupMembers(ctx, groupId, addList)
			if err != nil {
				return
			}
		}

		if len(removeList) > 0 {
			_, err = inst.Client.Iam.RemoveGroupMembers(ctx, groupId, removeList)
			if err != nil {
				return
			}
		}
	}

	if rd.HasChanges("policy_ids") {
		addList, removeList := common.GetAddRemoveItemsStringListFromSet(rd, "policy_ids")

		if len(addList) > 0 {
			_, err = inst.Client.Iam.AddGroupPolicies(ctx, groupId, addList)
			if err != nil {
				return
			}
		}

		if len(removeList) > 0 {
			_, err = inst.Client.Iam.RemoveGroupPolicies(ctx, groupId, removeList)
			if err != nil {
				return
			}
		}
	}

	return resourceGroupRead(ctx, rd, meta)
}

func resourceGroupDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.Iam.DeleteGroups(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// no waiting needed
	return nil
}

func validateGroupName(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	var err error = nil
	cnt := utf8.RuneCountInString(value) // cause we have hanguel here :)
	if cnt < 3 {
		err = fmt.Errorf("input must be longer than 3 characters")
	} else if cnt > 24 {
		err = fmt.Errorf("input must be shorter than 24 characters")
	}

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z0-9+=,.@\\-_ㄱ-ㅎ|ㅏ-ㅣ|가-힣]*$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

// iam api들은 wait가 필요 없음
/*
func waitForGroupStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {}
*/
