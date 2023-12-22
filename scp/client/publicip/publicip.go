package publicip

import (
	"context"
	_ "github.com/antihax/optional"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	publicip2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/public-ip2"
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

func (client *Client) GetPublicIps(ctx context.Context, param *publicip2.PublicIpOpenApiV3ControllerApiListPublicIpsV3Opts) (publicip2.ListResponseOfDetailPublicIpResponse, error) {
	result, _, err := client.sdkClient.PublicIpOpenApiV3ControllerApi.ListPublicIpsV3(ctx, client.config.ProjectId, param)
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

func (client *Client) CreatePublicIp(ctx context.Context, serviceZoneId string, uplinkType string, publicIpDescription string, tags map[string]interface{}) (publicip2.DetailPublicIpResponse, error) {
	result, _, err := client.sdkClient.PublicIpOpenApiV4ControllerApi.CreatePublicIpV4(ctx, client.config.ProjectId, publicip2.CreatePublicIpV4Request{
		ServiceZoneId:       serviceZoneId,
		UplinkType:          uplinkType,
		PublicIpDescription: publicIpDescription,
		Tags:                client.sdkClient.ToTagRequestList(tags),
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
