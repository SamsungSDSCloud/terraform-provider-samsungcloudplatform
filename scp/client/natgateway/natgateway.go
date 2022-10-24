package natgateway

import (
	sdk "github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/client"
	natgateway2 "github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/library/nat-gateway2"
	"golang.org/x/net/context"
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
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.CreateNatGatewayV2(ctx, client.config.ProjectId, natgateway2.NatGatewayCreateRequest{
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
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.DetailNatGatewayV2(ctx, client.config.ProjectId, natGatewayId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateNatGateway(ctx context.Context, natGatewayId string, description string) (natgateway2.NatGatewayDetailResponse, int, error) {
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.UpdateNatGatewayDescriptionV2(ctx, client.config.ProjectId, natGatewayId, natgateway2.NatGatewayDescriptionUpdateRequest{
		NatGatewayDescription: description,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteNatGateway(ctx context.Context, natGatewayId string) (natgateway2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.NatGatewayV2ControllerV2Api.DeleteNatGatewayV2(ctx, client.config.ProjectId, natGatewayId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

//func (client *Client) ListNatGateway(ctx context.Context) {
//	client.sdkClient.NatGatewayV2ControllerV2Api.ListNatGatewaysV2(ctx, client.config.ProjectId, &natgateway2.NatGatewayV2ControllerV2ApiListNatGatewaysV2Opts{
//		NatGatewayId:    optional.String{},
//		NatGatewayName:  optional.String{},
//		NatGatewayState: optional.String{},
//		SubnetId:        optional.String{},
//		VpcId:           optional.String{},
//		CreatedBy:       optional.String{},
//		Page:            optional.NewInt32(0),
//		Size:            optional.NewInt32(10000),
//		Sort:            optional.Interface{},
//	})
//}
