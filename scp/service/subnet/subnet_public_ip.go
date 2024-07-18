package subnet

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/subnet2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	scp.RegisterResource("scp_subnet_public_ip", ResourceSubnetPublicIp())
}

func ResourceSubnetPublicIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetPublicIpAttach,
		ReadContext:   resourceSubnetPublicIpRead,
		UpdateContext: resourceSubnetPublicIpUpdate,
		DeleteContext: resourceSubnetPublicIpDetach,
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
			"public_ip_address_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Public Ip Address Id (Reserved Public ip id)",
			},
		},
		Description: "Provides a Subnet Vip reserve resource.",
	}
}

func resourceSubnetPublicIpAttach(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get values from schema
	subnetId := rd.Get("subnet_id").(string)
	vipId := rd.Get("vip_id").(string)
	publicIpAddressId := rd.Get("public_ip_address_id").(string)

	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Subnet.GetSubnet(ctx, subnetId)
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = inst.Client.Subnet.GetSubnetVip(ctx, subnetId, vipId)
	if err != nil {
		return diag.FromErr(err)
	}

	//입력받은 public ip가 존재하는 경우, public ip 할당  / 입력받지 않은 경우 자동할당
	if len(publicIpAddressId) > 0 {
		_, _, err = inst.Client.PublicIp.GetPublicIp(ctx, publicIpAddressId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	result, err := inst.Client.Subnet.AttachSubnetPublicIp(ctx, subnetId, vipId, publicIpAddressId)

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
	err = rd.Set("public_ip_address_id", publicIpAddressId)
	if err != nil {
		return nil
	}

	return resourceSubnetPublicIpRead(ctx, rd, meta)
}

func resourceSubnetPublicIpRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	subnetId := rd.Get("subnet_id").(string)
	vipId := rd.Get("vip_id").(string)
	publicIpAddressId := rd.Get("public_ip_address_id").(string)

	// vip 목록 조회(VipState = ATTACHED 인 대상만) > 해당 ip_id 기준 으로 존재 여부 확인 (vip 상세 조회 시, subnetId와 vipId가 필요)
	requestParam := &subnet2.SubnetVipOpenApiControllerApiListSubnetVipsV2Opts{
		VipState: optional.NewString("ATTACHED"),
	}

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

	if vipInfo.VipState != "ATTACHED" {
		return diag.Errorf("Subnet Public IP Attached Failed")
	}

	//public ip 상세 조회
	if len(publicIpAddressId) > 0 {
		publicIpInfo, _, _ := inst.Client.PublicIp.GetPublicIp(ctx, publicIpAddressId)

		if publicIpInfo.IpAddressId != vipInfo.NatIpId {
			return diag.Errorf("Attached Nat IP Of Subnet Virtual IP And Attaching Public IP is different")
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

func resourceSubnetPublicIpUpdate(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return diag.Errorf("Update function is not supported!")
}

func resourceSubnetPublicIpDetach(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get values from schema
	subnetId := rd.Get("subnet_id").(string)
	vipId := rd.Get("vip_id").(string)

	inst := meta.(*client.Instance)

	_, err := inst.Client.Subnet.DetachSubnetPublicIp(ctx, subnetId, vipId)

	if err != nil {
		return diag.FromErr(err)
	}

	resourceSubnetPublicIpDetachRead(ctx, rd, meta)

	err = waitForSubnetStatus(ctx, inst.Client, subnetId, []string{"EDITING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSubnetPublicIpDetachRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	subnetId := rd.Get("subnet_id").(string)
	vipId := rd.Get("vip_id").(string)
	publicIpAddressId := rd.Get("public_ip_address_id").(string)

	// vip 목록 조회(VipState = RESERVED 인 대상만) > 해당 ip_id 기준 으로 존재 여부 확인 (vip 상세 조회 시, subnetId와 vipId가 필요)
	requestParam := &subnet2.SubnetVipOpenApiControllerApiListSubnetVipsV2Opts{
		VipState: optional.NewString("RESERVED"),
	}

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

	if vipInfo.VipState != "RESERVED" || vipInfo.NatIpId != "" || vipInfo.NatIpAddress != "" {
		return diag.Errorf("Subnet Public IP Detached Failed")
	}

	//입력받은 public ip가 존재하는 경우, public ip 할당  / 입력받지 않은 경우 자동할당
	if len(publicIpAddressId) > 0 {
		publicIpInfo, _, _ := inst.Client.PublicIp.GetPublicIp(ctx, publicIpAddressId)

		if publicIpInfo.PublicIpState == "ATTACHED" {
			return diag.Errorf("Released Public IP Failed")
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
func convertSubnetVipIdsToHclSet(subnetVips []subnet2.SubnetVirtualIpListItemResVo) (common.HclSetObject, []string) {
	var subnetVipList common.HclSetObject
	var vipIds []string
	for _, sb := range subnetVips {
		if len(subnetVips) == 0 {
			continue
		}
		vipIds = append(vipIds, sb.VipId)
		kv := common.HclKeyValueObject{
			"subnet_ip_id": sb.SubnetIpId,
			"vip_id":       sb.VipId,
		}
		subnetVipList = append(subnetVipList, kv)
	}
	return subnetVipList, vipIds
}

func existVipId(vipIds []string, virtualIpId string) bool {

	isExist := false //default

	for _, vipId := range vipIds {
		if vipId == virtualIpId {
			isExist = true
		}
	}

	return isExist
}
