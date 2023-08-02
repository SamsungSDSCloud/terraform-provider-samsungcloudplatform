package publicip

import (
	"context"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/client"
	publicip2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/public-ip2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *publicip2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: publicip2.NewAPIClient(config),
	}
}

/*
func (client *Client) GetAvailablePublicIpList(ctx context.Context, reservedIpPurpose string, uplinkType string, serviceZoneId string) (publicip2.ListResponseOfPublicIpAvailableResponse, error) {
	result, _, err := client.sdkClient.PublicIpOpenApiControllerApi.ListAvailableIpsV2(ctx, client.config.ProjectId, reservedIpPurpose, serviceZoneId, uplinkType)
	return result, err
}
*/

func (client *Client) GetPublicIpList(ctx context.Context, serviceZoneId string, param *publicip2.PublicIpOpenApiControllerApiListPublicIpsV2Opts) (publicip2.ListResponseOfDetailPublicIpResponse, error) {
	result, _, err := client.sdkClient.PublicIpOpenApiControllerApi.ListPublicIpsV2(ctx, client.config.ProjectId, param)
	return result, err
}

func (client *Client) GetPublicIp(ctx context.Context, publicIpAddressId string) (publicip2.DetailPublicIpResponse, int, error) {
	result, c, err := client.sdkClient.PublicIpOpenApiControllerApi.DetailPublicIpV2(ctx, client.config.ProjectId, publicIpAddressId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreatePublicIp(ctx context.Context, productGroupId string, publicIpPurpose string, serviceZoneId string, uplinkType string, publicIpDescription string) (publicip2.DetailPublicIpResponse, error) {
	result, _, err := client.sdkClient.PublicIpOpenApiControllerApi.CreatePublicIpV2(ctx, client.config.ProjectId, publicip2.CreatePublicIpRequest{
		ProductGroupId:             productGroupId,
		PublicIpPurpose:            publicIpPurpose,
		ServiceZoneId:              serviceZoneId,
		UplinkType:                 uplinkType,
		PublicIpAddressDescription: publicIpDescription,
	})
	return result, err
}

func (client *Client) UpdatePublicIp(ctx context.Context, publicIpAddressId string, publicIpDescription string) (error, error) {
	_, err := client.sdkClient.PublicIpOpenApiControllerApi.ChangePublicIpV2(ctx, client.config.ProjectId, publicIpAddressId, publicip2.ChangePublicIpRequest{
		PublicIpAddressDescription: publicIpDescription,
	})
	return err, nil
}

func (client *Client) DeletePublicIp(ctx context.Context, publicIpAddressId string) error {
	_, err := client.sdkClient.PublicIpOpenApiControllerApi.DeletePublicIpV2(ctx, client.config.ProjectId, publicIpAddressId)
	return err
}
