package mysql

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/mysql"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *mysql.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: mysql.NewAPIClient(config),
	}
}

func (client *Client) CreateMysqlCluster(ctx context.Context, request mysql.MysqlClusterCreateRequest, tags map[string]interface{}) (mysql.AsyncResponse, int, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, c, err := client.sdkClient.MysqlProvisioningApi.CreateMysqlCluster(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListMysqlClusters(ctx context.Context, request *mysql.MysqlSearchApiListMysqlClustersOpts) (mysql.ListResponseMysqlClusterListItemResponse, int, error) {
	result, c, err := client.sdkClient.MysqlSearchApi.ListMysqlClusters(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailMysqlCluster(ctx context.Context, mysqlClusterId string) (mysql.MysqlClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.MysqlSearchApi.DetailMysqlCluster(ctx, client.config.ProjectId, mysqlClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteMysqlCluster(ctx context.Context, mysqlClusterId string) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlManagementOpenApiV1ControllerApi.DeleteMysqlCluster(ctx, client.config.ProjectId, mysqlClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AttachMysqlClusterSecurityGroup(ctx context.Context, mysqlClusterId string, request mysql.DbClusterAttachSecurityGroupRequest) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlNetworkApi.AttachMysqlClusterSecurityGroup(ctx, client.config.ProjectId, mysqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachMysqlClusterSecurityGroup(ctx context.Context, mysqlClusterId string, request mysql.DbClusterDetachSecurityGroupRequest) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlNetworkApi.DetachMysqlClusterSecurityGroup(ctx, client.config.ProjectId, mysqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AddMysqlClusterBlockStorages(ctx context.Context, mysqlClusterId string, request mysql.MysqlClusterAddBlockStoragesRequest) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlInfraResourceApi.AddMysqlClusterBlockStorages(ctx, client.config.ProjectId, mysqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeMysqlClusterBlockStorages(ctx context.Context, mysqlClusterId string, request mysql.MysqlClusterResizeBlockStoragesRequest) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlInfraResourceApi.ResizeMysqlClusterBlockStorages(ctx, client.config.ProjectId, mysqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeMysqlClusterVirtualServers(ctx context.Context, mysqlClusterId string, request mysql.MysqlClusterResizeVirtualServersRequest) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlInfraResourceApi.ResizeMysqlClusterVirtualServers(ctx, client.config.ProjectId, mysqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StartMysqlCluster(ctx context.Context, mysqlClusterId string) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlManagementOpenApiV1ControllerApi.StartMysqlCluster(ctx, client.config.ProjectId, mysqlClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StopMysqlCluster(ctx context.Context, mysqlClusterId string) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlManagementOpenApiV1ControllerApi.StopMysqlCluster(ctx, client.config.ProjectId, mysqlClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyMysqlClusterContract(ctx context.Context, mysqlClusterId string, request mysql.DbClusterModifyContractRequest) (mysql.MysqlClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.MysqlPricingApi.ModifyMysqlClusterContract(ctx, client.config.ProjectId, mysqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyMysqlClusterNextContract(ctx context.Context, mysqlClusterId string, request mysql.DbClusterModifyNextContractRequest) (mysql.MysqlClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.MysqlPricingApi.ModifyMysqlClusterNextContract(ctx, client.config.ProjectId, mysqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateMysqlClusterFullBackupConfig(ctx context.Context, mysqlClusterId string, request mysql.DbClusterCreateFullBackupConfigRequest) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlBackupApi.CreateMysqlClusterFullBackupConfig(ctx, client.config.ProjectId, mysqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyMysqlClusterFullBackupConfig(ctx context.Context, mysqlClusterId string, request mysql.DbClusterModifyFullBackupConfigRequest) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlBackupApi.ModifyMysqlClusterFullBackupConfig(ctx, client.config.ProjectId, mysqlClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteMysqlClusterFullBackupConfig(ctx context.Context, mysqlClusterId string) (mysql.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MysqlBackupApi.DeleteMysqlClusterFullBackupConfig(ctx, client.config.ProjectId, mysqlClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
