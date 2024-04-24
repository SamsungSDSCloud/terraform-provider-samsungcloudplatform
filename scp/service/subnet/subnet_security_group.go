package subnet

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/subnet2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_subnet_security_group", ResourceSubnetSecurityGroup())
}

func ResourceSubnetSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetSecurityGroupAttach,
		ReadContext:   resourceSubnetSecurityGroupRead,
		UpdateContext: resourceSubnetSecurityGroupUpdate,
		DeleteContext: resourceSubnetSecurityGroupDetach,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target Subnet id",
			},
			"vip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "subnet Virtual ip id. (Reserved Virtual ip id)",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "securityGroup Id",
			},
		},
		Description: "Provides a Subnet Vip reserve resource.",
	}
}

func resourceSubnetSecurityGroupAttach(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get values from schema
	subnetId := rd.Get("subnet_id").(string)
	vipId := rd.Get("vip_id").(string)
	securityGroupId := rd.Get("security_group_id").(string)

	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Subnet.GetSubnet(ctx, subnetId)
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = inst.Client.Subnet.GetSubnetVip(ctx, subnetId, vipId)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := inst.Client.Subnet.AttachSubnetSecurityGroup(ctx, subnetId, vipId, securityGroupId)

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForSubnetStatus(ctx, inst.Client, subnetId, []string{"EDITING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.ResourceId) //vip_id
	err = rd.Set("subnet_id", subnetId)
	if err != nil {
		return nil
	}
	err = rd.Set("vip_id", vipId)
	if err != nil {
		return nil
	}
	err = rd.Set("security_group_id", securityGroupId)
	if err != nil {
		return nil
	}

	return resourceSubnetSecurityGroupRead(ctx, rd, meta)
}

func resourceSubnetSecurityGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	subnetId := rd.Get("subnet_id").(string)
	vipId := rd.Get("vip_id").(string)
	securityGroupId := rd.Get("security_group_id").(string)

	// vip 목록 조회 > 해당 ip_id 기준 으로 존재 여부 확인 (vip 상세 조회 시, subnetId와 vipId가 필요)
	requestParam := &subnet2.SubnetVipOpenApiControllerApiListSubnetVipsV2Opts{}

	subnetVips, _, err := inst.Client.Subnet.GetSubnetVipV2List(ctx, subnetId, requestParam)

	if err != nil {
		return diag.FromErr(err)
	}

	setSubnetVips, vipIds := convertSubnetVipIdsToHclSet(subnetVips.Contents)

	if len(setSubnetVips) == 0 {
		return diag.Errorf("Reserved Subnet Vip not found")
	}

	if !existVipId(vipIds, vipId) {
		return diag.Errorf("Input vip id is invalid (maybe not Reserved) : " + vipId)
	}

	// vip 상세 조회
	vipInfo, _, err := inst.Client.Subnet.GetSubnetVip(ctx, subnetId, vipId)

	if len(vipInfo.VipId) == 0 {
		return diag.Errorf("no matching Subnet Vip found, maybe vip Reservation Failed")
	}

	//입력받은 securityGroupId가 존재하는 경우
	if len(securityGroupId) > 0 {
		if !existSecurityGroupId(vipInfo.SecurityGroupIds, securityGroupId) {
			return diag.Errorf("Subnet SecurityGroupId Attached Failed")
		}
	}

	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	return nil
}

func resourceSubnetSecurityGroupUpdate(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return diag.Errorf("Update function is not supported!")
}

func resourceSubnetSecurityGroupDetach(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get values from schema
	subnetId := rd.Get("subnet_id").(string)
	vipId := rd.Get("vip_id").(string)
	securityGroupId := rd.Get("security_group_id").(string)
	inst := meta.(*client.Instance)

	_, err := inst.Client.Subnet.DetachSubnetSecurityGroup(ctx, subnetId, vipId, securityGroupId)

	if err != nil {
		return diag.FromErr(err)
	}

	resourceSubnetSecurityGroupDetachRead(ctx, rd, meta)

	err = waitForSubnetStatus(ctx, inst.Client, subnetId, []string{"EDITING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSubnetSecurityGroupDetachRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	subnetId := rd.Get("subnet_id").(string)
	vipId := rd.Get("vip_id").(string)
	securityGroupId := rd.Get("security_group_id").(string)

	// vip 목록 조회(VipState = RESERVED 인 대상만) > 해당 ip_id 기준 으로 존재 여부 확인 (vip 상세 조회 시, subnetId와 vipId가 필요)
	requestParam := &subnet2.SubnetVipOpenApiControllerApiListSubnetVipsV2Opts{}

	subnetVips, _, err := inst.Client.Subnet.GetSubnetVipV2List(ctx, subnetId, requestParam)

	if err != nil {
		return diag.FromErr(err)
	}

	setSubnetVips, ids := convertSubnetVipListToHclSet(subnetVips.Contents)

	if len(setSubnetVips) == 0 {
		return diag.Errorf("Reserved Subnet Vip not found")
	}

	if !existIpId(ids, vipId) {
		return diag.Errorf("Input vip id is invalid (maybe not Reserved) : " + vipId)
	}

	setVipId := getVipId(subnetVips.Contents, vipId)

	// vip 상세 조회
	vipInfo, _, err := inst.Client.Subnet.GetSubnetVip(ctx, subnetId, setVipId)

	if len(vipInfo.VipId) == 0 {
		return diag.Errorf("no matching Subnet Vip found, maybe vip Reservation Failed")
	}

	if existSecurityGroupId(vipInfo.SecurityGroupIds, securityGroupId) {
		return diag.Errorf("Subnet SecurityGroupId Detached Failed")
	}

	if err != nil {
		rd.SetId("")
		if common.IsDeleted(err) {
			return nil
		}

		return diag.FromErr(err)
	}

	return nil
}

func existSecurityGroupId(securityGroupIds []subnet2.SecurityGroupIdsResponse, securityGroupId string) bool {

	isExist := false //default

	for i := range securityGroupIds {
		if securityGroupIds[i].SecurityGroupId == securityGroupId {
			isExist = true
		}
	}

	return isExist
}
