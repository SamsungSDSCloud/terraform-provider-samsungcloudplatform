package placementgroup

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	placementgroup "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/placement-group"
	"github.com/antihax/optional"
)

type Client struct {
	config *sdk.Configuration
	sdk    *placementgroup.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config: config,
		sdk:    placementgroup.NewAPIClient(config),
	}
}

func (client *Client) CreatePlacementGroup(ctx context.Context, request CreateRequest) (placementgroup.PlacementGroupDetailResponse, error) {
	result, _, err := client.sdk.PlacementGroupV1Api.CreatePlacementGroup1(ctx, client.config.ProjectId, placementgroup.PlacementGroupCreateRequest{
		PlacementGroupName:        request.PlacementGroupName,
		ServiceZoneId:             request.ServiceZoneId,
		Tags:                      client.sdk.ToTagRequestList(request.Tags),
		VirtualServerType:         request.VirtualServerType,
		PlacementGroupDescription: request.PlacementGroupDescription,
		AvailabilityZoneName:      request.AvailabilityZoneName,
	})
	return result, err
}

func (client *Client) GetPlacementGroup(ctx context.Context, placemenGroupId string) (placementgroup.PlacementGroupDetailResponse, int, error) {
	result, c, err := client.sdk.PlacementGroupV1Api.DetailPlacementGroup1(ctx, client.config.ProjectId, placemenGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailPlacementGroup(ctx context.Context, id string) (placementgroup.PlacementGroupDetailResponse, error) {
	result, _, err := client.sdk.PlacementGroupV1Api.DetailPlacementGroup1(ctx, client.config.ProjectId, id)
	return result, err
}

func (client *Client) DeletePlacementGroup(ctx context.Context, id string) error {
	_, err := client.sdk.PlacementGroupV1Api.DeletePlacementGroup1(ctx, client.config.ProjectId, id)
	return err
}

func (client *Client) AddPlacementGroupMember(ctx context.Context, placementGroupId string, virtualServerId string) error {
	_, _, err := client.sdk.PlacementGroupOperateV1Api.AddPlacementGroupMember1(ctx, client.config.ProjectId, placementGroupId, virtualServerId)
	return err
}

func (client *Client) RemovePlacementGroupMember(ctx context.Context, placementGroupId string, virtualServerId string) error {
	_, _, err := client.sdk.PlacementGroupOperateV1Api.RemovePlacementGroupMember1(ctx, client.config.ProjectId, placementGroupId, virtualServerId)
	return err
}

func (client *Client) ListPlacementGroups(ctx context.Context, request ListPlacementGroupsRequestParam) (placementgroup.PageResponseV2PlacementGroupsResponse, error) {
	result, _, err := client.sdk.PlacementGroupV1Api.ListPlacementGroups1(ctx, client.config.ProjectId, &placementgroup.PlacementGroupV1ApiListPlacementGroups1Opts{
		PlacementGroupName:      optional.NewString(request.PlacementGroupName),
		PlacementGroupStateList: optional.NewInterface(request.PlacementGroupStateList),
		VirtualServerType:       optional.NewString(request.VirtualServerType),
		ServiceZoneId:           optional.NewString(request.ServiceZoneId),
		CreatedBy:               optional.NewString(request.CreatedBy),
		Page:                    optional.NewInt32(request.Page),
		Size:                    optional.NewInt32(request.Size),
		Sort:                    optional.NewInterface(request.Sort),
	})
	return result, err
}

func (client *Client) UpdatePlacementGroupDescription(ctx context.Context, placementGroupId string, request UpdatePlacementGroupDescriptionRequest) error {
	_, _, err := client.sdk.PlacementGroupV1Api.UpdatePlacementGroupDescription1(ctx, client.config.ProjectId, placementGroupId, placementgroup.PlacementGroupDescriptionUpdateRequest{
		PlacementGroupDescription: request.PlacementGroupDescription,
	})
	return err
}
