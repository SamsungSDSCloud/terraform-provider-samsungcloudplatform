package subnet

import (
	"context"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/subnet2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *subnet2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: subnet2.NewAPIClient(config),
	}
}

func (client *Client) CreateSubnet(ctx context.Context, vpcId string, cidrIpv4Block string, subnetType string, name string, description string) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.CreateSubnetV2(ctx, client.config.ProjectId, subnet2.CreateSubnetRequest{
		SubnetCidrBlock:   cidrIpv4Block,
		SubnetName:        name,
		SubnetType:        subnetType,
		VpcId:             vpcId,
		SubnetDescription: description,
	})
	return result, err
}

func (client *Client) GetSubnet(ctx context.Context, subnetId string) (subnet2.SubnetDetailResVo, int, error) {
	result, c, err := client.sdkClient.SubnetOpenApiControllerApi.DetailSubnetV2(ctx, client.config.ProjectId, subnetId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetSubnetList(ctx context.Context, request ListSubnetRequest) (subnet2.ListResponseOfSubnetListItemResVo, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.ListSubnetV21(ctx, client.config.ProjectId, &subnet2.SubnetOpenApiControllerApiListSubnetV21Opts{
		SubnetCidrBlock: optional.NewString(request.SubnetCidrBlock),
		SubnetId:        optional.NewString(request.SubnetId),
		SubnetName:      optional.NewString(request.SubnetName),
		//SubnetTypes:     optional.NewInterface(request.SubnetTypes), //TODO:
		VpcId:     optional.NewString(request.VpcId),
		CreatedBy: optional.NewString(request.CreatedBy),
		Page:      optional.NewInt32(request.Page),
		Size:      optional.NewInt32(request.Size),
	})
	return result, err
}

func (client *Client) GetSubnetResourcesV2List(ctx context.Context, request ListSubnetResourceRequest) (subnet2.ListResponseOfSubnetResourceIpListItemResVo, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.ListSubnetResourcesV2(ctx, client.config.ProjectId, request.SubnetId, &subnet2.SubnetOpenApiControllerApiListSubnetResourcesV2Opts{
		IpAddress:        optional.NewString(request.IpAddress),
		LinkedObjectType: optional.NewString(request.LinkedObjectType),
		Page:             optional.NewInt32(request.Page),
		Size:             optional.NewInt32(request.Size),
	})
	return result, err
}

func (client *Client) UpdateSubnetDescription(ctx context.Context, subnetId string, description string) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.UpdateSubnetDescriptionV2(ctx, client.config.ProjectId, subnetId, subnet2.UpdateSubnetDescriptionRequest{
		SubnetDescription: description,
	})
	return result, err
}

func (client *Client) UpdateSubnetType(ctx context.Context, subnetId string, subnetType string) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.UpdateSubnetTypeV2(ctx, client.config.ProjectId, subnetId, subnet2.UpdateSubnetTypeRequest{
		SubnetType: subnetType,
	})
	return result, err
}

func (client *Client) DeleteSubnet(ctx context.Context, subnetId string) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.DeleteSubnetV2(ctx, client.config.ProjectId, subnetId)
	return result, err
}

func (client *Client) CheckSubnetName(ctx context.Context, name string) (bool, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.CheckSubnetNameDuplicationV2(ctx, client.config.ProjectId, name)
	return result.Result, err
}

func (client *Client) CheckSubnetCidrIpv4(ctx context.Context, subnetCidrBlock string, vpcId string) (bool, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.CheckSubnetCidrBlockDuplicationV2(ctx, client.config.ProjectId, subnetCidrBlock, vpcId)
	return result.Result, err
}
