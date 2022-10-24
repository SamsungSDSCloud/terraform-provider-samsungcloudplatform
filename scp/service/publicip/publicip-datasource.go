package publicip

import (
	"context"

	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client/publicip"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

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
	inst := meta.(*client.Instance)

	requestParam := publicip.ListPublicIpRequest{
		ServiceZoneId:   rd.Get(common.ToSnakeCase("ServiceZoneId")).(string),
		IpAddress:       rd.Get(common.ToSnakeCase("IpAddress")).(string),
		IsBillable:      rd.Get(common.ToSnakeCase("IsBillable")).(bool),
		IsViewable:      rd.Get(common.ToSnakeCase("IsViewable")).(bool),
		PublicIpPurpose: rd.Get(common.ToSnakeCase("PublicIpPurpose")).(string),
		PublicIpState:   rd.Get(common.ToSnakeCase("PublicIpState")).(string),
		UplinkType:      rd.Get(common.ToSnakeCase("UplinkType")).(string),
		CreatedBy:       rd.Get(common.ToSnakeCase("CreatedBy")).(string),
		Page:            (int32)(rd.Get(common.ToSnakeCase("Page")).(int)),
		Size:            (int32)(rd.Get(common.ToSnakeCase("Size")).(int)),
	}

	responses, err := inst.Client.PublicIp.GetPublicIpListV2(ctx, requestParam)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourcePublicIpElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ProjectId"):                  {Type: schema.TypeString, Computed: true, Description: "Project id"},
			common.ToSnakeCase("AttachedObjectName"):         {Type: schema.TypeString, Computed: true, Description: "Name of object with public ip"},
			common.ToSnakeCase("BlockId"):                    {Type: schema.TypeString, Computed: true, Description: "Block id of this region"},
			common.ToSnakeCase("IpAddress"):                  {Type: schema.TypeString, Computed: true, Description: "Ip address"},
			common.ToSnakeCase("IpAddressId"):                {Type: schema.TypeString, Computed: true, Description: "Id of ip address"},
			common.ToSnakeCase("NetworkServiceType"):         {Type: schema.TypeString, Computed: true, Description: "Network service type"},
			common.ToSnakeCase("ProductGroupId"):             {Type: schema.TypeString, Computed: true, Description: "Product group id"},
			common.ToSnakeCase("ProjectName"):                {Type: schema.TypeString, Computed: true, Description: "Project name"},
			common.ToSnakeCase("PublicIpAddressId"):          {Type: schema.TypeString, Computed: true, Description: "Id of Public ip address"},
			common.ToSnakeCase("PublicIpPurpose"):            {Type: schema.TypeString, Computed: true, Description: "Purpose of public ip (NAT)"},
			common.ToSnakeCase("PublicIpState"):              {Type: schema.TypeString, Computed: true, Description: "Public ip status"},
			common.ToSnakeCase("Region"):                     {Type: schema.TypeString, Computed: true, Description: "The region name to create"},
			common.ToSnakeCase("ServiceZoneId"):              {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			common.ToSnakeCase("UplinkType"):                 {Type: schema.TypeString, Computed: true, Description: "Uplink type (INTERNET)"},
			common.ToSnakeCase("UserName"):                   {Type: schema.TypeString, Computed: true, Description: "User name"},
			common.ToSnakeCase("ZoneName"):                   {Type: schema.TypeString, Computed: true, Description: "Service zone name"},
			common.ToSnakeCase("PublicIpAddressDescription"): {Type: schema.TypeString, Computed: true, Description: "Description of public ip address "},
			common.ToSnakeCase("CreatedBy"):                  {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			common.ToSnakeCase("CreatedDt"):                  {Type: schema.TypeString, Computed: true, Description: "Creation time"},
			common.ToSnakeCase("ModifiedBy"):                 {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			common.ToSnakeCase("ModifiedDt"):                 {Type: schema.TypeString, Computed: true, Description: "Modification time"},
		},
	}
}
