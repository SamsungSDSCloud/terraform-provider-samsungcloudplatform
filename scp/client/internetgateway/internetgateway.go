package internetgateway

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	internetgateway2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/internet-gateway2"
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
func (client *Client) CreateInternetGateway(ctx context.Context, vpcId string, igwType string, description string, useFirewall bool, useFirewallLog bool) (internetgateway2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.InternetGatewayV4ControllerApi.CreateInternetGateway2(ctx, client.config.ProjectId, internetgateway2.InternetGatewayCreateV4Request{
		VpcId:                      vpcId,
		InternetGatewayType:        igwType,
		InternetGatewayDescription: description,
		FirewallEnabled:            &useFirewall,
		FirewallLoggable:           &useFirewallLog,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetInternetGateway(ctx context.Context, internetGatewayId string) (internetgateway2.InternetGatewayDetailResponse, int, error) {
	result, c, err := client.sdkClient.InternetGatewayV2ControllerV2Api.DetailInternetGateway(ctx, client.config.ProjectId, internetGatewayId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateInternetGateway(ctx context.Context, internetGatewayId string, description string) (internetgateway2.InternetGatewayDetailResponse, int, error) {
	result, c, err := client.sdkClient.InternetGatewayV2ControllerV2Api.ModifyInternetGatewayDescription(ctx, client.config.ProjectId, internetGatewayId, internetgateway2.InternetGatewayModifyDescriptionRequest{
		InternetGatewayDescription: description,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteInternetGateway(ctx context.Context, internetGatewayId string) (internetgateway2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.InternetGatewayV2ControllerV2Api.DeleteInternetGateway(ctx, client.config.ProjectId, internetGatewayId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetInternetGatewayList(ctx context.Context, request *internetgateway2.InternetGatewayV2ControllerV2ApiListInternetGateways1Opts) (internetgateway2.ListResponseOfInternetGatewayListItemResponse, int, error) {
	result, c, err := client.sdkClient.InternetGatewayV2ControllerV2Api.ListInternetGateways1(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
