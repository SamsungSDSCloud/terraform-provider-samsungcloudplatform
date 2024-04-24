package mariadb

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/mariadb"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *mariadb.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: mariadb.NewAPIClient(config),
	}
}

func (client *Client) CreateMariadbCluster(ctx context.Context, request mariadb.MariadbClusterCreateRequest, tags map[string]interface{}) (mariadb.AsyncResponse, int, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, c, err := client.sdkClient.MariadbProvisioningApi.CreateMariadbCluster(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListMariadbClusters(ctx context.Context, request *mariadb.MariadbSearchApiListMariadbClustersOpts) (mariadb.ListResponseMariadbClusterListItemResponse, int, error) {
	result, c, err := client.sdkClient.MariadbSearchApi.ListMariadbClusters(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailMariadbCluster(ctx context.Context, mariadbClusterId string) (mariadb.MariadbClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.MariadbSearchApi.DetailMariadbCluster(ctx, client.config.ProjectId, mariadbClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteMariadbCluster(ctx context.Context, mariadbClusterId string) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbOperationManagementApi.DeleteMariadbCluster(ctx, client.config.ProjectId, mariadbClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AttachMariadbClusterSecurityGroup(ctx context.Context, mariadbClusterId string, request mariadb.DbClusterAttachSecurityGroupRequest) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbNetworkApi.AttachMariadbClusterSecurityGroup(ctx, client.config.ProjectId, mariadbClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachMariadbClusterSecurityGroup(ctx context.Context, mariadbClusterId string, request mariadb.DbClusterDetachSecurityGroupRequest) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbNetworkApi.DetachMariadbClusterSecurityGroup(ctx, client.config.ProjectId, mariadbClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AddMariadbClusterBlockStorages(ctx context.Context, mariadbClusterId string, request mariadb.MariadbClusterAddBlockStoragesRequest) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbInfraResourceApi.AddMariadbClusterBlockStorages(ctx, client.config.ProjectId, mariadbClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeMariadbClusterBlockStorages(ctx context.Context, mariadbClusterId string, request mariadb.MariadbClusterResizeBlockStoragesRequest) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbInfraResourceApi.ResizeMariadbClusterBlockStorages(ctx, client.config.ProjectId, mariadbClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeMariadbClusterVirtualServers(ctx context.Context, mariadbClusterId string, request mariadb.MariadbClusterResizeVirtualServersRequest) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbInfraResourceApi.ResizeMariadbClusterVirtualServers(ctx, client.config.ProjectId, mariadbClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StartMariadbCluster(ctx context.Context, mariadbClusterId string) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbOperationManagementApi.StartMariadbCluster(ctx, client.config.ProjectId, mariadbClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StopMariadbCluster(ctx context.Context, mariadbClusterId string) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbOperationManagementApi.StopMariadbCluster(ctx, client.config.ProjectId, mariadbClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyMariadbClusterContract(ctx context.Context, mariadbClusterId string, request mariadb.DbClusterModifyContractRequest) (mariadb.MariadbClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.MariadbPricingApi.ModifyMariadbClusterContract(ctx, client.config.ProjectId, mariadbClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyMariadbClusterNextContract(ctx context.Context, mariadbClusterId string, request mariadb.DbClusterModifyNextContractRequest) (mariadb.MariadbClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.MariadbPricingApi.ModifyMariadbClusterNextContract(ctx, client.config.ProjectId, mariadbClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateMariadbClusterFullBackupConfig(ctx context.Context, mariadbClusterId string, request mariadb.DbClusterCreateFullBackupConfigRequest) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbBackupApi.CreateMariadbClusterFullBackupConfig(ctx, client.config.ProjectId, mariadbClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyMariadbClusterFullBackupConfig(ctx context.Context, mariadbClusterId string, request mariadb.DbClusterModifyFullBackupConfigRequest) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbBackupApi.ModifyMariadbClusterFullBackupConfig(ctx, client.config.ProjectId, mariadbClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteMariadbClusterFullBackupConfig(ctx context.Context, mariadbClusterId string) (mariadb.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariadbBackupApi.DeleteMariadbClusterFullBackupConfig(ctx, client.config.ProjectId, mariadbClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
