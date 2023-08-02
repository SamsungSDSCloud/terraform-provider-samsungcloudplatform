package filestorage

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/client"
	filestorage2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/file-storage2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *filestorage2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: filestorage2.NewAPIClient(config),
	}
}

func (client *Client) CheckFileStorage(ctx context.Context, request CheckFileStorageRequest) (filestorage2.CheckResponse, error) {
	result, _, err := client.sdkClient.FileStorageOpenApiV3Api.CheckFileStorageDuplication(ctx, client.config.ProjectId, request.FileStorageName)
	if err != nil {
		return filestorage2.CheckResponse{}, err
	}

	return result, err
}

func (client *Client) CreateFileStorage(ctx context.Context, request CreateFileStorageRequest) (filestorage2.AsyncResponse, error) {
	tags := make([]filestorage2.TagRequest, 0)
	for _, tag := range request.Tags {
		tags = append(tags, filestorage2.TagRequest{
			TagKey:   tag.TagKey,
			TagValue: tag.TagValue,
		})
	}

	result, _, err := client.sdkClient.FileStorageOpenApiV4Api.CreateFileStorageV4(ctx, client.config.ProjectId,
		filestorage2.CreateFileStorageV4Request{
			CifsPassword:          request.CifsPassword,
			DiskType:              request.DiskType,
			FileStorageName:       request.FileStorageName,
			FileStorageProtocol:   request.FileStorageProtocol,
			MultiAvailabilityZone: request.MultiAvailabilityZone,
			ProductNames:          request.ProductNames,
			ServiceZoneId:         request.ServiceZoneId,
			//SnapshotRetentionCount: request.SnapshotRetentionCount,
			//SnapshotSchedule:       (*filestorage2.SnapshotSchedule)(request.SnapshotSchedule),
			Tags: tags,
		})
	return result, err
}

func (client *Client) ReadFileStorage(ctx context.Context, fileStorageId string) (filestorage2.FileStorageDetailResponse, int, error) {
	result, c, err := client.sdkClient.FileStorageOpenApiV3Api.DetailFileStorage(ctx, client.config.ProjectId, fileStorageId)
	return result, c.StatusCode, err
}

func (client *Client) ReadFileStorageList(ctx context.Context, request filestorage2.FileStorageOpenApiV3ApiListFileStoragesOpts) (filestorage2.ListResponseOfFileStoragesResponse, error) {
	result, _, err := client.sdkClient.FileStorageOpenApiV3Api.ListFileStorages(ctx, client.config.ProjectId, &request)
	return result, err
}

func (client *Client) DeleteFileStorage(ctx context.Context, fileStorageId string) (filestorage2.AsyncResponse, error) {
	result, _, err := client.sdkClient.FileStorageOpenApiV3Api.DeleteFileStorage(ctx, client.config.ProjectId, fileStorageId)
	return result, err
}

func (client *Client) CreateFileStorageSnapshotSchedule(ctx context.Context, fileStorageId string, retentionCount int32, schedule *SnapshotSchedule) (filestorage2.AsyncResponse, error) {
	result, _, err := client.sdkClient.FileStorageOpenApiV3Api.CreateFileStorageSnapshotSchedule(ctx, client.config.ProjectId, fileStorageId, filestorage2.FsSnapshotScheduleRequest{
		SnapshotRetentionCount: &retentionCount,
		SnapshotSchedule:       (*filestorage2.SnapshotSchedule)(schedule),
	})

	return result, err
}

func (client *Client) ReadFileStorageSnapshotSchedule(ctx context.Context, fileStorageId string) (filestorage2.FileStorageSnapshotScheduleResponse, int, error) {
	result, c, err := client.sdkClient.FileStorageOpenApiV3Api.SearchFileStorageSnapshotSchedule(ctx, client.config.ProjectId, fileStorageId)
	return result, c.StatusCode, err
}

func (client *Client) UpdateFileStorageSnapshotSchedule(ctx context.Context, fileStorageId string, retentionCount int32, schedule *SnapshotSchedule) (filestorage2.AsyncResponse, error) {
	result, _, err := client.sdkClient.FileStorageOpenApiV3Api.UpdateFileStorageSnapshotSchedule(
		ctx,
		client.config.ProjectId,
		fileStorageId,
		filestorage2.FsSnapshotScheduleRequest{
			SnapshotRetentionCount: &retentionCount,
			SnapshotSchedule:       (*filestorage2.SnapshotSchedule)(schedule),
		})
	return result, err
}

func (client *Client) DeleteFileStorageSnapshotSchedule(ctx context.Context, fileStorageId string) (filestorage2.AsyncResponse, error) {
	result, _, err := client.sdkClient.FileStorageOpenApiV3Api.DeleteFileStorageSnapshotSchedule(ctx, client.config.ProjectId, fileStorageId)
	return result, err
}
