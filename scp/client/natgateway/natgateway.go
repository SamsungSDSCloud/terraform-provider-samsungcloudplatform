package natgateway

import (
	"context"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/client"
	natgateway2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/nat-gateway2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *natgateway2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: natgateway2.NewAPIClient(config),
	}
}

func (client *Client) CreateNatGateway(ctx context.Context, publicIpAddressId string, subnetId string, description string) (natgateway2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.CreateNatGateway(ctx, client.config.ProjectId, natgateway2.NatGatewayCreateRequest{
		PublicIpAddressId:     publicIpAddressId,
		SubnetId:              subnetId,
		NatGatewayDescription: description,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetNatGateway(ctx context.Context, natGatewayId string) (natgateway2.NatGatewayDetailResponse, int, error) {
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.DetailNatGateway(ctx, client.config.ProjectId, natGatewayId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateNatGateway(ctx context.Context, natGatewayId string, description string) (natgateway2.NatGatewayDetailResponse, int, error) {
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.UpdateNatGatewayDescription(ctx, client.config.ProjectId, natGatewayId, natgateway2.NatGatewayDescriptionUpdateRequest{
		NatGatewayDescription: description,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteNatGateway(ctx context.Context, natGatewayId string) (natgateway2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.DeleteNatGateway(ctx, client.config.ProjectId, natGatewayId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListNatGateway(ctx context.Context, request *natgateway2.NatGatewayV2ControllerV2ApiListNatGateways1Opts) (natgateway2.ListResponseOfNatGatewayListItemResponse, int, error) {
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.ListNatGateways1(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
