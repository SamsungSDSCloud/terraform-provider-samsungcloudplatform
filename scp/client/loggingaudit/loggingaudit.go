package loggingaudit

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/client"
	loggingaudit "github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/library/logging-audit"
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

func (client *Client) CreateTrail(ctx context.Context, request CreateTrailRequest) (loggingaudit.TrailResponseModel, error) {

	var users []loggingaudit.UserResponse
	for _, b := range request.LoggingTargetUsers {
		users = append(users, loggingaudit.UserResponse{
			UserId: b.UserId,
		})
	}

	result, _, err := client.sdkClient.TrailControllerApi.CreateTrail(ctx, client.config.ProjectId, loggingaudit.CreateTrailRequest{
		TrailName:                  request.TrailName,
		ObsBucketId:                request.ObsBucketId,
		IsLoggingTargetAllUser:     request.IsLoggingTargetAllUser,
		LoggingTargetUsers:         users,
		IsLoggingTargetAllResource: request.IsLoggingTargetAllResource,
		LoggingTargetResourceIds:   request.LoggingTargetResourceIds,
		TrailSaveType:              request.TrailSaveType,
		TrailDescription:           request.TrailDescription,
	})

	return result, err
}

func (client *Client) ReadTrail(ctx context.Context, trailId string) (loggingaudit.TrailResponseModel, int, error) {
	result, c, err := client.sdkClient.TrailControllerApi.DetailTrail(ctx, client.config.ProjectId, trailId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
