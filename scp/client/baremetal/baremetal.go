package baremetal

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/client"
	baremetal "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/bare-metal-server"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *baremetal.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: baremetal.NewAPIClient(config),
	}
}

func (client *Client) GetBareMetalServers(ctx context.Context, serverName, ipAddress string) (baremetal.ListResponseOfBareMetalServerResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalServerSimpleTaskOpenApiControllerApi.ListBareMetalServers(ctx, client.config.ProjectId, &baremetal.BareMetalServerSimpleTaskOpenApiControllerApiListBareMetalServersOpts{
		BareMetalServerName: optional.NewString(serverName),
		IpAddress:           optional.NewString(ipAddress),
		Page:                optional.NewInt32(0),
		Size:                optional.NewInt32(10000),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetBareMetalServerDetail(ctx context.Context, serverId string) (baremetal.BareMetalServerDetailResponse, int, error) {
	result, c, err := client.sdkClient.BareMetalServerSimpleTaskOpenApiControllerApi.DetailBareMetalServer(ctx, client.config.ProjectId, serverId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateBareMetalServer(ctx context.Context, request BMServerCreateRequest) (baremetal.AsyncResponse, error) {
	blockStorages := make([]baremetal.AdditionalBlockStorageCreateRequest, 0)
	for _, s := range request.ServerDetails[0].StorageDetails {
		ab := baremetal.AdditionalBlockStorageCreateRequest{
			BareMetalBlockStorageTypeId: s.BareMetalBlockStorageTypeId,
			BareMetalBlockStorageSize:   s.BareMetalBlockStorageSize,
			BareMetalBlockStorageType:   s.BareMetalBlockStorageType,
			BareMetalBlockStorageName:   s.BareMetalBlockStorageName,
			EncryptionEnabled:           &s.EncryptionEnabled,
		}
		blockStorages = append(blockStorages, ab)
	}
	serverDetails := make([]baremetal.BareMetalServerDetailsRequest, 0)
	serverDetail := baremetal.BareMetalServerDetailsRequest{
		BareMetalLocalSubnetEnabled:   &request.ServerDetails[0].BareMetalLocalSubnetEnabled,
		BareMetalLocalSubnetId:        request.ServerDetails[0].BareMetalLocalSubnetId,
		BareMetalLocalSubnetIpAddress: request.ServerDetails[0].BareMetalLocalSubnetIpAddress,
		BareMetalServerName:           request.ServerDetails[0].BareMetalServerName,
		DnsEnabled:                    &request.ServerDetails[0].DnsEnabled,
		IpAddress:                     request.ServerDetails[0].IpAddress,
		NatEnabled:                    &request.ServerDetails[0].NatEnabled,
		PublicIpAddressId:             request.ServerDetails[0].PublicIpAddressId,
		ServerTypeId:                  request.ServerDetails[0].ServerTypeId,
		StorageDetails:                blockStorages,
		UseHyperThreading:             request.ServerDetails[0].UseHyperThreading,
	}
	serverDetails = append(serverDetails, serverDetail)

	result, _, err := client.sdkClient.BareMetalServerCreateDeleteOpenApiControllerApi.CreateBareMetalServer(ctx, client.config.ProjectId, baremetal.BmServerCreateRequest{
		BlockId:                   request.BlockId,
		ContractId:                request.ContractId,
		DeletionProtectionEnabled: &request.DeletionProtectionEnabled,
		ImageId:                   request.ImageId,
		InitScript:                request.InitScript,
		OsUserId:                  request.OsUserId,
		OsUserPassword:            request.OsUserPassword,
		ProductGroupId:            request.ProductGroupId,
		ServerDetails:             serverDetails,
		ServiceZoneId:             request.ServiceZoneId,
		SubnetId:                  request.SubnetId,
		Tags:                      []baremetal.TagRequest{},
		VpcId:                     request.VpcId,
	})
	return result, err
}

func (client *Client) AttachBMLocalSubnet(ctx context.Context, serverId string, subnetId string, subnetIp string) (baremetal.AttachLocalSubnetResponse, error) {
	ipType := "AUTO"
	if subnetIp != "" {
		ipType = "MANUAL"
	}
	result, _, err := client.sdkClient.BareMetalServerLocalSubnetOpenApiControllerApi.AttachLocalSubnet(ctx, client.config.ProjectId, serverId, baremetal.AttachLocalSubnetRequest{
		IpAddress:              subnetIp,
		BareMetalLocalSubnetId: subnetId,
		IpType:                 ipType,
	})

	return result, err
}

func (client *Client) DetachBMLocalSubnet(ctx context.Context, serverId string) (baremetal.DetachLocalSubnetOutSo, error) {
	result, _, err := client.sdkClient.BareMetalServerLocalSubnetOpenApiControllerApi.DetachLocalSubnet(ctx, client.config.ProjectId, serverId)

	return result, err
}

func (client *Client) ChangeBMContract(ctx context.Context, serverId string, contractId string) (baremetal.BareMetalServerContractPeriodUpdateResponse, error) {
	result, _, err := client.sdkClient.BareMetalServerSimpleTaskOpenApiControllerApi.UpdateContractPeriod(ctx, client.config.ProjectId, serverId, baremetal.BmServerContractPeriodUpdateRequest{
		ContractId: contractId,
	})

	return result, err
}

func (client *Client) ChangeBMDeletePolicy(ctx context.Context, serverId string, deleteProtectionEnabled string) (baremetal.BareMetalServerDetailResponse, error) {
	result, _, err := client.sdkClient.BareMetalServerSimpleTaskOpenApiControllerApi.UpdateBareMetalServerDeletionProtectionEnabled(ctx, client.config.ProjectId, serverId, baremetal.BmServerDeletionProtectionEnabledUpdateRequest{
		DeletionProtectionEnabled: deleteProtectionEnabled,
	})

	return result, err
}

func (client *Client) DeleteBareMetalServer(ctx context.Context, serverId string) (baremetal.AsyncResponse, error) {
	result, _, err := client.sdkClient.BareMetalServerCreateDeleteOpenApiControllerApi.DeleteBareMetalServer(ctx, client.config.ProjectId, serverId)

	return result, err
}

func (client *Client) EnableBMNat(ctx context.Context, serverId string, natIpAddressType string, publicIpAddressId string) (baremetal.BareMetalServerPublicNatResponse, error) {
	result, _, err := client.sdkClient.BareMetalServerStaticNatOpenApiControllerApi.AssignPublicNat(ctx, client.config.ProjectId, serverId, baremetal.BareMetalServerAssignPublicNatRequest{
		NatIpAddress:      "",
		NatEnabled:        "Y",
		NatIpAddressType:  natIpAddressType,
		PublicIpAddressId: publicIpAddressId,
	})

	return result, err
}

func (client *Client) DisableBMNat(ctx context.Context, serverId string) (baremetal.BareMetalServerPublicNatResponse, error) {
	result, _, err := client.sdkClient.BareMetalServerStaticNatOpenApiControllerApi.DeletePublicNat(ctx, client.config.ProjectId, serverId)

	return result, err
}
