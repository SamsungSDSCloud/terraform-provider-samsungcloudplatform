package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAssessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: createAssessKey,
		ReadContext:   readAssessKey,
		UpdateContext: updateAssessKey,
		DeleteContext: deleteAssessKey,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"duration_days": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "",
			},
			"access_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"access_key_activated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func createAssessKey(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	response, err := inst.Client.Iam.CreateAccessKey(ctx, data.Get("project_id").(string), int32(data.Get("duration_days").(int)))

	if err != nil {
		if err.Error() == "400 Bad Request" {
			return diag.Errorf("400 Bad Request (Adding an encryption disk is only available on an encrypted Virtual Server.)")
		}
		return diag.FromErr(err)
	}

	if !data.Get("access_key_activated").(bool) {
		inst.Client.Iam.DactivateAccessKey(ctx, response.AccessKeyId)
	}

	data.SetId(response.AccessKeyId)
	data.Set("access_secret_key", response.AccessSecretKey)
	data.Set("access_key_activated", response.AccessKeyActivated)

	return nil
}

func readAssessKey(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func updateAssessKey(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if data.HasChanges("access_key_activated") {
		if data.Get("access_key_activated").(bool) {
			inst.Client.Iam.ActivateAccessKey(ctx, data.Id())
		} else {
			inst.Client.Iam.DactivateAccessKey(ctx, data.Id())
		}
	}
	return nil
}

func deleteAssessKey(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.Iam.DeleteAccessKey(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
