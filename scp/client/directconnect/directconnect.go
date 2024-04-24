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

func (client *Client) CreateDirectConnect(ctx context.Context, bandwidth int32, dcName string, serviceZoneId string, dcDescription string, tags map[string]interface{}) (directconnect2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.DirectConnectOpenApiV3ControllerApi.CreateDirectConnect1(ctx, client.config.ProjectId, directconnect2.DirectConnectCreateV3Request{
		BandwidthGbps:            bandwidth,
		DirectConnectName:        dcName,
		ServiceZoneId:            serviceZoneId,
		DirectConnectDescription: dcDescription,
		Tags:                     client.sdkClient.ToTagRequestList(tags),
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

func (client *Client) GetDirectConnectList(ctx context.Context, request *directconnect2.DirectConnectOpenApiControllerApiListDirectConnectsOpts) (directconnect2.ListResponseDirectConnectListItemResponse, int, error) {
	result, c, err := client.sdkClient.DirectConnectOpenApiControllerApi.ListDirectConnects(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

//------------Direct Connect Connection-------------------//

func (client *Client) CreateDconVpcConnection(ctx context.Context, approverProjectId string, approverVpcId string, connectionType string, firewallEnabled bool, requesterDcId string,
	requesterProjectId string, connectionDescription string, tags map[string]interface{}) (directconnect2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.DirectConnectConnectionOpenApiControllerApi.CreateDirectConnectConnection(ctx, client.config.ProjectId, directconnect2.DirectConnectConnectionCreateRequest{
		ApproverProjectId:                  approverProjectId,
		ApproverVpcId:                      approverVpcId,
		DirectConnectConnectionType:        connectionType,
		FirewallEnabled:                    &firewallEnabled,
		RequesterDirectConnectId:           requesterDcId,
		RequesterProjectId:                 requesterProjectId,
		DirectConnectConnectionDescription: connectionDescription,
		Tags:                               client.sdkClient.ToTagRequestList(tags),
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

func (client *Client) GetDconVpcConnectionList(ctx context.Context, request *directconnect2.DirectConnectConnectionOpenApiControllerApiListDirectConnectConnectionsOpts) (directconnect2.ListResponseDirectConnectConnectionListResponse, int, error) {
	result, c, err := client.sdkClient.DirectConnectConnectionOpenApiControllerApi.ListDirectConnectConnections(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
