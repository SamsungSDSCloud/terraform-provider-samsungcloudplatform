package kubernetesengine

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	kubernetesengine2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/kubernetes-engine2"
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

	privateAclResources := make([]kubernetesengine2.PrivateEndpointAccessControlResourceVo, 0)
	for _, resource := range request.PrivateAclResources {
		privateAclResources = append(privateAclResources, kubernetesengine2.PrivateEndpointAccessControlResourceVo{
			Id:    resource.Id,
			Type_: resource.Type,
			Value: resource.Value,
		})
	}

	result, response, err := client.sdk.K8sEngineV2Api.CreateKubernetesEngineV2(
		ctx,
		client.config.ProjectId,
		kubernetesengine2.KubernetesEngineCreateRequest{
			CloudLoggingEnabled:  &request.CloudLoggingEnabled,
			K8sVersion:           request.K8sVersion,
			KubernetesEngineName: request.KubernetesEngineName,
			LbId:                 request.LbId,
			PrivateAclResources:  privateAclResources,
			PublicAclIpAddress:   request.PublicAclIpAddress,
			SecurityGroupId:      request.SecurityGroupId,
			SubnetId:             request.SubnetId,
			VolumeId:             request.VolumeId,
			CifsVolumeId:         request.CifsVolumeId,
			VpcId:                request.VpcId,
			ZoneId:               request.ZoneId,
		})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadEngine(ctx context.Context, id string) (kubernetesengine2.ClusterV2Response, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.DetailKubernetesEngineV2(ctx, client.config.ProjectId, id)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateEngine(ctx context.Context, id string, request UpdateEngineRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.UpdateKubernetesEngineV2(ctx, client.config.ProjectId, id, kubernetesengine2.ClusterUpdateV2Request{
		K8sVersion:         request.K8sVersion,
		PublicAclIpAddress: request.PublicAclIpAddress,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateLoggingEngine(ctx context.Context, id string, request UpdateEngineLoggingRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.UpdateKubernetesEngineLoggingV2(ctx, client.config.ProjectId, id, kubernetesengine2.ClusterCloudLoggingUpdateV2Request{
		CloudLoggingEnabled: &request.CloudLoggingEnabled,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateLoadBalancerEngine(ctx context.Context, id string, request UpdateEngineLoadBalancerRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.UpdateKubernetesEngineLoadBalancerV2(ctx, client.config.ProjectId, id, kubernetesengine2.ClusterLoadBalancerUpdateV2Request{
		LoadBalancerEnabled: &request.LoadBalancerEnabled,
		LbId:                request.LbId,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateCifsVolumeEngine(ctx context.Context, id string, request UpdateEngineCifsVolumeRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.UpdateKubernetesEngineCifsVolumeV2(ctx, client.config.ProjectId, id, kubernetesengine2.ClusterCifsVolumeUpdateV2Request{
		CifsVolumeIdEnabled: &request.CifsVolumeIdEnabled,
		CifsVolumeId:        request.CifsVolumeId,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdatePrivateAclEngine(ctx context.Context, id string, request UpdateEnginePrivateAclRequest) (kubernetesengine2.AsyncResponse, int, error) {
	privateAclResources := make([]kubernetesengine2.PrivateEndpointAccessControlResourceVo, 0)
	for _, resource := range request.PrivateAclResources {
		privateAclResources = append(privateAclResources, kubernetesengine2.PrivateEndpointAccessControlResourceVo{
			Id:    resource.Id,
			Type_: resource.Type,
			Value: resource.Value,
		})
	}

	result, response, err := client.sdk.K8sEngineV2Api.UpdateKubernetesEnginePrivateAclV2(ctx, client.config.ProjectId, id, kubernetesengine2.ClusterPrivateEndpointAccessControlUpdateV2Request{
		PrivateAclResources: privateAclResources,
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

func (client *Client) GetEngineList(ctx context.Context, request *kubernetesengine2.K8sEngineV2ApiListKubernetesEnginesV2Opts) (kubernetesengine2.PageResponseOfClustersV2Response, int, error) {
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

func (client *Client) GetEngineVersionList(ctx context.Context, request *kubernetesengine2.K8sTemplateV2ApiListKubernetesVersionV21Opts) (kubernetesengine2.PageResponseOfK8sVersionWithProjectIdResponse, int, error) {
	result, response, err := client.sdk.K8sTemplateV2Api.ListKubernetesVersionV21(ctx, client.config.ProjectId, request)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateNodePool(ctx context.Context, engineId string, request CreateNodePoolRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.NodePoolV2Api.CreateNodePoolV2(
		ctx,
		client.config.ProjectId,
		engineId,
		kubernetesengine2.NodePoolCreateV2Request{
			AvailabilityZoneName: request.AvailabilityZoneName,
			AutoRecovery:         &request.AutoRecovery,
			AutoScale:            &request.AutoScale,
			ContractId:           request.ContractId,
			DesiredNodeCount:     request.DesiredNodeCount,
			ImageId:              request.ImageId,
			MaxNodeCount:         request.MaxNodeCount,
			MinNodeCount:         request.MinNodeCount,
			NodePoolName:         request.NodePoolName,
			ProductGroupId:       request.ProductGroupId,
			ScaleId:              request.ScaleId,
			ServiceLevelId:       request.ServiceLevelId,
			StorageId:            request.StorageId,
			StorageSize:          request.StorageSize,
		})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadNodePool(ctx context.Context, engineId string, nodePoolId string) (kubernetesengine2.PageResponseOfNodePoolsV2Response, int, error) {
	result, response, err := client.sdk.NodePoolV2Api.ListNodePoolsV2(ctx, client.config.ProjectId, engineId, &kubernetesengine2.NodePoolV2ApiListNodePoolsV2Opts{})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

/* 상세 API에 NodePoolName이 없음
func (client *Client) ReadNodePool(ctx context.Context, engineId string, nodePoolId string) (kubernetesengine2.ResourceNodePool, int, error) {
	result, response, err := client.sdk.NodePoolV2ControllerApi.DetailNodePoolV2(ctx, client.config.ProjectId, engineId, nodePoolId)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}
*/

func (client *Client) UpdateNodePool(ctx context.Context, engineId string, nodePoolId string, request NodePoolUpdateRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.NodePoolV2Api.UpdateNodePoolV2(ctx, client.config.ProjectId, engineId, nodePoolId, kubernetesengine2.NodePoolUpdateV2Request{
		AutoRecovery:     &request.AutoRecovery,
		AutoScale:        &request.AutoScale,
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
	response, err := client.sdk.NodePoolV2Api.DeleteNodePoolV2(ctx, client.config.ProjectId, engineId, nodePoolId)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return statusCode, err
}

func (client *Client) GetNodePoolList(ctx context.Context, kubernetesEngineId string, request *kubernetesengine2.NodePoolV2ApiListNodePoolsV2Opts) (kubernetesengine2.PageResponseOfNodePoolsV2Response, int, error) {
	result, response, err := client.sdk.NodePoolV2Api.ListNodePoolsV2(ctx, client.config.ProjectId, kubernetesEngineId, request)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}
