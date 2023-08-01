package filestorage

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/storage/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			/*"snapshot_retention_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Snapshot retention count",
				ValidateFunc: validation.All(
					validation.IntBetween(1, 128),
				),
			},
			"snapshotSchedule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Snapshot schedule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},*/
			//"tags": {
			//	Type:     schema.TypeList,
			//	Optional: true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"tag_key": {
			//				Type:             schema.TypeString,
			//				Required:         true,
			//				ValidateDiagFunc: common.ValidateName1to256DotDashUnderscore,
			//				Description:      "Tag Key",
			//			},
			//			"tag_value": {
			//				Type:             schema.TypeString,
			//				Optional:         true,
			//				ValidateDiagFunc: common.ValidateName1to256DotDashUnderscore,
			//				Description:      "Tag Value",
			//			},
			//		},
			//	},
			//	Description: "Tag list",
			//},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags",
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

	//snapshotRetentionCount := (int32)(rd.Get("snapshot_retention_count").(int))

	request := filestorage.CreateFileStorageRequest{
		DiskType:              diskType,
		FileStorageName:       fileStorageName,
		FileStorageProtocol:   fileStorageProtocol,
		MultiAvailabilityZone: &multiAvailabilityZone,
		ProductNames:          productNames,
		ServiceZoneId:         serviceZoneId,
		//SnapshotRetentionCount: &snapshotRetentionCount,
		Tags: getTagRequestArray(rd),
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

	/*snapshot, _, err := inst.Client.FileStorage.ReadFileStorageSnapshotSchedule(ctx, rd.Id())
	if err != nil {
		rd.SetId("")

		//not show error message for deleted resource
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}*/

	rd.Set("disk_type", info.DiskType)
	rd.Set("file_storage_name", info.FileStorageName)
	rd.Set("file_storage_protocol", info.FileStorageProtocol)
	rd.Set("multi_availability_zone", info.MultiAvailabilityZone)
	rd.Set("service_zone_id", info.ServiceZoneId)
	rd.Set("cifs_id", info.CifsId)

	/*if *snapshot.IsSnapshotPolicy {
		rd.Set("snapshot_retention_count", *snapshot.SnapshotRetentionCount)
		rd.Set("snapshotSchedule.day_of_week", snapshot.SnapshotSchedule.DayOfWeek)
		rd.Set("snapshotSchedule.hour", *snapshot.SnapshotSchedule.Hour)

		if snapshot.SnapshotSchedule.Frequency == "NONE" {
			rd.Set("snapshot_frequency", "") // set "" to avoid omit-empty, need to resolve this issue later in SDK.
		} else {
			rd.Set("snapshotSchedule.frequency", snapshot.SnapshotSchedule.Frequency)
		}
	}*/
	return nil
}

func updateFileStorage(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*inst := meta.(*client.Instance)

	if rd.HasChanges("snapshot_retention_count", "snapshot_day_of_week", "snapshot_frequency", "snapshot_hour") {
		newDayOfWeek := rd.Get("snapshot_day_of_week").(string)
		oldFrequencyTmp, newFrequencyTmp := rd.GetChange("snapshot_frequency")
		newHour := (int32)(rd.Get("snapshot_hour").(int))
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
		err := waitForFileStorageStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}*/
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

func getTagRequestArray(rd *schema.ResourceData) []filestorage.TagRequest {
	tags := rd.Get("tags").(map[string]interface{})
	tagsRequests := make([]filestorage.TagRequest, 0)
	for key, value := range tags {
		tagsRequests = append(tagsRequests, filestorage.TagRequest{
			TagKey:   key,
			TagValue: value.(string),
		})
	}
	return tagsRequests
}

func convertToStringArray(interfaceArray []interface{}) []string {
	stringArray := make([]string, 0)
	for _, interfaceElem := range interfaceArray {
		stringArray = append(stringArray, interfaceElem.(string))
	}
	return stringArray
}
