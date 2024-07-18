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
	uuid "github.com/satori/go.uuid"
	"log"
)

func init() {
	scp.RegisterDataSource("scp_subnet_vip_detail", DatasourceSubnetVip())
	scp.RegisterDataSource("scp_subnet_vips", DatasourceSubnetVips())
	scp.RegisterDataSource("scp_subnet_available_ips", DatasourceSubnetAvailableIps())
}

func DatasourceSubnetVip() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetVip,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id":      {Type: schema.TypeString, Required: true, Description: "Subnet id"},
			"vip_id":         {Type: schema.TypeString, Required: true, Description: "Subnet Virtual Ip id"},
			"nat_ip_address": {Type: schema.TypeString, Computed: true, Description: "Nat Ip address"},
			"nat_ip_id":      {Type: schema.TypeString, Computed: true, Description: "Nat Ip id"},
			"security_group_ids": {Type: schema.TypeList, Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id":           {Type: schema.TypeString, Computed: true, Description: "Security Group Id"},
						"security_group_member_state": {Type: schema.TypeString, Computed: true, Description: "Security Group Member State"},
					},
					Description: "Security-Group ids of this subnet virtual ip",
				},
			},
			"service_zone_id":   {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			"subnet_ip_address": {Type: schema.TypeString, Computed: true, Description: "Subnet Ip address"},
			"subnet_ip_id":      {Type: schema.TypeString, Computed: true, Description: "Subnet Ip id"},
			"vip_state":         {Type: schema.TypeString, Computed: true, Description: "Subnet Virtual Ip state"},
			"vip_description":   {Type: schema.TypeString, Computed: true, Description: "Description of Ip"},
			"created_by":        {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":        {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":       {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":       {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			"project_id":        {Type: schema.TypeString, Computed: true, Description: "Project id"},
		},
		Description: "Provides detail of subnet Vip",
	}
}

func dataSourceSubnetVip(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	subnetId := rd.Get("subnet_id").(string)
	vipId := rd.Get("vip_id").(string)

	vipInfo, _, err := inst.Client.Subnet.GetSubnetVip(ctx, subnetId, vipId)

	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	securityGroupIds := common.ConvertStructToMaps(vipInfo.SecurityGroupIds)

	rd.SetId(uuid.NewV4().String())
	rd.Set("nat_ip_address", vipInfo.NatIpAddress)
	rd.Set("nat_ip_id", vipInfo.NatIpId)
	rd.Set("security_group_ids", securityGroupIds)
	rd.Set("service_zone_id", vipInfo.ServiceZoneId)
	rd.Set("subnet_ip_address", vipInfo.SubnetIpAddress)
	rd.Set("subnet_ip_id", vipInfo.SubnetIpId)
	rd.Set("vip_state", vipInfo.VipState)
	rd.Set("vip_id", vipInfo.VipId)
	rd.Set("vip_description", vipInfo.VipDescription)
	rd.Set("vip_state", vipInfo.VipState)
	rd.Set("created_by", vipInfo.CreatedBy)
	rd.Set("created_dt", vipInfo.CreatedDt)
	rd.Set("modified_by", vipInfo.ModifiedBy)
	rd.Set("modified_dt", vipInfo.ModifiedDt)
	rd.Set("project_id", vipInfo.ProjectId)

	return nil
}

func DatasourceSubnetVips() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetVipList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id":         {Type: schema.TypeString, Required: true, Description: "Subnet id"},
			"subnet_ip_address": {Type: schema.TypeString, Optional: true, Description: "Subnet Virtual Ip address"},
			"vip_state":         {Type: schema.TypeString, Optional: true, Description: "Subnet Virtual Ip State"},
			"page":              {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":              {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":          {Type: schema.TypeList, Optional: true, Description: "Subnet resource list size", Elem: datasourceSubnetVipsElem()},
			"total_count":       {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of subnet Vip",
	}
}

