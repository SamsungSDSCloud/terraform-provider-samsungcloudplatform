package epas

import (
	"context"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	"github.com/ScpDevTerra/trf-sdk/library/epas2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *epas2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: epas2.NewAPIClient(config),
	}
}

func (client *Client) CreateEpas(ctx context.Context, request epas2.CreateEpasRequest) (epas2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasConfigurationControllerApi.CreateEpas(ctx, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListEpas(ctx context.Context, dbName string, serverGroupName string, virtualServerName string) (epas2.ListResponseOfDbServerGroupsResponse, error) {
	result, _, err := client.sdkClient.EpasConfigurationControllerApi.ListEpas(ctx, dbName, serverGroupName, virtualServerName, &epas2.EpasConfigurationControllerApiListEpasOpts{
		CreatedBy: optional.String{},
		Page:      optional.NewInt32(0),
		Size:      optional.NewInt32(1000),
		Sort:      optional.Interface{},
	})
	return result, err
}

func (client *Client) DeleteEpas(ctx context.Context, serverGroupId string) (epas2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DeleteDatabase1(ctx, serverGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetEpas(ctx context.Context, dbServerGroupId string) (epas2.DetailDatabaseResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DetailDatabase1(ctx, dbServerGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AddEpasBlock(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageType string, blockStorageSize int) (epas2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.AddDatabaseStorage1(ctx, dbServerGroupId, epas2.AddStorageRequest{
		VirtualServerId:  virtualServerId,
		BlockStorageType: blockStorageType,
		BlockStorageSize: int32(blockStorageSize),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateEpasBlockSize(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageId string, blockStorageSize int) (epas2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseStorage1(ctx, dbServerGroupId, epas2.ResizeStorageRequest{
		VirtualServerId:  virtualServerId,
		BlockStorageId:   blockStorageId,
		BlockStorageSize: int32(blockStorageSize),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateEpasScale(ctx context.Context, dbServerGroupId string, virtualServerId string, scaleId string) (epas2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseScale1(ctx, dbServerGroupId, epas2.ResizeScaleRequest{
		VirtualServerId: virtualServerId,
		ScaleProductId:  scaleId,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
