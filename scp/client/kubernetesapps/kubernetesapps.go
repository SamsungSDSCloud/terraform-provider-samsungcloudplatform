package kubernetesapps

import (
	"context"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	kubernetesapps "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/kubernetes-apps"
	"github.com/antihax/optional"
)

type Client struct {
	config *sdk.Configuration
	sdk    *kubernetesapps.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config: config,
		sdk:    kubernetesapps.NewAPIClient(config),
	}
}

func (client *Client) CreateApps(ctx context.Context, clusterId string, namespace string, imageId string, productGroupId string, name string, additionalParams map[string]interface{}, tags map[string]interface{}) (kubernetesapps.K8sAppsResponse, int, error) {
	result, response, err := client.sdk.ReleaseV1ControllerApi.CreateReleaseV1(ctx, client.config.ProjectId, clusterId, namespace, kubernetesapps.ReleaseCreateRequest{
		ProjectId:        client.config.ProjectId,
		ClusterId:        clusterId,
		ImageId:          imageId,
		NamespaceName:    namespace,
		ProductGroupId:   productGroupId,
		ReleaseName:      name,
		AdditionalParams: additionalParams,
		Tags:             client.sdk.ToTagRequestList(tags),
	})
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadApps(ctx context.Context, id string) (kubernetesapps.K8sAppsResponse, int, error) {
	result, response, err := client.sdk.K8sAppsV1ControllerApi.DetailK8sAppsV1(ctx, client.config.ProjectId, id)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteApps(ctx context.Context, clusterId string, namespace string, id string) (int, error) {
	response, err := client.sdk.ReleaseV1ControllerApi.DeleteReleaseV1(ctx, client.config.ProjectId, clusterId, namespace, id)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return statusCode, err
}

func (client *Client) ReadImage(ctx context.Context, id string) (kubernetesapps.ImagesResponse, int, error) {
	result, response, err := client.sdk.ImageV1ControllerApi.ListImagesV1(ctx, client.config.ProjectId, &kubernetesapps.ImageV1ControllerApiListImagesV1Opts{
		ImageId: optional.NewString(id),
	})

	if len(result.Contents) == 1 {
		return result.Contents[0], response.StatusCode, nil
	}

	return kubernetesapps.ImagesResponse{}, response.StatusCode, err
}

func (client *Client) ListImages(ctx context.Context) ([]kubernetesapps.ImagesResponse, int, error) {
	var images []kubernetesapps.ImagesResponse

	page := 0

	for {
		result, _, err := client.sdk.ImageV1ControllerApi.ListImagesV1(ctx, client.config.ProjectId,
			&kubernetesapps.ImageV1ControllerApiListImagesV1Opts{
				Page: optional.NewInt32(int32(page * 10)),
				Size: optional.NewInt32(10),
				Sort: optional.NewString("imageName:asc"),
			})

		if err != nil {
			return []kubernetesapps.ImagesResponse{}, -1, err
		}

		images = append(images, result.Contents[:]...)

		if len(images) == 0 || len(images) >= int(result.TotalCount) {
			break
		}

		page++
	}

	return images, 200, nil
}

func (client *Client) GetImageList(ctx context.Context, request ListStandardImageRequest) (kubernetesapps.PageResponseOfImagesResponse, error) {
	result, _, err := client.sdk.ImageV1ControllerApi.ListImagesV1(ctx, client.config.ProjectId,
		&kubernetesapps.ImageV1ControllerApiListImagesV1Opts{
			//Category:         optional.NewString(request.Category),
			//ImageId:          optional.NewString(request.ImageId),
			//ImageName:        optional.NewString(request.ImageName),
			//IsCarepack:       optional.NewString(request.IsCarepack),
			//IsNew:            optional.NewString(request.IsNew),
			//IsRecommended:    optional.NewString(request.IsRecommended),
			//PricePolicy:      optional.NewString(request.PricePolicy),
			//ProductGroupName: optional.NewString(request.ProductGroupName),
			Size: optional.NewInt32(request.Size),
			Page: optional.NewInt32(request.Page),
			Sort: optional.NewString("imageName:asc"),
		})

	return result, err
}
