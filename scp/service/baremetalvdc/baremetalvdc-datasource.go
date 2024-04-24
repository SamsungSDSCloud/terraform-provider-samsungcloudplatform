package baremetalvdc

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_bm_vdc_servers", DatasourceBareMetalServers())
	scp.RegisterDataSource("scp_bm_vdc_server", DatasourceBareMetalServer())
}

func DatasourceBareMetalServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: bareMetalServerList, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{ //스키마 정의
			"contents":    {Type: schema.TypeList, Optional: true, Description: "DNS list", Elem: bareMetalServerElem()},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides Baremetal Server VDC List",
	}
}

func bareMetalServerList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.BareMetalVdc.GetBareMetalServersVDC(ctx, "", "")
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func bareMetalServerElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ProjectId"):            {Type: schema.TypeString, Computed: true, Description: "Project Id"},
			common.ToSnakeCase("BareMetalServerId"):    {Type: schema.TypeString, Computed: true, Description: "Baremetal Server Id"},
			common.ToSnakeCase("BareMetalServerName"):  {Type: schema.TypeString, Computed: true, Description: "Baremetal Server Name"},
			common.ToSnakeCase("BareMetalServerState"): {Type: schema.TypeString, Computed: true, Description: "Baremetal Server State"},
			common.ToSnakeCase("BlockId"):              {Type: schema.TypeString, Computed: true, Description: "Block Id"},
			common.ToSnakeCase("ImageId"):              {Type: schema.TypeString, Computed: true, Description: "Image Id"},
			common.ToSnakeCase("IpAddress"):            {Type: schema.TypeString, Computed: true, Description: "Ip Address"},
			common.ToSnakeCase("ServerTypeId"):         {Type: schema.TypeString, Computed: true, Description: "Server Type Id"},
			common.ToSnakeCase("ServiceZoneId"):        {Type: schema.TypeString, Computed: true, Description: "Service Zone Id"},
			common.ToSnakeCase("CreatedBy"):            {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):            {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			common.ToSnakeCase("ModifiedBy"):           {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):           {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
	}
}

func DatasourceBareMetalServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: bareMetalServerDetail, //데이터 조회 함수
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ServerId"):          {Type: schema.TypeString, Required: true, Description: "Server Id"},
			common.ToSnakeCase("ProjectId"):         {Type: schema.TypeString, Computed: true, Description: "Project Id"},
			common.ToSnakeCase("AllMountedStorage"): {Type: schema.TypeBool, Computed: true, Description: "Is All Storage Mounted"},
			common.ToSnakeCase("BareMetalBlockStorageIds"): {Type: schema.TypeList, Computed: true, Description: "Baremetal Block Storage Ids",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			common.ToSnakeCase("BareMetalServerId"):         {Type: schema.TypeString, Computed: true, Description: "Baremetal Server Id"},
			common.ToSnakeCase("BareMetalServerName"):       {Type: schema.TypeString, Computed: true, Description: "Baremetal Server Name"},
			common.ToSnakeCase("BareMetalServerState"):      {Type: schema.TypeString, Computed: true, Description: "Baremetal Server State"},
			common.ToSnakeCase("BlockId"):                   {Type: schema.TypeString, Computed: true, Description: "Block Id"},
			common.ToSnakeCase("CheckCriticalError"):        {Type: schema.TypeBool, Computed: true, Description: "Check Critical Error"},
			common.ToSnakeCase("Contract"):                  {Type: schema.TypeString, Computed: true, Description: "Contract"},
			common.ToSnakeCase("ContractEndDate"):           {Type: schema.TypeString, Computed: true, Description: "Contract End Date"},
			common.ToSnakeCase("ContractId"):                {Type: schema.TypeString, Computed: true, Description: "Contract Id"},
			common.ToSnakeCase("ContractStartDate"):         {Type: schema.TypeString, Computed: true, Description: "Contract State Date"},
			common.ToSnakeCase("DeletionProtectionEnabled"): {Type: schema.TypeString, Computed: true, Description: "Delete Protection Enabled"},
			common.ToSnakeCase("DnsEnabled"):                {Type: schema.TypeString, Computed: true, Description: "Dns Enabled"},
			common.ToSnakeCase("ErrorCheck"):                {Type: schema.TypeBool, Computed: true, Description: "Error Check"},
			common.ToSnakeCase("ImageId"):                   {Type: schema.TypeString, Computed: true, Description: "Image Id"},
			common.ToSnakeCase("InitialScriptContent"):      {Type: schema.TypeString, Computed: true, Description: "Initial Script"},
			common.ToSnakeCase("IpAddress"):                 {Type: schema.TypeString, Computed: true, Description: "Ip Address"},
			common.ToSnakeCase("LocalSubnetStatus"):         {Type: schema.TypeString, Computed: true, Description: "Local Subnet State"},
			common.ToSnakeCase("NextContract"):              {Type: schema.TypeString, Computed: true, Description: "Next Contract"},
			common.ToSnakeCase("NextContractEndDate"):       {Type: schema.TypeString, Computed: true, Description: "Next Contract End Date"},
			common.ToSnakeCase("OsType"):                    {Type: schema.TypeString, Computed: true, Description: "Os Type"},
			common.ToSnakeCase("ProductGroupId"):            {Type: schema.TypeString, Computed: true, Description: "ProductGroup Id"},
			common.ToSnakeCase("ProductType"):               {Type: schema.TypeString, Computed: true, Description: "Product Type"},
			common.ToSnakeCase("ServerType"):                {Type: schema.TypeString, Computed: true, Description: "Server Type"},
			common.ToSnakeCase("ServerTypeId"):              {Type: schema.TypeString, Computed: true, Description: "Server Type Id"},
			common.ToSnakeCase("ServiceLevelId"):            {Type: schema.TypeString, Computed: true, Description: "Service Level Id"},
			common.ToSnakeCase("ServiceZoneId"):             {Type: schema.TypeString, Computed: true, Description: "Service Zone Id"},
			common.ToSnakeCase("SubnetId"):                  {Type: schema.TypeString, Computed: true, Description: "Subnet Id"},
			common.ToSnakeCase("Timezone"):                  {Type: schema.TypeString, Computed: true, Description: "TimeZone"},
			common.ToSnakeCase("UseHyperThreading"):         {Type: schema.TypeString, Computed: true, Description: "Use HyperThreading"},
			common.ToSnakeCase("VdcId"):                     {Type: schema.TypeString, Computed: true, Description: "Vdc Id"},
			common.ToSnakeCase("CreatedBy"):                 {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):                 {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			common.ToSnakeCase("ModifiedBy"):                {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):                {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
		Description: "Provides Baremetal Server VDC Detail",
	}
}

func bareMetalServerDetail(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	serverId := rd.Get("server_id").(string)

	response, _, err := inst.Client.BareMetalVdc.GetBareMetalServerDetailVDC(ctx, serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())

	rd.Set(common.ToSnakeCase("AllMountedStorage"), response.AllMountedStorage)
	rd.Set(common.ToSnakeCase("BareMetalBlockStorageIds"), response.BareMetalBlockStorageIds)
	rd.Set(common.ToSnakeCase("BareMetalServerId"), response.BareMetalServerId)
	rd.Set(common.ToSnakeCase("BareMetalServerName"), response.BareMetalServerName)
	rd.Set(common.ToSnakeCase("BareMetalServerState"), response.BareMetalServerState)
	rd.Set(common.ToSnakeCase("BlockId"), response.BlockId)
	rd.Set(common.ToSnakeCase("CheckCriticalError"), response.CheckCriticalError)
	rd.Set(common.ToSnakeCase("Contract"), response.Contract)
	rd.Set(common.ToSnakeCase("ContractEndDate"), response.ContractEndDate)
	rd.Set(common.ToSnakeCase("ContractId"), response.ContractId)
	rd.Set(common.ToSnakeCase("ContractStartDate"), response.ContractStartDate)
	rd.Set(common.ToSnakeCase("DeletionProtectionEnabled"), response.DeletionProtectionEnabled)
	rd.Set(common.ToSnakeCase("DnsEnabled"), response.DnsEnabled)
	rd.Set(common.ToSnakeCase("ErrorCheck"), response.ErrorCheck)
	rd.Set(common.ToSnakeCase("ImageId"), response.ImageId)
	rd.Set(common.ToSnakeCase("InitialScriptContent"), response.InitialScriptContent)
	rd.Set(common.ToSnakeCase("IpAddress"), response.IpAddress)
	rd.Set(common.ToSnakeCase("NextContract"), response.NextContract)
	rd.Set(common.ToSnakeCase("NextContractEndDate"), response.NextContractEndDate)
	rd.Set(common.ToSnakeCase("OsType"), response.OsType)
	rd.Set(common.ToSnakeCase("ProductGroupId"), response.ProductGroupId)
	rd.Set(common.ToSnakeCase("ProductType"), response.ProductType)
	rd.Set(common.ToSnakeCase("ServerTypeId"), response.ServerTypeId)
	rd.Set(common.ToSnakeCase("ServiceLevelId"), response.ServiceLevelId)
	rd.Set(common.ToSnakeCase("ServiceZoneId"), response.ServiceZoneId)
	rd.Set(common.ToSnakeCase("SubnetId"), response.SubnetId)
	rd.Set(common.ToSnakeCase("Timezone"), response.Timezone)
	rd.Set(common.ToSnakeCase("UseHyperThreading"), response.UseHyperThreading)
	rd.Set(common.ToSnakeCase("VdcId"), response.VdcId)
	rd.Set(common.ToSnakeCase("CreatedBy"), response.CreatedBy)
	rd.Set(common.ToSnakeCase("CreatedDt"), response.CreatedDt)
	rd.Set(common.ToSnakeCase("ModifiedBy"), response.ModifiedBy)
	rd.Set(common.ToSnakeCase("ModifiedDt"), response.ModifiedDt)

	return nil
}
