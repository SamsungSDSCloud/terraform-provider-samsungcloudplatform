package project

import (
	"context"

	sdk "github.com/ScpDevTerra/trf-sdk/client"
	"github.com/ScpDevTerra/trf-sdk/library/project"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *project.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: project.NewAPIClient(config),
	}
}

func (client *Client) GetProjectInfo(ctx context.Context) (project.ProjectDetailV2, error) {
	result, _, err := client.sdkClient.ProjectControllerV2Api.DetailProject(ctx, client.config.ProjectId, client.config.ProjectId)
	return result, err
}

func (client *Client) GetProjectList(ctx context.Context, request ListProjectRequest) (project.ProjectResponseOfProjectV2, error) {
	result, _, err := client.sdkClient.ProjectControllerV2Api.ListProjects(ctx, &project.ProjectControllerV2ApiListProjectsOpts{
		AccessLevel:         optional.NewString(request.AccessLevel),
		ActionName:          optional.NewString(request.ActionName),
		CmpServiceName:      optional.NewString(request.CmpServiceName),
		IsUserAuthorization: optional.NewBool(request.IsUserAuthorization),
	})
	return result, err
}

func (client *Client) GetAccountList(ctx context.Context, request ListAccountRequest) (project.ProjectResponseOfAccountV2, error) {
	result, _, err := client.sdkClient.ProjectControllerV2Api.ListAccountsByMyProject(ctx, &project.ProjectControllerV2ApiListAccountsByMyProjectOpts{
		AccessLevel:         optional.NewString(request.AccessLevel),
		ActionName:          optional.NewString(request.ActionName),
		CmpServiceName:      optional.NewString(request.CmpServiceName),
		IsUserAuthorization: optional.NewBool(request.IsUserAuthorization),
		MyProject:           optional.NewBool(request.MyProject),
	})
	return result, err
}
