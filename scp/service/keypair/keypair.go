package keypair

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/keypair"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_key_pair", ResourceKeyPair())
}

func ResourceKeyPair() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceKeyPairCreate,
		ReadContext:   resourceKeyPairRead,
		UpdateContext: resourceKeyPairUpdate,
		DeleteContext: resourceKeyPairDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key_pair_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateName3to28Dash,
				Description:      "Key Pair Name",
			},
			"private_key": {
				Type:             schema.TypeString,
				Computed:         true,
				ValidateDiagFunc: nil,
				Description:      "Private Key",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags",
			},
		},
	}
}

func resourceKeyPairCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	keyPairName := rd.Get("key_pair_name").(string)
	tags := rd.Get("tags").(map[string]interface{})
	tagsRequests := make([]keypair.TagRequest, 0)
	for key, value := range tags {
		tagsRequests = append(tagsRequests, keypair.TagRequest{
			TagKey:   key,
			TagValue: value.(string),
		})
	}

	response, err := inst.Client.KeyPair.CreateKeyPair(ctx, keypair.CreateRequest{
		KeyPairName: keyPairName,
		Tags:        tagsRequests,
	})
	if err != nil {
		return
	}

	err = WaitForKeyPairStatus(ctx, inst.Client, response.KeyPairId, []string{}, []string{common.ActiveState}, true)
	if err != nil {
		return
	}

	rd.SetId(response.KeyPairId)
	rd.Set("private_key", response.PrivateKey)
	rd.Set("key_pair_name", response.KeyPairName)

	return nil
}

func resourceKeyPairRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)
	keyPairInfo, _, err := inst.Client.KeyPair.GetKeyPair(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	rd.Set("key_pair_name", keyPairInfo.KeyPairName)

	return nil
}

func resourceKeyPairUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	return nil
}

func resourceKeyPairDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	inst := meta.(*client.Instance)
	err := inst.Client.KeyPair.DeleteKeyPair(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = WaitForKeyPairStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func WaitForKeyPairStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.KeyPair.GetKeyPair(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			if c == 403 && !errorOnNotFound {
				return "", common.DeletedState, nil
			}
			return nil, "", err
		}
		return info, info.KeyPairState, nil
	})
}
