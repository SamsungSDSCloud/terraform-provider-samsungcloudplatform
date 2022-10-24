package filestorage

import (
	"context"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	filestorage2 "github.com/ScpDevTerra/trf-sdk/library/file-storage2"
	"github.com/antihax/optional"
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
	fs := &filestorage2.FileStorageOpenApiV2ApiCheckFileStorageDuplicationV2Opts{
		FileStorageName: optional.NewString(request.FileStorageName),
	}

	if len(request.CifsId) != 0 {
		fs.CifsId = optional.NewString(request.CifsId)
	}

	result, _, err := client.sdkClient.FileStorageOpenApiV3Api.CheckFileStorageDuplicationV3(ctx, client.config.ProjectId, client.config.UserId, client.config.Email, client.config.LoginId, "", "", request.FileStorageName)

	return result, err
}

func (client *Client) CreateFileStorage(ctx context.Context, request CreateFileStorageRequest) (filestorage2.AsyncResponse, error) {
	result, _, err := client.sdkClient.FileStorageOpenApiV3Api.CreateFileStorageV3(
		ctx,
		client.config.ProjectId,
		client.config.UserId,
		client.config.LoginId,
		client.config.Email,
		"",
		"",
		filestorage2.CreateFileStorageRequest{
			//CifsId:                request.CifsId,
			CifsPassword:      request.CifsPassword,
			EncryptionEnabled: request.IsEncrypted,
			//FileStorageCapacityGb: request.FileStorageCapacityGb,
			FileStorageName:     request.FileStorageName,
			FileStorageProtocol: request.FileStorageProtocol,
			ProductGroupId:      request.ProductGroupId,
			ProductIds:          request.ProductIds,
			ServiceZoneId:       request.ServiceZoneId,
			//SnapshotCapacityRate:  request.SnapshotCapacityRate,
			DiskType:         request.DiskType,
			SnapshotSchedule: (*filestorage2.SnapshotSchedule)(request.SnapshotSchedule),
		})
	return result, err
}

func (client *Client) ReadFileStorage(ctx context.Context, fileStorageId string) (filestorage2.FileStorageDetailResponseV3, int, error) {
	result, c, err := client.sdkClient.FileStorageOpenApiV3Api.DetailFileStorageV3(ctx, client.config.ProjectId, client.config.UserId, client.config.LoginId, client.config.Email, "", "", fileStorageId)
	return result, c.StatusCode, err
}

func (client *Client) ReadFileStorageList(ctx context.Context, request ReadFileStorageRequest) (filestorage2.ListResponseOfFileStoragesResponse, error) {
	localVarOptionals := &filestorage2.FileStorageOpenApiV3ApiListFileStoragesV31Opts{
		FileStorageId:   optional.NewString(request.FileStorageId),
		FileStorageName: optional.NewString(request.FileStorageName),
		ServiceZoneId:   optional.NewString(request.ServiceZoneId),
		CreatedBy:       optional.NewString(request.CreatedBy),
		Page:            optional.NewInt32(request.Page),
		Size:            optional.NewInt32(request.Size),
		//Sort:            optional.NewString("modifiedDt:asc"),
	}

	if len(request.FileStorageProtocol) > 0 {
		localVarOptionals.FileStorageProtocol = optional.NewString(request.FileStorageProtocol)
	}

	result, _, err := client.sdkClient.FileStorageOpenApiV3Api.ListFileStoragesV31(
		ctx,
		client.config.ProjectId,
		client.config.UserId,
		client.config.LoginId,
		client.config.Email,
		"",
		"",
		localVarOptionals)

	return result, err
}

/* no api(cannot control file storage capacity GB in v3 api)
func (client *Client) IncreaseFileStorage(ctx context.Context, request UpdateFileStorageRequest) (filestorage2.AsyncResponse, error) {
	result, _, err := client.sdkClient.FileStorageOpenApiV2Api.IncreaseFileStorageCapacityV2(
		ctx,
		client.config.ProjectId,
		request.FileStorageId,
		filestorage2.IncreaseFileStorageCapacityRequest{
			FileStorageCapacityGb: request.FileStorageCapacityGb,
		})

	return result, err
}
*/

func (client *Client) DeleteFileStorage(ctx context.Context, fileStorageId string) (filestorage2.AsyncResponse, error) {
	result, _, err := client.sdkClient.FileStorageOpenApiV3Api.DeleteFileStorageV3(ctx, client.config.ProjectId, client.config.UserId, client.config.LoginId, client.config.Email, "", "", fileStorageId)
	return result, err
}
