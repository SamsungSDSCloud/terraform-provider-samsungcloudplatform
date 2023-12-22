package sqlserver

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/sqlserver2"
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

func (client *Client) CreateSqlServer(ctx context.Context, request sqlserver2.CreateSqlServerRequest, tags map[string]interface{}) (sqlserver2.AsyncResponse, int, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, c, err := client.sdkClient.MsSqlConfigurationControllerApi.CreateSqlserver(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetSqlServer(ctx context.Context, sqlServerId string) (sqlserver2.DetailDatabaseResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DetailDatabase10(ctx, client.config.ProjectId, sqlServerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteSqlServer(ctx context.Context, sqlServerId string) (sqlserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DeleteDatabase12(ctx, client.config.ProjectId, sqlServerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}

	return result, statusCode, err
}

func (client *Client) ListSqlServer(ctx context.Context, request *sqlserver2.MsSqlConfigurationControllerApiListSqlserverOpts) (sqlserver2.ListResponseOfDbServerGroupsResponse, int, error) {
	result, c, err := client.sdkClient.MsSqlConfigurationControllerApi.ListSqlserver(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateSqlServerScale(ctx context.Context, dbServerGroupId string, virtualServerId string, scaleId string) (sqlserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseScale11(ctx, client.config.ProjectId, dbServerGroupId, sqlserver2.ResizeScaleRequest{
		VirtualServerId: virtualServerId,
		ScaleProductId:  scaleId,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateSqlServerBlockSize(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageId string, blockStorageSize int) (sqlserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.DatabaseBlockStorageControllerApi.ResizeDatabaseStorage12(ctx, client.config.ProjectId, dbServerGroupId, sqlserver2.ResizeStorageRequest{
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

func (client *Client) UpdateBackupSetting(ctx context.Context, dbServerGroupId string, request sqlserver2.UpdateBackupSettingRequest) (sqlserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.DatabaseBackupControllerApi.UpdateBackupSetting9(ctx, client.config.ProjectId, dbServerGroupId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
