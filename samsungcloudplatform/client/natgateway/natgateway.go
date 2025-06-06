package natgateway

import (
	"context"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	natgateway2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/nat-gateway2"
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

func (client *Client) CreateNatGateway(ctx context.Context, publicIpAddressId string, subnetId string, description string, tags map[string]interface{}) (natgateway2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.CreateNatGateway(ctx, client.config.ProjectId, natgateway2.NatGatewayCreateRequest{
		PublicIpAddressId:     publicIpAddressId,
		SubnetId:              subnetId,
		NatGatewayDescription: description,
		Tags:                  client.sdkClient.ToTagRequestList(tags),
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

func (client *Client) ListNatGateway(ctx context.Context, request *natgateway2.NatGatewayV2ControllerV2ApiListNatGatewaysOpts) (natgateway2.ListResponseNatGatewayListItemResponse, int, error) {
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.ListNatGateways(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
