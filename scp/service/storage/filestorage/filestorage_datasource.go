package filestorage

import (
	"context"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	filestorage2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/file-storage2"
	"github.com/antihax/optional"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_file_storages", DatasourceFileStorages())
}

func DatasourceFileStorages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"block_id":              {Type: schema.TypeString, Optional: true, Description: "Block ID"},
			"file_storage_id":       {Type: schema.TypeString, Optional: true, Description: "File Storage ID"},
			"file_storage_name":     {Type: schema.TypeString, Optional: true, Description: "File Storage Name"},
			"file_storage_protocol": {Type: schema.TypeString, Optional: true, Description: "File Storage Protocol"},
			"file_storage_state":    {Type: schema.TypeString, Optional: true, Description: "File Storage State"},
			"file_storage_states":   {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "File Storage State"}, Description: "File Storage States"},
			"service_zone_id":       {Type: schema.TypeString, Optional: true, Description: "Service Zone ID"},
			"created_by":            {Type: schema.TypeString, Optional: true, Description: "Created By"},
			"page":                  {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                  {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"sort":                  {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, Description: "Sort"}, Description: "Sort"},
			"contents":              {Type: schema.TypeList, Computed: true, Description: "File Storage List", Elem: datasourceElem()},
			"total_count":           {Type: schema.TypeInt, Computed: true, Description: "Total List Size"},
		},
		Description: "Provides list of file storages",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	var blockId optional.String
	var fileStorageId optional.String
	var fileStorageName optional.String
	var fileStorageProtocol optional.String
	var fileStorageState optional.String
	var fileStorageStates optional.Interface
	var serviceZoneId optional.String
	var createdBy optional.String
	var page optional.Int32
	var size optional.Int32
	var sort optional.Interface

	if v, ok := rd.GetOk("block_id"); ok {
		blockId = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("file_storage_id"); ok {
		fileStorageId = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("file_storage_name"); ok {
		fileStorageName = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("file_storage_protocol"); ok {
		fileStorageProtocol = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("file_storage_state"); ok {
		fileStorageState = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("file_storage_states"); ok {
		fileStorageStates = optional.NewInterface(v.([]interface{}))
	}
	if v, ok := rd.GetOk("service_zone_id"); ok {
		serviceZoneId = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("created_by"); ok {
		createdBy = optional.NewString(v.(string))
	}
	if v, ok := rd.GetOk("page"); ok {
		page = optional.NewInt32(int32(v.(int)))
	}
	if v, ok := rd.GetOk("size"); ok {
		size = optional.NewInt32(int32(v.(int)))
	}
	if v, ok := rd.GetOk("sort"); ok {
		sort = optional.NewInterface(v.([]interface{}))
	}

	responses, err := inst.Client.FileStorage.ReadFileStorageList(ctx, filestorage2.FileStorageOpenApiV3ApiListFileStoragesOpts{
		BlockId:             blockId,
		FileStorageId:       fileStorageId,
		FileStorageName:     fileStorageName,
		FileStorageProtocol: fileStorageProtocol,
		FileStorageState:    fileStorageState,
		FileStorageStates:   fileStorageStates,
		ServiceZoneId:       serviceZoneId,
		CreatedBy:           createdBy,
		Page:                page,
		Size:                size,
		Sort:                sort,
	})

	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", *responses.TotalCount)

	return nil
}

func datasourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":            {Type: schema.TypeString, Computed: true, Description: "Project ID"},
			"block_id":              {Type: schema.TypeString, Computed: true, Description: "Block ID"},
			"disk_type":             {Type: schema.TypeString, Computed: true, Description: "Disk Type"},
			"encryption_enabled":    {Type: schema.TypeBool, Computed: true, Description: "Encryption enabled"},
			"file_storage_id":       {Type: schema.TypeString, Computed: true, Description: "File Storage ID"},
			"file_storage_name":     {Type: schema.TypeString, Computed: true, Description: "File Storage Name"},
			"file_storage_protocol": {Type: schema.TypeString, Computed: true, Description: "File Storage Protocol"},
			"file_storage_purpose":  {Type: schema.TypeString, Computed: true, Description: "File Storage Purpose"},
			"file_storage_state":    {Type: schema.TypeString, Computed: true, Description: "File Storage State"},
			"linked_object_count":   {Type: schema.TypeInt, Computed: true, Description: "Linked Object Count"},
			"product_group_id":      {Type: schema.TypeString, Computed: true, Description: "Product Group ID"},
			"service_zone_id":       {Type: schema.TypeString, Computed: true, Description: "Service Zone ID"},
			"tiering_enabled":       {Type: schema.TypeBool, Computed: true, Description: "Tiering enabled"},
			"created_by":            {Type: schema.TypeString, Computed: true, Description: "Created By"},
			"created_dt":            {Type: schema.TypeString, Computed: true, Description: "Created Date"},
			"modified_by":           {Type: schema.TypeString, Computed: true, Description: "Modified By"},
			"modified_dt":           {Type: schema.TypeString, Computed: true, Description: "Modified Date"},
		},
	}
}
