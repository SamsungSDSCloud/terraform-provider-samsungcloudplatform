package virtualserver

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/client"
	virtualserver2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/virtual-server2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *virtualserver2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: virtualserver2.NewAPIClient(config),
	}
}

func (client *Client) GetVirtualServer(ctx context.Context, virtualServerId string) (virtualserver2.DetailVirtualServerResponse, int, error) {
	result, c, err := client.sdkClient.VirtualServerOpenApiControllerApi.DetailVirtualServer(ctx, client.config.ProjectId, virtualServerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetVirtualServerList(ctx context.Context, virtualServerName string) (virtualserver2.ListResponseOfVirtualServersResponse, int, error) {
	var optVirtualServerName optional.String
	if len(virtualServerName) > 0 {
		optVirtualServerName = optional.NewString(virtualServerName)
	}
	result, c, err := client.sdkClient.VirtualServerOpenApiControllerApi.ListVirtualServers(ctx, client.config.ProjectId, &virtualserver2.VirtualServerOpenApiControllerApiListVirtualServersOpts{
		AutoscalingEnabled:   optional.Bool{},
		ServerGroupId:        optional.String{},
		ServicedForList:      optional.Interface{},
		ServicedGroupForList: optional.Interface{},
		VirtualServerName:    optVirtualServerName,
		Page:                 optional.NewInt32(0),
		Size:                 optional.NewInt32(10000),
		Sort:                 optional.NewInterface([]string{"createdDt:desc"}),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateVirtualServer(ctx context.Context, request CreateRequest) (virtualserver2.AsyncResponse, error) {

	var extraBlockStorages []virtualserver2.VirtualServerCreateBlockStorageRequest
	for _, b := range request.ExtraBlockStorages {
		extraBlockStorages = append(extraBlockStorages, virtualserver2.VirtualServerCreateBlockStorageRequest{
			BlockStorageName: b.BlockStorageName,
			DiskSize:         b.DiskSize,
			EncryptEnabled:   b.EncryptEnabled,
			ProductId:        b.ProductId,
		})
	}

	result, _, err := client.sdkClient.VirtualServerCreateDeleteOpenApiControllerApi.CreateVirtualServer(ctx, client.config.ProjectId, virtualserver2.VirtualServerCreateRequest{
		BlockStorage: &virtualserver2.VirtualServerCreateBlockStorageRequest{
			BlockStorageName: request.BlockStorage.BlockStorageName,
			DiskSize:         request.BlockStorage.DiskSize,
			EncryptEnabled:   request.BlockStorage.EncryptEnabled,
			ProductId:        request.BlockStorage.ProductId,
		},
		ContractId:                request.ContractId,
		DeletionProtectionEnabled: request.DeletionProtectionEnabled,
		DnsEnabled:                request.DnsEnabled,
		ExtraBlockStorages:        extraBlockStorages,
		ImageId:                   request.ImageId,
		InitialScript: &virtualserver2.VirtualServerCreateInitialScriptRequest{
			EncodingType:         request.InitialScript.EncodingType,
			InitialScriptContent: request.InitialScript.InitialScriptContent,
			InitialScriptShell:   request.InitialScript.InitialScriptShell,
			InitialScriptType:    request.InitialScript.InitialScriptType,
		},
		LocalSubnet: &virtualserver2.VirtualServerCreateLocalSubnetNicRequest{
			LocalSubnetIpAddress: request.LocalSubnet.LocalSubnetIpAddress,
			SubnetId:             request.LocalSubnet.SubnetId,
		},
		Nic: &virtualserver2.VirtualServerCreateNicRequest{
			InternalIpAddress: request.Nic.InternalIpAddress,
			NatEnabled:        request.Nic.NatEnabled,
			PublicIpId:        request.Nic.PublicIpAddressId,
			SubnetId:          request.Nic.SubnetId,
		},
		OsAdmin: &virtualserver2.VirtualServerCreateOsCredentialRequest{
			OsUserId:       request.OsAdmin.OsUserId,
			OsUserPassword: request.OsAdmin.OsUserPassword,
		},
		ProductGroupId:    request.ProductGroupId,
		SecurityGroupIds:  request.SecurityGroupIds,
		ServerGroupId:     request.ServerGroupId,
		ServerTypeId:      request.ServerTypeId,
		ServiceLevelId:    request.ServiceLevelId,
		ServiceZoneId:     request.ServiceZoneId,
		Timezone:          request.Timezone,
		VirtualServerName: request.VirtualServerName,
	})
	return result, err
}

func (client *Client) DeleteVirtualServer(ctx context.Context, virtualServerId string) (virtualserver2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VirtualServerCreateDeleteOpenApiControllerApi.DeleteVirtualServer1(ctx, client.config.ProjectId, virtualServerId)
	return result, err
}

func (client *Client) GetNicList(ctx context.Context, virtualServerId string) (virtualserver2.ListResponseOfNicResponse, error) {
	result, _, err := client.sdkClient.VirtualServerNicOpenApiControllerApi.ListNics(ctx, client.config.ProjectId, virtualServerId)
	return result, err
}

func (client *Client) AttachLocalSubnet(ctx context.Context, virtualServerId string, localSubnetId string, localSubnetIp string) (virtualserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.VirtualServerNicOpenApiControllerApi.AttachLocalSubnetNic(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerCreateLocalSubnetNicRequest{
		LocalSubnetIpAddress: localSubnetIp,
		SubnetId:             localSubnetId,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachLocalSubnet(ctx context.Context, virtualServerId string, nicId string) (virtualserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.VirtualServerNicOpenApiControllerApi.DetachLocalSubnetNic(ctx, client.config.ProjectId, nicId, virtualServerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateDeleteProtectionEnabled(ctx context.Context, virtualServerId string, isDeleteProtectionEnabled bool) (virtualserver2.DetailVirtualServerResponse, int, error) {
	result, c, err := client.sdkClient.VirtualServerOpenApiControllerApi.UpdateVirtualServerDeletionProtectionEnabled(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerDeletionProtectionEnabledUpdateRequest{
		DeletionProtectionEnabled: isDeleteProtectionEnabled,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateScale(ctx context.Context, virtualServerId string, serverTypeId string) (virtualserver2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VirtualServerOperateOpenApiControllerApi.ResizeVirtualServer(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerResizeRequest{
		ServerTypeId: serverTypeId,
	})
	return result, err
}

func (client *Client) DeleteSecurityGroup(ctx context.Context, virtualServerId string, securityGroupId string) (virtualserver2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VirtualServerOperateOpenApiControllerApi.DetachVirtualServerSecurityGroup(ctx, client.config.ProjectId, securityGroupId, virtualServerId)
	return result, err
}

func (client *Client) AddSecurityGroup(ctx context.Context, virtualServerId string, securityGroupId string) (virtualserver2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VirtualServerOperateOpenApiControllerApi.AttachVirtualServerSecurityGroup(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerSecurityGroupUpdateRequest{
		SecurityGroupId: securityGroupId,
	})
	return result, err
}

func (client *Client) AttachPublicIp(ctx context.Context, virtualServerId string, nicId string, publicIpId string) (virtualserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.VirtualServerNicOpenApiControllerApi.AttachNatIp(ctx, client.config.ProjectId, nicId, virtualServerId, virtualserver2.VirtualServerAttachNatIpRequest{
		PublicIpId: publicIpId,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachPublicIp(ctx context.Context, virtualServerId string, nicId string) (virtualserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.VirtualServerNicOpenApiControllerApi.DetachNatIp(ctx, client.config.ProjectId, nicId, virtualServerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
