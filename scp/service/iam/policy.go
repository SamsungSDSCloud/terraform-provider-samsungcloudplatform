package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: createPolicy,
		ReadContext:   readPolicy,
		UpdateContext: updatePolicy,
		DeleteContext: deletePolicy,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"policy_json": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"principals": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"principal_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"principal_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func convertPricipal(list common.HclListObject) ([]iam.PolicyPrincipalRequest, error) {
	var result []iam.PolicyPrincipalRequest
	for _, l := range list {
		itemObject := l.(common.HclKeyValueObject)
		info := iam.PolicyPrincipalRequest{}
		if principal_id, ok := itemObject["principal_id"]; ok {
			info.PrincipalId = principal_id.(string)
		}
		if principal_type, ok := itemObject["principal_type"]; ok {
			info.PrincipalType = principal_type.(string)
		}

		result = append(result, info)
	}
	return result, nil
}

func createPolicy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, statusCode, err := inst.Client.Iam.ValidPolicyJson(ctx, data.Get("policy_json").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if statusCode == 400 {
		return diag.Errorf("JSON format Wrong.")
	}

	principals, err := convertPricipal(data.Get("principals").(common.HclListObject))
	if err != nil {
		return nil
	}

	response, err := inst.Client.Iam.CreatePolicy(ctx, data.Get("policy_name").(string), data.Get("policy_json").(string), principals, data.Get("description").(string))

	if err != nil {
		if err.Error() == "400 Bad Request" {
			return diag.Errorf("400 Bad Request")
		}
		return diag.FromErr(err)
	}

	data.SetId(response.PolicyId)

	return readPolicy(ctx, data, meta)
}

func readPolicy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	result, err := inst.Client.Iam.DetailPolicy(ctx, data.Id())

	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	data.Set("policy_name", result.PolicyName)
	data.Set("policy_json", result.PolicyJson)
	data.Set("description", result.Description)

	return nil
}
func updatePolicy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, statusCode, err := inst.Client.Iam.ValidPolicyJson(ctx, data.Get("policy_json").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if statusCode == 400 {
		return diag.Errorf("JSON format Wrong.")
	}

	principals, err := convertPricipal(data.Get("principals").(common.HclListObject))
	if err != nil {
		return nil
	}

	inst.Client.Iam.UpdatePolciy(ctx, data.Id(), data.Get("policy_json").(string), data.Get("policy_name").(string), principals, data.Get("description").(string))

	return nil
}
func deletePolicy(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	_, err := inst.Client.Iam.DeletePolicy(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
