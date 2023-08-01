package backup

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/backup2"
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

	tags := make([]backup2.TagRequest, 0)
	for _, tag := range request.Tags {
		tags = append(tags, backup2.TagRequest{
			TagKey:   tag.TagKey,
			TagValue: tag.TagValue,
		})
	}

	result, _, err := client.sdkClient.BackupOpenApiV3Api.CreateBackup1(ctx, client.config.ProjectId,
		backup2.BackupCreateV3Request{
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
			Schedules:                  backupSchedules,
			ServiceZoneId:              request.ServiceZoneId,
			Tags:                       tags,
		})
	return result, err
}

func (client *Client) ReadBackup(ctx context.Context, backupId string) (backup2.DetailBackupResponse, int, error) {
	result, c, err := client.sdkClient.BackupSearchOpenApiV3Api.DetailBackup1(ctx, client.config.ProjectId, backupId)
	return result, c.StatusCode, err
}

func (client *Client) ReadBackupList(ctx context.Context, request backup2.BackupSearchOpenApiV2ApiListBackupsOpts) (backup2.ListResponseOfBackupResponse, error) {
	result, _, err := client.sdkClient.BackupSearchOpenApiV2Api.ListBackups(ctx, client.config.ProjectId, &request)
	return result, err
}

func (client *Client) DeleteBackup(ctx context.Context, backupId string) (backup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.BackupOpenApiV2Api.DeleteBackup(ctx, client.config.ProjectId, backupId)
	return result, err
}

func (client *Client) ReadBackupScheduleList(ctx context.Context, backupId string, request backup2.BackupSearchOpenApiV2ApiListSchedulesOpts) (backup2.ListResponseOfBackupSchedulesResponse, error) {
	result, _, err := client.sdkClient.BackupSearchOpenApiV2Api.ListSchedules(ctx, client.config.ProjectId, backupId, &request)
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

	result, _, err := client.sdkClient.BackupOpenApiV2Api.UpdateBackupSchedule(
		ctx,
		client.config.ProjectId,
		backupId,
		backup2.BackupScheduleUpdateRequest{
			Schedules: backupSchedules,
		})
	return result, err
}
