package kubernetesengine

import (
	"context"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	kubernetesengine2 "github.com/ScpDevTerra/trf-sdk/library/kubernetes-engine2"
	"io/ioutil"
)

type Client struct {
	config *sdk.Configuration
	sdk    *kubernetesengine2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config: config,
		sdk:    kubernetesengine2.NewAPIClient(config),
	}
}

func (client *Client) CreateEngine(ctx context.Context, request CreateEngineRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.CreateKubernetesEngineV2(
		ctx,
		client.config.ProjectId,
		kubernetesengine2.KubernetesEngineCreateRequest{
			CloudLoggingEnabled:  request.CloudLoggingEnabled,
			K8sVersion:           request.K8sVersion,
			KubernetesEngineName: request.KubernetesEngineName,
			LbId:                 request.LbId,
			PublicAclIpAddress:   request.PublicAclIpAddress,
			SecurityGroupId:      request.SecurityGroupId,
			SubnetId:             request.SubnetId,
			VolumeId:             request.VolumeId,
			//CifsVolumeId:         request.CifsVolumeId,
			VpcId:  request.VpcId,
			ZoneId: request.ZoneId,
		})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadEngine(ctx context.Context, id string) (kubernetesengine2.ResponseKubernetesEngine, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.DetailKubernetesEngineV2(ctx, client.config.ProjectId, id)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateEngine(ctx context.Context, id string, request UpdateEngineRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.UpdateKubernetesEngineV2(ctx, client.config.ProjectId, id, kubernetesengine2.KubernetesEngineUpdateRequest{
		//CloudLoggingEnabled: request.CloudLoggingEnabled,
		K8sVersion:         request.K8sVersion,
		PublicAclIpAddress: request.PublicAclIpAddress,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteEngine(ctx context.Context, id string) (int, error) {
	_, response, err := client.sdk.K8sEngineV2Api.DeleteKubernetesEngineV2(ctx, client.config.ProjectId, id)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return statusCode, err
}

func (client *Client) GetEngineList(ctx context.Context, request *kubernetesengine2.K8sEngineV2ApiListKubernetesEnginesV2Opts) (kubernetesengine2.PageResponseOfKubernetesEnginesResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.ListKubernetesEnginesV2(ctx, client.config.ProjectId, request)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetKubeConfig(ctx context.Context, id string) (string, int, error) {
	response, err := client.sdk.K8sEngineV2Api.DownloadKubernetesEngineConfigV2(ctx, client.config.ProjectId, id, "public")
	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", -1, err
	}
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return string(data), statusCode, err
}

func (client *Client) CheckUsableSubnet(ctx context.Context, subnetId string, vpcId string) (kubernetesengine2.CheckResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.CheckUsableSubnetV2(ctx, client.config.ProjectId, subnetId, vpcId)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetEngineVersionList(ctx context.Context, request *kubernetesengine2.K8sTemplateV2ApiListKubernetesVersionV2Opts) (kubernetesengine2.PageResponseOfK8sVersionsResponse, int, error) {
	result, response, err := client.sdk.K8sTemplateV2Api.ListKubernetesVersionV2(ctx, client.config.ProjectId, request)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateNodePool(ctx context.Context, engineId string, request CreateNodePoolRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.NodePoolV2ControllerApi.CreateNodePoolV2(
		ctx,
		client.config.ProjectId,
		engineId,
		kubernetesengine2.NodePoolCreateRequest{
			AutoRecovery:     request.AutoRecovery,
			AutoScale:        request.AutoScale,
			ContractId:       request.ContractId,
			DesiredNodeCount: request.DesiredNodeCount,
			ImageId:          request.ImageId,
			MaxNodeCount:     request.MaxNodeCount,
			MinNodeCount:     request.MinNodeCount,
			NodePoolName:     request.NodePoolName,
			ProductGroupId:   request.ProductGroupId,
			ScaleId:          request.ScaleId,
			ServiceLevelId:   request.ServiceLevelId,
			StorageId:        request.StorageId,
			StorageSize:      request.StorageSize,
		})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadNodePool(ctx context.Context, engineId string, nodePoolId string) (kubernetesengine2.ResourceNodePool, int, error) {
	result, response, err := client.sdk.NodePoolV2ControllerApi.DetailNodePoolV2(ctx, client.config.ProjectId, engineId, nodePoolId)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateNodePool(ctx context.Context, engineId string, nodePoolId string, request NodePoolUpdateRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.NodePoolV2ControllerApi.UpdateNodePoolV2(ctx, client.config.ProjectId, engineId, nodePoolId, kubernetesengine2.NodePoolUpdateRequest{
		AutoRecovery:     request.AutoRecovery,
		AutoScale:        request.AutoScale,
		ContractId:       request.ContractId,
		DesiredNodeCount: request.DesiredNodeCount,
		ImageId:          request.ImageId,
		MaxNodeCount:     request.MaxNodeCount,
		MinNodeCount:     request.MinNodeCount,
		NodePoolName:     request.NodePoolName,
		ProductGroupId:   request.ProductGroupId,
		ScaleId:          request.ScaleId,
		ServiceLevelId:   request.ServiceLevelId,
		StorageId:        request.StorageId,
		StorageSize:      request.StorageSize,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteNodePool(ctx context.Context, engineId string, nodePoolId string) (int, error) {
	response, err := client.sdk.NodePoolV2ControllerApi.DeleteNodePoolV2(ctx, client.config.ProjectId, engineId, nodePoolId)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return statusCode, err
}

func (client *Client) GetNodePoolList(ctx context.Context, kubernetesEngineId string, request *kubernetesengine2.NodePoolV2ControllerApiListNodePoolsV2Opts) (kubernetesengine2.PageResponseOfNodePoolsResponse, int, error) {
	result, response, err := client.sdk.NodePoolV2ControllerApi.ListNodePoolsV2(ctx, client.config.ProjectId, kubernetesEngineId, request)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}
