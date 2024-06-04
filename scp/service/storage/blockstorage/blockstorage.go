package blockstorage

import (
	"context"
	"fmt"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/blockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/virtualserver"
	blockstorage2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/block-storage2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func init() {
	scp.RegisterResource("scp_block_storage", ResourceBlockStorage())
}

func ResourceBlockStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: createBlockStorage,
		ReadContext:   readBlockStorage,
		UpdateContext: updateBlockStorage,
		DeleteContext: deleteBlockStorage,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The block storage name to create. (3 to 28 characters with -)",
				ValidateDiagFunc: common.ValidateName3to28Dash,
			},
			"storage_size_gb": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The storage size(GB) of the block storage to create. (4 to  12288 GB)",
				/*ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics

					// Get attribute key
					attr := path[len(path)-1].(cty.GetAttrStep)
					attrKey := attr.Name

					// Get value
					value := (int32)(v.(int))

					if value < 4 || value > 16384 {
						diags = append(diags, diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, fmt.Errorf("capacity size is out of bounds. (4 to 12288 GB) ")),
							AttributePath: path,
						})
					}
					return diags
				},*/
			},
			//"virtual_server_id": {
			//	Type:        schema.TypeString,
			//	Required:    true,
			//	ForceNew:    true,
			//	Description: "Virtual server ID to which you want to assign the block storage.",
			//},
			"virtual_server_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Virtual server ID to which you want to assign the block storage.",
			},
			"virtual_server_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Virtual server IDs to which you want to assign the block storage.",
			},
			"product_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "You can use by selecting SSD or HDD based storage.",
			},
			"shared_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "You can use by selecting DEDICATED or SHARED",
			},
			"encrypt_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "The block storage whether to use encryption. This can be enabled when the virtual server is encryption enabled.",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Block Storage resource.",
	}
}

func createBlockStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	virtualServerId := data.Get("virtual_server_id").(string)

	virtualServerIds := make([]string, 0)
	for _, virtualServerId := range data.Get("virtual_server_ids").([]interface{}) {
		virtualServerIds = append(virtualServerIds, virtualServerId.(string))
	}

	if len(virtualServerId) == 0 && len(virtualServerIds) == 0 {
		return diag.Errorf("You should input virtual server Id !!")
	}

	var finalVirtualServerId string
	if len(virtualServerId) > 0 {
		finalVirtualServerId = virtualServerId
	}
	if len(virtualServerIds) > 0 {
		finalVirtualServerId = virtualServerIds[0]
	}

	encryptEnable := data.Get("encrypt_enable").(bool) // TODO : (add Validation) Virtual Server 암호화 True -> EncryptEnable도 True가능
	sharedType := data.Get("shared_type").(string)

	err := virtualserver.WaitForVirtualServerStatus(ctx, inst.Client, finalVirtualServerId, common.VirtualServerProcessingStates(), []string{common.RunningState, common.StoppedState}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := inst.Client.BlockStorage.CreateBlockStorage(ctx, blockstorage.CreateBlockStorageRequest{
		BlockStorageName: data.Get("name").(string),
		BlockStorageSize: (int32)(data.Get("storage_size_gb").(int)),
		EncryptEnabled:   encryptEnable,
		DiskType:         data.Get("product_name").(string),
		SharedType:       sharedType,
		VirtualServerId:  finalVirtualServerId,
	}, data.Get("tags").(map[string]interface{}))

	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForBlockStorageStatus(ctx, inst.Client, response.ResourceId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.ResourceId)

	return readBlockStorage(ctx, data, meta)
}

func readBlockStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	info, _, err := inst.Client.BlockStorage.ReadBlockStorage(ctx, data.Id())
	if err != nil {
		data.SetId("")
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	data.Set("name", info.BlockStorageName)
	data.Set("storage_size_gb", info.BlockStorageSize)
	data.Set("encrypt_enable", info.EncryptEnabled)
	data.Set("product_id", info.ProductId)
	data.Set("shared_type", info.SharedType)
	virtualServerIds := getVirtualServerIds(info)
	data.Set("virtual_server_ids", virtualServerIds)

	tfTags.SetTags(ctx, data, meta, data.Id())

	return nil
}

//func getTags(tagInfo tag.PageResponseV2OfTagResponse) map[string]string {
//	tags := make(map[string]string)
//	for _, content := range tagInfo.Contents {
//		tags[content.TagKey] = content.TagValue
//	}
//	return tags
//}

