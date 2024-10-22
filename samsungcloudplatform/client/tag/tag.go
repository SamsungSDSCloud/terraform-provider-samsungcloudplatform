package tag

import (
	"context"
	"errors"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/tag"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *tag.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: tag.NewAPIClient(config),
	}
}

type Filter struct {
	TagKey    string
	TagValues []string
}

func (client *Client) ListResourceTags(ctx context.Context, resourceId string) (tag.PageResponseV2TagResponse, int, error) {
	result, c, err := client.sdkClient.ResourceTagControllerApi.ListResourceTags(ctx, client.config.ProjectId, resourceId, &tag.ResourceTagControllerApiListResourceTagsOpts{
		Page: optional.NewInt32(0),
		Size: optional.NewInt32(10000),
		Sort: optional.Interface{},
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListResources(ctx context.Context, resourceIds []string, resourceTypeFilters []string, filters []Filter) (tag.PageResponseV2TagsResponse, int, error) {
	tagFilters := make([]tag.TagFilter, 0)
	if len(filters) > 0 {
		for _, f := range filters {
			tagFilters = append(tagFilters, tag.TagFilter{
				TagKey:    f.TagKey,
				TagValues: f.TagValues,
			})
		}
	}

	result, c, err := client.sdkClient.ResourceTagControllerApi.ListResources(ctx, client.config.ProjectId, tag.ResourceSearchCriteria{
		ResourceIds:         resourceIds,
		ResourceTypeFilters: resourceTypeFilters,
		TagFilters:          tagFilters,
		Page:                0,
		Size:                10000,
		Sort:                nil,
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateResourceTag(ctx context.Context, resourceId string, tags []tag.TagRequest) (tag.TagsResponse, int, error) {
	if len(tags) > 50 {
		return tag.TagsResponse{}, 400, errors.New("the number of tags to be updated cannot exceed 50")
	}

	result, c, err := client.sdkClient.ResourceTagControllerApi.UpdateResourceTag(ctx, client.config.ProjectId, resourceId, tag.ResourceTagUpdateRequest{
		Tags: tags,
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AttachResourceTag(ctx context.Context, resourceId string, tags []tag.TagRequest) (tag.TagsResponse, int, error) {
	if len(tags) > 50 {
		return tag.TagsResponse{}, 400, errors.New("the number of tags to be updated cannot exceed 50")
	}

	result, c, err := client.sdkClient.ResourceTagControllerApi.CreateResourceTag(ctx, client.config.ProjectId, resourceId, tag.ResourceTagCreateRequest{
		Tags: tags,
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetachResourceTag(ctx context.Context, resourceId string, tagKey string) (int, error) {
	c, err := client.sdkClient.ResourceTagControllerApi.DeleteResourceTag(ctx, client.config.ProjectId, resourceId, tagKey)

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return statusCode, err
}
