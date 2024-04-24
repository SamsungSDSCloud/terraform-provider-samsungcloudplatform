package redis

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

func (client *Client) CreateRedis(ctx context.Context, request redis.RedisCreateRequest, tags map[string]interface{}) (redis.AsyncResponse, int, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, c, err := client.sdkClient.RedisProvisioningApi.CreateRedis(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListRedis(ctx context.Context, request *redis.RedisSearchApiListRedisOpts) (redis.ListResponseRedisListItemResponse, int, error) {
	result, c, err := client.sdkClient.RedisSearchApi.ListRedis(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailRedis(ctx context.Context, redisClusterId string) (redis.RedisDetailResponse, int, error) {
	result, c, err := client.sdkClient.RedisSearchApi.DetailRedis(ctx, client.config.ProjectId, redisClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteRedis(ctx context.Context, redisClusterId string) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisOperationManagementApi.DeleteRedis(ctx, client.config.ProjectId, redisClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AttachRedisSecurityGroup(ctx context.Context, redisClusterId string, request redis.DbClusterAttachSecurityGroupRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisNetworkApi.AttachRedisSecurityGroup(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachRedisSecurityGroup(ctx context.Context, redisClusterId string, request redis.DbClusterDetachSecurityGroupRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisNetworkApi.DetachRedisSecurityGroup(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeRedisBlockStorages(ctx context.Context, redisClusterId string, request redis.RedisResizeBlockStoragesRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisInfraResourceApi.ResizeRedisBlockStorages(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ResizeRedisVirtualServers(ctx context.Context, redisClusterId string, request redis.RedisResizeVirtualServersRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisInfraResourceApi.ResizeRedisVirtualServers(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StartRedis(ctx context.Context, redisClusterId string) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisOperationManagementApi.StartRedis(ctx, client.config.ProjectId, redisClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StopRedis(ctx context.Context, redisClusterId string) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisOperationManagementApi.StopRedis(ctx, client.config.ProjectId, redisClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyRedisContract(ctx context.Context, redisClusterId string, request redis.DbClusterModifyContractRequest) (redis.RedisDetailResponse, int, error) {
	result, c, err := client.sdkClient.RedisPricingApi.ModifyRedisContract(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyRedisNextContract(ctx context.Context, redisClusterId string, request redis.DbClusterModifyNextContractRequest) (redis.RedisDetailResponse, int, error) {
	result, c, err := client.sdkClient.RedisPricingApi.ModifyRedisNextContract(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateRedisFullBackupConfig(ctx context.Context, redisClusterId string, request redis.RedisCreateFullBackupConfigRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisBackupApi.CreateRedisFullBackupConfig(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ModifyRedisFullBackupConfig(ctx context.Context, redisClusterId string, request redis.RedisModifyFullBackupConfigRequest) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisBackupApi.ModifyRedisFullBackupConfig(ctx, client.config.ProjectId, redisClusterId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteRedisFullBackupConfig(ctx context.Context, redisClusterId string) (redis.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.RedisBackupApi.DeleteRedisFullBackupConfig(ctx, client.config.ProjectId, redisClusterId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