func getVirtualServerIds(info blockstorage2.BlockStorageResponse) []string {
	virtualServerIds := make([]string, 0)
	for _, virtualServer := range info.VirtualServers {
		virtualServerIds = append(virtualServerIds, virtualServer.VirtualServerId)
	}
	return virtualServerIds
}

// Block Storage Resize
func updateBlockStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if data.HasChanges("storage_size_gb") {
		info, _, err := inst.Client.BlockStorage.ReadBlockStorage(ctx, data.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		if (int32)(data.Get("storage_size_gb").(int)) < info.BlockStorageSize {
			return diag.Errorf("Only capacity expansion is possible")
		}

		virtualServerIds := make([]string, 0)
		for _, virtualServerId := range data.Get("virtual_server_ids").([]interface{}) {
			virtualServerIds = append(virtualServerIds, virtualServerId.(string))
		}
		_, err = inst.Client.BlockStorage.ResizeBlockStorage(ctx, blockstorage.UpdateBlockStorageRequest{
			BlockStorageId:   data.Id(),
			BlockStorageSize: (int32)(data.Get("storage_size_gb").(int)),
		})
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForBlockStorageStatus(ctx, inst.Client, data.Id(), []string{}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChanges("virtual_server_ids") {
		oldVmIds, newVmIds := getOldAndNewVmIds(data)
		deletedVmIds := getDeletedVmIds(oldVmIds, newVmIds)
		addedVmIds := getAddedVmIds(oldVmIds, newVmIds)

		for _, deletedVmId := range deletedVmIds {
			_, err := inst.Client.BlockStorage.DetachBlockStorage(ctx, data.Id(), blockstorage.BlockStorageDetachRequest{
				VirtualServerId: deletedVmId,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			err = waitForBlockStorageStatus(ctx, inst.Client, data.Id(), []string{}, []string{"ACTIVE"}, true)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		for _, addedVmId := range addedVmIds {
			//log.Println("AttachBlockStorage !!!!!")
			_, err := inst.Client.BlockStorage.AttachBlockStorage(ctx, data.Id(), blockstorage.BlockStorageAttachRequest{
				VirtualServerId: addedVmId,
			})
			if err != nil {
				return diag.FromErr(err)
			}

			err = waitForBlockStorageStatus(ctx, inst.Client, data.Id(), []string{}, []string{"ACTIVE"}, true)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		//log.Println("deletedVmIds : ", deletedVmIds)
		//log.Println("addedVmIds : ", addedVmIds)

	}

	err := tfTags.UpdateTags(ctx, data, meta, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return readBlockStorage(ctx, data, meta)
}

func getAddedVmIds(oldVmIds []string, newVmIds []string) []string {
	addedVmIds := make([]string, 0)
	for _, newVmId := range newVmIds {
		var i int
		for i = 0; i < len(oldVmIds); i++ {
			if strings.Compare(newVmId, oldVmIds[i]) == 0 {
				break
			}
		}
		if i == len(oldVmIds) {
			addedVmIds = append(addedVmIds, newVmId)
		}
	}
	return addedVmIds
}

func getDeletedVmIds(oldVmIds []string, newVmIds []string) []string {
	deletedVmIds := make([]string, 0)

	for _, oldVmId := range oldVmIds {
		var i int
		for i = 0; i < len(newVmIds); i++ {
			if strings.Compare(newVmIds[i], oldVmId) == 0 {
				break
			}
		}
		if i == len(newVmIds) {
			deletedVmIds = append(deletedVmIds, oldVmId)
		}
	}
	return deletedVmIds
}

func getOldAndNewVmIds(data *schema.ResourceData) ([]string, []string) {
	oldValue, newValue := data.GetChange("virtual_server_ids")
	oldValues := oldValue.([]interface{})
	newValues := newValue.([]interface{})
	oldVmIds := make([]string, len(oldValues))
	newVmIds := make([]string, len(newValues))
	for i, oldVmId := range oldValues {
		oldVmIds[i] = oldVmId.(string)
	}
	for i, newVmId := range newValues {
		newVmIds[i] = newVmId.(string)
	}
	//log.Println("oldVmIds : ", oldVmIds)
	//log.Println("newVmIds : ", newVmIds)
	return oldVmIds, newVmIds
}

func deleteBlockStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.BlockStorage.DeleteBlockStorage(ctx, data.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitForBlockStorageStatus(ctx, inst.Client, data.Id(), []string{}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForBlockStorageStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.BlockStorage.ReadBlockStorage(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}

			return nil, "", err
		}
		if info.BlockStorageId != id {
			return nil, "", fmt.Errorf("invalid resource status")
		}
		return info, info.BlockStorageState, nil
	})
}
