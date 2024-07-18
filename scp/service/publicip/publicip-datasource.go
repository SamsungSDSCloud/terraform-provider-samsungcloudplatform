package publicip

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	publicip2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/public-ip2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_public_ip", DatasourceVpcPublicIp())
	scp.RegisterDataSource("scp_public_ips", DatasourcePublicIps())
}
func DatasourceVpcPublicIp() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceVpcPublicIpRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ServiceZoneId"):              {Type: schema.TypeString, Optional: true, Description: "Service zone id"},
			common.ToSnakeCase("IpAddress"):                  {Type: schema.TypeString, Optional: true, Description: "Ip address"},
			common.ToSnakeCase("IpAddressId"):                {Type: schema.TypeString, Optional: true, Description: "Ip address Id"},
			common.ToSnakeCase("publicIpId"):                 {Type: schema.TypeString, Required: true, Description: "Public ip Id"},
			common.ToSnakeCase("PublicIpPurpose"):            {Type: schema.TypeString, Optional: true, Description: "The reason to make public ip"},
			common.ToSnakeCase("PublicIpAddressDescription"): {Type: schema.TypeString, Optional: true, Description: "Public ip Description"},
			common.ToSnakeCase("PublicIpState"):              {Type: schema.TypeString, Optional: true, Description: "Public ip status"},
			common.ToSnakeCase("UplinkType"):                 {Type: schema.TypeString, Optional: true, Description: "Uplink type"},
			common.ToSnakeCase("AttachedObjectName"):         {Type: schema.TypeString, Optional: true, Description: "Name of object with public ip"},
			common.ToSnakeCase("CreatedBy"):                  {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
		},
		Description: "Provide Detail of public ip",
	}
}

func datasourceVpcPublicIpRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	publicIpId := rd.Get(common.ToSnakeCase("publicIpId")).(string)
	info, _, err := inst.Client.PublicIp.GetPublicIp(ctx, publicIpId)

	if err != nil {
		rd.SetId("")
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set(common.ToSnakeCase("ServiceZoneId"), info.ServiceZoneId)
	rd.Set(common.ToSnakeCase("IpAddress"), info.IpAddress)
	rd.Set(common.ToSnakeCase("IpAddressId"), info.IpAddressId)
	rd.Set(common.ToSnakeCase("publicIpId"), info.PublicIpAddressId)
	rd.Set(common.ToSnakeCase("PublicIpAddressDescription"), info.PublicIpAddressDescription)
	rd.Set(common.ToSnakeCase("uplinkType"), info.UplinkType)
	rd.Set(common.ToSnakeCase("PublicIpState"), info.PublicIpState)
	rd.Set(common.ToSnakeCase("PublicIpPurpose"), info.PublicIpPurpose)
	rd.Set(common.ToSnakeCase("AttachedObjectName"), info.AttachedObjectName)
	rd.Set(common.ToSnakeCase("CreatedBy"), info.CreatedBy)

	return nil
}

func DatasourcePublicIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicIpList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ServiceZoneId"):   {Type: schema.TypeString, Optional: true, Description: "Service zone id"},
			common.ToSnakeCase("IpAddress"):       {Type: schema.TypeString, Optional: true, Description: "Ip address"},
			common.ToSnakeCase("IsBillable"):      {Type: schema.TypeBool, Optional: true, Description: "Enable bill"},
			common.ToSnakeCase("IsViewable"):      {Type: schema.TypeBool, Optional: true, Description: "Enable view"},
			common.ToSnakeCase("PublicIpPurpose"): {Type: schema.TypeString, Optional: true, Description: "The reason to make public ip"},
			common.ToSnakeCase("PublicIpState"):   {Type: schema.TypeString, Optional: true, Description: "Public ip status"},
			common.ToSnakeCase("UplinkType"):      {Type: schema.TypeString, Optional: true, Description: "Uplink type"},
			common.ToSnakeCase("VpcId"):           {Type: schema.TypeString, Optional: true, Description: "VPC id"},
			common.ToSnakeCase("CreatedBy"):       {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			common.ToSnakeCase("Page"):            {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			common.ToSnakeCase("Size"):            {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":                            {Type: schema.TypeList, Optional: true, Description: "Public ip list size", Elem: datasourcePublicIpElem()},
			"total_count":                         {Type: schema.TypeInt, Computed: true},
		},
		Description: "Provides list of public ips",
	}
}

func dataSourcePublicIpList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	publicIpList, err := inst.Client.PublicIp.GetPublicIps(ctx, &publicip2.PublicIpOpenApiV3ControllerApiListPublicIpsV3Opts{
		ServiceZoneId: common.GetKeyString(rd, common.ToSnakeCase("ServiceZoneId")),
		IpAddress:     common.GetKeyString(rd, common.ToSnakeCase("IpAddress")),
		PublicIpState: common.GetKeyString(rd, common.ToSnakeCase("PublicIpState")),
		VpcId:         common.GetKeyString(rd, common.ToSnakeCase("VpcId")),
		UplinkType:    common.GetKeyString(rd, common.ToSnakeCase("UplinkType")),
		CreatedBy:     common.GetKeyString(rd, common.ToSnakeCase("CreatedBy")),
		Page:          optional.NewInt32(0),
		Size:          optional.NewInt32(10000),
		Sort:          optional.NewInterface([]string{"createdDt:desc"}),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(publicIpList.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", publicIpList.TotalCount)

	return nil
}

func datasourcePublicIpElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ProjectId"):                  {Type: schema.TypeString, Computed: true, Description: "Project id"},
			common.ToSnakeCase("AttachedObjectName"):         {Type: schema.TypeString, Computed: true, Description: "Name of object with public ip"},
			common.ToSnakeCase("IpAddress"):                  {Type: schema.TypeString, Computed: true, Description: "Ip address"},
			common.ToSnakeCase("IpAddressId"):                {Type: schema.TypeString, Computed: true, Description: "Id of ip address"},
			common.ToSnakeCase("ProductGroupId"):             {Type: schema.TypeString, Computed: true, Description: "Product group id"},
			common.ToSnakeCase("PublicIpAddressId"):          {Type: schema.TypeString, Computed: true, Description: "Id of Public ip address"},
			common.ToSnakeCase("PublicIpPurpose"):            {Type: schema.TypeString, Computed: true, Description: "Purpose of public ip (NAT)"},
			common.ToSnakeCase("PublicIpState"):              {Type: schema.TypeString, Computed: true, Description: "Public ip status"},
			common.ToSnakeCase("ServiceZoneId"):              {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			common.ToSnakeCase("UplinkType"):                 {Type: schema.TypeString, Computed: true, Description: "Uplink type (INTERNET)"},
			common.ToSnakeCase("PublicIpAddressDescription"): {Type: schema.TypeString, Computed: true, Description: "Description of public ip address "},
			common.ToSnakeCase("CreatedBy"):                  {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			common.ToSnakeCase("CreatedDt"):                  {Type: schema.TypeString, Computed: true, Description: "Creation time"},
			common.ToSnakeCase("ModifiedBy"):                 {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			common.ToSnakeCase("ModifiedDt"):                 {Type: schema.TypeString, Computed: true, Description: "Modification time"},
		},
	}
}
