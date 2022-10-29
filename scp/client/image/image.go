package image

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/image2"
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

func (client *Client) GetStandardImageList(ctx context.Context, zoneId string, imageState string, servicedGroupFor string, servicedFor string) (image2.ListResponseOfStandardImageResponse, error) {
	result, _, err := client.sdkClient.StandardImageControllerApi.ListStandardImages(ctx, zoneId, &image2.StandardImageControllerApiListStandardImagesOpts{
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
	result, _, err := client.sdkClient.StandardImageControllerApi.DetailStandardImage1(ctx, standardImageId)
	return result, err
}
