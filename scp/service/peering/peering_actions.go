package peering

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_vpc_peering_cancel", ResourceVpcPeeringCancel())
	scp.RegisterResource("scp_vpc_peering_reject", ResourceVpcPeeringReject())
	scp.RegisterResource("scp_vpc_peering_approve", ResourceVpcPeeringApprove())
}

func ResourceVpcPeeringApprove() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceVpcPeeringApproveCreate,
		ReadContext:   resourceVpcPeeringActionRead,
		UpdateContext: resourceVpcPeeringApproveCreate,
		DeleteContext: resourceVpcPeeringActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("VpcPeeringId"):    {Type: schema.TypeString, Required: true, Description: "Vpc Peering Id"},
			common.ToSnakeCase("FirewallEnabled"): {Type: schema.TypeBool, Required: true, Description: "Firewall Enabled"},
			common.ToSnakeCase("VpcPeeringState"): {Type: schema.TypeString, Computed: true, Description: "Vpc Peering Id"},
		},
		Description: "Approve Peering Request.",
	}
}

func resourceVpcPeeringApproveCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcPeeringId := rd.Get(common.ToSnakeCase("VpcPeeringId")).(string)
	firewallEnabled := rd.Get(common.ToSnakeCase("FirewallEnabled")).(bool)

	inst := meta.(*client.Instance)

	result, err := inst.Client.Peering.ApproveVpcPeering(ctx, vpcPeeringId, firewallEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	if !(*result.Success) {
		// check when false //
	}

	err = waitVpcPeeringCreating(ctx, inst.Client, vpcPeeringId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.VpcPeeringId)
	return resourceVpcPeeringActionRead(ctx, rd, meta)
}

func ResourceVpcPeeringReject() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceVpcPeeringRejectCreate,
		ReadContext:   resourceVpcPeeringActionRead,
		UpdateContext: resourceVpcPeeringActionUpdate,
		DeleteContext: resourceVpcPeeringActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("VpcPeeringId"):    {Type: schema.TypeString, Required: true, Description: "Vpc Peering Id"},
			common.ToSnakeCase("VpcPeeringState"): {Type: schema.TypeString, Computed: true, Description: "Vpc Peering Id"},
		},
		Description: "Reject Peering Request.",
	}
}

func resourceVpcPeeringRejectCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcPeeringId := rd.Get(common.ToSnakeCase("VpcPeeringId")).(string)

	inst := meta.(*client.Instance)

	result, err := inst.Client.Peering.RejectVpcPeering(ctx, vpcPeeringId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !(*result.Success) {
		// check when false //
	}
	rd.SetId(result.VpcPeeringId)
	return resourceVpcPeeringActionRead(ctx, rd, meta)
}

func ResourceVpcPeeringCancel() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceVpcPeeringCancelCreate,
		ReadContext:   resourceVpcPeeringActionRead,
		UpdateContext: resourceVpcPeeringActionUpdate,
		DeleteContext: resourceVpcPeeringActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("VpcPeeringId"):    {Type: schema.TypeString, Required: true, Description: "Vpc Peering Id"},
			common.ToSnakeCase("VpcPeeringState"): {Type: schema.TypeString, Computed: true, Description: "Vpc Peering Id"},
		},
		Description: "Reject Peering Request.",
	}
}

func resourceVpcPeeringCancelCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcPeeringId := rd.Get(common.ToSnakeCase("VpcPeeringId")).(string)

	inst := meta.(*client.Instance)

	result, err := inst.Client.Peering.CancelVpcPeering(ctx, vpcPeeringId)
	if err != nil {
		return diag.FromErr(err)
	}

	if !(*result.Success) {
		// check when false //
	}
	rd.SetId(result.VpcPeeringId)
	return resourceVpcPeeringActionRead(ctx, rd, meta)
}

func resourceVpcPeeringActionRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	ruleInfo, _, err := inst.Client.Peering.GetVpcPeeringDetail(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}
	rd.Set(common.ToSnakeCase("VpcPeeringState"), ruleInfo.VpcPeeringState)

	return nil
}

func resourceVpcPeeringActionUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("Update function is not supported!")
}

func resourceVpcPeeringActionDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//return diag.Errorf("Delete function is not supported!")
	rd.SetId("")
	return nil
}

func waitVpcPeeringCreating(ctx context.Context, scpClient *client.SCPClient, peeringId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"CREATING"}, []string{"ACTIVE"}, func() (interface{}, string, error) {
		return scpClient.Peering.GetVpcPeeringDetail(ctx, peeringId)
	})
}
