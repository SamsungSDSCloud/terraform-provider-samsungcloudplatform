package epas

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/epas"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *epas.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: epas.NewAPIClient(config),
	}
}

func (client *Client) CreateEpasCluster(ctx context.Context, request epas.EpasClusterCreateRequest, tags map[string]interface{}) (epas.AsyncResponse, int, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, c, err := client.sdkClient.EpasProvisioningApi.CreateEpasCluster(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListEpasClusters(ctx context.Context, request *epas.EpasSearchApiListEpasClustersOpts) (epas.ListResponseEpasClusterListItemResponse, int, error) {
	result, c, err := client.sdkClient.EpasSearchApi.ListEpasClusters(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailEpasCluster(ctx context.Context, epasClusterId string) (epas.EpasClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.EpasSearchApi.DetailEpasCluster(ctx, client.config.ProjectId, epasClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteEpasCluster(ctx context.Context, epasClusterId string) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasOperationManagementApi.DeleteEpasCluster(ctx, client.config.ProjectId, epasClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AttachEpasClusterSecurityGroup(ctx context.Context, epasClusterId string, request epas.DbClusterAttachSecurityGroupRequest) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasNetworkApi.AttachEpasClusterSecurityGroup(ctx, client.config.ProjectId, epasClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachEpasClusterSecurityGroup(ctx context.Context, epasClusterId string, request epas.DbClusterDetachSecurityGroupRequest) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasNetworkApi.DetachEpasClusterSecurityGroup(ctx, client.config.ProjectId, epasClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AddEpasClusterBlockStorages(ctx context.Context, epasClusterId string, request epas.EpasClusterAddBlockStoragesRequest) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasInfraResourceApi.AddEpasClusterBlockStorages(ctx, client.config.ProjectId, epasClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeEpasClusterBlockStorages(ctx context.Context, epasClusterId string, request epas.EpasClusterResizeBlockStoragesRequest) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasInfraResourceApi.ResizeEpasClusterBlockStorages(ctx, client.config.ProjectId, epasClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeEpasClusterVirtualServers(ctx context.Context, epasClusterId string, request epas.EpasClusterResizeVirtualServersRequest) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasInfraResourceApi.ResizeEpasClusterVirtualServers(ctx, client.config.ProjectId, epasClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StartEpasCluster(ctx context.Context, epasClusterId string) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasOperationManagementApi.StartEpasCluster(ctx, client.config.ProjectId, epasClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StopEpasCluster(ctx context.Context, epasClusterId string) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasOperationManagementApi.StopEpasCluster(ctx, client.config.ProjectId, epasClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyEpasClusterContract(ctx context.Context, epasClusterId string, request epas.DbClusterModifyContractRequest) (epas.EpasClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.EpasPricingApi.ModifyEpasClusterContract(ctx, client.config.ProjectId, epasClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyEpasClusterNextContract(ctx context.Context, epasClusterId string, request epas.DbClusterModifyNextContractRequest) (epas.EpasClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.EpasPricingApi.ModifyEpasClusterNextContract(ctx, client.config.ProjectId, epasClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateEpasClusterFullBackupConfig(ctx context.Context, epasClusterId string, request epas.DbClusterCreateFullBackupConfigRequest) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasBackupApi.CreateEpasClusterFullBackupConfig(ctx, client.config.ProjectId, epasClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyEpasClusterFullBackupConfig(ctx context.Context, epasClusterId string, request epas.DbClusterModifyFullBackupConfigRequest) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasBackupApi.ModifyEpasClusterFullBackupConfig(ctx, client.config.ProjectId, epasClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteEpasClusterFullBackupConfig(ctx context.Context, epasClusterId string) (epas.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.EpasBackupApi.DeleteEpasClusterFullBackupConfig(ctx, client.config.ProjectId, epasClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
