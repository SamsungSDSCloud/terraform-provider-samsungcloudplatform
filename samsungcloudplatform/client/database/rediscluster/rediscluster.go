package rediscluster

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/redis"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *redis.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: redis.NewAPIClient(config),
	}
}

func (client *Client) CreateRedisCluster(ctx context.Context, request redis.RedisClusterCreateRequest, tags map[string]interface{}) (redis.AsyncResponse, int, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, c, err := client.sdkClient.RedisClusterProvisioningApi.CreateRedisCluster(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListRedisCluster(ctx context.Context, request *redis.RedisClusterSearchApiListRedisClusterOpts) (redis.ListResponseRedisListItemResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterSearchApi.ListRedisCluster(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailRedisCluster(ctx context.Context, redisClusterId string) (redis.RedisClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterSearchApi.DetailRedisCluster(ctx, client.config.ProjectId, redisClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteRedisCluster(ctx context.Context, redisClusterId string) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterOperationManagementApi.DeleteRedisCluster(ctx, client.config.ProjectId, redisClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AttachRedisClusterSecurityGroup(ctx context.Context, redisClusterId string, request redis.DbClusterAttachSecurityGroupRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterNetworkApi.AttachRedisClusterSecurityGroup(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachRedisClusterSecurityGroup(ctx context.Context, redisClusterId string, request redis.DbClusterDetachSecurityGroupRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterNetworkApi.DetachRedisClusterSecurityGroup(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeRedisClusterBlockStorages(ctx context.Context, redisClusterId string, request redis.RedisResizeBlockStoragesRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterInfraResourceApi.ResizeRedisClusterBlockStorages(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeRedisClusterVirtualServers(ctx context.Context, redisClusterId string, request redis.RedisResizeVirtualServersRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterInfraResourceApi.ResizeRedisClusterVirtualServers(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StartRedisCluster(ctx context.Context, redisClusterId string) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterOperationManagementApi.StartRedisCluster(ctx, client.config.ProjectId, redisClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StopRedisCluster(ctx context.Context, redisClusterId string) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterOperationManagementApi.StopRedisCluster(ctx, client.config.ProjectId, redisClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyRedisClusterContract(ctx context.Context, redisClusterId string, request redis.DbClusterModifyContractRequest) (redis.RedisClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterPricingApi.ModifyRedisClusterContract(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyRedisClusterNextContract(ctx context.Context, redisClusterId string, request redis.DbClusterModifyNextContractRequest) (redis.RedisClusterDetailResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterPricingApi.ModifyRedisClusterNextContract(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateRedisClusterFullBackupConfig(ctx context.Context, redisClusterId string, request redis.RedisCreateFullBackupConfigRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterBackupApi.CreateRedisClusterFullBackupConfig(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyRedisClusterFullBackupConfig(ctx context.Context, redisClusterId string, request redis.RedisModifyFullBackupConfigRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterBackupApi.ModifyRedisClusterFullBackupConfig(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteRedisClusterFullBackupConfig(ctx context.Context, redisClusterId string) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisClusterBackupApi.DeleteRedisClusterFullBackupConfig(ctx, client.config.ProjectId, redisClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
