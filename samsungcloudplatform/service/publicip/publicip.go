package publicip

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	samsungcloudplatform.RegisterResource("samsungcloudplatform_public_ip", ResourceVpcPublicIp())
}

func ResourceVpcPublicIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcPublicIpCreate,
		ReadContext:   resourceVpcPublicIpRead,
		UpdateContext: resourceVpcPublicIpUpdate,
		DeleteContext: resourceVpcPublicIpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Description of public IP",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region name",
			},
			"ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of public IP",
			},
			"uplink_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Public IP uplinkType ('INTERNET'|'DEDICATED_INTERNET'|'SHARED_GROUP'|'SECURE_INTERNET')",
				ValidateFunc: validation.StringInSlice([]string{
					"INTERNET",
					"DEDICATED_INTERNET",
					"SHARED_GROUP",
					"SECURE_INTERNET",
				}, false),
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Public IP resource.",
	}
}

func resourceVpcPublicIpCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	description := rd.Get("description").(string)
	location := rd.Get("region").(string)
	uplinkType := rd.Get("uplink_type").(string)
	tags := rd.Get("tags").(map[string]interface{})
	inst := meta.(*client.Instance)

	serviceZoneId, err := client.FindServiceZoneId(ctx, inst.Client, location)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := inst.Client.PublicIp.CreatePublicIp(ctx, serviceZoneId, uplinkType, description, tags)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForVpcPublicIpStatus(ctx, inst.Client, result.PublicIpAddressId, []string{}, []string{common.ReservedState}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.PublicIpAddressId)

	return resourceVpcPublicIpRead(ctx, rd, meta)
}

func resourceVpcPublicIpRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)
	info, _, err := inst.Client.PublicIp.GetPublicIp(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("ipv4", info.IpAddress)
	rd.Set("description", info.PublicIpAddressDescription)

	location, err := client.FindLocationName(ctx, inst.Client, info.ServiceZoneId)
	if err != nil {
		tflog.Warn(ctx, "Failed to get service zone information")
	}

	rd.Set("region", location)
	rd.Set("uplinkType", info.UplinkType)

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceVpcPublicIpUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	if rd.HasChanges("description") {
		_, err := inst.Client.PublicIp.UpdatePublicIp(ctx, rd.Id(), rd.Get("description").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVpcPublicIpRead(ctx, rd, meta)
}

func resourceVpcPublicIpDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	err := inst.Client.PublicIp.DeletePublicIp(ctx, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForVpcPublicIpStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"DELETED", "FREE"}, false)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForVpcPublicIpStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.PublicIp.GetPublicIp(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.PublicIpState, nil
	})
}
