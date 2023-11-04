package customimage

import (
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	image "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/image2"
	"github.com/antihax/optional"
	"golang.org/x/net/context"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *image.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: image.NewAPIClient(config),
	}
}

func (client *Client) CreateCustomImage(ctx context.Context, createCustomImageRequest image.CustomImageCreateRequest) (image.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.CustomImageV2Api.CreateCustomImage(ctx, client.config.ProjectId, createCustomImageRequest)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
func (client *Client) GetCustomImage(ctx context.Context, imageId string) (image.CustomImageResponse, int, error) {
	result, c, err := client.sdkClient.CustomImageV2Api.DetailCustomImage1(ctx, client.config.ProjectId, imageId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetCustomImageList(ctx context.Context, request image.CustomImageV2ApiListCustomImagesOpts) (image.ListResponseOfCustomImageResponse, error) {
	if request.Size == optional.NewInt32(0) {
		request.Size = optional.NewInt32(20)
	}
	result, _, err := client.sdkClient.CustomImageV2Api.ListCustomImages(ctx, client.config.ProjectId, &request)
	return result, err
}

func (client *Client) UpdateCustomImageDescription(ctx context.Context, imageId string, imageDescription string) (image.CustomImageResponse, error) {
	result, _, err := client.sdkClient.CustomImageV2Api.UpdateCustomImage(ctx, client.config.ProjectId, imageId, image.CustomImageUpdateRequest{
		ImageDescription: imageDescription,
	})
	return result, err
}

func (client *Client) DeleteCustomImage(ctx context.Context, imageId string) error {
	_, _, err := client.sdkClient.CustomImageV2Api.DeleteCustomImage(ctx, client.config.ProjectId, imageId)
	return err
}
