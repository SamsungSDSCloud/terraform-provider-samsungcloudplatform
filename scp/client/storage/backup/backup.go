package backup

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/backup2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *backup2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: backup2.NewAPIClient(config),
	}
}

func (client *Client) CreateBackup(ctx context.Context, request CreateBackupRequest) (backup2.AsyncResponse, error) {
	var backupSchedules []backup2.BackupScheduleInfo
	for _, backupSchedule := range request.Schedules {
		backupSchedules = append(backupSchedules, backup2.BackupScheduleInfo{
			ScheduleFrequency:       backupSchedule.ScheduleFrequency,
			ScheduleFrequencyDetail: backupSchedule.ScheduleFrequencyDetail,
			ScheduleType:            backupSchedule.ScheduleType,
			StartTime:               backupSchedule.StartTime,
		})
	}

	result, _, err := client.sdkClient.BackupOperateOpenApiApi.CreateBackup2(ctx, client.config.ProjectId,
		backup2.BackupCreateV5Request{
			AzCode:                     request.AzCode,
			BackupDrZoneId:             request.BackupDrZoneId,
			BackupName:                 request.BackupName,
			BackupPolicyTypeCategory:   request.BackupPolicyTypeCategory,
			BackupRepository:           request.BackupRepository,
			DrAzCode:                   request.DrAzCode,
			FileSystemBackupSelections: request.FileSystemBackupSelections,
			IsBackupDrEnabled:          request.IsBackupDrEnabled,
			ObjectId:                   request.ObjectId,
			ObjectType:                 request.ObjectType,
			PolicyType:                 request.PolicyType,
			ProductNames:               request.ProductNames,
			RetentionPeriod:            request.RetentionPeriod,
			IncrementalRetentionPeriod: request.IncrementalRetentionPeriod,
			Schedules:                  backupSchedules,
			ServiceZoneId:              request.ServiceZoneId,
			Tags:                       client.sdkClient.ToTagRequestList(request.Tags),
		})
	return result, err
}

func (client *Client) ReadBackup(ctx context.Context, backupId string) (backup2.DetailBackupV4Response, int, error) {
	result, c, err := client.sdkClient.BackupSearchOpenApiApi.DetailBackup(ctx, client.config.ProjectId, backupId)
	return result, c.StatusCode, err
}

func (client *Client) UpdateBackupDr(ctx context.Context, rd *schema.ResourceData, backupDrId string) (backup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.BackupDrOpenApiV2Api.CancelBackupDrRelationship(ctx, client.config.ProjectId, rd.Get("backup_dr_id").(string))
	return result, err
}
func (client *Client) ReadBackupList(ctx context.Context, request backup2.BackupSearchOpenApiApiListBackupsOpts) (backup2.ListResponseBackupV3Response, error) {
	result, _, err := client.sdkClient.BackupSearchOpenApiApi.ListBackups(ctx, client.config.ProjectId, &request)
	return result, err
}

func (client *Client) DeleteBackup(ctx context.Context, backupId string) (backup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.BackupOperateOpenApiApi.DeleteBackup(ctx, client.config.ProjectId, backupId)
	return result, err
}

func (client *Client) ReadBackupScheduleList(ctx context.Context, backupId string, request backup2.BackupSearchOpenApiApiListSchedulesOpts) (backup2.ListResponseBackupSchedulesResponse, error) {
	result, _, err := client.sdkClient.BackupSearchOpenApiApi.ListSchedules(ctx, client.config.ProjectId, backupId, &request)
	return result, err
}

func (client *Client) UpdateBackupSchedule(ctx context.Context, backupId string, request UpdateBackupScheduleRequest) (backup2.AsyncResponse, error) {
	var backupSchedules []backup2.BackupScheduleInfo
	for _, backupSchedule := range request.Schedules {
		backupSchedules = append(backupSchedules, backup2.BackupScheduleInfo{
			ScheduleFrequency:       backupSchedule.ScheduleFrequency,
			ScheduleFrequencyDetail: backupSchedule.ScheduleFrequencyDetail,
			ScheduleType:            backupSchedule.ScheduleType,
			StartTime:               backupSchedule.StartTime,
		})
	}

	result, _, err := client.sdkClient.BackupOperateOpenApiApi.UpdateBackupSchedule1(
		ctx,
		client.config.ProjectId,
		backupId,
		backup2.BackupScheduleUpdateV3Request{
			RetentionPeriod:            request.RetentionPeriod,
			IncrementalRetentionPeriod: request.IncrementalRetentionPeriod,
			Schedules:                  backupSchedules,
		})
	return result, err
}

func (client *Client) DeleteBackupDr(ctx context.Context, backupDrId string) (backup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.BackupDrOpenApiV2Api.DeleteBackupDr(ctx, client.config.ProjectId, backupDrId)
	return result, err
}
