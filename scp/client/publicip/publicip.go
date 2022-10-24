package publicip

import (
	"context"

	"github.com/ScpDevTerra/trf-provider/scp/common"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	publicip2 "github.com/ScpDevTerra/trf-sdk/library/public-ip2"
	"github.com/antihax/optional"
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

func (client *Client) GetPublicIpList(ctx context.Context, serviceZoneId string, param *publicip2.PublicIpOpenApiControllerApiListPublicIpsV21Opts) (publicip2.ListResponseOfDetailPublicIpResponse, error) {
	result, _, err := client.sdkClient.PublicIpOpenApiControllerApi.ListPublicIpsV21(ctx, client.config.ProjectId, param)
	return result, err
}

func (client *Client) GetPublicIpListV2(ctx context.Context, request ListPublicIpRequest) (publicip2.ListResponseOfDetailPublicIpResponse, error) {
	result, _, err := client.sdkClient.PublicIpOpenApiControllerApi.ListPublicIpsV21(ctx, client.config.ProjectId, &publicip2.PublicIpOpenApiControllerApiListPublicIpsV21Opts{
		//IpAddress:       optional.NewString(request.IpAddress),
		//IsBillable:      optional.NewBool(request.IsBillable),
		//IsViewable:      optional.NewBool(request.IsViewable),
		//PublicIpPurpose: optional.NewString(request.PublicIpPurpose),
		//PublicIpState:   optional.NewString(request.PublicIpState),
		//UplinkType:      optional.NewString(request.UplinkType),
		//CreatedBy:       optional.NewString(request.CreatedBy),
		//Size:            optional.NewInt32(request.Size),
		//Page:            optional.NewInt32(request.Page),

		IpAddress:       optional.String{},
		IsBillable:      optional.NewBool(true),
		IsViewable:      optional.NewBool(true),
		PublicIpPurpose: optional.NewString(common.VpcPublicIpPurpose),
		PublicIpState:   optional.String{},
		UplinkType:      optional.NewString(common.VpcPublicIpUplinkType),
		CreatedBy:       optional.String{},
		Page:            optional.NewInt32(0),
		Size:            optional.NewInt32(10000),
		Sort:            optional.NewInterface([]string{"createdDt:desc"}),
	})
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
