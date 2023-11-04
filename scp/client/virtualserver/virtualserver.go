package virtualserver

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	virtualserver2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/virtual-server2"
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

func (client *Client) GetVirtualServer(ctx context.Context, virtualServerId string) (virtualserver2.DetailVirtualServerV3Response, int, error) {
	result, c, err := client.sdkClient.VirtualServerV3Api.DetailVirtualServer1(ctx, client.config.ProjectId, virtualServerId)

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
	result, c, err := client.sdkClient.VirtualServerV2Api.ListVirtualServers2(ctx, client.config.ProjectId, &virtualserver2.VirtualServerV2ApiListVirtualServers2Opts{
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

	var extraBlockStorages []virtualserver2.VirtualServerCreateBlockStorageV3Request
	for _, b := range request.ExtraBlockStorages {
		extraBlockStorages = append(extraBlockStorages, virtualserver2.VirtualServerCreateBlockStorageV3Request{
			BlockStorageName: b.BlockStorageName,
			DiskSize:         b.DiskSize,
			EncryptEnabled:   &b.EncryptEnabled,
			DiskType:         b.DiskType,
		})
	}

	tags := make([]virtualserver2.TagRequest, 0)
	for _, tag := range request.Tags {
		tags = append(tags, virtualserver2.TagRequest{
			TagKey:   tag.TagKey,
			TagValue: tag.TagValue,
		})
	}
	result, _, err := client.sdkClient.VirtualServerCreateV3Api.CreateVirtualServer1(ctx, client.config.ProjectId, virtualserver2.VirtualServerCreateV3Request{
		BlockStorage: &virtualserver2.VirtualServerCreateBlockStorageV3Request{
			BlockStorageName: request.BlockStorage.BlockStorageName,
			DiskSize:         request.BlockStorage.DiskSize,
			EncryptEnabled:   &request.BlockStorage.EncryptEnabled,
			DiskType:         request.BlockStorage.DiskType,
		},
		ContractDiscount:          request.ContractDiscount,
		DeletionProtectionEnabled: &request.DeletionProtectionEnabled,
		DnsEnabled:                &request.DnsEnabled,
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
			NatEnabled:        &request.Nic.NatEnabled,
			PublicIpId:        request.Nic.PublicIpAddressId,
			SubnetId:          request.Nic.SubnetId,
		},
		OsAdmin: &virtualserver2.VirtualServerCreateOsCredentialRequest{
			OsUserId:       request.OsAdmin.OsUserId,
			OsUserPassword: request.OsAdmin.OsUserPassword,
		},
		SecurityGroupIds: request.SecurityGroupIds,
		ServerGroupId:    request.ServerGroupId,
		ServerType:       request.ServerType,
		//ServiceLevelId:    request.ServiceLevelId,	// deprecated
		ServiceZoneId: request.ServiceZoneId,
		//Timezone:          request.Timezone,			// deprecated
		VirtualServerName:    request.VirtualServerName,
		Tags:                 tags,
		AvailabilityZoneName: request.AvailabilityZoneName,
	})
	return result, err
}

func (client *Client) CreateVirtualServerV4(ctx context.Context, request CreateRequest) (virtualserver2.AsyncResponse, error) {

	var extraBlockStorages []virtualserver2.VirtualServerCreateBlockStorageV3Request
	for _, b := range request.ExtraBlockStorages {
		extraBlockStorages = append(extraBlockStorages, virtualserver2.VirtualServerCreateBlockStorageV3Request{
			BlockStorageName: b.BlockStorageName,
			DiskSize:         b.DiskSize,
			EncryptEnabled:   &b.EncryptEnabled,
			DiskType:         b.DiskType,
		})
	}

	tags := make([]virtualserver2.TagRequest, 0)
	for _, tag := range request.Tags {
		tags = append(tags, virtualserver2.TagRequest{
			TagKey:   tag.TagKey,
			TagValue: tag.TagValue,
		})
	}
	result, _, err := client.sdkClient.VirtualServerCreateV4Api.CreateVirtualServer2(ctx, client.config.ProjectId, virtualserver2.VirtualServerCreateV4Request{
		BlockStorage: &virtualserver2.VirtualServerCreateBlockStorageV3Request{
			BlockStorageName: request.BlockStorage.BlockStorageName,
			DiskSize:         request.BlockStorage.DiskSize,
			EncryptEnabled:   &request.BlockStorage.EncryptEnabled,
			DiskType:         request.BlockStorage.DiskType,
		},
		ContractDiscount:          request.ContractDiscount,
		DeletionProtectionEnabled: &request.DeletionProtectionEnabled,
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
			NatEnabled:        &request.Nic.NatEnabled,
			PublicIpId:        request.Nic.PublicIpAddressId,
			SubnetId:          request.Nic.SubnetId,
		},
		KeyPairId:        request.KeyPairId,
		PlacementGroupId: request.PlacementGroupId,
		SecurityGroupIds: request.SecurityGroupIds,
		ServerGroupId:    request.ServerGroupId,
		ServerType:       request.ServerType,
		//ServiceLevelId:    request.ServiceLevelId,	// deprecated
		ServiceZoneId: request.ServiceZoneId,
		//Timezone:          request.Timezone,			// deprecated
		VirtualServerName:    request.VirtualServerName,
		Tags:                 tags,
		AvailabilityZoneName: request.AvailabilityZoneName,
	})
	return result, err
}

