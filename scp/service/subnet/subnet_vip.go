package subnet

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/subnet2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	scp.RegisterResource("scp_subnet_vip", ResourceSubnetVip())
}

func ResourceSubnetVip() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetVipReserve,
		ReadContext:   resourceSubnetVipRead,
		UpdateContext: resourceSubnetVipUpdate,
		DeleteContext: resourceSubnetVipRelease,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target Subnet id",
			},
			"subnet_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "subnet ip id. (Available ip id)",
			},
			"vip_description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Subnet vip description. (Up to 50 characters)",
				ValidateFunc: validation.StringLenBetween(0, 50),
			},
		},
		Description: "Provides a Subnet Vip reserve resource.",
	}
}

func resourceSubnetVipReserve(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get values from schema
	subnetId := rd.Get("subnet_id").(string)
	subnetIpId := rd.Get("subnet_ip_id").(string)
	vipDescription := rd.Get("vip_description").(string)

	inst := meta.(*client.Instance)

	_, _, err := inst.Client.Subnet.GetSubnet(ctx, subnetId)
	if err != nil {
		return diag.FromErr(err)
	}

	option := subnet2.SubnetVipOpenApiControllerApiListAvailableVipsV2Opts{}

	subnetAvailableVips, err := inst.Client.Subnet.GetSubnetAvailableVipV2List(ctx, subnetId, &option)
	if err != nil {
		return diag.FromErr(err)
	}

	setSubnetAvailableVips, ids := convertSubnetAvailableVipListToHclSet(subnetAvailableVips.Contents)

	if len(setSubnetAvailableVips) == 0 {
		return diag.Errorf("Available subnet ip not found")
	}

	if !existIpId(ids, subnetIpId) {
		return diag.Errorf("Input subnet id is invalid (maybe not Available state) : " + subnetIpId)
	}

	result, err := inst.Client.Subnet.ReserveSubnetVipsV2(ctx, subnetId, subnetIpId, vipDescription)

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForSubnetStatus(ctx, inst.Client, subnetId, []string{"EDITING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(result.ResourceId) //ip_id
	rd.Set("subnet_id", subnetId)
	rd.Set("subnet_ip_id", result.ResourceId)

	return resourceSubnetVipRead(ctx, rd, meta)
}

func resourceSubnetVipRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	subnetId := rd.Get("subnet_id").(string)
	subnetIpId := rd.Get("subnet_ip_id").(string)

	// vip 목록조회 > 해당 ip_id 기준으로 존재여부 확인 (vip 상세조회 시, subnetId와 vipId가 필요)
	requestParam := &subnet2.SubnetVipOpenApiControllerApiListSubnetVipsV2Opts{}

	subnetVips, _, err := inst.Client.Subnet.GetSubnetVipV2List(ctx, subnetId, requestParam)

	if err != nil {
		return diag.FromErr(err)
	}

	setSubnetVips, ids := convertSubnetVipListToHclSet(subnetVips.Contents)

	if len(setSubnetVips) == 0 {
		return diag.Errorf("Reserved Subnet Vip not found")
	}

	if !existIpId(ids, subnetIpId) {
		return diag.Errorf("Input subnet ip id is invalid (maybe not Reserved) : " + subnetIpId)
	}

	setVipId := getVipId(subnetVips.Contents, subnetIpId)

	// vip 상세조회
	vipInfo, _, err := inst.Client.Subnet.GetSubnetVip(ctx, subnetId, setVipId)

	if len(vipInfo.VipId) == 0 {
		return diag.Errorf("no matching Subnet Vip found, maybe vip Reservation Failed")
	}

	if vipInfo.VipState != "RESERVED" {
		return diag.Errorf("Subnet Vip Reservation Failed")
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

func resourceSubnetVipUpdate(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("Update function is not supported!")
}

func resourceSubnetVipRelease(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get values from schema
	subnetId := rd.Get("subnet_id").(string)
	subnetIpId := rd.Get("subnet_ip_id").(string)

	inst := meta.(*client.Instance)

	_, err := inst.Client.Subnet.ReleaseSubnetVipsV2(ctx, subnetId, subnetIpId)

	if err != nil {
		return diag.FromErr(err)
	}

	resourceSubnetVipReleaseRead(ctx, rd, meta)

	err = waitForSubnetStatus(ctx, inst.Client, subnetId, []string{"EDITING"}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSubnetVipReleaseRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	subnetId := rd.Get("subnet_id").(string)
	subnetIpId := rd.Get("subnet_ip_id").(string)

	// vip 목록조회 > 해당 ip_id 기준으로 존재여부 확인 (vip 상세조회 시, subnetId와 vipId가 필요)
	requestParam := &subnet2.SubnetVipOpenApiControllerApiListSubnetVipsV2Opts{}

	subnetVips, _, err := inst.Client.Subnet.GetSubnetVipV2List(ctx, subnetId, requestParam)

	if err != nil {
		return diag.FromErr(err)
	}

	//목록조회 시, 0건이면 모두 반납 / 0건이 아닌경우는 반납한 ip_id 존재확인
	if subnetVips.TotalCount != 0 {
		_, ids := convertSubnetVipListToHclSet(subnetVips.Contents)

		if existIpId(ids, subnetIpId) {
			return diag.Errorf("Input subnet ip id is invalid (maybe not Release) : " + subnetIpId)
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

func convertSubnetAvailableVipListToHclSet(subnetAvaVips []subnet2.SubnetVirtualIpAvailableListItemResVo) (common.HclSetObject, []string) {
	var subnetAvailableVipList common.HclSetObject
	var ids []string
	for _, sb := range subnetAvaVips {
		if len(subnetAvaVips) == 0 {
			continue
		}
		ids = append(ids, sb.IpId)
		kv := common.HclKeyValueObject{
			"subnet_ip_id": sb.IpId,
			"subnet_ip":    sb.SubnetIpAddress,
			"vip_state":    sb.VipState,
			"vip_desc":     sb.VipDescription,
		}
		subnetAvailableVipList = append(subnetAvailableVipList, kv)
	}
	return subnetAvailableVipList, ids
}

func convertSubnetVipListToHclSet(subnetVips []subnet2.SubnetVirtualIpListItemResVo) (common.HclSetObject, []string) {
	var subnetVipList common.HclSetObject
	var ids []string
	for _, sb := range subnetVips {
		if len(subnetVips) == 0 {
			continue
		}
		ids = append(ids, sb.SubnetIpId)
		kv := common.HclKeyValueObject{
			"subnet_ip_id": sb.SubnetIpId,
			"vip_id":       sb.VipId,
		}
		subnetVipList = append(subnetVipList, kv)
	}
	return subnetVipList, ids
}

func getVipId(subnetVips []subnet2.SubnetVirtualIpListItemResVo, subnetIpId string) string {
	var vipId string

	for i := range subnetVips {
		if subnetVips[i].SubnetIpId == subnetIpId {
			vipId = subnetVips[i].VipId
			continue
		}
	}

	return vipId
}

func existIpId(ids []string, subnetIpId string) bool {

	isExist := false //default

	for _, ipId := range ids {
		if ipId == subnetIpId {
			isExist = true
		}
	}

	return isExist
}
