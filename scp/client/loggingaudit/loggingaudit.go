package loggingaudit

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	loggingaudit "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/logging-audit"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *loggingaudit.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: loggingaudit.NewAPIClient(config),
	}
}

func ynStringFromBool(flag bool) string {
	if flag {
		return "Y"
	} else {
		return "N"
	}
}

func (client *Client) CreateTrail(ctx context.Context, tags []interface{}, request CreateTrailRequest) (loggingaudit.TrailResponse, error) {
	var users []loggingaudit.UserResponse
	for _, b := range request.LoggingTargetUsers {
		users = append(users, loggingaudit.UserResponse{UserId: b})
	}

	result, _, err := client.sdkClient.TrailControllerApi.CreateTrail(ctx, client.config.ProjectId, loggingaudit.CreateTrailRequest{
		TagCreateRequests:          toTagRequestList(tags),
		TrailName:                  request.TrailName,
		ObsBucketId:                request.ObsBucketId,
		IsLoggingTargetAllUser:     ynStringFromBool(request.IsLoggingTargetAllUser),
		LoggingTargetUsers:         users,
		IsLoggingTargetAllResource: ynStringFromBool(request.IsLoggingTargetAllResource),
		LoggingTargetResourceIds:   request.LoggingTargetResourceIds,
		IsLoggingTargetAllRegion:   ynStringFromBool(request.IsLoggingTargetAllRegion),
		TrailSaveType:              request.TrailSaveType,
		LoggingTargetRegions:       request.LoggingTargetRegions,
		ValidationYn:               ynStringFromBool(request.UseVerification),
		TrailDescription:           request.TrailDescription,
	})

	return result, err
}

func (client *Client) ReadTrail(ctx context.Context, trailId string) (loggingaudit.TrailResponse, int, error) {
	result, c, err := client.sdkClient.TrailControllerApi.DetailTrail(ctx, client.config.ProjectId, trailId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListTrails(ctx context.Context, isMy *bool, regions []string, resourceIds []string, state string, name string) (loggingaudit.PageResponseV2OfTrailResponse, error) {
	result, _, err := client.sdkClient.TrailControllerApi.ListTrails(ctx, client.config.ProjectId, loggingaudit.TrailSearchCriteria{
		IsMy:                     isMy,
		LoggingTargetRegions:     regions,
		LoggingTargetResourceIds: resourceIds,
		Page:                     0,
		Size:                     10000,
		State:                    state,
		TrailName:                name,
	})

	return result, err
}

func (client *Client) DeleteTrail(ctx context.Context, trailId string) (loggingaudit.TrailStateResponse, int, error) {
	result, c, err := client.sdkClient.TrailControllerApi.DeleteTrail(ctx, client.config.ProjectId, trailId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateTrail(ctx context.Context, trailId string, request UpdateTrailRequest) (loggingaudit.TrailUpdateResponse, error) {
	var users []loggingaudit.UserResponse
	for _, b := range request.LoggingTargetUsers {
		users = append(users, loggingaudit.UserResponse{UserId: b})
	}

	result, _, error := client.sdkClient.TrailControllerApi.UpdateTrail(ctx, client.config.ProjectId, trailId, loggingaudit.UpdateTrailRequest{
		TrailUpdateType:            request.TrailUpdateType,
		IsLoggingTargetAllUser:     request.IsLoggingTargetAllUser,
		LoggingTargetUsers:         users,
		IsLoggingTargetAllResource: request.IsLoggingTargetAllResource,
		LoggingTargetResourceIds:   request.LoggingTargetResourceIds,
		TrailSaveType:              request.TrailSaveType,
		TrailDescription:           request.TrailDescription,
		IsLoggingTargetAllRegion:   request.IsLoggingTargetAllRegion,
		LoggingTargetRegions:       request.LoggingTargetRegions,
		ValidationYn:               ynStringFromBool(request.UseVerification),
	})
	return result, error
}

func (client *Client) CheckTrailName(ctx context.Context, trailName string) (bool, error) {
	result, _, err := client.sdkClient.TrailControllerApi.CheckTrailName(ctx, client.config.ProjectId, client.config.ProjectId, trailName)
	return result, err
}

func (client *Client) StopTrail(ctx context.Context, trailId string) (loggingaudit.TrailStateResponse, int, error) {
	result, c, err := client.sdkClient.TrailControllerApi.StopTrail(ctx, client.config.ProjectId, trailId, loggingaudit.TrailStateRequest{})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) StartTrail(ctx context.Context, trailId string) (loggingaudit.TrailStateResponse, int, error) {
	result, c, err := client.sdkClient.TrailControllerApi.StartTrail(ctx, client.config.ProjectId, trailId, loggingaudit.TrailStateRequest{})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailLogging(ctx context.Context, loggingId string) (loggingaudit.LoggingsResponse, int, error) {
	result, c, err := client.sdkClient.LoggingControllerApi.DetailLogging(ctx, client.config.ProjectId, loggingId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}

	return result, statusCode, err
}

func (client *Client) ListLoggings(ctx context.Context, request loggingaudit.LoggingSearchCriteria) (loggingaudit.PageResponseV2OfLoggingsResponse, int, error) {
	response, c, err := client.sdkClient.LoggingControllerApi.ListLoggings(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func (client *Client) ListUsers(ctx context.Context, userName string) (loggingaudit.PageResponseV2OfMembersResponse, int, error) {
	response, c, err := client.sdkClient.TrailControllerApi.ListUsers(ctx, client.config.ProjectId, &loggingaudit.TrailControllerApiListUsersOpts{
		UserName: optional.NewString(userName),
		Page:     optional.NewInt32(0),
		Size:     optional.NewInt32(10000),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return response, statusCode, err
}

func toTagRequestList(list []interface{}) []loggingaudit.TagRequest {
	if len(list) == 0 {
		return nil
	}
	var result []loggingaudit.TagRequest

	for _, val := range list {
		kv := val.(common.HclKeyValueObject)
		result = append(result, loggingaudit.TagRequest{
			TagKey:   kv["tag_key"].(string),
			TagValue: kv["tag_value"].(string),
		})
	}
	return result
}