func (client *Client) DeleteVirtualServer(ctx context.Context, virtualServerId string) (virtualserver2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VirtualServerCreateDeleteV2Api.DeleteVirtualServer2(ctx, client.config.ProjectId, virtualServerId)
	return result, err
}

func (client *Client) GetNicList(ctx context.Context, virtualServerId string) (virtualserver2.ListResponseOfNicResponse, error) {
	result, _, err := client.sdkClient.VirtualServerNicV2Api.ListNics2(ctx, client.config.ProjectId, virtualServerId)
	return result, err
}

func (client *Client) AttachLocalSubnet(ctx context.Context, virtualServerId string, localSubnetId string, localSubnetIp string) (virtualserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.VirtualServerNicV2Api.AttachLocalSubnetNic(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerCreateLocalSubnetNicRequest{
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
	result, c, err := client.sdkClient.VirtualServerNicV2Api.DetachLocalSubnetNic(ctx, client.config.ProjectId, nicId, virtualServerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateDeleteProtectionEnabled(ctx context.Context, virtualServerId string, isDeleteProtectionEnabled bool) (virtualserver2.DetailVirtualServerV3Response, int, error) {
	result, c, err := client.sdkClient.VirtualServerV3Api.UpdateVirtualServerDeletionProtectionEnabled1(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerDeletionProtectionEnabledUpdateRequest{
		DeletionProtectionEnabled: &isDeleteProtectionEnabled,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateScale(ctx context.Context, virtualServerId string, serverType string) (virtualserver2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VirtualServerOperateV3Api.ResizeVirtualServer1(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerResizeV3Request{
		ServerType: serverType,
	})
	return result, err
}

func (client *Client) DeleteSecurityGroup(ctx context.Context, virtualServerId string, securityGroupId string) (virtualserver2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VirtualServerOperateV2Api.DetachVirtualServerSecurityGroup(ctx, client.config.ProjectId, securityGroupId, virtualServerId)
	return result, err
}

func (client *Client) AddSecurityGroup(ctx context.Context, virtualServerId string, securityGroupId string) (virtualserver2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VirtualServerOperateV2Api.AttachVirtualServerSecurityGroup(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerSecurityGroupUpdateRequest{
		SecurityGroupId: securityGroupId,
	})
	return result, err
}

func (client *Client) AttachPublicIp(ctx context.Context, virtualServerId string, nicId string, publicIpId string) (virtualserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.VirtualServerNicV2Api.AttachNatIp1(ctx, client.config.ProjectId, nicId, virtualServerId, virtualserver2.VirtualServerAttachNatIpRequest{
		PublicIpId: publicIpId,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachPublicIp(ctx context.Context, virtualServerId string, nicId string) (virtualserver2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.VirtualServerNicV2Api.DetachNatIp1(ctx, client.config.ProjectId, nicId, virtualServerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailVirtualServer(ctx context.Context, virtualServerId string) (virtualserver2.DetailVirtualServerV3Response, error) {
	result, _, err := client.sdkClient.VirtualServerV3Api.DetailVirtualServer1(ctx, client.config.ProjectId, virtualServerId)
	return result, err
}

func (client *Client) ListVirtualServers(ctx context.Context, request ListVirtualServersRequestParam) (virtualserver2.ListResponseOfVirtualServersResponse, error) {

	result, _, err := client.sdkClient.VirtualServerV2Api.ListVirtualServers2(ctx, client.config.ProjectId,
		&virtualserver2.VirtualServerV2ApiListVirtualServers2Opts{
			AutoscalingEnabled:   optional.NewBool(request.AutoscalingEnabled),
			ServerGroupId:        optional.NewString(request.ServerGroupId),
			ServicedForList:      optional.NewInterface(request.ServicedForList),
			ServicedGroupForList: optional.NewInterface(request.ServicedGroupForList),
			VirtualServerName:    optional.NewString(request.VirtualServerName),
			Page:                 optional.NewInt32(request.Page),
			Size:                 optional.NewInt32(request.Size),
			Sort:                 optional.NewInterface(request.Sort),
		})
	return result, err
}

func (client *Client) UpdateVirtualServerSubnetIp(ctx context.Context, virtualServerId string, request VirtualServerSubnetIpUpdateRequest) (virtualserver2.AsyncResponse, error) {

	result, _, err := client.sdkClient.VirtualServerNicV2Api.UpdateVirtualServerSubnetIp(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerSubnetIpUpdateRequest{
		InternalIpAddress: request.InternalIpAddress,
		SubnetId:          request.SubnetId,
	})
	return result, err
}

func (client *Client) UpdateVirtualServerContract(ctx context.Context, virtualServerId string, request VirtualServerContractUpdateRequest) (virtualserver2.DetailVirtualServerV3Response, error) {

	result, _, err := client.sdkClient.VirtualServerV4Api.UpdateVirtualServerContract2(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerContractUpdateV4Request{
		ContractDiscount: request.ContractDiscount,
	})

	return result, err
}

func (client *Client) UpdateVirtualServerNextContract(ctx context.Context, virtualServerId string, request VirtualServerContractUpdateRequest) (virtualserver2.DetailVirtualServerV3Response, error) {
	result, _, err := client.sdkClient.VirtualServerV4Api.UpdateVirtualServerNextContract1(ctx, client.config.ProjectId, virtualServerId, virtualserver2.VirtualServerContractUpdateV4Request{
		ContractDiscount: request.ContractDiscount,
	})
	return result, err
}

func (client *Client) ListNics(ctx context.Context, virtualServerId string) (virtualserver2.ListResponseOfNicResponse, error) {
	result, _, err := client.sdkClient.VirtualServerNicV2Api.ListNics2(ctx, client.config.ProjectId, virtualServerId)
	return result, err
}

func (client *Client) StopVirtualServer(ctx context.Context, virtualServerId string) (virtualserver2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VirtualServerOperateV2Api.StopVirtualServer2(ctx, client.config.ProjectId, virtualServerId)
	return result, err
}

func (client *Client) StartVirtualServer(ctx context.Context, virtualServerId string) (virtualserver2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VirtualServerOperateV2Api.StartVirtualServer2(ctx, client.config.ProjectId, virtualServerId)
	return result, err
}
