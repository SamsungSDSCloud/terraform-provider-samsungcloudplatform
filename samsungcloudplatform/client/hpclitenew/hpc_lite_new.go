package hpclitenew

import (
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	hpclitenew "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/hpc-lite-new"
	"golang.org/x/net/context"
	"net/http"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *hpclitenew.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: hpclitenew.NewAPIClient(config),
	}
}

func (client *Client) GetHpcLiteNewDetail(ctx context.Context, serverId string) (hpclitenew.HpcLitePlusOpenApiDetailResponseDto, int, error) {
	responseVo, httpResponse, err := client.sdkClient.HpcLitePlusOpenAPIV1ControllerApi.DetailHpcLitePlusV1(ctx, client.config.ProjectId, serverId)
	if err != nil {
		return *new(hpclitenew.HpcLitePlusOpenApiDetailResponseDto), httpResponse.StatusCode, err
	}

	return responseVo, httpResponse.StatusCode, err
}

func (client *Client) CreateHpcLiteNew(ctx context.Context, request HpcLiteNewCreateRequest) (hpclitenew.AsyncListResponse, int, error) {
	serverDetails := toServerDetailsVoList(request.ServerDetails)
	requestVo := hpclitenew.HpcLitePlusOpenApiCreateRequestVo{
		CoServiceZoneId:       request.CoServiceZoneId,
		Contract:              request.Contract,
		HyperThreadingEnabled: request.HyperThreadingEnabled,
		ImageId:               request.ImageId,
		InitScript:            request.InitScript,
		OsUserId:              request.OsUserId,
		OsUserPassword:        request.OsUserPassword,
		ProductGroupId:        request.ProductGroupId,
		ResourcePoolId:        request.ResourcePoolId,
		ServerDetails:         serverDetails,
		ServerType:            request.ServerType,
		ServiceZoneId:         request.ServiceZoneId,
		Tags:                  client.sdkClient.ToTagRequestList(request.Tags),
		VlanPoolCidr:          request.VlanPoolCidr,
	}

	result, c, err := client.sdkClient.HpcLitePlusOpenAPIV1ControllerApi.CreateHpcLitePlusV1(ctx, client.config.ProjectId, requestVo)
	statusCode := getStatusCode(c)

	return result, statusCode, err
}

func (client *Client) DeleteHpcLiteNew(ctx context.Context, request HpcLiteNewDeleteRequest) (hpclitenew.AsyncListResponse, int, error) {
	result, c, err := client.sdkClient.HpcLitePlusOpenAPIV1ControllerApi.DeleteHpcLitePlusV1(ctx, client.config.ProjectId, hpclitenew.HpcLitePlusOpenApiDeleteRequestVo{
		ServerIds:     request.ServerIds,
		ServiceZoneId: request.ServiceZoneId,
	})

	statusCode := getStatusCode(c)
	return result, statusCode, err
}

func toServerDetailsVoList(serverDetailList []ServerDetailRequest) []hpclitenew.ServerDetailRequestVo {
	var ret = []hpclitenew.ServerDetailRequestVo{}
	for _, v := range serverDetailList {
		ret = append(ret, hpclitenew.ServerDetailRequestVo{
			ServerName: v.ServerName,
			IpAddress:  v.IpAddress,
		})
	}
	return ret
}

func getStatusCode(response *http.Response) int {
	var ret int
	if response != nil {
		ret = response.StatusCode
	}
	return ret
}
