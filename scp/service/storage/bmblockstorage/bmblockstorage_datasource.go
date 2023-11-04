package bmblockstorage

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	baremetalblockstorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/bare-metal-block-storage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_bm_block_storage", DataSourceBmBlockStorage())
	scp.RegisterDataSource("scp_bm_block_storages", DatasourceBmBlockStorages())
}

func DatasourceBmBlockStorages() *schema.Resource {
	return &schema.Resource{
		ReadContext: bmBlockStorageList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{ //스키마 정의
			"contents":    {Type: schema.TypeList, Optional: true, Description: "BareMetal Block Storages", Elem: bmBlockStoragesElem()},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides Block Storage(BM) List",
	}
}

func bmBlockStorageList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	responses, _, err := inst.Client.BareMetalBlockStorage.GetBareMetalBlockStorages(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)
	for i, v := range responses.Contents {
		contents[i][common.ToSnakeCase("BareMetalServerIds")] = v.BareMetalServerIds
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func bmBlockStoragesElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("BareMetalBlockStorageId"):      {Type: schema.TypeString, Computed: true, Description: "Baremetal Block Storage Id"},
			common.ToSnakeCase("BareMetalBlockStorageName"):    {Type: schema.TypeString, Computed: true, Description: "Baremetal Block Storage Name"},
			common.ToSnakeCase("BareMetalBlockStoragePurpose"): {Type: schema.TypeString, Computed: true, Description: "Baremetal Block Storage Purpose"},
			common.ToSnakeCase("BareMetalBlockStorageSize"):    {Type: schema.TypeInt, Computed: true, Description: "Baremetal Block Storage Size"},
			common.ToSnakeCase("BareMetalBlockStorageState"):   {Type: schema.TypeString, Computed: true, Description: "Baremetal Block Storage State"},
			common.ToSnakeCase("BareMetalBlockStorageTypeId"):  {Type: schema.TypeString, Computed: true, Description: "Baremetal Block Storage Type"},
			common.ToSnakeCase("BareMetalServerIds"):           {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Baremetal Server Ids"},
			common.ToSnakeCase("BlockId"):                      {Type: schema.TypeString, Computed: true, Description: "Block Id"},
			common.ToSnakeCase("EncryptionEnabled"):            {Type: schema.TypeBool, Computed: true, Description: "Encryption Enabled"},
			common.ToSnakeCase("Location"):                     {Type: schema.TypeString, Computed: true, Description: "Location"},
			common.ToSnakeCase("ServiceZoneId"):                {Type: schema.TypeString, Computed: true, Description: "Service Zone Id"},
			common.ToSnakeCase("CreatedBy"):                    {Type: schema.TypeString, Computed: true, Description: "Created By"},
			common.ToSnakeCase("CreatedDt"):                    {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			common.ToSnakeCase("ModifiedBy"):                   {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			common.ToSnakeCase("ModifiedDt"):                   {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
	}
}

func DataSourceBmBlockStorage() *schema.Resource {
	return &schema.Resource{
		ReadContext: bmBlockStorageDetail,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("StorageId"):                     {Type: schema.TypeString, Required: true, Description: "baremetal_block_storage_id"},
			common.ToSnakeCase("ProjectId"):                     {Type: schema.TypeString, Computed: true, Description: "project_id"},
			common.ToSnakeCase("BackupBareMetalBlockStorageId"): {Type: schema.TypeString, Computed: true, Description: "backup_baremetal_block_storage_id"},
			common.ToSnakeCase("BareMetalBlockStorageName"):     {Type: schema.TypeString, Computed: true, Description: "baremetal_block_storage_name"},
			common.ToSnakeCase("BareMetalBlockStoragePurpose"):  {Type: schema.TypeString, Computed: true, Description: "baremetal_block_storage_purpose"},
			common.ToSnakeCase("BareMetalBlockStorageSize"):     {Type: schema.TypeInt, Computed: true, Description: "baremetal_block_storage_size"},
			common.ToSnakeCase("BareMetalBlockStorageState"):    {Type: schema.TypeString, Computed: true, Description: "baremetal_block_storage_state"},
			common.ToSnakeCase("BlockId"):                       {Type: schema.TypeString, Computed: true, Description: "block_id"},
			common.ToSnakeCase("DrBareMetalBlockStorageId"):     {Type: schema.TypeString, Computed: true, Description: "dr_baremetal_block_storage_id"},
			common.ToSnakeCase("EncryptionEnabled"):             {Type: schema.TypeBool, Computed: true, Description: "encryption_enabled"},
			common.ToSnakeCase("ErrorCheck"):                    {Type: schema.TypeBool, Computed: true, Description: "error_check"},
			common.ToSnakeCase("IscsiTargetIp"):                 {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "iscsi_target_ip"},
			common.ToSnakeCase("OriginBareMetalBlockStorage"):   {Type: schema.TypeList, Computed: true, Elem: originBaremetalBlockStorageElem(), Description: "origin_baremetal_block_storage"},
			common.ToSnakeCase("ProductId"):                     {Type: schema.TypeString, Computed: true, Description: "product_id"},
			common.ToSnakeCase("Servers"):                       {Type: schema.TypeList, Computed: true, Elem: serversElem(), Description: "servers"},
			common.ToSnakeCase("ServiceZoneId"):                 {Type: schema.TypeString, Computed: true, Description: "service_zone_id"},
			common.ToSnakeCase("CreatedBy"):                     {Type: schema.TypeString, Computed: true, Description: "created_by"},
			common.ToSnakeCase("CreatedDt"):                     {Type: schema.TypeString, Computed: true, Description: "created_dt"},
			common.ToSnakeCase("ModifiedBy"):                    {Type: schema.TypeString, Computed: true, Description: "modified_by"},
			common.ToSnakeCase("ModifiedDt"):                    {Type: schema.TypeString, Computed: true, Description: "modified_dt"},
		},
		Description: "Provide Block Storage(BM) Detail",
	}
}

func bmBlockStorageDetail(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	bsId := rd.Get(common.ToSnakeCase("StorageId")).(string)

	response, _, err := inst.Client.BareMetalBlockStorage.GetBareMetalBlockStorageDetail(ctx, bsId)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set(common.ToSnakeCase("ProjectId"), response.ProjectId)
	rd.Set(common.ToSnakeCase("BackupBareMetalBlockStorageId"), response.BackupBareMetalBlockStorageId)
	rd.Set(common.ToSnakeCase("StorageId"), response.BareMetalBlockStorageId)
	rd.Set(common.ToSnakeCase("BareMetalBlockStorageName"), response.BareMetalBlockStorageName)
	rd.Set(common.ToSnakeCase("BareMetalBlockStoragePurpose"), response.BareMetalBlockStoragePurpose)
	rd.Set(common.ToSnakeCase("BareMetalBlockStorageSize"), response.BareMetalBlockStorageSize)
	rd.Set(common.ToSnakeCase("BareMetalBlockStorageState"), response.BareMetalBlockStorageState)
	rd.Set(common.ToSnakeCase("BlockId"), response.BlockId)
	rd.Set(common.ToSnakeCase("DrBareMetalBlockStorageId"), response.DrBareMetalBlockStorageId)
	rd.Set(common.ToSnakeCase("EncryptionEnabled"), response.EncryptionEnabled)
	rd.Set(common.ToSnakeCase("ErrorCheck"), response.ErrorCheck)
	rd.Set(common.ToSnakeCase("IscsiTargetIp"), response.IscsiTargetIp)
	rd.Set(common.ToSnakeCase("OriginBareMetalBlockStorage"), toElemList(response.OriginBareMetalBlockStorage))
	rd.Set(common.ToSnakeCase("ProductId"), response.ProductId)
	rd.Set(common.ToSnakeCase("Servers"), toElemList(response.Servers))
	rd.Set(common.ToSnakeCase("ServiceZoneId"), response.ServiceZoneId)
	rd.Set(common.ToSnakeCase("CreatedBy"), response.CreatedBy)
	rd.Set(common.ToSnakeCase("CreatedDt"), response.CreatedDt.String())
	rd.Set(common.ToSnakeCase("ModifiedBy"), response.ModifiedBy)
	rd.Set(common.ToSnakeCase("ModifiedDt"), response.ModifiedDt.String())

	return nil
}

func toElemList[T baremetalblockstorage.OriginBareMetalBlockStorageResponse | baremetalblockstorage.BlockStorageServerResponse](list []T) []map[string]interface{} {
	retList := make([]map[string]interface{}, 0)
	for _, v := range list {
		retList = append(retList, common.ToMap(v))
	}
	return retList
}

func originBaremetalBlockStorageElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("BareMetalBlockStorageId"):   {Type: schema.TypeString, Computed: true, Description: "baremetal_block_storage_id"},
			common.ToSnakeCase("BareMetalBlockStorageName"): {Type: schema.TypeString, Computed: true, Description: "baremetal_block_storage_name"},
			common.ToSnakeCase("BlockId"):                   {Type: schema.TypeString, Computed: true, Description: "block_id"},
			common.ToSnakeCase("ServiceZoneId"):             {Type: schema.TypeString, Computed: true, Description: "service_zone_id"},
		},
	}
}

func serversElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("ServerId"):   {Type: schema.TypeString, Computed: true, Description: "server_id"},
			common.ToSnakeCase("ServerType"): {Type: schema.TypeString, Computed: true, Description: "server_type"},
		},
	}
}
