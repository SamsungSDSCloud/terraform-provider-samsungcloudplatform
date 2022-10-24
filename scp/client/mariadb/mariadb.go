package mariadb

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/client"
	"github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/library/mariadb2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *mariadb2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: mariadb2.NewAPIClient(config),
	}
}

func (client *Client) CreateMariadb(ctx context.Context, request mariadb2.CreateMariaDbRequest) (mariadb2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.MariaDbConfigurationControllerApi.CreateMariadb(ctx, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListMariadb(ctx context.Context, dbName string, serverGroupName string, virtualServerName string) (mariadb2.ListResponseOfDbServerGroupsResponse, error) {
	result, _, err := client.sdkClient.MariaDbConfigurationControllerApi.ListMariadb(ctx, dbName, serverGroupName, virtualServerName, &mariadb2.MariaDbConfigurationControllerApiListMariadbOpts{
		CreatedBy: optional.String{},
		Page:      optional.NewInt32(0),
		Size:      optional.NewInt32(1000),
		Sort:      optional.Interface{},
	})
	return result, err
}

func (client *Client) GetMariadb(ctx context.Context, dbServerGroupId string) (mariadb2.DetailDatabaseResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.DetailDatabase4(ctx, dbServerGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateMariadbBlockSize(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageId string, blockStorageSize int) (mariadb2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseScale4(ctx, dbServerGroupId, mariadb2.ResizeScaleRequest{
		VirtualServerId: virtualServerId,
		ScaleProductId:  blockStorageId, //need to change
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AddPostgresqlBlock(ctx context.Context, dbServerGroupId string, virtualServerId string, blockStorageType string, blockStorageSize int) (mariadb2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.AddDatabaseStorage3(ctx, dbServerGroupId, mariadb2.AddStorageRequest{
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

func (client *Client) UpdatePostgresqlScale(ctx context.Context, dbServerGroupId string, virtualServerId string, scaleId string) (mariadb2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.ConfigurationControllerApi.ResizeDatabaseStorage4(ctx, dbServerGroupId, mariadb2.ResizeStorageRequest{
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
