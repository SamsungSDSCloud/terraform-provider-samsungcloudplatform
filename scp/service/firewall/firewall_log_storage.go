package firewall

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/firewall2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_firewall_logstorage", resourceFirewallLogStorage())
}

func resourceFirewallLogStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallLogStorageCreate,
		ReadContext:   resourceFirewallLogStorageRead,
		UpdateContext: resourceFirewallLogStorageUpdate,
		DeleteContext: resourceFirewallLogStorageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC id",
			},
			"obs_bucket_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Object storage bucket id to save firewall log",
			},
		},
		Description: "Set up firewall log storage.",
	}
}

func resourceFirewallLogStorageCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	vpcId := rd.Get("vpc_id").(string)
	obsBucketId := rd.Get("obs_bucket_id").(string)

	firewallLogStorageCreatRequest := firewall2.FirewallLogStorageCreatRequest{
		LogStorageType: "FIREWALL",
		ObsBucketId:    obsBucketId,
		VpcId:          vpcId,
	}
	logStorage, _, err := inst.Client.Firewall.CreateFirewallLogStorage(ctx, firewallLogStorageCreatRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(logStorage.LogStorageId)

	return resourceFirewallLogStorageRead(ctx, rd, meta)
}

func resourceFirewallLogStorageRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	logStageInfo, _, err := inst.Client.Firewall.GetFirewallLogStorage(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.Set("vpc_id", logStageInfo.VpcId)
	rd.Set("obs_bucket_id", logStageInfo.ObsBucketId)

	return nil
}

func resourceFirewallLogStorageUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("obs_bucket_id") {

		obsBucketId := rd.Get("obs_bucket_id").(string)

		if _, _, err := inst.Client.Firewall.UpdateFirewallLogStorage(ctx, rd.Id(), obsBucketId); err != nil {
			return diag.FromErr(err)
		}

	}

	return resourceFirewallLogStorageRead(ctx, rd, meta)
}

func resourceFirewallLogStorageDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	if err := inst.Client.Firewall.DeleteFirewallLogStorage(ctx, rd.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
