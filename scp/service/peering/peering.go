package peering

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/peering"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	scp.RegisterResource("scp_vpc_peering", ResourceVpcPeering())
}

func ResourceVpcPeering() *schema.Resource {
	return &schema.Resource{
		//CRUD
		CreateContext: resourceVpcPeeringCreate,
		ReadContext:   resourceVpcPeeringRead,
		UpdateContext: resourceVpcPeeringUpdate,
		DeleteContext: resourceVpcPeeringDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"approver_vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Approver VPC Id",
				ValidateFunc: validation.StringLenBetween(3, 100),
			},
			"firewall_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Firewall Enabled",
			},
			"requester_vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Requester VPC Id",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"vpc_peering_description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "VPC Peering Description",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			common.ToSnakeCase("VpcPeeringState"): {Type: schema.TypeString, Computed: true, Description: "Vpc Peering State"},
			"tags":                                tfTags.TagsSchema(),
		},
		Description: "Provides a VPC Peering Rule.",
	}
}

func resourceVpcPeeringCreate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get values from schema
	approverVpcId := rd.Get("approver_vpc_id").(string)
	firewallEnabled := rd.Get("firewall_enabled").(bool)
	requesterVpcId := rd.Get("requester_vpc_id").(string)
	vpcPeeringDescription := rd.Get("vpc_peering_description").(string)

	inst := meta.(*client.Instance)

	approverVpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, approverVpcId)

	// from vpc requesterVpcInfo get project-id & product-group-id
	requesterVpcInfo, _, err := inst.Client.Vpc.GetVpcInfo(ctx, requesterVpcId)
	if err != nil {
		return diag.FromErr(err)
	}

	request := peering.VpcPeeringCreateRequest{
		ApproverProjectId:     approverVpcInfo.ProjectId,
		ApproverVpcId:         approverVpcId,
		FirewallEnabled:       firewallEnabled,
		RequesterProjectId:    requesterVpcInfo.ProjectId,
		RequesterVpcId:        requesterVpcId,
		VpcPeeringDescription: vpcPeeringDescription,
		Tags:                  rd.Get("tags").(map[string]interface{}),
	}

	tflog.Debug(ctx, "Try create vpc peering : "+approverVpcId+", "+requesterVpcId)

	result, err := inst.Client.Peering.CreateVpcPeering(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.VpcPeeringId)

	// Refresh
	return resourceVpcPeeringRead(ctx, rd, meta)
}

func resourceVpcPeeringRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	ruleInfo, _, err := inst.Client.Peering.GetVpcPeeringDetail(ctx, rd.Id())
	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		//except resource deleted error
		return diag.FromErr(err)
	}

	rd.Set("approver_project_id", ruleInfo.ApproverProjectId)
	rd.Set("approver_vpc_id", ruleInfo.ApproverVpcId)
	rd.Set("firewall_enabled", ruleInfo.RequesterVpcFirewallEnabled)
	rd.Set("requester_project_id", ruleInfo.RequesterProjectId)
	rd.Set("requester_vpc_id", ruleInfo.RequesterVpcId)
	rd.Set("vpc_peering_description", ruleInfo.VpcPeeringDescription)
	rd.Set("vpc_peering_state", ruleInfo.VpcPeeringState)

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func resourceVpcPeeringUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	if rd.HasChanges("approver_project_id", "approver_vpc_id", "firewall_enabled", "requester_vpc_id", "vpc_peering_type") {
		return diag.Errorf("Only description can be changed")
	}
	if rd.HasChanges("vpc_peering_description") {
		_, err := inst.Client.Peering.UpdateVpcPeeringDescription(ctx, rd.Id(), rd.Get("vpc_peering_description").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVpcPeeringRead(ctx, rd, meta)
}

func resourceVpcPeeringDelete(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	state := rd.Get("vpc_peering_state").(string)
	if state == "REQUESTING" {
		if _, err := inst.Client.Peering.CancelVpcPeering(ctx, rd.Id()); err != nil {
			return diag.FromErr(err)
		}
	}

	err := inst.Client.Peering.DeleteVpcPeering(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitVpcPeeringDeleting(ctx, inst.Client, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitVpcPeeringDeleting(ctx context.Context, scpClient *client.SCPClient, peeringId string) error {
	return client.WaitForStatus(ctx, scpClient, []string{"DELETING"}, []string{"DELETED"}, func() (interface{}, string, error) {
		return scpClient.Peering.GetVpcPeeringForDelete(ctx, peeringId)
	})
}
