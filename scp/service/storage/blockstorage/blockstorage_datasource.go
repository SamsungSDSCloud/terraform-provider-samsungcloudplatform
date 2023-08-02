package blockstorage

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/client/storage/blockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	scp.RegisterDataSource("scp_block_storages", DatasourceBlockStorages())
}

func DatasourceBlockStorages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			//"block_storage_name":  {Type: schema.TypeString, Optional: true, Description: "Block Storage Name"},
			"block_storage_id":  {Type: schema.TypeString, Optional: true, Description: "block_storage_id"},
			"virtual_server_id": {Type: schema.TypeString, Optional: true, Description: "Virtual server id"},
			//"virtual_server_name": {Type: schema.TypeString, Optional: true, Description: "Virtual Server Name"},
			"created_by":  {Type: schema.TypeString, Optional: true, Description: "The person who created the resource"},
			"page":        {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Page start number from which to get the list"},
			"size":        {Type: schema.TypeInt, Optional: true, Default: 20, Description: "Size to get list"},
			"contents":    {Type: schema.TypeList, Optional: true, Description: "Block Storage list", Elem: datasourceElem()},
			"total_count": {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of block storages",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if len(rd.Get("block_storage_id").(string)) > 0 {

		response, _, err := inst.Client.BlockStorage.ReadBlockStorage(ctx, rd.Get("block_storage_id").(string))
		if err != nil {
			diag.FromErr(err)
		}

		contents := make([]map[string]interface{}, 0)

		content := common.ToMap(response)
		mapVirtualServerList := make([]map[string]interface{}, 0)
		for _, virtualServer := range response.VirtualServers {
			mapVirtualServer := common.ToMap(virtualServer)
			mapVirtualServerList = append(mapVirtualServerList, mapVirtualServer)
		}
		content["virtual_servers"] = mapVirtualServerList
		contents = append(contents, content)

		rd.SetId(uuid.NewV4().String())
		rd.Set("contents", contents)
		rd.Set("total_count", 1)
	} else {

		requestParam := blockstorage.ReadBlockStorageRequest{
			//BlockStorageName: rd.Get("block_storage_name").(string),
			VirtualServerId: rd.Get("virtual_server_id").(string),
			//VirtualServerName: rd.Get("virtual_server_name").(string),
			CreatedBy: rd.Get("created_by").(string),
			Page:      (int32)(rd.Get("page").(int)),
			Size:      (int32)(rd.Get("size").(int)),
		}

		responses, err := inst.Client.BlockStorage.ReadBlockStorageList(ctx, requestParam)
		if err != nil {
			return diag.FromErr(err)
		}

		contents := common.ConvertStructToMaps(responses.Contents)

		for i, resContent := range responses.Contents {
			mapVirtualServerList := make([]map[string]interface{}, 0)
			for _, virtualServer := range resContent.VirtualServers {
				mapVirtualServer := common.ToMap(virtualServer)
				mapVirtualServerList = append(mapVirtualServerList, mapVirtualServer)
			}
			contents[i]["virtual_servers"] = mapVirtualServerList
		}

		rd.SetId(uuid.NewV4().String())
		rd.Set("contents", contents)
		rd.Set("total_count", responses.TotalCount)
	}
	return nil
}

func datasourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id":          {Type: schema.TypeString, Computed: true, Description: "Project id"},
			"block_id":            {Type: schema.TypeString, Computed: true, Description: "Block id of this region"},
			"block_storage_id":    {Type: schema.TypeString, Computed: true, Description: "Block Storage id"},
			"block_storage_name":  {Type: schema.TypeString, Computed: true, Description: "Block storage name to create."},
			"block_storage_size":  {Type: schema.TypeInt, Computed: true, Description: "Storage size(GB)"},
			"block_storage_state": {Type: schema.TypeString, Computed: true, Description: "Block storage status"},
			"block_storage_uuid":  {Type: schema.TypeString, Computed: true, Description: "Block Storage uuid"},
			"device_node":         {Type: schema.TypeString, Computed: true, Description: "Device node"},
			"encrypt_enabled":     {Type: schema.TypeBool, Computed: true, Description: "Enable encryption feature in storage"},
			"is_boot_disk":        {Type: schema.TypeBool, Computed: true, Description: "Check whether it is OS(Boot) disk or not"},
			"mount_path":          {Type: schema.TypeString, Computed: true, Description: "Mount path"},
			"product_id":          {Type: schema.TypeString, Computed: true, Description: "Product id of block storage"},
			"service_zone_id":     {Type: schema.TypeString, Computed: true, Description: "Service zone id"},
			"shared_type":         {Type: schema.TypeString, Computed: true, Description: "Shared type of block storage"},
			"virtual_server_id":   {Type: schema.TypeString, Computed: true, Description: "Virtual server id to assign the block storage."},
			"virtual_servers": {Type: schema.TypeList, Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_state":       {Type: schema.TypeString, Computed: true, Description: "Mount State"},
						"virtual_server_id": {Type: schema.TypeString, Computed: true, Description: "Virtual Server Id"},
					}},
				Description: "Mounted Virtual Servers",
			},
			"created_by":  {Type: schema.TypeString, Computed: true, Description: "Person who created the resource"},
			"created_dt":  {Type: schema.TypeString, Computed: true, Description: "Creation time"},
			"modified_by": {Type: schema.TypeString, Computed: true, Description: "Person who modified the resource"},
			"modified_dt": {Type: schema.TypeString, Computed: true, Description: "Modification time"},
		},
	}
}
