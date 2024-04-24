package iam

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	scp.RegisterResource("scp_iam_role", ResourceRole())
}

func ResourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDestroy,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Role name",
				ValidateDiagFunc: common.ValidateNameHangeulAlphabetSomeSpecials3to64,
			},
			"trust_principals": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Performing subjects",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Project IDs",
						},
						"user_srns": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "User SRNs",
						},
					},
				},
			},
			"tags": tfTags.TagsSchema(),
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 1000)),
				Description:      "Description",
			},
			"policy_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of policy IDs",
			},

			"project_id":        {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"role_policy_count": {Type: schema.TypeInt, Computed: true, Description: "Role's policy count"},
			"role_srn":          {Type: schema.TypeString, Computed: true, Description: "Role's SRN"},
			"session_time":      {Type: schema.TypeInt, Computed: true, Description: "Session time"},
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

func convertTrustPrincipal(rd *schema.ResourceData) ([]string, []string, error) {
	tpSet := rd.Get("trust_principals").(*schema.Set).List()

	var resultProjectIds []string
	var resultUserSrns []string

	for _, tp := range tpSet {
		currentTp := tp.(map[string]interface{})

		if v, ok := currentTp["project_ids"]; ok {
			ipidList := v.([]interface{})
			for _, pid := range ipidList {
				resultProjectIds = append(resultProjectIds, pid.(string))
			}
		}

		if v, ok := currentTp["user_srns"]; ok {
			isrnList := v.([]interface{})
			for _, srn := range isrnList {
				resultUserSrns = append(resultUserSrns, srn.(string))
			}
		}
	}
	return resultProjectIds, resultUserSrns, nil
}

func resourceRoleCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	projectIds, userSrns, _ := convertTrustPrincipal(rd)

	roleName := rd.Get("role_name").(string)
	tags := rd.Get("tags").(map[string]interface{})
	desc := rd.Get("description").(string)

	response, _, err := inst.Client.Iam.CreateRole(ctx, roleName, projectIds, userSrns, tags, desc)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.RoleId)

	policyIdsRaw := rd.Get("policy_ids").(*schema.Set).List()
	if len(policyIdsRaw) > 0 {
		policyIds := common.ToStringList(policyIdsRaw)
		_, err = inst.Client.Iam.AddRolePolicies(ctx, response.RoleId, policyIds)
		if err != nil {
			rd.SetId("")
			return diag.FromErr(err)
		}
	}
	return resourceRoleRead(ctx, rd, meta)
}

func resourceRoleRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	result, err := inst.Client.Iam.DetailRole(ctx, rd.Id())

	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rolePolicies, err := inst.Client.Iam.ListRolePolicies(ctx, rd.Id(), "", "")
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	policyIds := make([]string, rolePolicies.TotalCount)
	for i, policy := range rolePolicies.Contents {
		policyIds[i] = policy.PolicyId
	}

	rd.Set("policy_ids", policyIds)
	rd.Set("project_id", result.ProjectId)
	rd.Set("role_policy_count", result.RolePolicyCount)
	rd.Set("role_srn", result.RoleSrn)
	rd.Set("session_time", result.SessionTime)
	rd.Set("created_by", result.CreatedBy)
	rd.Set("created_by_name", result.CreatedByName)
	rd.Set("created_by_email", result.CreatedByEmail)
	rd.Set("created_dt", result.CreatedDt.String())
	rd.Set("modified_by", result.ModifiedBy)
	rd.Set("modified_by_name", result.ModifiedByName)
	rd.Set("modified_by_email", result.ModifiedByEmail)
	rd.Set("modified_dt", result.ModifiedDt.String())

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	var principals common.HclSetObject
	if result.TrustPrincipals != nil {
		kv := common.HclKeyValueObject{
			"project_ids": result.TrustPrincipals.ProjectIds,
			"user_srns":   result.TrustPrincipals.UserSrns,
		}
		principals = append(principals, kv)
	}
	rd.Set("trust_principals", principals)
	return nil
}

func resourceRoleUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("role_name", "trust_principals", "description") {
		roleName := rd.Get("role_name").(string)
		desc := rd.Get("description").(string)
		projectIds, userSrns, _ := convertTrustPrincipal(rd)

		_, err := inst.Client.Iam.UpdateRole(ctx, rd.Id(), roleName, projectIds, userSrns, desc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if rd.HasChanges("policy_ids") {
		addIds, removeIds := common.GetAddRemoveItemsStringListFromSet(rd, "policy_ids")

		if len(addIds) > 0 {
			_, err := inst.Client.Iam.AddRolePolicies(ctx, rd.Id(), addIds)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(removeIds) > 0 {
			_, err := inst.Client.Iam.RemoveRolePolicies(ctx, rd.Id(), removeIds)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceRoleRead(ctx, rd, meta)
}

func resourceRoleDestroy(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.Iam.DeleteRole(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
