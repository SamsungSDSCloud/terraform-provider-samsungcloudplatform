package natgateway

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/service/tag"
	publicip2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/public-ip2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	samsungcloudplatform.RegisterResource("samsungcloudplatform_nat_gateway", ResourceNATGateway())
}

func ResourceNATGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNATGatewayCreate,
		ReadContext:   resourceNATGatewayRead,
		UpdateContext: resourceNATGatewayUpdate,
		DeleteContext: resourceNATGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target VPC id",
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "NAT-Gateway public IP. If not set, it will be auto generated.",
			},
			"public_ipv4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NAT-Gateway public IP.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "NAT-Gateway description. (Up to 50 characters)",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a NAT Gateway resource.",
	}
}

func resourceNATGatewayCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			diagnostics = diag.FromErr(err)
		}
	}()

	subnetId := rd.Get("subnet_id").(string)
	publicIpId := rd.Get("public_ip_id").(string)
	description := rd.Get("description").(string)
	tags := rd.Get("tags").(map[string]interface{})

	inst := meta.(*client.Instance)

	result, _, err := inst.Client.NatGateway.CreateNatGateway(ctx, publicIpId, subnetId, description, tags)
	if err != nil {
		return
	}

	err = waitForNATStatus(ctx, inst.Client, result.ResourceId, []string{"CREATING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return
	}

	rd.SetId(result.ResourceId)

	return resourceNATGatewayRead(ctx, rd, meta)
}

func resourceNATGatewayRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) (diagnostics diag.Diagnostics) {
	var err error = nil
	defer func() {
		if err != nil {
			rd.SetId("")
			diagnostics = diag.FromErr(err)
		}
	}()

	inst := meta.(*client.Instance)

	info, _, err := inst.Client.NatGateway.GetNatGateway(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	publicIpInfos, err := inst.Client.PublicIp.GetPublicIps(ctx, &publicip2.PublicIpOpenApiV3ControllerApiListPublicIpsV3Opts{
		IpAddress:     optional.NewString(info.NatGatewayIpAddress),
		VpcId:         optional.NewString(info.VpcId),
		PublicIpState: optional.String{},
		UplinkType:    optional.String{},
		CreatedBy:     optional.String{},
		Page:          optional.NewInt32(0),
		Size:          optional.NewInt32(10000),
		Sort:          optional.Interface{},
	})
	if err != nil {
		return
	}
	var publicIpId string
	for _, publicIpInfo := range publicIpInfos.Contents {
		if publicIpInfo.IpAddress == info.NatGatewayIpAddress {
			publicIpId = publicIpInfo.PublicIpAddressId
			break
		}
	}
	//if len(publicIpId) == 0 {
	//	err = fmt.Errorf("target public ip address id not found")
	//	return
	//}

	rd.Set("subnet_id", info.SubnetId)
	rd.Set("public_ip_id", publicIpId)
	rd.Set("public_ipv4", info.NatGatewayIpAddress)
	rd.Set("description", info.NatGatewayDescription)

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceNATGatewayUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	if rd.HasChanges("description") {
		_, _, err := inst.Client.NatGateway.UpdateNatGateway(ctx, rd.Id(), rd.Get("description").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForNATStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNATGatewayRead(ctx, rd, meta)
}

func resourceNATGatewayDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	_, _, err := inst.Client.NatGateway.DeleteNatGateway(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitForNATStatus(ctx, inst.Client, rd.Id(), []string{"DELETING"}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForNATStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.NatGateway.GetNatGateway(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.NatGatewayState, nil
	})
}
