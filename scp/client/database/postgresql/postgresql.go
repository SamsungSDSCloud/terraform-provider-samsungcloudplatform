package postgresql

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/postgresql"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *postgresql.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: postgresql.NewAPIClient(config),
	}
}

func (client *Client) CreatePostgresqlCluster(ctx context.Context, request postgresql.PostgresqlClusterCreateRequest, tags map[string]interface{}) (postgresql.AsyncResponse, int, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, c, err := client.sdkClient.PostgresqlProvisioningApi.CreatePostgresqlCluster(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListPostgresqlClusters(ctx context.Context, request *postgresql.PostgresqlSearchApiListPostgresqlClustersOpts) (postgresql.ListResponseOfPostgresqlClusterListItemResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlSearchApi.ListPostgresqlClusters(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailPostgresqlCluster(ctx context.Context, postgresqlClusterId string) (postgresql.PostgresqlClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlSearchApi.DetailPostgresqlCluster(ctx, client.config.ProjectId, postgresqlClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeletePostgresqlCluster(ctx context.Context, postgresqlClusterId string) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlOperationManagementApi.DeletePostgresqlCluster(ctx, client.config.ProjectId, postgresqlClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AttachPostgresqlClusterSecurityGroup(ctx context.Context, postgresqlClusterId string, request postgresql.DbClusterAttachSecurityGroupRequest) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlNetworkApi.AttachPostgresqlClusterSecurityGroup(ctx, client.config.ProjectId, postgresqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachPostgresqlClusterSecurityGroup(ctx context.Context, postgresqlClusterId string, request postgresql.DbClusterDetachSecurityGroupRequest) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlNetworkApi.DetachPostgresqlClusterSecurityGroup(ctx, client.config.ProjectId, postgresqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AddPostgresqlClusterBlockStorages(ctx context.Context, postgresqlClusterId string, request postgresql.PostgresqlClusterAddBlockStoragesRequest) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlInfraResourceApi.AddPostgresqlClusterBlockStorages(ctx, client.config.ProjectId, postgresqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizePostgresqlClusterBlockStorages(ctx context.Context, postgresqlClusterId string, request postgresql.PostgresqlClusterResizeBlockStoragesRequest) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlInfraResourceApi.ResizePostgresqlClusterBlockStorages(ctx, client.config.ProjectId, postgresqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizePostgresqlClusterVirtualServers(ctx context.Context, postgresqlClusterId string, request postgresql.PostgresqlClusterResizeVirtualServersRequest) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlInfraResourceApi.ResizePostgresqlClusterVirtualServers(ctx, client.config.ProjectId, postgresqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StartPostgresqlCluster(ctx context.Context, postgresqlClusterId string) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlOperationManagementApi.StartPostgresqlCluster(ctx, client.config.ProjectId, postgresqlClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StopPostgresqlCluster(ctx context.Context, postgresqlClusterId string) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlOperationManagementApi.StopPostgresqlCluster(ctx, client.config.ProjectId, postgresqlClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyPostgresqlClusterContract(ctx context.Context, postgresqlClusterId string, request postgresql.DbClusterModifyContractRequest) (postgresql.PostgresqlClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlPricingApi.ModifyPostgresqlClusterContract(ctx, client.config.ProjectId, postgresqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyPostgresqlClusterNextContract(ctx context.Context, postgresqlClusterId string, request postgresql.DbClusterModifyNextContractRequest) (postgresql.PostgresqlClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlPricingApi.ModifyPostgresqlClusterNextContract(ctx, client.config.ProjectId, postgresqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreatePostgresqlClusterFullBackupConfig(ctx context.Context, postgresqlClusterId string, request postgresql.DbClusterCreateFullBackupConfigRequest) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlBackupApi.CreatePostgresqlClusterFullBackupConfig(ctx, client.config.ProjectId, postgresqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyPostgresqlClusterFullBackupConfig(ctx context.Context, postgresqlClusterId string, request postgresql.DbClusterModifyFullBackupConfigRequest) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlBackupApi.ModifyPostgresqlClusterFullBackupConfig(ctx, client.config.ProjectId, postgresqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeletePostgresqlClusterFullBackupConfig(ctx context.Context, postgresqlClusterId string) (postgresql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.PostgresqlBackupApi.DeletePostgresqlClusterFullBackupConfig(ctx, client.config.ProjectId, postgresqlClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
