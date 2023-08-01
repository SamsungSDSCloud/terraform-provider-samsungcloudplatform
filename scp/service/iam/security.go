package iam

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurity() *schema.Resource {
	return &schema.Resource{
		CreateContext: createSecurity,
		ReadContext:   readSecurity,
		UpdateContext: updateSecurity,
		DeleteContext: deleteSecurity,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ip_acl_activated": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "",
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Source ip addresses list",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"mfa_activated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func getIps(rd *schema.ResourceData) []string {
	ipAddrs := rd.Get("ip_addresses").([]interface{})
	ips := make([]string, len(ipAddrs))
	for i, valueIpv4 := range ipAddrs {
		ips[i] = valueIpv4.(string)
	}
	return ips
}

func createSecurity(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId("1")
	return readSecurity(ctx, data, meta)
}

func readSecurity(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*
			result, err := inst.Client.Iam.GetSecurityInfo(ctx)

			if err != nil {
				data.SetId("")
				return diag.FromErr(err)
			}

			data.Set("ip_acl_activated", result.IpAclActivated)
			data.Set("ip_addresses", result.IpAddresses)
			data.Set("mfa_activated", result.MfaActivated)
		inst := meta.(*client.Instance)
	*/

	return nil
}

func updateSecurity(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.Iam.UpdateSecurityInfo(ctx, data.Get("ip_acl_activated").(bool), getIps(data), data.Get("mfa_activated").(bool))
	if err != nil {
		return diag.FromErr(err)
	}
	return readSecurity(ctx, data, meta)
}

func deleteSecurity(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
