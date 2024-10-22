package sqlserver

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/sqlserver"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *sqlserver.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: sqlserver.NewAPIClient(config),
	}
}

func (client *Client) CreateSqlserverCluster(ctx context.Context, request sqlserver.SqlserverClusterCreateRequest, tags map[string]interface{}) (sqlserver.AsyncResponse, int, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, c, err := client.sdkClient.SqlserverProvisioningApi.CreateSqlserverCluster(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListSqlserverClusters(ctx context.Context, request *sqlserver.SqlserverSearchApiListSqlserverClustersOpts) (sqlserver.ListResponseSqlserverClusterListItemResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverSearchApi.ListSqlserverClusters(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailSqlserverCluster(ctx context.Context, sqlserverClusterId string) (sqlserver.SqlserverClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverSearchApi.DetailSqlserverCluster(ctx, client.config.ProjectId, sqlserverClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StartSqlserverCluster(ctx context.Context, sqlserverClusterId string) (sqlserver.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverOperationManagementApi.StartSqlserverCluster(ctx, client.config.ProjectId, sqlserverClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}

	return result, statusCode, err
}

func (client *Client) StopSqlserverCluster(ctx context.Context, sqlserverClusterId string) (sqlserver.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverOperationManagementApi.StopSqlserverCluster(ctx, client.config.ProjectId, sqlserverClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}

	return result, statusCode, err
}

func (client *Client) DeleteSqlserverCluster(ctx context.Context, sqlserverClusterId string) (sqlserver.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverOperationManagementApi.DeleteSqlserverCluster(ctx, client.config.ProjectId, sqlserverClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}

	return result, statusCode, err
}

func (client *Client) ModifySqlserverClusterContract(ctx context.Context, sqlserverClusterId string, request sqlserver.DbClusterModifyContractRequest) (sqlserver.SqlserverClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverPricingApi.ModifySqlserverClusterContract(ctx, client.config.ProjectId, sqlserverClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifySqlserverClusterNextContract(ctx context.Context, sqlserverClusterId string, request sqlserver.DbClusterModifyNextContractRequest) (sqlserver.SqlserverClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverPricingApi.ModifySqlserverClusterNextContract(ctx, client.config.ProjectId, sqlserverClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeSqlserverClusterVirtualServers(ctx context.Context, sqlserverClusterId string, request sqlserver.SqlserverClusterResizeVirtualServersRequest) (sqlserver.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverInfraResourceApi.ResizeSqlserverClusterVirtualServers(ctx, client.config.ProjectId, sqlserverClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeSqlserverClusterBlockStorages(ctx context.Context, sqlserverClusterId string, request sqlserver.SqlserverClusterResizeBlockStoragesRequest) (sqlserver.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverInfraResourceApi.ResizeSqlserverClusterBlockStorages(ctx, client.config.ProjectId, sqlserverClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AttachSqlserverClusterSecurityGroup(ctx context.Context, sqlserverClusterId string, request sqlserver.DbClusterAttachSecurityGroupRequest) (sqlserver.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverNetworkApi.AttachSqlserverClusterSecurityGroup(ctx, client.config.ProjectId, sqlserverClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachSqlserverClusterSecurityGroup(ctx context.Context, sqlserverClusterId string, request sqlserver.DbClusterDetachSecurityGroupRequest) (sqlserver.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverNetworkApi.DetachSqlserverClusterSecurityGroup(ctx, client.config.ProjectId, sqlserverClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateSqlserverClusterFullBackupConfig(ctx context.Context, sqlserverClusterId string, request sqlserver.SqlserverCreateFullBackupConfigRequest) (sqlserver.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverBackupApi.CreateSqlserverClusterFullBackupConfig(ctx, client.config.ProjectId, sqlserverClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifySqlserverClusterFullBackupConfig(ctx context.Context, sqlserverClusterId string, request sqlserver.SqlserverModifyFullBackupConfigRequest) (sqlserver.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverBackupApi.ModifySqlserverClusterFullBackupConfig(ctx, client.config.ProjectId, sqlserverClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteSqlserverClusterFullBackupConfig(ctx context.Context, sqlserverClusterId string) (sqlserver.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.SqlserverBackupApi.DeleteSqlserverClusterFullBackupConfig(ctx, client.config.ProjectId, sqlserverClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
