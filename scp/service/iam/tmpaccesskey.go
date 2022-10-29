package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTmpAssessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: createTmpAccessKey,
		ReadContext:   readTmpAccessKey,
		UpdateContext: updateTmpAccessKey,
		DeleteContext: deleteTmpAccessKey,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"otp": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"duration_minutes": {
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

func createTmpAccessKey(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	response, err := inst.Client.Iam.CreateTemporaryAccessKey(ctx, data.Get("otp").(string), int32(data.Get("duration_minutes").(int)))

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

func readTmpAccessKey(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func updateTmpAccessKey(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if data.HasChanges("access_key_activated") {
		if data.Get("access_key_activated").(bool) {
			inst.Client.Iam.ActivateTmpAccessKey(ctx, data.Id())
		} else {
			inst.Client.Iam.DactivateTmpAccessKey(ctx, data.Id())
		}
	}
	return nil
}

func deleteTmpAccessKey(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.Iam.DeleteTmpAccessKey(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
