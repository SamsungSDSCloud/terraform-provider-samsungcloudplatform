package backup

import (
	"context"
	"fmt"
	scp "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/backup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	tfTags "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/service/tag"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/backup2"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func init() {
	scp.RegisterResource("scp_backup", ResourceBackup())
}

func ResourceBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: createBackup,
		ReadContext:   readBackup,
		UpdateContext: updateBackup,
		DeleteContext: deleteBackup,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"az_code": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Multi AZ Code",
			},
			"backup_dr_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Backup(DR) Service Zone Id",
			},
			"backup_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "Backup Name",
				ValidateDiagFunc: common.ValidateName3to30AlphaNumeric,
			},
			"backup_policy_type_category": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Backup Policy Type Category (VM, FILESYSTEM)",
			},
			"backup_repository": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Backup Repository (SD_STORAGE)",
			},
			"dr_az_code": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Multi AZ(DR) Code",
			},
			/*"file_system_backup_selections": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description:      "Target Filesystem (INSTANCE, BAREMETAL)",
			},*/
			"is_backup_dr_enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backup(DR) Activation (If 'Y', Backup(DR) will be activated)",
			},
			"is_backup_dr_deleted": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Is Backup DR Deleted.",
			},
			"is_backup_dr_destroy_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "IF 'Y', Destroy DR replica together.",
			},
			"backup_dr_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backup DR ID",
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Backup Object ID",
			},
			"object_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Backup Object Type",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Backup Policy Type",
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
			"retention_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Full Backup Retention Period",
			},
			"incremental_retention_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Incremental Backup Retention Period",
			},
			"schedules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Backup Schedules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_frequency": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Backup Schedule Frequency (MONTHLY, WEEKLY, DAYS)",
						},
						"schedule_frequency_detail": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Backup Schedule Frequency details",
						},
						"schedule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup Schedule ID",
						},
						"schedule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup Schedule Name",
						},
						"schedule_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Backup Schedule Type (FULL, INCREMENTAL)",
						},
						"start_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Backup Start Time (format:HH:mm±hh:mm)",
						},
					},
				},
			},
			"service_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Service Zone ID",
			},
			"tags": tfTags.TagsSchema(),
		},
		Description: "Provides a Backup resource.",
	}
}

