package mysql

import (
	"context"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	"github.com/ScpDevTerra/trf-sdk/library/mysql2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *mysql2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: mysql2.NewAPIClient(config),
	}
}

func (client *Client) CreateMariadb(ctx context.Context, request mysql2.CreateMySqlRequest) (mysql2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MySqlConfigurationControllerApi.CreateMysql(ctx, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListMariadb(ctx context.Context, dbName string, serverGroupName string, virtualServerName string) (mysql2.ListResponseOfDbServerGroupsResponse, error) {
	result, _, err := client.sdkClient.MySqlConfigurationControllerApi.ListMysql(ctx, dbName, serverGroupName, virtualServerName, &mysql2.MySqlConfigurationControllerApiListMysqlOpts{
		CreatedBy: optional.String{},
		Page:      optional.NewInt32(0),
		Size:      optional.NewInt32(1000),
		Sort:      optional.Interface{},
	})
	return result, err
}

func (client *Client) GetMariadb(ctx context.Context, dbServerGroupId string) (mysql2.DetailDatabaseResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DetailDatabase5(ctx, dbServerGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateMariadbBlockSize(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageId string, blockStorageSize int) (mysql2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseStorage5(ctx, dbServerGroupId, mysql2.ResizeStorageRequest{
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

func (client *Client) AddPostgresqlBlock(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageType string, blockStorageSize int) (mysql2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.AddDatabaseStorage4(ctx, dbServerGroupId, mysql2.AddStorageRequest{
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

func (client *Client) UpdatePostgresqlScale(ctx context.Context, dbServerGroupId string, virtualServerId string, scaleId string) (mysql2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseScale5(ctx, dbServerGroupId, mysql2.ResizeScaleRequest{
		VirtualServerId: virtualServerId,
		ScaleProductId:  scaleId,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
