package filestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

func init() {
	scp.RegisterResource("scp_file_storage", ResourceFileStorage())
}

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
			"cifs_password": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Description:      "CIFS Password is only available for CIFS Protocol. (6 to 20 alphabet and numeric characters without following special characters ($, %, {, }, [, ], \", \\)",
				ValidateDiagFunc: common.ValidateName6to20AlphaAndNumericWithoutSomeSpecials, // 영문+숫자+특수문자 6-20 ($ % { } [ ] " \ 제외)
			},
			"disk_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "File Storage Disk Type (HDD, SSD, HP_SSD)",
			},
			"file_storage_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "File Storage Name (3 to 21 lower alphabet and numeric characters with '_' symbol are available, but it must be started with lower alphabet)",
				ValidateDiagFunc: common.ValidateName3to21LowerAlphaAndNumericWithUnderscoreStartsWithLowerAlpha,
			},
			"file_storage_name_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "File Storage Name with UUID (10 to 28 lower alphabet and numeric characters with '_' symbol are available, but it must be started with lower alphabet)",
			},
			"file_storage_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "File Storage Protocol (NFS, CIFS)",
			},
			"multi_availability_zone": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Multi AZ (If null, default value is false)",
			},
			"product_names": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Product Names",
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service Zone ID",
			},
			"cifs_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CIFS ID",
			},
			"snapshot_retention_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Snapshot retention count",
			},
			"snapshot_schedule": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Snapshot schedule",
			},
			"frequency": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Snapshot schedule frequency must be one of \"DAILY\" or \"WEEKLY\"",
			},
			"day_of_week": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Snapshot schedule dayOfWeek must be one of \"SUN\", \"MON\", \"TUE\", \"WED\", \"THU\", \"FRI\" or \"SAT\"",
			},
			"hour": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Snapshot schedule hour (0 to 23)",
			},
			"tags": tfTags.TagsSchema(),
			"file_unit_recovery_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "File Unit Recovery",
			},
			"link_objects": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Link Objects",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"link_object_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Link object ID",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type",
						},
					},
				},
			},
			"unlink_objects": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Unlink Objects",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"link_object_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Link object ID",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type",
						},
					},
				},
			},
		},
		Description: "Provides a File Storage resource.",
	}
}

