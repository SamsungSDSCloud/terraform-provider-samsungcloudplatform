package gslb

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	gslb2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/gslb2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *gslb2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: gslb2.NewAPIClient(config),
	}
}

func (client *Client) CreateGslb(ctx context.Context, request CreateGslbRequest) (gslb2.AsyncResponse, int, error) {
	gslbResources := make([]gslb2.GslbResourceMappingRequest, 0)
	for _, gslbResource := range request.GslbResources {
		gslbResources = append(gslbResources, gslb2.GslbResourceMappingRequest{
			GslbDestination:         gslbResource.GslbDestination,
			GslbRegion:              gslbResource.GslbRegion,
			GslbResourceWeight:      gslbResource.GslbResourceWeight,
			GslbResourceDescription: gslbResource.GslbResourceDescription,
		})
	}

	result, c, err := client.sdkClient.GslbOpenApiV3ControllerApi.CreateGslb1(ctx, client.config.ProjectId, gslb2.CreateGslbServiceV3Request{
		GslbName:      request.GslbName,
		GslbEnvUsage:  request.GslbEnvUsage,
		GslbAlgorithm: request.GslbAlgorithm,
		GslbHealthCheck: &gslb2.GslbHealthCheckReqVo1{
			Protocol:                    request.GslbHealthCheck.Protocol,
			GslbHealthCheckInterval:     request.GslbHealthCheck.GslbHealthCheckInterval,
			GslbHealthCheckTimeout:      request.GslbHealthCheck.GslbHealthCheckTimeout,
			ProbeTimeout:                request.GslbHealthCheck.ProbeTimeout,
			ServicePort:                 request.GslbHealthCheck.ServicePort,
			GslbHealthCheckUserId:       request.GslbHealthCheck.GslbHealthCheckUserId,
			GslbHealthCheckUserPassword: request.GslbHealthCheck.GslbHealthCheckUserPassword,
			GslbSendString:              request.GslbHealthCheck.GslbSendString,
			GslbResponseString:          request.GslbHealthCheck.GslbResponseString,
		},
		GslbResources: gslbResources,
		Tags:          client.sdkClient.ToTagRequestList(request.Tags),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetGslb(ctx context.Context, gslbId string) (gslb2.GslbServiceDetailResponse, int, error) {
	result, c, err := client.sdkClient.GslbOpenApiV2ControllerApi.DetailGslb(ctx, client.config.ProjectId, gslbId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetGslbResource(ctx context.Context, gslbId string) (gslb2.ListResponseOfGslbResourceMappingResponse, error) {
	result, _, err := client.sdkClient.GslbOpenApiV2ControllerApi.ListGslbResources(ctx, client.config.ProjectId, gslbId, nil)
	return result, err
}

func (client *Client) UpdateGslbAlgorithm(ctx context.Context, gslbId string, algorithm string) (gslb2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.GslbOpenApiV2ControllerApi.UpdateGslbAlgorithm(ctx, client.config.ProjectId, gslbId, gslb2.ChangeGslbAlgorithmRequest{
		GslbAlgorithm: algorithm,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateGslbHealthCheck(ctx context.Context, gslbId string, request gslb2.ChangeGslbHealthCheckRequest) (gslb2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.GslbOpenApiV2ControllerApi.UpdateGslbHealthCheck(ctx, client.config.ProjectId, gslbId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateGslbResources(ctx context.Context, gslbId string, gslbResources []gslb2.GslbResourceMappingRequest) (gslb2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.GslbOpenApiV2ControllerApi.UpdateGslbResources(ctx, client.config.ProjectId, gslbId, gslb2.ChangeGslbResourceMappingRequest{
		GslbResources: gslbResources,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteGslb(ctx context.Context, gslbId string) (gslb2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.GslbOpenApiV2ControllerApi.DeleteGslb(ctx, client.config.ProjectId, gslbId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetGslbList(ctx context.Context, request *gslb2.GslbOpenApiV2ControllerApiListGslbsOpts) (gslb2.ListResponseOfGslbServiceListItemResponse, int, error) {
	result, c, err := client.sdkClient.GslbOpenApiV2ControllerApi.ListGslbs(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
