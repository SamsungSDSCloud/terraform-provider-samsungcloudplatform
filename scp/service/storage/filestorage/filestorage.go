package filestorage

import (
	"context"
	"fmt"
	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/ScpDevTerra/trf-provider/scp/client/storage/filestorage"
	"github.com/ScpDevTerra/trf-provider/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFileStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: createFileStorage,
		ReadContext:   readFileStorage,
		UpdateContext: updateFileStorage,
		DeleteContext: deleteFileStorage,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The file storage name to create. (3 to 28 lowercase characters with _)",
				ValidateDiagFunc: common.ValidateName3to28UnderscoreLowercase,
			},
			"disk_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HDD / SSD / HP_SSD(Only Multi-node GPU Clusters)",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The file storage protocol type to create (NFS, CIFS)",
			},
			/* auto generated
			"storage_size_gb": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The storage size(GB) of the file storage to create. (1 to 102,400 GB)",
				ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics

					// Get attribute key
					attr := path[len(path)-1].(cty.GetAttrStep)
					attrKey := attr.Name

					// Get value
					value := (int32)(v.(int))

					if value < 1 || value > 102400 {
						diags = append(diags, diag.Diagnostic{
							Severity:      diag.Error,
							Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, fmt.Errorf("storage size is out of bounds. (1 < size < 102400 gb) ")),
							AttributePath: path,
						})
					}
					return diags
				},
			},
			*/
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The file storage region name to create.",
			},
			/* auto generated
			"cifs_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Description:      "Cifs ID is only available for CIFS protocol. (4 to 20 characters without specials)",
				ValidateDiagFunc: common.ValidateName4to20NoSpecialsLowercase, // 영소문자시작, 영소문자+숫자 4-20
			},
			*/
			"cifs_password": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Description:      "Cifs password is only available for CIFS protocol. (6 to 20 characters without following special characters ($, %, {, }, [, ], \", \\)",
				ValidateDiagFunc: common.ValidateName6to20, // 영문+숫자+특수문자 6-20 ($ % { } [ ] " \ 제외)
			},
			/* deleted in v3 api
			"snapshot_capacity_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The capacity rate of the snapshot to create.",
			},
			*/
			"snapshot_day_of_week": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Snapshot creation cycle, It is only available when you use a Snapshot",
			},
			"snapshot_frequency": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Snapshot creation frequency, It is only available when you use a Snapshot",
			},
			"snapshot_hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Snapshot creation hour, It is only available when you use a Snapshot",
			},
			"retention_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Snapshot archiving count",
			},
			"is_encrypted": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "The file storage whether to use encryption.",
			},
		},
		Description: "Provides a File Storage resource.",
	}
}

func createFileStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	protocol := data.Get("protocol").(string)
	//cifsId := data.Get("cifs_id").(string)
	cifsPw := data.Get("cifs_password").(string)
	/*
		if protocol == "CIFS" {
			if len(cifsId) == 0 {
				return diag.Errorf("The cifsId value is required.")
			}
			if len(cifsPw) == 0 {
				return diag.Errorf("The cifsPw value is required.")
			}
		} else {
			if len(cifsId) > 0 {
				return diag.Errorf("If it is not CIFS protocol, cifsId accepts null only")
			}
			if len(cifsPw) > 0 {
				return diag.Errorf("If it is not CIFS protocol, cifsPassword accepts null only")
			}
		}*/

	name := data.Get("name").(string)
	region := data.Get("region").(string)
	serviceZoneId, productGroupId, err := client.FindServiceZoneIdAndProductGroupId(ctx, inst.Client, region, "", common.FileStorageProductName)
	if err != nil {
		return diag.FromErr(err)
	}

	isNameInvalid, err := inst.Client.FileStorage.CheckFileStorage(ctx, filestorage.CheckFileStorageRequest{
		ServiceZoneId:   serviceZoneId,
		FileStorageName: name,
		//CifsId:          cifsId,
	})
	if isNameInvalid.Result {
		return diag.Errorf("Input storage name is invalid (maybe duplicated) : " + name)
	}

	//snapshotRate := (int32)(data.Get("snapshot_capacity_rate").(int))
	snapshotDayOfWeek := data.Get("snapshot_day_of_week").(string)
	snapshotFreq := data.Get("snapshot_frequency").(string)
	snapshotHour := (int32)(data.Get("snapshot_hour").(int))

	/*
		if snapshotRate > 0 {
			if snapshotRate < 1 || snapshotRate > 50 {
				return diag.Errorf("Snapshot capacity rate is out of bounds. (1 < rate < 50) ")
			}
		} else {
			if len(snapshotDayOfWeek) > 0 {
				return diag.Errorf("If snapshot capacity rate is 0, snapshot day of week accepts null only")
			}
			if len(snapshotFreq) > 0 {
				return diag.Errorf("If snapshot capacity rate is 0, snapshot frequency accepts null only")
			}
			if snapshotHour != 0 {
				return diag.Errorf("If snapshot capacity rate is 0, snapshot hour accepts null only")
			}
		}

	*/

	productIds, _ := client.FindProductIdByType(ctx, inst.Client, productGroupId, common.ProductTypeDisk)
	if len(productIds) == 0 {
		return diag.Errorf("matching available productIds not found")
	}

	var productId []string
	productId = append(productId, productIds[0]) // 일단 무조건 첫 번째 인덱스만 넘기는 걸로.

	request := filestorage.CreateFileStorageRequest{
		FileStorageName:     name,
		DiskType:            data.Get("disk_type").(string),
		FileStorageProtocol: protocol,
		ProductGroupId:      productGroupId,
		ProductIds:          productId,
		ServiceZoneId:       serviceZoneId,
		//SnapshotCapacityRate: snapshotRate,
		RetentionCount: (int32)(data.Get("retention_count").(int)),
		IsEncrypted:    data.Get("is_encrypted").(bool),
	}
	// 빈 값으로 데이터 넘기면 500 Error
	if protocol == "CIFS" {
		request.CifsPassword = cifsPw
	}
	//if snapshotRate > 0 {
	// snapshotRate > 0 이어도 schedule 설정안할 수 있음
	if len(snapshotDayOfWeek) > 0 {
		request.SnapshotSchedule = &filestorage.SnapshotSchedule{
			DayOfWeek: snapshotDayOfWeek,
			Frequency: snapshotFreq,
			Hour:      snapshotHour,
		}
	}
	//}

	response, err := inst.Client.FileStorage.CreateFileStorage(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForFileStorageStatus(ctx, inst.Client, response.ResourceId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.ResourceId)

	return readFileStorage(ctx, data, meta)
}

func readFileStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	info, _, err := inst.Client.FileStorage.ReadFileStorage(ctx, data.Id())

	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}

	data.Set("name", info.FileStorageName)
	data.Set("protocol", info.FileStorageProtocol)
	data.Set("is_encrypted", info.EncryptionEnabled)
	data.Set("cifs_id", info.CifsId)

	return nil
}

func updateFileStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/* capacityGb deleted in v3 Api
	inst := meta.(*client.Instance)

		if data.HasChanges("storage_size_gb") {
			info, _, err := inst.Client.FileStorage.ReadFileStorage(ctx, data.Id())
			if err != nil {
				return diag.FromErr(err)
			}


			if (int32)(data.Get("storage_size_gb").(int)) < info.FileStorageCapacityGb {
				return diag.Errorf("Only capacity expansion is possible")
			}

			_, err = inst.Client.FileStorage.IncreaseFileStorage(ctx, filestorage.UpdateFileStorageRequest{
				FileStorageId:         data.Id(),
				FileStorageCapacityGb: (int32)(data.Get("storage_size_gb").(int)),
			})
			if err != nil {
				return diag.FromErr(err)
			}


			err = waitForFileStorageStatus(ctx, inst.Client, data.Id(), []string{}, []string{"ACTIVE"}, true)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	*/
	return readFileStorage(ctx, data, meta)
}

func deleteFileStorage(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.FileStorage.DeleteFileStorage(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForFileStorageStatus(ctx, inst.Client, data.Id(), []string{}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForFileStorageStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.FileStorage.ReadFileStorage(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		if info.FileStorageId != id {
			return nil, "", fmt.Errorf("invalid resource status")
		}
		return info, info.FileStorageState, nil
	})
}