func createFileStorage(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	fileStorageName := rd.Get("file_storage_name").(string)
	serviceZoneId := rd.Get("service_zone_id").(string)

	isNameInvalid, err := inst.Client.FileStorage.CheckFileStorage(ctx, filestorage.CheckFileStorageRequest{
		ServiceZoneId:   serviceZoneId,
		FileStorageName: fileStorageName,
	})
	if err != nil {
		return diag.FromErr(err)
	} else if isNameInvalid.Result == nil || *isNameInvalid.Result {
		return diag.Errorf("File Storage Name is duplicated")
	}

	cifsPassword := rd.Get("cifs_password").(string)
	diskType := rd.Get("disk_type").(string)
	fileStorageProtocol := rd.Get("file_storage_protocol").(string)
	multiAvailabilityZone := rd.Get("multi_availability_zone").(bool)
	productNames := convertToStringArray(rd.Get("product_names").([]interface{}))
	snapshotRetentionCount := (int32)(rd.Get("snapshot_retention_count").(int))

	request := filestorage.CreateFileStorageRequest{
		DiskType:               diskType,
		FileStorageName:        fileStorageName,
		FileStorageProtocol:    fileStorageProtocol,
		MultiAvailabilityZone:  &multiAvailabilityZone,
		ProductNames:           productNames,
		ServiceZoneId:          serviceZoneId,
		SnapshotRetentionCount: &snapshotRetentionCount,
		SnapshotSchedule:       getSnapshotSchedule(rd),
		Tags:                   rd.Get("tags").(map[string]interface{}),
	}

	// 빈 값으로 데이터 넘기면 500 Error
	if fileStorageProtocol == "CIFS" {
		request.CifsPassword = cifsPassword
	}

	response, err := inst.Client.FileStorage.CreateFileStorage(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForFileStorageStatus(ctx, inst.Client, response.ResourceId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ResourceId)

	// if file_unit_recovery_enabled true//
	if rd.Get("file_unit_recovery_enabled").(bool) {
		fileUnitRecoveryEnabled := rd.Get("file_unit_recovery_enabled").(bool)
		if _, err := inst.Client.FileStorage.UpdateFileStorageFileRecoveryEnabled(ctx, rd.Id(), fileUnitRecoveryEnabled); err != nil {
			return diag.FromErr(err)
		}

		errUpdateRecovery := waitForFileStorageStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
		if errUpdateRecovery != nil {
			return diag.FromErr(errUpdateRecovery)
		}

	}
	// TODO: if attach objects have
	linkObjects, ok := rd.Get("link_objects").(map[string]interface{})
	if ok {
		if len(linkObjects) > 0 {
			if _, err := inst.Client.FileStorage.UpdateFileStorageObjectsLink(ctx, rd.Id(), filestorage.LinkFileStorageObjectRequest{
				LinkObjects: getLinkObjectsArray(rd),
			}); err != nil {
				return diag.FromErr(err)
			}
			errUpdateRecovery := waitForFileStorageStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
			if errUpdateRecovery != nil {
				return diag.FromErr(errUpdateRecovery)
			}
		}
	}
	return readFileStorage(ctx, rd, meta)
}

func readFileStorage(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	info, _, err := inst.Client.FileStorage.ReadFileStorage(ctx, rd.Id())
	if err != nil {
		rd.SetId("")

		//not show error message for deleted resource
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}
	snapshot, _, err := inst.Client.FileStorage.ReadFileStorageSnapshotSchedule(ctx, rd.Id())

	linkedObjects := common.HclSetObject{}
	for _, linkedObject := range info.LinkedObjects {
		linkedObjects = append(linkedObjects, common.HclKeyValueObject{
			"link_object_id": linkedObject.LinkedObjectId,
			"type":           linkedObject.LinkedObjectType,
		})
	}

	if err != nil {
		rd.SetId("")

		//not show error message for deleted resource
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}
	fileStorageName := info.FileStorageName[:len(info.FileStorageName)-7]
	rd.Set("disk_type", info.DiskType)
	rd.Set("file_storage_name", fileStorageName)
	rd.Set("file_storage_name_uuid", info.FileStorageName)
	rd.Set("file_storage_protocol", info.FileStorageProtocol)
	rd.Set("multi_availability_zone", info.MultiAvailabilityZone)
	rd.Set("service_zone_id", info.ServiceZoneId)
	rd.Set("cifs_id", info.CifsId)
	rd.Set("file_unit_recovery_enabled", info.FileUnitRecoveryEnabled)
	if len(linkedObjects) > 0 {
		rd.Set("link_objects", linkedObjects)
	}
	if *snapshot.IsSnapshotPolicy {
		rd.Set("snapshot_retention_count", *snapshot.SnapshotRetentionCount)
		rd.Set("snapshot_schedule.day_of_week", snapshot.SnapshotSchedule.DayOfWeek)
		rd.Set("snapshot_schedule.hour", *snapshot.SnapshotSchedule.Hour)

		if snapshot.SnapshotSchedule.Frequency == "NONE" {
			rd.Set("snapshot_schedule.frequency", "") // set "" to avoid omit-empty, need to resolve this issue later in SDK.
		} else {
			rd.Set("snapshot_schedule.frequency", snapshot.SnapshotSchedule.Frequency)
		}
	}

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func updateFileStorage(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {

	inst := meta.(*client.Instance)

	if rd.HasChanges("link_objects", "unlink_objects") {
		unlinkObjects := rd.Get("unlink_objects").([]interface{})
		linkObjects := rd.Get("link_objects").([]interface{})
		if (len(linkObjects) > 0) || (len(unlinkObjects) > 0) {
			if _, err := inst.Client.FileStorage.UpdateFileStorageObjectsLink(ctx, rd.Id(), filestorage.LinkFileStorageObjectRequest{
				LinkObjects:   getLinkObjectsArray(rd),
				UnlinkObjects: getUnLinkObjectsArray(rd),
			}); err != nil {
				return diag.FromErr(err)
			}
			err := waitForFileStorageStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if rd.HasChanges("file_unit_recovery_enabled") {
		fileUnitRecoveryEnabled := rd.Get("file_unit_recovery_enabled").(bool)
		if _, err := inst.Client.FileStorage.UpdateFileStorageFileRecoveryEnabled(ctx, rd.Id(), fileUnitRecoveryEnabled); err != nil {
			return diag.FromErr(err)
		}
		err := waitForFileStorageStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChanges("snapshot_retention_count", "snapshot_schedule.day_of_week", "snapshot_schedule.frequency", "snapshot_schedule.hour") {
		newDayOfWeek := rd.Get("snapshot_schedule.day_of_week").(string)
		oldFrequencyTmp, newFrequencyTmp := rd.GetChange("snapshot_schedule.frequency")
		beforeHour := rd.Get("snapshot_schedule.hour").(string)
		hourInt, numerr := strconv.Atoi(beforeHour)
		if numerr != nil {
			// 변환 오류 처리
			log.Printf("Error converting hour string to int: %v\n", numerr)
		}
		newHour := int32(hourInt)
		oldRetentionCountTmp, newRetentionCountTmp := rd.GetChange("snapshot_retention_count")
		oldFrequency := oldFrequencyTmp.(string)
		newFrequency := newFrequencyTmp.(string)
		oldRetentionCount := (int32)(oldRetentionCountTmp.(int))
		newRetentionCount := (int32)(newRetentionCountTmp.(int))

		state := "UPDATE"

		if len(oldFrequency) == 0 && len(newFrequency) > 0 && oldRetentionCount == 0 && newRetentionCount > 0 {
			state = "CREATE"
		} else if len(oldFrequency) > 0 && len(newFrequency) == 0 && oldRetentionCount > 0 && newRetentionCount == 0 {
			state = "DELETE"
		}

		if state == "CREATE" {
			if _, err := inst.Client.FileStorage.CreateFileStorageSnapshotSchedule(ctx, rd.Id(), newRetentionCount, &filestorage.SnapshotSchedule{
				DayOfWeek: newDayOfWeek,
				Frequency: newFrequency,
				Hour:      &newHour,
			}); err != nil {
				return diag.FromErr(err)
			}
		} else if state == "DELETE" {
			if _, err := inst.Client.FileStorage.DeleteFileStorageSnapshotSchedule(ctx, rd.Id()); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if _, err := inst.Client.FileStorage.UpdateFileStorageSnapshotSchedule(ctx, rd.Id(), newRetentionCount, &filestorage.SnapshotSchedule{
				DayOfWeek: newDayOfWeek,
				Frequency: newFrequency,
				Hour:      &newHour,
			}); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	err := waitForFileStorageStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	err = tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return readFileStorage(ctx, rd, meta)
}

func deleteFileStorage(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	_, err := inst.Client.FileStorage.DeleteFileStorage(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitForFileStorageStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"DELETED"}, false)
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

func getLinkObjectsArray(rd *schema.ResourceData) []filestorage.LinkObjectRequest {
	linkObjects := rd.Get("link_objects").([]interface{})
	linkObjectRequest := make([]filestorage.LinkObjectRequest, 0)
	for _, linkObject := range linkObjects {
		if linkObject == nil {
			continue
		}
		linkObjectMap, ok := linkObject.(map[string]interface{})
		if !ok {
			fmt.Println("linkObject가 map[string]interface{} 타입이 아닙니다.")
			continue
		}
		linkObjectRequest = append(linkObjectRequest, filestorage.LinkObjectRequest{
			ObjectId: linkObjectMap["link_object_id"].(string),
			Type:     linkObjectMap["type"].(string),
		})
	}
	return linkObjectRequest
}

func getUnLinkObjectsArray(rd *schema.ResourceData) []filestorage.LinkObjectRequest {
	linkObjects := rd.Get("unlink_objects").([]interface{})
	unlinkObjectRequest := make([]filestorage.LinkObjectRequest, 0)
	for _, linkObject := range linkObjects {
		if linkObject == nil {
			continue
		}
		linkObjectMap, ok := linkObject.(map[string]interface{})
		if !ok {
			fmt.Println("linkObject가 map[string]interface{} 타입이 아닙니다.")
			continue
		}
		unlinkObjectRequest = append(unlinkObjectRequest, filestorage.LinkObjectRequest{
			ObjectId: linkObjectMap["link_object_id"].(string),
			Type:     linkObjectMap["type"].(string),
		})
	}
	return unlinkObjectRequest
}

func getSnapshotSchedule(rd *schema.ResourceData) filestorage.SnapshotSchedule {
	frequency := rd.Get("snapshot_schedule.frequency").(string)
	day_of_week := rd.Get("snapshot_schedule.day_of_week").(string)
	beforeHour := rd.Get("snapshot_schedule.hour").(string)
	hourInt, err := strconv.Atoi(beforeHour)
	if err != nil {
		// 변환 오류 처리
		log.Printf("Error converting hour string to int: %v\n", err)
	}
	hour := int32(hourInt)
	hourPtr := &hour
	snapshotScheduleRequest := filestorage.SnapshotSchedule{
		DayOfWeek: day_of_week,
		Frequency: frequency,
		Hour:      hourPtr,
	}
	return snapshotScheduleRequest
}

func convertToStringArray(interfaceArray []interface{}) []string {
	stringArray := make([]string, 0)
	for _, interfaceElem := range interfaceArray {
		stringArray = append(stringArray, interfaceElem.(string))
	}
	return stringArray
}
