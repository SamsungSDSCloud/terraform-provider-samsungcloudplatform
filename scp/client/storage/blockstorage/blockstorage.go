package blockstorage

import (
	"context"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/client"
	blockstorage2 "github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/library/block-storage2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *blockstorage2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: blockstorage2.NewAPIClient(config),
	}
}

func (client *Client) CreateBlockStorage(ctx context.Context, request CreateBlockStorageRequest) (blockstorage2.AsyncResponse, error) {
	result, _, err := client.sdkClient.BlockStorageControllerApi.CreateBlockStorage(
		ctx,
		client.config.ProjectId,
		blockstorage2.BlockStorageCreateRequest{
			BlockStorageName: request.BlockStorageName,
			BlockStorageSize: request.BlockStorageSize,
			EncryptEnabled:   request.EncryptEnabled,
			ProductId:        request.ProductId,
			VirtualServerId:  request.VirtualServerId,
		})

	return result, err
}

func (client *Client) ReadBlockStorage(ctx context.Context, blockStorageId string) (blockstorage2.BlockStorageResponse, int, error) {
	result, c, err := client.sdkClient.BlockStorageControllerApi.DetailBlockStorage(ctx, client.config.ProjectId, blockStorageId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadBlockStorageList(ctx context.Context, request ReadBlockStorageRequest) (blockstorage2.ListResponseOfBlockStorageResponse, error) {
	result, _, err := client.sdkClient.BlockStorageControllerApi.ListBlockStorages(
		ctx,
		client.config.ProjectId,
		&blockstorage2.BlockStorageControllerApiListBlockStoragesOpts{
			BlockStorageName:  optional.NewString(request.BlockStorageName),
			VirtualServerId:   optional.NewString(request.VirtualServerId),
			VirtualServerName: optional.NewString(request.BlockStorageName),
			CreatedBy:         optional.NewString(request.CreatedBy),
			Page:              optional.NewInt32(request.Page),
			Size:              optional.NewInt32(request.Size),
			//Sort: optional.NewInterface([]string{"modifiedDt:asc"}),
		})

	return result, err
}

func (client *Client) ResizeBlockStorage(ctx context.Context, request UpdateBlockStorageRequest) (blockstorage2.AsyncResponse, error) {
	result, _, err := client.sdkClient.BlockStorageControllerApi.ResizeBlockStorage(
		ctx,
		client.config.ProjectId,
		request.BlockStorageId,
		blockstorage2.BlockStorageResizeRequest{
			BlockStorageSize: request.BlockStorageSize,
			ProductId:        request.ProductId,
		})

	return result, err
}

func (client *Client) DeleteBlockStorage(ctx context.Context, blockStorageId string) (blockstorage2.AsyncResponse, error) {
	result, _, err := client.sdkClient.BlockStorageControllerApi.DeleteBlockStorage(ctx, client.config.ProjectId, blockStorageId)
	return result, err
}
