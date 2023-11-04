package transitgateway

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"time"
)

func init() {
	scp.RegisterResource("scp_transit_gateway", ResourceTransitGateway())
}

func ResourceTransitGateway() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceTransitGatewayCreate,
		ReadContext:   resourceTransitGatewayRead,
		UpdateContext: resourceTransitGatewayUpdate,
		DeleteContext: resourceTransitGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"transit_gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Transit Gateway Name. ( 3 to 20 characters consist of alphabets and numbers)",
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "must contain only alphabets and numbers"),
				),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region Name",
			},
			"uplink_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Option for Uplink",
			},
			"bandwidth_gbps": {
				Type:             schema.TypeInt,
				Required:         true,
				ForceNew:         true,
				Description:      "Bandwidth Gbps(UplinkEnable=false : 1, UplinkEnable=true : 1/10, Reserved for designated: 20/40)",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{1, 10})),
			},
			"transit_gateway_description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Transit Gateway description. (Up to 50 characters)",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
		},
		Description: "Provides  a TransitGateway resource.",
	}
}

func resourceTransitGatewayCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := rd.Get("transit_gateway_name").(string)
	uplinkEnabled := rd.Get("uplink_enabled").(bool)
	bandwidthGbps := int32(rd.Get("bandwidth_gbps").(int))
	location := rd.Get("region").(string)
	description := rd.Get("transit_gateway_description").(string)

	inst := meta.(*client.Instance)

	serviceZoneId, err := client.FindServiceZoneId(ctx, inst.Client, location)

	response, _, err := inst.Client.TransitGateway.CreateTransitGateway(ctx, bandwidthGbps, serviceZoneId, name, uplinkEnabled, description)

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitTransitGatewayCreating(ctx, inst.Client, response.ResourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ResourceId)
	return resourceTransitGatewayRead(ctx, rd, meta)

}
func resourceTransitGatewayRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	tgwInfo, _, err := inst.Client.TransitGateway.GetTransitGatewayInfo(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	rd.Set("transit_gateway_name", tgwInfo.TransitGatewayName)
	rd.Set("transit_gateway_description", tgwInfo.TransitGatewayDescription)

	location, err := client.FindLocationName(ctx, inst.Client, tgwInfo.ServiceZoneId)
	if err != nil {
		tflog.Warn(ctx, "Failed to get service zone information from Location")
	}
	rd.Set("region", location)
	return nil
}

func resourceTransitGatewayUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	if rd.HasChanges("uplink_enabled") {
		_, _, err := inst.Client.TransitGateway.UpdateTransitGatewayUplinkEnable(ctx, rd.Id(), rd.Get("uplink_enabled").(bool))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceTransitGatewayRead(ctx, rd, meta)
}

func resourceTransitGatewayDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	err := inst.Client.TransitGateway.DeleteTransitGateway(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	time.Sleep(10 * time.Second)

	err = waitTransitGatewayDeleting(ctx, inst.Client, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitTransitGatewayCreating(ctx context.Context, scpClient *client.SCPClient, transitGatewayId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		tgwDetail, _, err := scpClient.TransitGateway.GetTransitGatewayInfo(ctx, transitGatewayId)
		if err != nil {
			return nil, "", err
		}
		return tgwDetail, tgwDetail.TransitGatewayState, nil
	})
}

func waitTransitGatewayDeleting(ctx context.Context, scpClient *client.SCPClient, transitGatewayId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		tgwDetail, c, err := scpClient.TransitGateway.GetTransitGatewayInfo(ctx, transitGatewayId)
		if err != nil {
			if c == 404 {
				return "", "DELETED", nil
			}
			if c == 403 {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return tgwDetail, tgwDetail.TransitGatewayState, nil
	})
}
