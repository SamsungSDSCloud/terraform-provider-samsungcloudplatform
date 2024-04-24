package baremetalvdc

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	baremetalvdc "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/bare-metal-server-vdc"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *baremetalvdc.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: baremetalvdc.NewAPIClient(config),
	}
}

// 생성 요청(V1)
func (client *Client) CreateBareMetalServerVDC(ctx context.Context, request BMVDCServerCreateRequest, tags map[string]interface{}) (baremetalvdc.AsyncResponse, error) {
	blockStorages := make([]baremetalvdc.VxLanAdditionalBlockStorageCreateRequest, 0)
	for _, s := range request.ServerDetails[0].StorageDetails {
		blockStorage := baremetalvdc.VxLanAdditionalBlockStorageCreateRequest{
			BareMetalBlockStorageName:   s.BareMetalBlockStorageName,
			BareMetalBlockStorageSize:   s.BareMetalBlockStorageSize,
			BareMetalBlockStorageType:   s.BareMetalBlockStorageType,
			BareMetalBlockStorageTypeId: s.BareMetalBlockStorageTypeId,
			EncryptionEnabled:           &s.EncryptionEnabled,
		}
		blockStorages = append(blockStorages, blockStorage)
	}

	serverDetails := make([]baremetalvdc.VxLanBmServerDetailsRequest, 0)

	for index, _ := range request.ServerDetails {

		serverDetail := baremetalvdc.VxLanBmServerDetailsRequest{
			BareMetalServerName: request.ServerDetails[index].BareMetalServerName,
			DnsEnabled:          &request.ServerDetails[index].DnsEnabled,
			IpAddress:           request.ServerDetails[index].IpAddress,
			ServerTypeId:        request.ServerDetails[index].ServerTypeId,
			StorageDetails:      blockStorages,
			UseHyperThreading:   request.ServerDetails[index].UseHyperThreading,
		}
		serverDetails = append(serverDetails, serverDetail)
	}

	result, _, err := client.sdkClient.VxLanBareMetalServerCreateDeleteOpenApiControllerApi.CreateVdcBareMetalServers(ctx, client.config.ProjectId, baremetalvdc.VxLanBmServerCreateRequest{
		BlockId:                   request.BlockId,
		ContractId:                request.ContractId,
		DeletionProtectionEnabled: &request.DeletionProtectionEnabled,
		ImageId:                   request.ImageId,
		InitScript:                request.InitScript,
		ProductGroupId:            request.ProductGroupId,
		OsUserId:                  request.OsUserId,
		OsUserPassword:            request.OsUserPassword,
		ServerDetails:             serverDetails,
		ServiceZoneId:             request.ServiceZoneId,
		SubnetId:                  request.SubnetId,
		Tags:                      client.sdkClient.ToTagRequestList(tags),
		VdcId:                     request.VdcId,
	})

	return result, err
}

// 다건 해지 요청(V1)
func (client *Client) DeleteBareMetalServersVDC(ctx context.Context, serverIds []string) (baremetalvdc.AsyncResponse, error) {

	result, _, err := client.sdkClient.VxLanBareMetalServerCreateDeleteOpenApiControllerApi.DeleteVdcBareMetalServers(ctx, client.config.ProjectId, baremetalvdc.BaremetalServersTerminateRequest{
		BareMetalServerIds: serverIds,
	})

	return result, err
}

// 상세 조회(V1)
func (client *Client) GetBareMetalServerDetailVDC(ctx context.Context, serverId string) (baremetalvdc.VxLanBareMetalServerDetailResponse, int, error) {
	result, c, err := client.sdkClient.VxLanBareMetalServerSimpleTaskOpenApiControllerApi.DetailVdcBareMetalServer(ctx, client.config.ProjectId, serverId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

// 목록 조회(V1)
func (client *Client) GetBareMetalServersVDC(ctx context.Context, serverName, ipAddress string) (baremetalvdc.ListResponseVxLanBmServerGridResponse, int, error) {
	result, c, err := client.sdkClient.VxLanBareMetalServerSimpleTaskOpenApiControllerApi.ListVdcBareMetalServers(ctx, client.config.ProjectId, &baremetalvdc.VxLanBareMetalServerSimpleTaskOpenApiControllerApiListVdcBareMetalServersOpts{
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
