package endpoint

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/endpoint2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *endpoint2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: endpoint2.NewAPIClient(config),
	}
}

func (client *Client) GetEndpoint(ctx context.Context, endpointId string) (endpoint2.EndpointDetailResponse, int, error) {
	result, c, err := client.sdkClient.EndpointOpenApiControllerApi.DetailEndpoint(ctx, client.config.ProjectId, endpointId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateEndpoint(ctx context.Context, endpointId string, description string) (endpoint2.EndpointDetailResponse, error) {
	result, _, err := client.sdkClient.EndpointOpenApiControllerApi.ModifyEndpointDescription(ctx, client.config.ProjectId, endpointId, endpoint2.EndpointModifyDescriptionRequest{
		EndpointDescription: description,
	})
	return result, err
}

func (client *Client) GetEndpointList(ctx context.Context, request *endpoint2.EndpointOpenApiControllerApiListEndpointOpts) (endpoint2.ListResponseEndpointResponse, int, error) {
	result, c, err := client.sdkClient.EndpointOpenApiControllerApi.ListEndpoint(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

/*func (client *Client) CreateEndpoint(ctx context.Context, request CreateEndpointRequest) (endpoint2.AsyncResponse, error) {
	tags := make([]endpoint2.TagRequest, 0)
	for _, tag := range request.Tags {
		tags = append(tags, endpoint2.TagRequest{
			TagKey:   tag.TagKey,
			TagValue: tag.TagValue,
		})
	}

	result, _, err := client.sdkClient.EndpointOpenApiV3ControllerApi.CreateEndpointV3(ctx, client.config.ProjectId, endpoint2.EndpointCreateV3Request{
		EndpointIpAddress:   request.EndpointIpAddress,
		EndpointName:        request.EndpointName,
		EndpointType:        request.EndpointType,
		ObjectId:            request.ObjectId,
		ServiceZoneId:       request.ServiceZoneId,
		VpcId:               request.VpcId,
		EndpointDescription: request.EndpointDescription,
		Tags:                tags,
	})
	return result, err
}*/

func (client *Client) CreateEndpoint(ctx context.Context, endpointIpAddress string, endpointName string, endpointType string, objectId string, vpcId string, endpointDescription string, serviceZoneId string, tags map[string]interface{}) (endpoint2.AsyncResponse, error) {
	result, _, err := client.sdkClient.EndpointOpenApiV3ControllerApi.CreateEndpointV3(ctx, client.config.ProjectId, endpoint2.EndpointCreateV3Request{
		EndpointIpAddress:   endpointIpAddress,
		EndpointName:        endpointName,
		EndpointType:        endpointType,
		ObjectId:            objectId,
		ServiceZoneId:       serviceZoneId,
		VpcId:               vpcId,
		EndpointDescription: endpointDescription,
		Tags:                client.sdkClient.ToTagRequestList(tags),
	})
	return result, err
}

func (client *Client) DeleteEndpoint(ctx context.Context, endpointId string) error {
	_, _, err := client.sdkClient.EndpointOpenApiControllerApi.DeleteEndpoint(ctx, client.config.ProjectId, endpointId)
	return err
}
