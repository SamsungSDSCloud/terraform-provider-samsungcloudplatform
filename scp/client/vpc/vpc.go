package vpc

import (
	"context"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/client"
	"github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/library/vpc2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *vpc2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: vpc2.NewAPIClient(config),
	}
}

func (client *Client) GetVpcInfo(ctx context.Context, vpcId string) (vpc2.DetailVpcResponse, int, error) {
	result, c, err := client.sdkClient.VpcOpenApiControllerApi.DetailVpcV2(ctx, client.config.ProjectId, vpcId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateVpc(ctx context.Context, vpcId string, description string) (vpc2.DetailVpcResponse, error) {
	result, _, err := client.sdkClient.VpcOpenApiControllerApi.UpdateVpcDescriptionV2(ctx, client.config.ProjectId, vpcId, vpc2.ModifyVpcDescriptionRequest{
		VpcDescription: description,
	})
	return result, err
}

func (client *Client) GetVpcList(ctx context.Context) (vpc2.ListResponseOfVpcResponse, error) {
	result, _, err := client.sdkClient.VpcOpenApiControllerApi.ListVpcV21(ctx, client.config.ProjectId, &vpc2.VpcOpenApiControllerApiListVpcV21Opts{
		Size: optional.NewInt32(20),
		Page: optional.NewInt32(0),
	})
	return result, err
}

func (client *Client) GetVpcListV2(ctx context.Context, request ListVpcRequest) (vpc2.ListResponseOfVpcResponse, error) {
	result, _, err := client.sdkClient.VpcOpenApiControllerApi.ListVpcV21(ctx, client.config.ProjectId, &vpc2.VpcOpenApiControllerApiListVpcV21Opts{
		ServiceZoneId: optional.NewString(request.ServiceZoneId),
		VpcId:         optional.NewString(request.VpcId),
		VpcName:       optional.NewString(request.VpcName),
		VpcStates:     optional.NewInterface(request.VpcStates),
		CreatedBy:     optional.NewString(request.CreatedBy),
		Size:          optional.NewInt32(request.Size),
		Page:          optional.NewInt32(request.Page),
	})
	return result, err
}

func (client *Client) CheckVpcName(ctx context.Context, vpcName string) (bool, error) {
	result, _, err := client.sdkClient.VpcOpenApiControllerApi.CheckDuplicationVpcV2(ctx, client.config.ProjectId, vpcName)
	return result.Result, err
}

func (client *Client) CreateVpc(ctx context.Context, vpcName string, vpcDescription string, productGroupId string, serviceZoneId string) (vpc2.AsyncResponse, error) {

	result, _, err := client.sdkClient.VpcOpenApiControllerApi.CreateVpcV2(ctx, client.config.ProjectId, vpc2.CreateVpcRequest{
		VpcName:        vpcName,
		VpcDescription: vpcDescription,
		ServiceZoneId:  serviceZoneId,
		ProductGroupId: productGroupId,
	})
	return result, err
}

func (client *Client) DeleteVpc(ctx context.Context, vpcId string) error {
	_, _, err := client.sdkClient.VpcOpenApiControllerApi.DeleteVpcV2(ctx, client.config.ProjectId, vpcId)
	return err
}
