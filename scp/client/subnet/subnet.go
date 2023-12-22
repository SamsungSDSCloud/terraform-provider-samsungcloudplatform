package subnet

import (
	"context"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/subnet2"
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

func (client *Client) CreateSubnet(ctx context.Context, vpcId string, cidrIpv4Block string, subnetType string, name string, description string, tags map[string]interface{}) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.CreateSubnetV2(ctx, client.config.ProjectId, subnet2.CreateSubnetRequest{
		SubnetCidrBlock:   cidrIpv4Block,
		SubnetName:        name,
		SubnetType:        subnetType,
		VpcId:             vpcId,
		SubnetDescription: description,
		Tags:              client.sdkClient.ToTagRequestList(tags),
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

func (client *Client) GetSubnetList(ctx context.Context, request *subnet2.SubnetOpenApiControllerApiListSubnetV2Opts) (subnet2.ListResponseOfSubnetListItemResVo, int, error) {
	result, c, err := client.sdkClient.SubnetOpenApiControllerApi.ListSubnetV2(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetSubnetResourcesV2List(ctx context.Context, subnetId string, request *subnet2.SubnetVipOpenApiControllerApiListSubnetResourcesV2Opts) (subnet2.ListResponseOfSubnetResourceIpListItemResVo, int, error) {
	result, c, err := client.sdkClient.SubnetVipOpenApiControllerApi.ListSubnetResourcesV2(ctx, client.config.ProjectId, subnetId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetSubnetVipV2List(ctx context.Context, subnetId string, request *subnet2.SubnetVipOpenApiControllerApiListSubnetVipsV2Opts) (subnet2.ListResponseOfSubnetVirtualIpListItemResVo, int, error) {
	result, c, err := client.sdkClient.SubnetVipOpenApiControllerApi.ListSubnetVipsV2(ctx, client.config.ProjectId, subnetId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetSubnetVip(ctx context.Context, subnetId string, vipId string) (subnet2.SubnetVirtualIpDetailResVo, int, error) {
	result, c, err := client.sdkClient.SubnetVipOpenApiControllerApi.DetailSubnetVipsV2(ctx, client.config.ProjectId, subnetId, vipId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateSubnetDescription(ctx context.Context, subnetId string, description string) (subnet2.SubnetDetailResVo, error) {
	result, _, err := client.sdkClient.SubnetOpenApiV3ControllerApi.UpdateSubnetDescriptionV3(ctx, client.config.ProjectId, subnetId, subnet2.UpdateSubnetDescriptionRequest{
		SubnetDescription: description,
	})
	return result, err
}

func (client *Client) UpdateSubnetType(ctx context.Context, subnetId string, subnetType string) (subnet2.SubnetDetailResVo, error) {
	result, _, err := client.sdkClient.SubnetOpenApiV3ControllerApi.UpdateSubnetTypeV3(ctx, client.config.ProjectId, subnetId, subnet2.UpdateSubnetTypeRequest{
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
	if result.Result == nil {
		return false, err
	}
	return *result.Result, err
}

func (client *Client) CheckSubnetCidrIpv4(ctx context.Context, subnetCidrBlock string, vpcId string) (bool, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.CheckSubnetCidrBlockDuplicationV2(ctx, client.config.ProjectId, subnetCidrBlock, vpcId)
	if result.Result == nil {
		return false, err
	}
	return *result.Result, err
}
func (client *Client) CheckAvailableSubnetIp(ctx context.Context, subnetId string, ipAddress string) (subnet2.CheckResponse, error) {
	result, _, err := client.sdkClient.SubnetOpenApiControllerApi.CheckAvailableSubnetIpV2(ctx, client.config.ProjectId, subnetId, ipAddress)

	if err != nil {
		return subnet2.CheckResponse{}, err
	}

	return result, err
}

func (client *Client) GetSubnetAvailableVipV2List(ctx context.Context, subnetId string, request *subnet2.SubnetVipOpenApiControllerApiListAvailableVipsV2Opts) (subnet2.ListResponseOfSubnetVirtualIpAvailableListItemResVo, error) {
	result, _, err := client.sdkClient.SubnetVipOpenApiControllerApi.ListAvailableVipsV2(ctx, client.config.ProjectId, subnetId, request)
	return result, err
}

func (client *Client) ReserveSubnetVipsV2(ctx context.Context, subnetId string, subnetIpId string, vipDescription string) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetVipOpenApiControllerApi.ReserveSubnetVipsV2(ctx, client.config.ProjectId, subnetId, subnetIpId, subnet2.SubnetVirtualIpReserveRequest{
		VipDescription: vipDescription,
	})
	return result, err
}

func (client *Client) ReleaseSubnetVipsV2(ctx context.Context, subnetId string, subnetIpId string) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetVipOpenApiControllerApi.ReleaseSubnetVipsV2(ctx, client.config.ProjectId, subnetId, subnetIpId)
	return result, err
}

func (client *Client) AttachSubnetPublicIp(ctx context.Context, subnetId string, vipId string, publicIpAddressId string) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetVipOpenApiControllerApi.AttachSubnetPublicIpV2(ctx, client.config.ProjectId, subnetId, vipId, subnet2.AttachSubnetPublicIpRequest{
		PublicIpAddressId: publicIpAddressId,
	})
	return result, err
}

func (client *Client) DetachSubnetPublicIp(ctx context.Context, subnetId string, vipId string) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetVipOpenApiControllerApi.DetachSubnetPublicIpV2(ctx, client.config.ProjectId, subnetId, vipId)
	return result, err
}

func (client *Client) AttachSubnetSecurityGroup(ctx context.Context, subnetId string, vipId string, securityGroupId string) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetVipOpenApiControllerApi.AttachSubnetSecurityGroupV2(ctx, client.config.ProjectId, subnetId, vipId, subnet2.AttachSubnetSecurityGroupRequest{
		SecurityGroupId: securityGroupId,
	})
	return result, err
}

func (client *Client) DetachSubnetSecurityGroup(ctx context.Context, subnetId string, vipId string, securityGroupId string) (subnet2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SubnetVipOpenApiControllerApi.DetachSubnetSecurityGroupV2(ctx, client.config.ProjectId, securityGroupId, subnetId, vipId)
	return result, err
}
