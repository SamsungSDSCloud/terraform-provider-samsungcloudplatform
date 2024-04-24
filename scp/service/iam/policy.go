package iam

import (
	"context"
	"fmt"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	iam2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/iam"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"unicode/utf8"
)

func init() {
	scp.RegisterResource("scp_iam_policy", ResourcePolicy())
}
func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"policy_json": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 65535)),
				Description:      "Policy json statement",
			},
			"policy_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateNameHangeulAlphabetSomeSpecials3to64,
				Description:      "Policy name",
			},
			"principals": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"principal_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 60)),
							Description:      "Principal ID",
						},
						"principal_type": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 50)),
							Description:      "Principal type",
						},
					},
				},
				Description: "Policy principal list",
			},
			"tags": tfTags.TagsSchema(),
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 1000)),
				Description:      "Description",
			},

			"project_id":             {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"policy_id":              {Type: schema.TypeString, Computed: true, Description: "Policy ID"},
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
		},
		Description: "Provides IAM policy resource.",
	}
}

func resourcePolicyCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	policyJson := rd.Get("policy_json").(string)
	_, statusCode, err := inst.Client.Iam.ValidatePolicyJson(ctx, policyJson)
	if err != nil {
		return diag.FromErr(err)
	}
	if statusCode == 400 {
		return diag.Errorf("Invalid json format")
	}

	policyName := rd.Get("policy_name").(string)
	principals := toPrincipalRequestList(rd.Get("principals").([]interface{}))

	response, err := inst.Client.Iam.CreatePolicy(ctx, policyName, policyJson, principals, rd.Get("tags").(map[string]interface{}), rd.Get("description").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.PolicyId)
	return resourcePolicyRead(ctx, rd, meta)
}

func resourcePolicyRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	result, err := inst.Client.Iam.DetailPolicy(ctx, rd.Id())

	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.Set("policy_name", result.PolicyName)
	rd.Set("policy_json", result.PolicyJson)
	rd.Set("description", result.Description)
	rd.Set("project_id", result.ProjectId)
	rd.Set("policy_id", result.PolicyId)
	rd.Set("policy_principal_count", result.PolicyPrincipalCount)
	rd.Set("policy_srn", result.PolicySrn)
	rd.Set("policy_type", result.PolicyType)
	rd.Set("policy_version", result.PolicyVersion)
	rd.Set("created_by", result.CreatedBy)
	rd.Set("created_by_name", result.CreatedByName)
	rd.Set("created_by_email", result.CreatedByEmail)
	rd.Set("created_dt", result.CreatedDt.String())
	rd.Set("modified_by", result.ModifiedBy)
	rd.Set("modified_by_name", result.ModifiedByName)
	rd.Set("modified_by_email", result.ModifiedByEmail)
	rd.Set("modified_dt", result.ModifiedDt.String())

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}
func resourcePolicyUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("policy_json", "principals", "policy_name", "description") {
		policyJson := rd.Get("policy_json").(string)
		_, statusCode, err := inst.Client.Iam.ValidatePolicyJson(ctx, policyJson)
		if err != nil {
			return diag.FromErr(err)
		}
		if statusCode == 400 {
			return diag.Errorf("Invalid json format.")
		}

		policyName := rd.Get("policy_name").(string)
		principals := toPrincipalRequestList(rd.Get("principals").([]interface{}))
		desc := rd.Get("description").(string)

		_, err = inst.Client.Iam.UpdatePolicy(ctx, rd.Id(), policyJson, policyName, principals, desc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyRead(ctx, rd, meta)
}

func resourcePolicyDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.Iam.DeletePolicy(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func validatePolicyName(v interface{}, path cty.Path) diag.Diagnostics {
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
	} else if cnt > 64 {
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

func toPrincipalRequestList(list []interface{}) []iam2.PolicyPrincipalRequest {
	if len(list) == 0 {
		return nil
	}
	var result []iam2.PolicyPrincipalRequest

	for _, val := range list {
		kv := val.(common.HclKeyValueObject)
		result = append(result, iam2.PolicyPrincipalRequest{
			PrincipalId:   kv["principal_id"].(string),
			PrincipalType: kv["principal_type"].(string),
		})
	}
	return result
}
