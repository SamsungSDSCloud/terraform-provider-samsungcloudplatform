package filestorage

import (
	"context"

	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client/storage/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func DatasourceFileStorages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"file_storage_id":       {Type: schema.TypeString, Optional: true, Description: "File Storage id"},
			"file_storage_name":     {Type: schema.TypeString, Optional: true, Description: "File Storage name"},
			"disk_type":             {Type: schema.TypeString, Optional: true, Default: "SSD", Description: "Disk type(HDD / SSD / HP_SSD)"},
			"file_storage_protocol": {Type: schema.TypeString, Optional: true, Description: "File Storage protocol(NFS, CIFS)"},
			"file_storage_state":    {Type: schema.TypeString, Optional: true, Description: "File Storage status"},
			"service_zone_id":       {Type: schema.TypeString, Optional: true, Description: "Service zone id"},
			"created_by":            {Type: schema.TypeString, Optional: true, Description: "Person who created the resource"},
			"page":                  {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":                  {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":              {Type: schema.TypeList, Optional: true, Description: "File Storage list", Elem: datasourceElem()},
			"total_count":           {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of file storages",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	requestParam := filestorage.ReadFileStorageRequest{
		FileStorageId:       rd.Get("file_storage_id").(string),
		FileStorageName:     rd.Get("file_storage_name").(string),
		FileStorageProtocol: rd.Get("file_storage_protocol").(string),
		DiskType:            rd.Get("disk_type").(string),
		ServiceZoneId:       rd.Get("service_zone_id").(string),
		CreatedBy:           rd.Get("created_by").(string),
		Page:                (int32)(rd.Get("page").(int)),
		Size:                (int32)(rd.Get("size").(int)),
	}

	responses, err := inst.Client.FileStorage.ReadFileStorageList(ctx, requestParam)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":            {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"block_id":              {Type: schema.TypeString, Computed: true, Description: "Block id of this region"},
			"disk_type":             {Type: schema.TypeString, Computed: true, Description: "Disk type(HDD / SSD / HP_SSD)"},
			"encryption_enabled":    {Type: schema.TypeBool, Computed: true, Description: "Enable encryption feature in storage"},
			"file_storage_id":       {Type: schema.TypeString, Computed: true, Description: "File Storage id"},
			"file_storage_name":     {Type: schema.TypeString, Computed: true, Description: "File Storage name"},
			"file_storage_protocol": {Type: schema.TypeString, Computed: true, Description: "File Storage protocol type(NFS, CIFS)"},
			"file_storage_purpose":  {Type: schema.TypeString, Computed: true, Description: "Purpose of file storage"},
			"file_storage_state":    {Type: schema.TypeString, Computed: true, Description: "File Storage status"},
			//"linked_object_state": {},
			"product_group_id": {Type: schema.TypeString, Computed: true, Description: "File Storage group id"},
			"service_zone_id":  {Type: schema.TypeString, Computed: true, Description: "Service zone id "},
			"created_by":       {Type: schema.TypeString, Computed: true, Description: "The person who created the resource"},
			"created_dt":       {Type: schema.TypeString, Computed: true, Description: "Creation time"},
			"modified_by":      {Type: schema.TypeString, Computed: true, Description: "The person who modified the resource"},
			"modified_dt":      {Type: schema.TypeString, Computed: true, Description: "Modification time"},
		},
	}
}
