package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	scp.RegisterResource("scp_iam_access_key", ResourceAssessKey())
}
func ResourceAssessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessKeyCreate,
		ReadContext:   resourceAccessKeyRead,
		UpdateContext: resourceAccessKeyUpdate,
		DeleteContext: resourceAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID",
			},
			"duration_days": {
				Type:             schema.TypeInt,
				Required:         true,
				Description:      "Expiration time (days), set to zero to get permanent key",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 365)),
			},
			"access_key_activated": {Type: schema.TypeBool, Optional: true, Computed: true, Description: "Access key activation"},

			"access_key":        {Type: schema.TypeString, Computed: true, Description: "Access key"},
			"access_key_id":     {Type: schema.TypeString, Computed: true, Description: "Access key ID"},
			"access_key_state":  {Type: schema.TypeString, Computed: true, Description: "Access key state"},
			"access_secret_key": {Type: schema.TypeString, Computed: true, Description: "Access secret key"},
			"expired_dt":        {Type: schema.TypeString, Computed: true, Description: "Expired date"},
			"project_name":      {Type: schema.TypeString, Computed: true, Description: "Project name"},
		},
		Description: "Provides IAM access key resource.",
	}
}

func resourceAccessKeyCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	response, err := inst.Client.Iam.CreateAccessKey(ctx, rd.Get("project_id").(string), int32(rd.Get("duration_days").(int)))

	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.AccessKeyId)
	rd.Set("access_secret_key", response.AccessSecretKey) // can only be set here
	if *response.AccessKeyActivated {

	}

	return resourceAccessKeyRead(ctx, rd, meta)
}

func resourceAccessKeyRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	emptyString := ""
	listKeys, err := inst.Client.Iam.ListAccessKeys(ctx, emptyString, emptyString, emptyString, optional.Bool{}, emptyString)
	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	curId := rd.Id()
	for _, key := range listKeys.Contents {
		if key.AccessKeyId == curId {
			found = true

			rd.Set("project_id", key.ProjectId)
			rd.Set("access_key", key.AccessKey)
			rd.Set("access_key_activated", key.AccessKeyActivated)
			rd.Set("access_key_state", key.AccessKeyState)
			rd.Set("expired_dt", key.ExpiredDt)
			rd.Set("project_name", key.ProjectName)
		}
	}

	if !found {
		return diag.Errorf("access key ID %s was not found", curId)
	}

	return nil
}

func resourceAccessKeyUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("access_key_activated") {
		if rd.Get("access_key_activated").(bool) {
			_, err := inst.Client.Iam.ActivateAccessKey(ctx, rd.Id())
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			_, err := inst.Client.Iam.DeactivateAccessKey(ctx, rd.Id())
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return resourceAccessKeyRead(ctx, rd, meta)
}

func resourceAccessKeyDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	// deactivate 먼저
	_, err := inst.Client.Iam.DeactivateAccessKey(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = inst.Client.Iam.DeleteAccessKey(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
