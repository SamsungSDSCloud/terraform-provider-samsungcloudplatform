package iam

import (
	"context"
	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: createMember,
		UpdateContext: updateMember,
		ReadContext:   readMember,
		DeleteContext: deleteMember,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func createMember(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readMember(ctx, data, meta)
}

func updateMember(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readMember(ctx, data, meta)
}

func deleteMember(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return readMember(ctx, data, meta)
}

func readMember(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	result, err := inst.Client.Iam.DetailMember(ctx, data.Get("user_id").(string))

	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	data.SetId(result.UserId)
	data.Set("user_id", result.UserId)
	data.Set("email", result.Email)
	data.Set("user_name", result.UserName)

	return nil
}
