package directconnect

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	directconnect2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/direct-connect2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *directconnect2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: directconnect2.NewAPIClient(config),
	}
}

//------------Direct Connect -------------------//

func (client *Client) GetDirectConnectInfo(ctx context.Context, directConnectId string) (directconnect2.DirectConnectDetailResponse, int, error) {
	result, c, err := client.sdkClient.DirectConnectOpenApiControllerApi.DetailDirectConnect(ctx, client.config.ProjectId, directConnectId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateDirectConnect(ctx context.Context, bandwidth int32, dcName string, serviceZoneId string, dcDescription string) (directconnect2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.DirectConnectOpenApiV3ControllerApi.CreateDirectConnect1(ctx, client.config.ProjectId, directconnect2.DirectConnectCreateV3Request{
		BandwidthGbps:            bandwidth,
		DirectConnectName:        dcName,
		ServiceZoneId:            serviceZoneId,
		DirectConnectDescription: dcDescription,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteDirectConnect(ctx context.Context, dcId string) error {
	_, _, err := client.sdkClient.DirectConnectOpenApiControllerApi.DeleteDirectConnect(ctx, client.config.ProjectId, dcId)
	return err
}

func (client *Client) GetDirectConnectList(ctx context.Context, request *directconnect2.DirectConnectOpenApiControllerApiListDirectConnects1Opts) (directconnect2.ListResponseOfDirectConnectListItemResponse, int, error) {
	result, c, err := client.sdkClient.DirectConnectOpenApiControllerApi.ListDirectConnects1(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

//------------Direct Connect Connection-------------------//

func (client *Client) CreateDconVpcConnection(ctx context.Context, approverProjectId string, approverVpcId string, connectionType string, firewallEnabled bool, requesterDcId string,
	requesterProjectId string, connectionDescription string) (directconnect2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.DirectConnectConnectionOpenApiControllerApi.CreateDirectConnectConnection(ctx, client.config.ProjectId, directconnect2.DirectConnectConnectionCreateRequest{
		ApproverProjectId:                  approverProjectId,
		ApproverVpcId:                      approverVpcId,
		DirectConnectConnectionType:        connectionType,
		FirewallEnabled:                    &firewallEnabled,
		RequesterDirectConnectId:           requesterDcId,
		RequesterProjectId:                 requesterProjectId,
		DirectConnectConnectionDescription: connectionDescription,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteDconVpcConnection(ctx context.Context, dcId string) error {
	_, _, err := client.sdkClient.DirectConnectConnectionOpenApiControllerApi.DeleteDirectConnectConnection(ctx, client.config.ProjectId, dcId)
	return err
}

func (client *Client) GetDconVpcConnectionInfo(ctx context.Context, directConnectId string) (directconnect2.DirectConnectConnectionDetailResponse, int, error) {
	result, c, err := client.sdkClient.DirectConnectConnectionOpenApiControllerApi.DetailDirectConnectConnections(ctx, client.config.ProjectId, directConnectId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetDconVpcConnectionList(ctx context.Context, request *directconnect2.DirectConnectConnectionOpenApiControllerApiListDirectConnectConnections1Opts) (directconnect2.ListResponseOfDirectConnectConnectionListResponse, int, error) {
	result, c, err := client.sdkClient.DirectConnectConnectionOpenApiControllerApi.ListDirectConnectConnections1(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
