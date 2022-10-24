package sqlserver

import (
	"context"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	"github.com/ScpDevTerra/trf-sdk/library/sqlserver2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *sqlserver2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: sqlserver2.NewAPIClient(config),
	}
}

func (client *Client) CreateSqlServer(ctx context.Context, request sqlserver2.CreateSqlServerRequest) (sqlserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MsSqlConfigurationControllerApi.CreateSqlserver(ctx, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListSqlServer(ctx context.Context, dbName string, serverGroupName string, virtualServerName string) (sqlserver2.ListResponseOfDbServerGroupsResponse, error) {
	result, _, err := client.sdkClient.MsSqlConfigurationControllerApi.ListSqlserver(ctx, dbName, serverGroupName, virtualServerName, &sqlserver2.MsSqlConfigurationControllerApiListSqlserverOpts{
		CreatedBy: optional.String{},
		Page:      optional.NewInt32(0),
		Size:      optional.NewInt32(1000),
		Sort:      optional.Interface{},
	})
	return result, err
}

func (client *Client) DeleteSqlServer(ctx context.Context, serverGroupId string) (sqlserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DeleteDatabase10(ctx, serverGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetSqlServer(ctx context.Context, dbServerGroupId string) (sqlserver2.DetailDatabaseResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DetailDatabase8(ctx, dbServerGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateSqlServerBlockSize(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageId string, blockStorageSize int) (sqlserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseStorage8(ctx, dbServerGroupId, sqlserver2.ResizeStorageRequest{
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

func (client *Client) UpdateSqlServerScale(ctx context.Context, dbServerGroupId string, virtualServerId string, scaleId string) (sqlserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseStorage8(ctx, dbServerGroupId, sqlserver2.ResizeStorageRequest{
		VirtualServerId:  virtualServerId,
		BlockStorageId:   scaleId, //need to change
		BlockStorageSize: 10,      //need to change
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
