package transitgateway

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	scp.RegisterResource("scp_transit_gateway_connection", ResourceTransitGatewayConnection())
	scp.RegisterResource("scp_transit_gateway_connection_approve", ResourceTransitGatewayConnectionApprove())
}

func ResourceTransitGatewayConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTransitGatewayConnectionCreate,
		ReadContext:   resourceTransitGatewayConnectionRead,
		UpdateContext: resourceTransitGatewayConnectionUpdate,
		DeleteContext: resourceTransitGatewayConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"requester_transit_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Requester TGW ID",
				ValidateFunc: validation.StringLenBetween(3, 60),
			},
			"approver_vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Approver VPC ID",
				ValidateFunc: validation.StringLenBetween(3, 60),
			},
			"transit_gateway_connection_description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "TGW - VPC Connection description",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
			"firewall_enable": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Activate Firewall or not",
			},
			"firewall_loggable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Activate Firewall Logging or not",
			},
			"transit_gateway_connection_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Transit Gateway Connection State",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a TGW- VPC connection resource.",
	}
}

func ResourceTransitGatewayConnectionApprove() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTransitGatewayConnectionApprove,
		ReadContext:   resourceTransitGatewayConnectionRead,
		UpdateContext: resourceTransitGatewayConnectionUpdate,
		DeleteContext: resourceTransitGatewayConnectionPseudoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"transit_gateway_connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "TGW-VPC Connection ID",
			},
			"transit_gateway_connection_description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "TGW - VPC Connection description",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
			"transit_gateway_connection_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Transit Gateway Connection State",
			},
		},
		Description: "Approve TGW-VPC Connection",
	}
}

func resourceTransitGatewayConnectionCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	requesterTransitGatewayId := rd.Get("requester_transit_gateway_id").(string)
	approverVpcId := rd.Get("approver_vpc_id").(string)
	firewallEnabled := rd.Get("firewall_enable").(bool)
	firewallLogging := rd.Get("firewall_loggable").(bool)
	tgwConnectionDescription := rd.Get("transit_gateway_connection_description").(string)
	tgwConnectionType := "INTERNAL"
	tags := rd.Get("tags").(map[string]interface{})

	inst := meta.(*client.Instance)

	tgwInfo, _, err := inst.Client.TransitGateway.GetTransitGatewayInfo(ctx, requesterTransitGatewayId)
	if err != nil {
		return diag.FromErr(err)
	}

	vpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, approverVpcId)
	if err != nil {
		return diag.FromErr(err)
	}

	response, _, err := inst.Client.TransitGateway.CreateTransitGatewayConnection(ctx, requesterTransitGatewayId, approverVpcId, tgwInfo.ProjectId, vpcInfo.ProjectId, tgwConnectionDescription, firewallEnabled, firewallLogging, tgwConnectionType, tags)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.TransitGatewayConnectionId)

	return resourceTransitGatewayConnectionRead(ctx, rd, meta)

}

func resourceTransitGatewayConnectionRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.TransitGateway.GetTransitGatewayConnectionInfo(ctx, rd.Id())

	if err != nil {
		return diag.FromErr(err)
	}
	rd.Set("transit_gateway_connection_state", info.TransitGatewayConnectionState)
	rd.Set("transit_gateway_connection_description", info.TransitGatewayConnectionDescription)

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceTransitGatewayConnectionApprove(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	tgwConnectionId := rd.Get("transit_gateway_connection_id").(string)

	result, _, err := inst.Client.TransitGateway.ApproveTransitGatewayConnection(ctx, tgwConnectionId)

	if err != nil {
		return diag.FromErr(err)
	}

	if !*result.Success {
		diag.Errorf("Approve TGW -VPC Connection was failed. Approval Client call was failed.")
	}

	rd.SetId(tgwConnectionId)

	err = waitTransitGatewayConnectionCreating(ctx, inst.Client, rd.Id())

	if err != nil {
		return diag.FromErr(err)

	}

	return resourceTransitGatewayConnectionRead(ctx, rd, meta)
}

func resourceTransitGatewayConnectionUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("transit_gateway_connection_description") {
		_, _, err := inst.Client.TransitGateway.UpdateTransitGatewayConnectionDescription(ctx, rd.Id(), rd.Get("transit_gateway_connection_description").(string))
		if err != nil {
			diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTransitGatewayConnectionRead(ctx, rd, meta)

}

func resourceTransitGatewayConnectionDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.Get("transit_gateway_connection_state").(string) == "REQUESTING" {
		result, _, err := inst.Client.TransitGateway.CancelTransitGatewayConnection(ctx, rd.Id())
		if err != nil {
			diag.FromErr(err)
		}
		if !*result.Success {
			diag.Errorf("Cancel TGW -VPC Connection was failed. Approval Client call was failed.")
		}
	}

	err := inst.Client.TransitGateway.DeleteTransitGatewayConnection(ctx, rd.Id())

	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}
	err = waitTransitGatewayConnectionDeleting(ctx, inst.Client, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil

}

func resourceTransitGatewayConnectionPseudoDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rd.SetId("")
	return nil
}

func waitTransitGatewayConnectionCreating(ctx context.Context, scpClient *client.SCPClient, transitGatewayConnectionId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING", "PROGRESSING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		info, _, err := scpClient.TransitGateway.GetTransitGatewayConnectionInfo(ctx, transitGatewayConnectionId)
		if err != nil {
			return nil, "", err
		}
		return info, info.TransitGatewayConnectionState, nil
	})
}

func waitTransitGatewayConnectionDeleting(ctx context.Context, scpClient *client.SCPClient, transitGatewayConnectionId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		info, c, err := scpClient.TransitGateway.GetTransitGatewayConnectionInfo(ctx, transitGatewayConnectionId)
		if err != nil {
			if c == 404 {
				return "", "DELETED", nil
			}
			if c == 403 {
				return nil, "", err
			}
		}
		return info, info.TransitGatewayConnectionState, nil
	})
}
