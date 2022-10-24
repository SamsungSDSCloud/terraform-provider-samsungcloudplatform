package blockstorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/client/storage/blockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/common"
	"github.com/SamsungSDSCloud/terraform-provider-SamsungCloudPlatform/scp/service/virtualserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
			"virtual_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Virtual server ID to which you want to assign the block storage.",
			},
			"product_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "You can use by selecting SSD or HDD based storage.",
			},
			"encrypt_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "The block storage whether to use encryption. This can be enabled when the virtual server is encryption enabled.",
			},
		},
		Description: "Provides a Block Storage resource.",
	}
}

func createBlockStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	virtualServerId := data.Get("virtual_server_id").(string)
	serverInfo, _, err := inst.Client.VirtualServer.GetVirtualServer(ctx, virtualServerId)
	if err != nil {
		return diag.FromErr(err)
	}

	productId, _ := client.FindProductId(ctx, inst.Client, serverInfo.ProductGroupId, common.ProductTypeDisk, data.Get("product_name").(string))
	if len(productId) == 0 {
		return diag.Errorf("Matching available productId not found")
	}

	encryptEnable := data.Get("encrypt_enable").(bool) // TODO : (add Validation) Virtual Server 암호화 True -> EncryptEnable도 True가능

	err = virtualserver.WaitForVirtualServerStatus(ctx, inst.Client, virtualServerId, common.VirtualServerProcessingStates(), []string{common.RunningState, common.StoppedState}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := inst.Client.BlockStorage.CreateBlockStorage(ctx, blockstorage.CreateBlockStorageRequest{
		BlockStorageName: data.Get("name").(string),
		BlockStorageSize: (int32)(data.Get("storage_size_gb").(int)),
		EncryptEnabled:   encryptEnable,
		ProductId:        productId,
		VirtualServerId:  virtualServerId,
	})

	if err != nil {
		if err.Error() == "400 Bad Request" {
			return diag.Errorf("400 Bad Request (Adding an encryption disk is only available on an encrypted Virtual Server.)")
		}
		return diag.FromErr(err)
	}

	err = waitForBlockStorageStatus(ctx, inst.Client, response.ResourceId, []string{}, []string{"MOUNTED"}, true)
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
		return diag.FromErr(err)
	}

	data.Set("name", info.BlockStorageName)
	data.Set("storage_size_gb", info.BlockStorageSize)
	data.Set("encrypt_enable", info.EncryptEnabled)
	data.Set("product_id", info.ProductId)
	data.Set("virtual_server_id", info.VirtualServerId)

	return nil
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

		virtualServerId := data.Get("virtual_server_id").(string)
		serverInfo, _, err := inst.Client.VirtualServer.GetVirtualServer(ctx, virtualServerId)
		if err != nil {
			return diag.FromErr(err)
		}

		productId, _ := client.FindProductId(ctx, inst.Client, serverInfo.ProductGroupId, common.ProductTypeDisk, data.Get("product_name").(string))
		if len(productId) == 0 {
			return diag.Errorf("Matching available productId not found")
		}

		_, err = inst.Client.BlockStorage.ResizeBlockStorage(ctx, blockstorage.UpdateBlockStorageRequest{
			BlockStorageId:   data.Id(),
			BlockStorageSize: (int32)(data.Get("storage_size_gb").(int)),
			ProductId:        productId,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForBlockStorageStatus(ctx, inst.Client, data.Id(), []string{}, []string{"MOUNTED"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readBlockStorage(ctx, data, meta)
}

func deleteBlockStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.BlockStorage.DeleteBlockStorage(ctx, data.Id())
	if err != nil {
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
