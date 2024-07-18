package resourcegroup

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/resourcegroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_resource_group", ResourceResourceGroup())
}

func ResourceResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceGroupCreate,
		ReadContext:   resourceResourceGroupRead,
		UpdateContext: resourceResourceGroupUpdate,
		DeleteContext: resourceResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name":                 {Type: schema.TypeString, Required: true, Description: "Resource group name"},
			"target_resource_tags": tfTags.TagsSchema(),
			"resource_group_name":  {Type: schema.TypeString, Computed: true, Description: "Resource group name"},
			"target_resource_tag": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateName1to256DotDashUnderscore,
							Description:      "Tag key",
						},
						"tag_value": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: common.ValidateName1to256DotDashUnderscore,
							Description:      "Tag value",
						},
					},
				},
				Description: "Tag list",
			},
			"target_resource_types":      {Type: schema.TypeList, Optional: true, Computed: true, Description: "Resource group types", Elem: &schema.Schema{Type: schema.TypeString, Description: "type"}},
			"resource_group_description": {Type: schema.TypeString, Optional: true, Computed: true, Description: "Resource group description"},
			"created_by_id":              {Type: schema.TypeString, Optional: true, Description: "The user id which created the resource group"},
			"created_by_name":            {Type: schema.TypeString, Optional: true, Description: "The user name which created the resource group"},
			"created_by_email":           {Type: schema.TypeString, Computed: true, Description: "The user email which created the resource group"},
			"created_dt":                 {Type: schema.TypeString, Computed: true, Description: "The created date of the resource group"},
			"modified_by_id":             {Type: schema.TypeString, Optional: true, Description: "The user id which modified the resource group"},
			"modified_by_name":           {Type: schema.TypeString, Optional: true, Description: "The user name which modified the resource group"},
			"modified_by_email":          {Type: schema.TypeString, Computed: true, Description: "The user email which modified the resource group"},
			"modified_dt":                {Type: schema.TypeString, Computed: true, Description: "The modified date of the resource group"},
		},
	}
}

func resourceResourceGroupCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	result, _, err := inst.Client.ResourceGroup.CreateResourceGroup(ctx, resourcegroup.ResourceGroupRequest{
		ResourceGroupName:        rd.Get("name").(string),
		TargetResourceTags:       rd.Get("target_resource_tags").([]interface{}),
		TargetResourceTypes:      rd.Get("target_resource_types").([]interface{}),
		ResourceGroupDescription: rd.Get("resource_group_description").(string),
	})

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.ResourceGroupId)
	return resourceResourceGroupRead(ctx, rd, meta)
}

func resourceResourceGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)

	result, err := inst.Client.ResourceGroup.GetResourceGroup(ctx, rd.Id())
	if err != nil {
		return
	}

	rd.Set("resource_group_name", result.ResourceGroupName)
	rd.Set("target_resource_types", result.TargetResourceTypes)
	rd.Set("resource_group_description", result.ResourceGroupDescription)
	rd.Set("created_by_id", result.CreatedById)
	rd.Set("created_by_name", result.CreatedByName)
	rd.Set("created_by_email", result.CreatedByEmail)
	rd.Set("created_dt", result.CreatedDt)
	rd.Set("modified_by_id", result.ModifiedById)
	rd.Set("modified_by_name", result.ModifiedByName)
	rd.Set("modified_by_email", result.ModifiedByEmail)
	rd.Set("modified_dt", result.ModifiedDt)

	var tags common.HclSetObject
	for _, tag := range result.TargetResourceTag {
		kv := common.HclKeyValueObject{
			"tag_key":   tag.TagKey,
			"tag_value": tag.TagValue,
		}
		tags = append(tags, kv)
	}
	rd.Set("target_resource_tag", tags)

	return nil
}

func resourceResourceGroupUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	resourceGroupId := rd.Id()

	if rd.HasChanges("name", "target_resource_tags", "target_resource_types", "resource_group_description") {
		_, err = inst.Client.ResourceGroup.UpdateResourceGroup(ctx, resourceGroupId, resourcegroup.ResourceGroupRequest{
			ResourceGroupName:        rd.Get("name").(string),
			TargetResourceTags:       rd.Get("target_resource_tags").([]interface{}),
			TargetResourceTypes:      rd.Get("target_resource_types").([]interface{}),
			ResourceGroupDescription: rd.Get("resource_group_description").(string),
		})
		if err != nil {
			return
		}
	}

	if rd.HasChanges("target_resource_tags") {
		o, n := rd.GetChange("target_resource_tags")
		oldMap := o.(map[string]interface{})
		newMap := n.(map[string]interface{})

		err := client.UpdateResourceTag(ctx, inst.Client, rd.Id(), oldMap, newMap)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceResourceGroupRead(ctx, rd, meta)
}

func resourceResourceGroupDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	err := inst.Client.ResourceGroup.DeleteResourceGroup(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// no waiting needed
	return nil
}