func dataSourceSubnetVipList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	subnetId := rd.Get(common.ToSnakeCase("SubnetId")).(string)

	requestParam := &subnet2.SubnetVipOpenApiControllerApiListSubnetVipsV2Opts{
		SubnetIpAddress: optional.NewString(rd.Get("subnet_ip_address").(string)),
		VipState:        optional.NewString(rd.Get("vip_state").(string)),
		Page:            optional.NewInt32((int32)(rd.Get("page").(int))),
		Size:            optional.NewInt32((int32)(rd.Get("size").(int))),
		Sort:            optional.Interface{},
	}
	responses, _, err := inst.Client.Subnet.GetSubnetVipV2List(ctx, subnetId, requestParam)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	log.Println("=================================================")
	log.Println("contents: ", contents)
	log.Println("=================================================")

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceSubnetVipsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nat_ip_address": {Type: schema.TypeString, Computed: true, Description: "Nat Ip address"},
			"nat_ip_id":      {Type: schema.TypeString, Computed: true, Description: "Nat Ip id"},
			"security_group_ids": {Type: schema.TypeList, Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id":           {Type: schema.TypeString, Computed: true, Description: "Security Group Id"},
						"security_group_member_state": {Type: schema.TypeString, Computed: true, Description: "Security Group Member State"},
					},
					Description: "Security-Group ids of this subnet virtual ip",
				},
			},
			"service_zone_id":   {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			"subnet_ip_address": {Type: schema.TypeString, Computed: true, Description: "Subnet Ip address"},
			"subnet_ip_id":      {Type: schema.TypeString, Computed: true, Description: "Subnet Ip id"},
			"vip_state":         {Type: schema.TypeString, Computed: true, Description: "Subnet Virtual Ip state"},
			"vip_id":            {Type: schema.TypeString, Computed: true, Description: "Subnet Virtual Ip id"},
			"vip_description":   {Type: schema.TypeString, Computed: true, Description: "Description of Ip"},
			"created_by":        {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":        {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":       {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":       {Type: schema.TypeString, Computed: true, Description: "Modification date"},
			"project_id":        {Type: schema.TypeString, Computed: true, Description: "Project id"},
		},
	}
}

func DatasourceSubnetAvailableIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetAvailableIpList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"subnet_id":         {Type: schema.TypeString, Required: true, Description: "Subnet id"},
			"subnet_ip_address": {Type: schema.TypeString, Optional: true, Description: "Subnet Virtual Ip address"},
			"page":              {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":              {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":          {Type: schema.TypeList, Optional: true, Description: "Subnet resource list size", Elem: datasourceSubnetAvailableIpsElem()},
			"total_count":       {Type: schema.TypeInt, Computed: true},
		},
		Description: "Available ip addresses list in subnet",
	}
}

func dataSourceSubnetAvailableIpList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	SubnetId := rd.Get(common.ToSnakeCase("SubnetId")).(string)

	option := subnet2.SubnetVipOpenApiControllerApiListAvailableVipsV2Opts{
		SubnetIpAddress: optional.String{},
		Page:            optional.Int32{},
		Size:            optional.NewInt32(10000),
	}

	responses, err := inst.Client.Subnet.GetSubnetAvailableVipV2List(ctx, SubnetId, &option)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	log.Println("=================================================")
	log.Println("contents: ", contents)
	log.Println("=================================================")

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceSubnetAvailableIpsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"subnet_ip_address": {Type: schema.TypeString, Computed: true, Description: "Subnet Ip address"},
			"ip_id":             {Type: schema.TypeString, Computed: true, Description: "Ip id"},
			"vip_state":         {Type: schema.TypeString, Computed: true, Description: "Subnet Virtual Ip state"},
			"vip_description":   {Type: schema.TypeString, Computed: true, Description: "Description of Ip"},
			"created_by":        {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":        {Type: schema.TypeString, Computed: true, Description: "Creation date"},
			"modified_by":       {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":       {Type: schema.TypeString, Computed: true, Description: "Modification date"},
		},
	}
}
