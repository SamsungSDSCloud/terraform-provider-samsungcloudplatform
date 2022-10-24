package postgresql

import (
	"context"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	"github.com/ScpDevTerra/trf-sdk/library/postgresql2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *postgresql2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: postgresql2.NewAPIClient(config),
	}
}

func (client *Client) CreatePostgresql(ctx context.Context, request postgresql2.CreatePostgreSqlRequest) (postgresql2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgreSqlConfigurationControllerApi.CreatePostgresql(ctx, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListPostgresql(ctx context.Context, dbName string, serverGroupName string, virtualServerName string) (postgresql2.ListResponseOfDbServerGroupsResponse, error) {
	result, _, err := client.sdkClient.PostgreSqlConfigurationControllerApi.ListPostgresql(ctx, dbName, serverGroupName, virtualServerName, &postgresql2.PostgreSqlConfigurationControllerApiListPostgresqlOpts{
		CreatedBy: optional.String{},
		Page:      optional.NewInt32(0),
		Size:      optional.NewInt32(1000),
		Sort:      optional.Interface{},
	})
	return result, err
}

func (client *Client) DeletePostgresql(ctx context.Context, serverGroupId string) (postgresql2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DeleteDatabase8(ctx, serverGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetPostgresql(ctx context.Context, dbServerGroupId string) (postgresql2.DetailDatabaseResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DetailDatabase6(ctx, dbServerGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AddPostgresqlBlock(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageType string, blockStorageSize int) (postgresql2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.AddDatabaseStorage5(ctx, dbServerGroupId, postgresql2.AddStorageRequest{
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

func (client *Client) UpdatePostgresqlBlockSize(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageId string, blockStorageSize int) (postgresql2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseStorage6(ctx, dbServerGroupId, postgresql2.ResizeStorageRequest{
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

func (client *Client) UpdatePostgresqlScale(ctx context.Context, dbServerGroupId string, virtualServerId string, scaleId string) (postgresql2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseScale6(ctx, dbServerGroupId, postgresql2.ResizeScaleRequest{
		VirtualServerId: virtualServerId,
		ScaleProductId:  scaleId,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
