package publicip

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"
	publicip2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/public-ip2"
	"github.com/antihax/optional"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_public_ips", DatasourcePublicIps())
}

func DatasourcePublicIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicIpList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ServiceZoneId"):   {Type: schema.TypeString, Required: true, Description: "Service zone id"},
			common.ToSnakeCase("IpAddress"):       {Type: schema.TypeString, Optional: true, Description: "Ip address"},
			common.ToSnakeCase("IsBillable"):      {Type: schema.TypeBool, Optional: true, Description: "Enable bill"},
			common.ToSnakeCase("IsViewable"):      {Type: schema.TypeBool, Optional: true, Description: "Enable view"},
			common.ToSnakeCase("PublicIpPurpose"): {Type: schema.TypeString, Optional: true, Description: "The reason to make public ip"},
			common.ToSnakeCase("PublicIpState"):   {Type: schema.TypeString, Optional: true, Description: "Public ip status"},
			common.ToSnakeCase("UplinkType"):      {Type: schema.TypeString, Optional: true, Description: "Uplink type"},
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

	serviceZoneId := rd.Get("service_zone_id").(string)

	inst := meta.(*client.Instance)

	publicIpList, err := inst.Client.PublicIp.GetPublicIpList(ctx, serviceZoneId, &publicip2.PublicIpOpenApiControllerApiListPublicIpsV2Opts{
		ServiceZoneId:   optional.String{},
		IpAddress:       optional.String{},
		IsBillable:      optional.NewBool(true),
		IsViewable:      optional.NewBool(true),
		PublicIpPurpose: optional.NewString(common.VpcPublicIpPurpose),
		PublicIpState:   optional.String{},
		UplinkType:      optional.NewString(common.VpcPublicIpUplinkType),
		CreatedBy:       optional.String{},
		Page:            optional.NewInt32(0),
		Size:            optional.NewInt32(10000),
		Sort:            optional.NewInterface([]string{"createdDt:desc"}),
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
