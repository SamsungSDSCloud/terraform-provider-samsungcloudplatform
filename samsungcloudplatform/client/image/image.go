package image

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/image2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *image2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: image2.NewAPIClient(config),
	}
}

func (client *Client) GetStandardImageList(ctx context.Context, zoneId string, imageState string, servicedGroupFor string, servicedFor string) (image2.ListResponseStandardImageResponse, error) {
	result, _, err := client.sdkClient.StandardImageV2Api.ListStandardImages(ctx, client.config.ProjectId, zoneId, &image2.StandardImageV2ApiListStandardImagesOpts{
		ImageState:       optional.NewString(imageState),
		ServicedFor:      optional.NewString(servicedFor),
		ServicedGroupFor: optional.NewString(servicedGroupFor),
		Page:             optional.NewInt32(0),
		Size:             optional.NewInt32(10000),
		Sort:             optional.NewInterface([]string{"imageName:asc"}),
	})
	return result, err
}

func (client *Client) GetStandardImageInfo(ctx context.Context, standardImageId string) (image2.StandardImageResponse, error) {
	result, _, err := client.sdkClient.StandardImageV2Api.DetailStandardImage1(ctx, client.config.ProjectId, standardImageId)
	return result, err
}

func (client *Client) GetImageType(ctx context.Context, imageId string) (string, error) {
	result, _, err := client.sdkClient.CommonImageV2Api.DetailImageType(ctx, client.config.ProjectId, imageId)
	return result.ImageType, err // STANDARD, CUSTOM, MIGRATION
}
