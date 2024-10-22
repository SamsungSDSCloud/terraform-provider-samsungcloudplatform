package migrationimage

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/image2"
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

func (client *Client) GetMigrationImageList(ctx context.Context, request image2.MigrationImageV2ApiListMigrationImagesOpts) (image2.ListResponseMigrationImageResponse, error) {
	result, _, err := client.sdkClient.MigrationImageV2Api.ListMigrationImages(ctx, client.config.ProjectId, &request)
	return result, err
}

func (client *Client) GetMigrationImageInfo(ctx context.Context, migrationImageId string) (image2.MigrationImageResponse, int, error) {
	result, c, err := client.sdkClient.MigrationImageV2Api.DetailMigrationImage(ctx, client.config.ProjectId, migrationImageId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateMigrationImage(ctx context.Context, request image2.MigrationImageCreateRequest, tags map[string]interface{}) (image2.AsyncResponse, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, _, err := client.sdkClient.MigrationImageV2Api.CreateMigrationImage(ctx, client.config.ProjectId, request)
	return result, err
}

func (client *Client) UpdateMigrationImage(ctx context.Context, migrationImageId string, imageDescription string) (image2.MigrationImageResponse, error) {
	result, _, err := client.sdkClient.MigrationImageV2Api.UpdateMigrationImage(ctx, client.config.ProjectId, migrationImageId, image2.MigrationImageUpdateRequest{
		ImageDescription: imageDescription,
	})
	return result, err
}

func (client *Client) DeleteMigrationImage(ctx context.Context, migrationImageId string) (image2.AsyncResponse, error) {
	result, _, err := client.sdkClient.MigrationImageV2Api.DeleteMigrationImage(ctx, client.config.ProjectId, migrationImageId)
	return result, err
}
