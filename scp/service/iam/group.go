package iam

import (
	"context"
	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: createGroup,
		ReadContext:   readGroup,
		UpdateContext: updateGroup,
		DeleteContext: deleteGroup,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
		},
	}
}

func createGroup(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	response, err := inst.Client.Iam.CreatGroup(ctx, data.Get("group_name").(string), data.Get("description").(string))

	if err != nil {
		if err.Error() == "400 Bad Request" {
			return diag.Errorf("400 Bad Request")
		}
		return diag.FromErr(err)
	}

	data.SetId(response.GroupId)

	return readGroup(ctx, data, meta)
}

func readGroup(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	result, err := inst.Client.Iam.DetailGroup(ctx, data.Id())

	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	data.Set("group_name", result.GroupName)
	data.Set("description", result.Description)

	return nil
}
func updateGroup(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if data.HasChanges("group_name") || data.HasChanges("description") {
		inst.Client.Iam.UpdateGroup(ctx, data.Id(), data.Get("group_name").(string), data.Get("description").(string))
	}
	return nil
}
func deleteGroup(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.Iam.DeletePolicy(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
