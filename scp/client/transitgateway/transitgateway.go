package transitgateway

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	transitgateway2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/transit-gateway2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *transitgateway2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: transitgateway2.NewAPIClient(config),
	}
}

// Transit Gateway --------------------------->

func (client *Client) GetTransitGatewayInfo(ctx context.Context, transitGatewayId string) (transitgateway2.TransitGatewayDetailResponse, int, error) {
	result, c, err := client.sdkClient.TransitGatewayOpenApiControllerApi.DetailTransitGateway2(ctx, client.config.ProjectId, transitGatewayId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetTransitGatewayList(ctx context.Context, request *transitgateway2.TransitGatewayOpenApiControllerApiListTransitGateways3Opts) (transitgateway2.ListResponseOfTransitGatewayListItemResponse, int, error) {
	result, c, err := client.sdkClient.TransitGatewayOpenApiControllerApi.ListTransitGateways3(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateTransitGateway(ctx context.Context, bandwidthGbps int32, serviceZoneId string, name string, uplinkEnabled bool, description string) (transitgateway2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.TransitGatewayOpenApiV3ControllerApi.CreateTransitGateway1(ctx, client.config.ProjectId, transitgateway2.TransitGatewayCreateV3Request{
		BandwidthGbps:             bandwidthGbps,
		ServiceZoneId:             serviceZoneId,
		TransitGatewayName:        name,
		UplinkEnabled:             &uplinkEnabled,
		TransitGatewayDescription: description,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateTransitGatewayUplinkEnable(ctx context.Context, transitGatewayId string, uplinkEnabled bool) (transitgateway2.TransitGatewayDetailResponse, int, error) {
	result, c, err := client.sdkClient.TransitGatewayOpenApiControllerApi.UpdateTransitGatewayUplinkEnabled(ctx, client.config.ProjectId, transitGatewayId, transitgateway2.TransitGatewayUplinkUpdateRequest{
		UplinkEnabled: &uplinkEnabled,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteTransitGateway(ctx context.Context, transitGatewayId string) error {
	_, _, err := client.sdkClient.TransitGatewayOpenApiControllerApi.DeleteTransitGateway(ctx, client.config.ProjectId, transitGatewayId)
	return err
}

// <-------------------------------

//  Transit Gateway - VPC Connections  -------------------------->

func (client *Client) GetTransitGatewayConnectionInfo(ctx context.Context, transitGatewayConnectionId string) (transitgateway2.TransitGatewayConnectionDetailResponse, int, error) {
	result, c, err := client.sdkClient.TransitGatewayConnectionOpenApiControllerApi.DetailTransitGatewayConnection(ctx, client.config.ProjectId, transitGatewayConnectionId)

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err

}

func (client *Client) GetTransitGatewayConnectionList(ctx context.Context, request *transitgateway2.TransitGatewayConnectionOpenApiControllerApiListTransitGatewayConnections1Opts) (transitgateway2.ListResponseOfTransitGatewayConnectionListItemResponse, int, error) {
	result, c, err := client.sdkClient.TransitGatewayConnectionOpenApiControllerApi.ListTransitGatewayConnections1(ctx, client.config.ProjectId, request)

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}

	return result, statusCode, err
}

func (client *Client) CreateTransitGatewayConnection(ctx context.Context, transitGatewayId string, vpcId string, requesterProjectId string, approverProjectId string, description string, firewallEnabled bool, firewallLoggable bool, connectionType string, tags map[string]interface{}) (transitgateway2.TransitGatewayConnectionApprovalResponse, int, error) {
	result, c, err := client.sdkClient.TransitGatewayConnectionOpenApiControllerApi.CreateTransitGatewayConnection(ctx, client.config.ProjectId, transitgateway2.TransitGatewayConnectionCreateRequest{
		RequesterProjectId:                  requesterProjectId,
		RequesterTransitGatewayId:           transitGatewayId,
		ApproverProjectId:                   approverProjectId,
		ApproverVpcId:                       vpcId,
		TransitGatewayConnectionDescription: description,
		FirewallEnabled:                     &firewallEnabled,
		FirewallLoggable:                    &firewallLoggable,
		TransitGatewayConnectionType:        connectionType,
		Tags:                                client.sdkClient.ToTagRequestList(tags),
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ApproveTransitGatewayConnection(ctx context.Context, transitGatewayConnectionId string) (transitgateway2.TransitGatewayConnectionApprovalResponse, int, error) {
	result, c, err := client.sdkClient.TransitGatewayConnectionOpenApiControllerApi.ApproveTransitGatewayConnection(ctx, client.config.ProjectId, transitGatewayConnectionId)

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CancelTransitGatewayConnection(ctx context.Context, transitGatewayConnectionId string) (transitgateway2.TransitGatewayConnectionApprovalResponse, int, error) {
	result, c, err := client.sdkClient.TransitGatewayConnectionOpenApiControllerApi.CancelTransitGatewayConnection(ctx, client.config.ProjectId, transitGatewayConnectionId)

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateTransitGatewayConnectionDescription(ctx context.Context, transitGatewayConnectionId string, description string) (transitgateway2.TransitGatewayConnectionDetailResponse, int, error) {
	result, c, err := client.sdkClient.TransitGatewayConnectionOpenApiControllerApi.UpdateTransitGatewayConnectionDescription(ctx, client.config.ProjectId, transitGatewayConnectionId, transitgateway2.TransitGatewayConnectionDescriptionUpdateRequest{
		TransitGatewayConnectionDescription: description,
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteTransitGatewayConnection(ctx context.Context, transitGatewayConnectionId string) error {

	_, _, err := client.sdkClient.TransitGatewayConnectionOpenApiControllerApi.DeleteTransitGatewayConnection(ctx, client.config.ProjectId, transitGatewayConnectionId)

	return err

}

// <-------------------------------
