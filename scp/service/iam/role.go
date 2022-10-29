package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: createRole,
		ReadContext:   readRole,
		UpdateContext: updateRole,
		DeleteContext: deleteRole,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			/*"trust_principals": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    false,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_ids": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "",
						},
						"user_srns": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "",
						},
					},
				},
			},*/
			"trust_principals": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_ids": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "",
						},
						"user_srns": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func convertTrustPricipal(itemObject common.HclKeyValueObject) (iam.TrustPrincipalsResponse, error) {
	result := iam.TrustPrincipalsResponse{}
	if projectIds, ok := itemObject["project_ids"]; ok {
		result.ProjectIds = projectIds.([]string)
	}

	if userSrns, ok := itemObject["user_srns"]; ok {
		result.UserSrns = userSrns.([]string)
	}
	return result, nil
}

func expandStringArray(rd *schema.ResourceData) []string {
	addressesIpv4List := rd.Get("addresses_ipv4").([]interface{})
	addressesIpv4 := make([]string, len(addressesIpv4List))
	for i, valueIpv4 := range addressesIpv4List {
		addressesIpv4[i] = valueIpv4.(string)
	}
	return addressesIpv4
}

func convertTrustPricipal1(rd *schema.ResourceData) (iam.TrustPrincipalsResponse, error) {
	servicesSet := rd.Get("trust_principals").(*schema.Set).List()
	services := make([]iam.TrustPrincipalsResponse, len(servicesSet))
	for i, valueService := range servicesSet {
		s := valueService.(map[string]interface{})

		if t, ok := s["project_ids"]; ok {
			project_ids := t.([]interface{})
			list := make([]string, len(project_ids))
			for i, valueIpv4 := range project_ids {
				list[i] = valueIpv4.(string)
			}
			services[i].ProjectIds = list
		}

		if v, ok := s["user_srns"]; ok {
			user_srns := v.([]interface{})
			list := make([]string, len(user_srns))
			for i, valueIpv4 := range user_srns {
				list[i] = valueIpv4.(string)
			}
			services[i].UserSrns = list
		}
	}
	return services[0], nil
}

func createRole(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	principals, err := convertTrustPricipal1(data)
	if err != nil {
		return nil
	}

	response, err := inst.Client.Iam.CreateRole(ctx, data.Get("role_name").(string), principals, data.Get("description").(string))

	if err != nil {
		if err.Error() == "400 Bad Request" {
			return diag.Errorf("400 Bad Request")
		}
		return diag.FromErr(err)
	}

	data.SetId(response.RoleId)

	return readPolicy(ctx, data, meta)
}

func readRole(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	result, err := inst.Client.Iam.DetailRole(ctx, data.Id())

	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	data.Set("role_name", result.RoleName)
	data.Set("trust_principals", result.TrustPrincipals)
	data.Set("description", result.Description)

	return nil
}

func updateRole(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	principals, err := convertTrustPricipal(data.Get("trust_principals").(common.HclKeyValueObject))
	if err != nil {
		return nil
	}

	inst.Client.Iam.UpdateRole(ctx, data.Id(), data.Get("role_name").(string), principals, data.Get("description").(string))

	return nil
}

func deleteRole(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	_, err := inst.Client.Iam.DeleteRole(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
