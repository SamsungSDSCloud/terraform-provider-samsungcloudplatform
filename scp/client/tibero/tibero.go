package tibero

import (
	"context"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	"github.com/ScpDevTerra/trf-sdk/library/tibero2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *tibero2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: tibero2.NewAPIClient(config),
	}
}

func (client *Client) CreateTibero(ctx context.Context, request tibero2.CreateTiberoRequest) (tibero2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.TiberoConfigurationControllerApi.CreateTibero(ctx, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListTibero(ctx context.Context, dbName string, serverGroupName string, virtualServerName string) (tibero2.ListResponseOfDbServerGroupsResponse, error) {
	result, _, err := client.sdkClient.TiberoConfigurationControllerApi.ListTibero(ctx, dbName, serverGroupName, virtualServerName, &tibero2.TiberoConfigurationControllerApiListTiberoOpts{
		CreatedBy: optional.String{},
		Page:      optional.NewInt32(0),
		Size:      optional.NewInt32(1000),
		Sort:      optional.Interface{},
	})
	return result, err
}

func (client *Client) DeleteTibero(ctx context.Context, serverGroupId string) (tibero2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DeleteDatabase11(ctx, serverGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetTibero(ctx context.Context, dbServerGroupId string) (tibero2.DetailDatabaseResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DetailDatabase9(ctx, dbServerGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AddTiberoBlock(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageType string, blockStorageSize int) (tibero2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.AddDatabaseStorage6(ctx, dbServerGroupId, tibero2.AddStorageRequest{
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

func (client *Client) UpdateTiberoBlockSize(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageId string, blockStorageSize int) (tibero2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseScale8(ctx, dbServerGroupId, tibero2.ResizeScaleRequest{
		VirtualServerId: virtualServerId,
		ScaleProductId:  blockStorageId, //need to change
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateTiberoScale(ctx context.Context, dbServerGroupId string, virtualServerId string, scaleId string) (tibero2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseScale8(ctx, dbServerGroupId, tibero2.ResizeScaleRequest{
		VirtualServerId: virtualServerId,
		ScaleProductId:  scaleId,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