func createBackup(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	azCode := rd.Get("az_code").(string)
	backupDrZoneId := rd.Get("backup_dr_zone_id").(string)
	backupName := rd.Get("backup_name").(string)
	backupPolicyTypeCategory := rd.Get("backup_policy_type_category").(string)
	backupRepository := rd.Get("backup_repository").(string)
	drAzCode := rd.Get("dr_az_code").(string)
	//fileSystemBackupSelections := rd.Get("file_system_backup_selections").([]string)
	isBackupDrEnabled := rd.Get("is_backup_dr_enabled").(string)
	objectId := rd.Get("object_id").(string)
	objectType := rd.Get("object_type").(string)
	policyType := rd.Get("policy_type").(string)
	productNames := convertToStringArray(rd.Get("product_names").([]interface{}))
	retentionPeriod := rd.Get("retention_period").(string)
	incrementalRetentionPeriod := rd.Get("incremental_retention_period").(string)
	serviceZoneId := rd.Get("service_zone_id").(string)
	scheduleList := rd.Get("schedules").(common.HclListObject)
	scheduleInfoList := convertSchedules(scheduleList)

	request := backup.CreateBackupRequest{
		AzCode:                   azCode,
		BackupDrZoneId:           backupDrZoneId,
		BackupName:               backupName,
		BackupPolicyTypeCategory: backupPolicyTypeCategory,
		BackupRepository:         backupRepository,
		DrAzCode:                 drAzCode,
		//FileSystemBackupSelections: fileSystemBackupSelections,
		IsBackupDrEnabled:          isBackupDrEnabled,
		ObjectId:                   objectId,
		ObjectType:                 objectType,
		PolicyType:                 policyType,
		ProductNames:               productNames,
		RetentionPeriod:            retentionPeriod,
		IncrementalRetentionPeriod: incrementalRetentionPeriod,
		Schedules:                  scheduleInfoList,
		ServiceZoneId:              serviceZoneId,
		Tags:                       rd.Get("tags").(map[string]interface{}),
	}

	response, err := inst.Client.Backup.CreateBackup(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForBackupStatus(ctx, inst.Client, response.ResourceId, []string{}, []string{"ACTIVE"}, true)
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(response.ResourceId)

	return readBackup(ctx, rd, meta)
}

func readBackup(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)
	info, _, err := inst.Client.Backup.ReadBackup(ctx, rd.Id())
	if err != nil {
		rd.SetId("")

		//not show error message for deleted resource
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	rd.Set("az_code", info.AzCode)
	rd.Set("backup_dr_zone_id", info.BackupDrZoneId)
	rd.Set("backup_name", info.BackupName)
	rd.Set("backup_policy_type_category", info.BackupPolicyTypeCategory)
	rd.Set("backup_repository", info.BackupRepository)
	rd.Set("dr_az_code", info.DrAzCode)
	rd.Set("is_backup_dr_enabled", info.IsBackupDrEnabled)
	rd.Set("object_id", info.ObjectId)
	rd.Set("object_type", info.ObjectType)
	rd.Set("policy_type", info.PolicyType)
	rd.Set("is_backup_dr_deleted", info.IsBackupDrDeleted)
	/*productNames := common.HclListObject{}
	for _, schedule := range rd.ProductNames {
		scheduleIds = append(scheduleIds, schedule)
	}
	rd.Set("product_names")*/

	rd.Set("retention_period", info.RetentionPeriod)
	rd.Set("incremental_retention_period", info.IncrementalRetentionPeriod)
	rd.Set("service_zone_id", info.ServiceZoneId)
	rd.Set("backup_dr_id", info.BackupDrId)
	if _, ok := rd.GetOk("is_backup_dr_destroy_enabled"); !ok {
		rd.Set("is_backup_dr_destroy_enabled", false)
	}

	backupScheduleList, err := inst.Client.Backup.ReadBackupScheduleList(ctx, rd.Id(), backup2.BackupSearchOpenApiApiListSchedulesOpts{
		Page: optional.Int32{},
		Size: optional.Int32{},
		Sort: optional.Interface{},
	})
	if err != nil {
		rd.SetId("")

		//not show error message for deleted resource
		if common.IsDeleted(err) {
			return nil
		}
		return diag.FromErr(err)
	}

	backupSchedules := convertBackupScheduleResponseListToHclSetObject(backupScheduleList.Contents)
	rd.Set("schedules", backupSchedules)

	tfTags.SetTags(ctx, rd, meta, rd.Id())

	return nil
}

func convertBackupScheduleResponseListToHclSetObject(backupScheduleList []backup2.BackupSchedulesResponse) common.HclSetObject {
	var backupSchedules common.HclSetObject

	for _, schedule := range backupScheduleList {
		kv := common.HclKeyValueObject{
			"schedule_frequency":        schedule.ScheduleFrequency,
			"schedule_frequency_detail": schedule.ScheduleFrequencyDetail,
			"schedule_id":               schedule.ScheduleId,
			"schedule_name":             schedule.ScheduleName,
			"schedule_type":             schedule.ScheduleType,
			"start_time":                schedule.StartTime,
		}
		backupSchedules = append(backupSchedules, kv)
	}
	return backupSchedules
}

func updateBackup(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.HasChanges("schedules") || rd.HasChanges("retention_period") || rd.HasChanges("incremental_retention_period") {
		scheduleList := rd.Get("schedules").(common.HclListObject)
		retentionPeriod := rd.Get("retention_period").(string)
		incrementalRetentionPeriod := rd.Get("incremental_retention_period").(string)

		scheduleInfoList := convertSchedules(scheduleList)

		request := backup.UpdateBackupScheduleRequest{
			Schedules:                  scheduleInfoList,
			IncrementalRetentionPeriod: incrementalRetentionPeriod,
			RetentionPeriod:            retentionPeriod,
		}

		_, err := inst.Client.Backup.UpdateBackupSchedule(ctx, rd.Id(), request)
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForBackupStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if rd.HasChanges("is_backup_dr_enabled") && rd.Get("is_backup_dr_enabled").(string) == "N" {
		_, err := inst.Client.Backup.UpdateBackupDr(ctx, rd, rd.Get("backup_dr_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForBackupStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err := tfTags.UpdateTags(ctx, rd, meta, rd.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return readBackup(ctx, rd, meta)
}

func deleteBackup(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	if rd.Get("is_backup_dr_destroy_enabled").(bool) && strings.EqualFold(rd.Get("is_backup_dr_deleted").(string), "N") {
		_, err := inst.Client.Backup.DeleteBackupDr(ctx, rd.Get("backup_dr_id").(string))
		err = waitForBackupStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"ACTIVE"}, true)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	_, err := inst.Client.Backup.DeleteBackup(ctx, rd.Id())
	if err != nil && !common.IsDeleted(err) {
		return diag.FromErr(err)
	}

	err = waitForBackupStatus(ctx, inst.Client, rd.Id(), []string{}, []string{"DELETED"}, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForBackupStatus(ctx context.Context, scpClient *client.SCPClient, id string, pendingStates []string, targetStates []string, errorOnNotFound bool) error {
	return client.WaitForStatus(ctx, scpClient, pendingStates, targetStates, func() (interface{}, string, error) {
		info, c, err := scpClient.Backup.ReadBackup(ctx, id)
		if err != nil {
			if c == 404 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			if c == 403 && !errorOnNotFound {
				return "", "DELETED", nil
			}
			return nil, "", err
		}
		if info.BackupId != id {
			return nil, "", fmt.Errorf("invalid resource status")
		}
		return info, info.BackupState, nil
	})
}

func getTagRequestArray(rd *schema.ResourceData) []backup.TagRequest {
	tags := rd.Get("tags").([]interface{})
	tagsRequests := make([]backup.TagRequest, 0)
	for _, tag := range tags {
		tagMap := tag.(map[string]interface{})
		tagsRequests = append(tagsRequests, backup.TagRequest{
			TagKey:   tagMap["tag_key"].(string),
			TagValue: tagMap["tag_value"].(string),
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

func convertSchedules(schedules common.HclListObject) []backup.BackupScheduleInfo {
	var result []backup.BackupScheduleInfo
	for _, itemObject := range schedules {
		item := itemObject.(common.HclKeyValueObject)
		var info backup.BackupScheduleInfo
		if v, ok := item["schedule_frequency"]; ok {
			info.ScheduleFrequency = v.(string)
		}
		if v, ok := item["schedule_frequency_detail"]; ok {
			info.ScheduleFrequencyDetail = v.(string)
		}
		if v, ok := item["schedule_type"]; ok {
			info.ScheduleType = v.(string)
		}
		if v, ok := item["start_time"]; ok {
			info.StartTime = v.(string)
		}
		result = append(result, info)
	}
	return result
}
