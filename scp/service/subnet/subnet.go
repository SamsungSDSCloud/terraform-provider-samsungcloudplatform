package subnet

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"strings"
)

func init() {
	scp.RegisterResource("scp_subnet", ResourceSubnet())
}

func ResourceSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetCreate,
		ReadContext:   resourceSubnetRead,
		UpdateContext: resourceSubnetUpdate,
		DeleteContext: resourceSubnetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target VPC id",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet name. (3 to 20 characters without specials)",
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "must contain only alphanumeric characters"),
				),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Subnet description. (Up to 50 characters)",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet type ('PUBLIC'|'PRIVATE'|'BM'|'VM')",
				ValidateFunc: validation.StringInSlice([]string{
					"PUBLIC",
					"PRIVATE",
					"BM",
					"VM",
				}, false),
			},
			"cidr_ipv4": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Subnet cidr ipv4",
				ValidateFunc: validation.IsCIDR,
			},
		},
		Description: "Provides a Subnet resource.",
	}
}

func resourceSubnetCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	vpcId := rd.Get("vpc_id").(string)
	name := rd.Get("name").(string)
	description := rd.Get("description").(string)
	cidrIpv4 := rd.Get("cidr_ipv4").(string)
	subnetType := strings.ToUpper(rd.Get("type").(string))

	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Vpc.GetVpcInfo(ctx, vpcId)
	if err != nil {
		return diag.FromErr(err)
	}

	isNameInvalid, err := inst.Client.Subnet.CheckSubnetName(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}
	if isNameInvalid {
		return diag.Errorf("Input subnet name is invalid (maybe duplicated) : " + name)
	}

	isCidrInvalid, err := inst.Client.Subnet.CheckSubnetCidrIpv4(ctx, cidrIpv4, vpcId)
	if err != nil {
		return diag.FromErr(err)
	}
	if isCidrInvalid {
		return diag.Errorf("Input cidr is invalid (maybe duplicated) : " + cidrIpv4)
	}

	result, err := inst.Client.Subnet.CreateSubnet(ctx, vpcId, cidrIpv4, subnetType, name, description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForSubnetStatus(ctx, inst.Client, result.ResourceId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.ResourceId)

	return resourceSubnetRead(ctx, rd, meta)
}

func resourceSubnetRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	info, _, err := inst.Client.Subnet.GetSubnet(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	rd.Set("vpc_id", info.VpcId)
	rd.Set("cidr_ipv4", info.SubnetCidrBlock)
	rd.Set("name", info.SubnetName)
	rd.Set("description", info.SubnetDescription)
	rd.Set("type", info.SubnetType)

	return nil
}

func resourceSubnetUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("description") {
		_, err := inst.Client.Subnet.UpdateSubnetDescription(ctx, rd.Id(), rd.Get("description").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if rd.HasChanges("type") {
		_, err := inst.Client.Subnet.UpdateSubnetType(ctx, rd.Id(), rd.Get("type").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSubnetRead(ctx, rd, meta)
}

func resourceSubnetDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	_, err := inst.Client.Subnet.DeleteSubnet(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}
	err = waitForSubnetStatus(ctx, inst.Client, rd.Id(), []string{"TERMINATING"}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForSubnetStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.Subnet.GetSubnet(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		return info, info.SubnetState, nil
	})
}
