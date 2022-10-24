package internetgateway

import (
	"context"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	internetgateway2 "github.com/ScpDevTerra/trf-sdk/library/internet-gateway2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *internetgateway2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: internetgateway2.NewAPIClient(config),
	}
}
func (client *Client) CreateInternetGateway(ctx context.Context, vpcId string, serviceZoneId string, description string, useFirewall bool) (internetgateway2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.InternetGatewayV2ControllerV2Api.CreateInternetGatewayV2(ctx, client.config.ProjectId, internetgateway2.InternetGatewayCreateRequest{
		FirewallEnabled:            useFirewall,
		ServiceZoneId:              serviceZoneId,
		VpcId:                      vpcId,
		InternetGatewayDescription: description,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetInternetGateway(ctx context.Context, internetGatewayId string) (internetgateway2.InternetGatewayDetailResponse, int, error) {
	result, c, err := client.sdkClient.InternetGatewayV2ControllerV2Api.DetailInternetGatewayV2(ctx, client.config.ProjectId, internetGatewayId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateInternetGateway(ctx context.Context, internetGatewayId string, description string) (internetgateway2.InternetGatewayDetailResponse, int, error) {
	result, c, err := client.sdkClient.InternetGatewayV2ControllerV2Api.ModifyInternetGatewayDescriptionV2(ctx, client.config.ProjectId, internetGatewayId, internetgateway2.InternetGatewayModifyDescriptionRequest{
		InternetGatewayDescription: description,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteInternetGateway(ctx context.Context, internetGatewayId string) (internetgateway2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.InternetGatewayV2ControllerV2Api.DeleteInternetGatewayV2(ctx, client.config.ProjectId, internetGatewayId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
