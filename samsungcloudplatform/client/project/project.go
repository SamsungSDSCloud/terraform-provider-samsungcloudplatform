package project

import (
	"context"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/project"
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

func (client *Client) GetProjectInfo(ctx context.Context) (project.ProjectDetailResponseV3, error) {
	result, _, err := client.sdkClient.ProjectV3ControllerApi.DetailProject(ctx, client.config.ProjectId, client.config.ProjectId)
	return result, err
}

func (client *Client) GetProjectList(ctx context.Context, request ListProjectRequest) (project.PageResponseV2ProjectResponseV3, error) {
	result, _, err := client.sdkClient.ProjectV3ControllerApi.ListProjects(ctx, &project.ProjectV3ControllerApiListProjectsOpts{
		AccountName:          optional.NewString(request.AccountName),
		BillYearMonth:        optional.NewString(request.BillYearMonth),
		IsBillingInfoDemand:  optional.NewBool(request.IsBillingInfoDemand),
		IsResourceInfoDemand: optional.NewBool(request.IsResourceInfoDemand),
		IsUserInfoDemand:     optional.NewBool(request.IsUserInfoDemand),
		ProjectName:          optional.NewString(request.ProjectName),
		CreatedByEmail:       optional.NewString(request.CreatedByEmail),
		Page:                 optional.NewInt32(0),
		Size:                 optional.NewInt32(10000),
	})

	return result, err
}

func (client *Client) GetAccountList(ctx context.Context) (project.ListResponseV2AccountResponseV3, error) {
	result, _, err := client.sdkClient.AccountV3ControllerApi.ListAccountsByMyProject(ctx)
	return result, err
}

func (client *Client) GetProductResourceList(ctx context.Context, productCategoryId optional.String) (project.ProjectResponseProductCategoryResource, error) {
	result, _, err := client.sdkClient.ProjectControllerV2Api.ListProductResources(ctx, &project.ProjectControllerV2ApiListProductResourcesOpts{
		ProductCategoryId: productCategoryId,
	})

	return result, err
}

func (client *Client) GetProjectZoneList(ctx context.Context, projectId string) (project.ListResponseV2ZoneResponseV3, error) {
	result, _, err := client.sdkClient.ZoneV3ControllerApi.ListServiceZonesOfProject(ctx, client.config.ProjectId, projectId)

	return result, err
}

func (client *Client) GetProjectProductsList(ctx context.Context, projectId string, code optional.String) (project.ProjectResponseProductCategoryV2, error) {
	result, _, err := client.sdkClient.ProjectControllerV2Api.ListProjectProducts(ctx, client.config.ProjectId, projectId, &project.ProjectControllerV2ApiListProjectProductsOpts{
		LanguageCode: code,
	})

	return result, err
}

func (client *Client) GetProjectProductResourcesList(ctx context.Context, projectId string, productCategoryId optional.String) (project.ProjectResponseProductCategoryResource, error) {
	result, _, err := client.sdkClient.ProjectControllerV2Api.ListProjectProductResources(ctx, client.config.ProjectId, projectId, &project.ProjectControllerV2ApiListProjectProductResourcesOpts{
		ProductCategoryId: productCategoryId,
	})

	return result, err
}
