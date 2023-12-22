package bmblockstorage

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	baremetalblockstorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/bare-metal-block-storage"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *baremetalblockstorage.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: baremetalblockstorage.NewAPIClient(config),
	}
}

func (client *Client) GetBareMetalBlockStorageDetail(ctx context.Context, blockStorageId string) (baremetalblockstorage.BmBlockStorageDetailResponse, int, error) {
	result, c, err := client.sdkClient.BmBlockStorageControllerApi.DetailBareMetalBlockStorage(ctx, client.config.ProjectId, blockStorageId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetBareMetalBlockStorages(ctx context.Context) (baremetalblockstorage.ListResponseOfBmBlockStorageResponse, int, error) {
	result, c, err := client.sdkClient.BmBlockStorageControllerApi.ListBareMetalBlockStorages(ctx, client.config.ProjectId, &baremetalblockstorage.BmBlockStorageControllerApiListBareMetalBlockStoragesOpts{
		Page: optional.NewInt32(0),
		Size: optional.NewInt32(10000),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateBareMetalBlockStorage(ctx context.Context, request BmBlockStorageCreateRequest) (baremetalblockstorage.AsyncResponse, int, error) {
	snapshotSchedule := baremetalblockstorage.SnapshotSchedule{
		DayOfWeek: request.SnapshotSchedule.DayOfWeek,
		Frequency: request.SnapshotSchedule.Frequency,
		Hour:      request.SnapshotSchedule.Hour,
	}
	result, c, err := client.sdkClient.BmBlockStorageControllerApi.CreateBareMetalBlockStorage(ctx, client.config.ProjectId, baremetalblockstorage.BmBlockStorageCreateRequest{
		BareMetalBlockStorageName: request.BareMetalBlockStorageName,
		BareMetalBlockStorageSize: request.BareMetalBlockStorageSize,
		BareMetalServerIds:        request.BareMetalServerIds,
		EncryptionEnabled:         &request.EncryptionEnabled,
		IsSnapshotPolicy:          &request.IsSnapshotPolicy,
		ProductId:                 request.ProductId,
		ServiceZoneId:             request.ServiceZoneId,
		SnapshotCapacityRate:      request.SnapshotCapacityRate,
		SnapshotSchedule:          &snapshotSchedule,
		Tags:                      client.sdkClient.ToTagRequestList(request.Tags),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AttachBareMetalBlockStorage(ctx context.Context, storageId string, serverIds []string) (baremetalblockstorage.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.BmBlockStorageControllerApi.AttachBareMetalBlockStorage(ctx, client.config.ProjectId, storageId, baremetalblockstorage.BmBlockStorageAttachRequest{
		BareMetalServerIds: serverIds,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachBareMetalBlockStorage(ctx context.Context, storageId string, serverIds []string) (baremetalblockstorage.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.BmBlockStorageControllerApi.DetachBareMetalBlockStorage(ctx, client.config.ProjectId, storageId, baremetalblockstorage.BmBlockStorageDetachRequest{
		BareMetalServerIds: serverIds,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteBareMetalBlockStorage(ctx context.Context, storageId string) (baremetalblockstorage.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.BmBlockStorageControllerApi.TerminatedBareMetalBlockStorage(ctx, client.config.ProjectId, storageId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetBareMetalBlockStorageScheduleList(ctx context.Context, blockStorageId string) (baremetalblockstorage.ListResponseOfBmBlockStorageSnapshotScheduleResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalBlockStorageSnapshotScheduleOpenApiV1Api.ListBareMetalBlockStorageSnapshotScheduleV1(ctx, client.config.ProjectId, blockStorageId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateBareMetalBlockStorageSchedule(ctx context.Context, blockStorageId string, schedule SnapshotSchedule) (baremetalblockstorage.BmBlockStorageSnapshotScheduleResponse, int, error) {
	snapshotSchedule := baremetalblockstorage.SnapshotSchedule{
		DayOfWeek: schedule.DayOfWeek,
		Frequency: schedule.Frequency,
		Hour:      schedule.Hour,
	}
	result, c, err := client.sdkClient.BareMetalBlockStorageSnapshotScheduleOpenApiV1Api.CreateBareMetalBlockStorageSnapshotScheduleV1(ctx, client.config.ProjectId, blockStorageId, baremetalblockstorage.BmBlockStorageSnapshotScheduleRequest{
		SnapshotSchedule: &snapshotSchedule,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateBareMetalBlockStorageSchedule(ctx context.Context, blockStorageId string, schedule SnapshotSchedule) (baremetalblockstorage.BmBlockStorageSnapshotScheduleResponse, int, error) {
	snapshotSchedule := baremetalblockstorage.SnapshotSchedule{
		DayOfWeek: schedule.DayOfWeek,
		Frequency: schedule.Frequency,
		Hour:      schedule.Hour,
	}
	result, c, err := client.sdkClient.BareMetalBlockStorageSnapshotScheduleOpenApiV1Api.UpdateBareMetalBlockStorageSnapshotScheduleV1(ctx, client.config.ProjectId, blockStorageId, baremetalblockstorage.BmBlockStorageSnapshotScheduleRequest{
		SnapshotSchedule: &snapshotSchedule,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteBareMetalBlockStorageSchedule(ctx context.Context, blockStorageId string) (baremetalblockstorage.BmBlockStorageSnapshotScheduleResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalBlockStorageSnapshotScheduleOpenApiV1Api.DeleteBareMetalBlockStorageSnapshotScheduleV1(ctx, client.config.ProjectId, blockStorageId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetBareMetalBlockStorageSnapshotList(ctx context.Context, blockStorageId string) (baremetalblockstorage.ListResponseOfBmBlockStorageSnapshotsResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalBlockStorageOpenApiV1Api.ListBareMetalBlockStorageSnapshotsV1(ctx, client.config.ProjectId, blockStorageId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateBareMetalBlockStorageSnapshot(ctx context.Context, blockStorageId string) (baremetalblockstorage.BmBlockStorageSnapshotCreateResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalBlockStorageOpenApiV1Api.CreateBareMetalBlockStorageSnapshot(ctx, client.config.ProjectId, blockStorageId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateBareMetalBlockStorageSnapshotAttribute(ctx context.Context, blockStorageId string, isSnapshotPolicy string, snapshotCapacityRate int32) (baremetalblockstorage.BmBlockStorageSnapshotsAttributeResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalBlockStorageOpenApiV1Api.CreateBareMetalBlockStorageSnapshotAttribute(ctx, client.config.ProjectId, blockStorageId, baremetalblockstorage.BmBlockStorageSnapshotAttributeRequest{
		IsSnapshotPolicy:     isSnapshotPolicy,
		SnapshotCapacityRate: snapshotCapacityRate,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) RestoreBareMetalBlockStorageSnapshot(ctx context.Context, blockStorageId string, snapshotId string) (baremetalblockstorage.BmBlockStorageSnapshotRestoreResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalBlockStorageOpenApiV1Api.RestoreBareMetalBlockStorageSnapshotV1(ctx, client.config.ProjectId, blockStorageId, snapshotId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateBareMetalBlockStorageSnapshotAttribute(ctx context.Context, blockStorageId string, isSnapshotPolicy string, snapshotCapacityRate int32) (baremetalblockstorage.BmBlockStorageSnapshotsAttributeResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalBlockStorageOpenApiV1Api.UpdateBareMetalBlockStorageSnapshotAttribute(ctx, client.config.ProjectId, blockStorageId, baremetalblockstorage.BmBlockStorageSnapshotAttributeRequest{
		IsSnapshotPolicy:     isSnapshotPolicy,
		SnapshotCapacityRate: snapshotCapacityRate,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteBareMetalBlockStorageSnapshot(ctx context.Context, blockStorageId string, snapshotId string) (baremetalblockstorage.BmBlockStorageSnapshotDeleteResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalBlockStorageOpenApiV1Api.DeleteBareMetalBlockStorageSnapshotV1(ctx, client.config.ProjectId, blockStorageId, snapshotId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteBareMetalBlockStorageSnapshotAttribute(ctx context.Context, blockStorageId string, isSnapshotPolicy string) (baremetalblockstorage.BmBlockStorageSnapshotsAttributeResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalBlockStorageOpenApiV1Api.DeleteBareMetalBlockStorageSnapshotAttribute(ctx, client.config.ProjectId, blockStorageId, baremetalblockstorage.BmBlockStorageSnapshotAttributeDeleteRequest{
		IsSnapshotPolicy: isSnapshotPolicy,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
