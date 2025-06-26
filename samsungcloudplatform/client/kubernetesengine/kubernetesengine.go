package kubernetesengine

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	kubernetesengine2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/kubernetes-engine2"
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
		kubernetesengine2.ClusterCreateV2Request{
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
			Tags:                 client.sdk.ToTagRequestList(request.Tags),
		})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadEngine(ctx context.Context, id string) (kubernetesengine2.ClusterV4Response, int, error) {
	result, response, err := client.sdk.K8sEngineV4Api.DetailKubernetesEngineV4(ctx, client.config.ProjectId, id)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdatePublicEndpointAccessControlEngine(ctx context.Context, id string, request UpdatePublicEndpointAccessControlRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV3Api.UpdateKubernetesEnginePublicEndpointAccessControlV3(ctx, client.config.ProjectId, id, kubernetesengine2.ClusterPublicEndpointAccessV3UpdateRequest{
		PublicEndpointAccessControlIp: request.PublicAclIpAddress,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpgradeEngine(ctx context.Context, id string, request UpgradeRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV3Api.UpgradeKubernetesEngineVersionV3(ctx, client.config.ProjectId, id, kubernetesengine2.ClusterUpgradeV3Request{
		UpgradeK8sVersion: request.K8sVersion,
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
	result, response, err := client.sdk.K8sEngineV3Api.UpdateKubernetesEngineLoadBalancerV3(ctx, client.config.ProjectId, id, kubernetesengine2.ClusterLoadBalancerUpdateV3Request{
		LoadBalancerId: request.LbId,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateCifsVolumeEngine(ctx context.Context, id string, request UpdateEngineCifsVolumeRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV3Api.UpdateKubernetesEngineCifsVolumeV3(ctx, client.config.ProjectId, id, kubernetesengine2.ClusterCifsVolumeUpdateV3Request{
		CifsVolumeId: request.CifsVolumeId,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdatePrivateAclEngine(ctx context.Context, id string, request UpdateEnginePrivateAclRequest) (kubernetesengine2.AsyncResponse, int, error) {
	privateAclResources := make([]kubernetesengine2.PrivateEndpointAccessControlResourceV3Vo, 0)
	for _, resource := range request.PrivateAclResourcesToUpdate {
		privateAclResources = append(privateAclResources, kubernetesengine2.PrivateEndpointAccessControlResourceV3Vo{
			ResourceId:   resource.ResourceId,
			ResourceName: resource.ResourceValue,
			ResourceType: resource.ResourceType,
		})
	}

	result, response, err := client.sdk.K8sEngineV3Api.UpdateKubernetesEnginePrivateAclV3(ctx, client.config.ProjectId, id, kubernetesengine2.ClusterPrivateEndpointAccessControlUpdateV3Request{
		PrivateEndpointAccessControlResourceList: privateAclResources,
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

func (client *Client) GetEngineList(ctx context.Context, request *kubernetesengine2.K8sEngineV2ApiListKubernetesEnginesV2Opts) (kubernetesengine2.PageResponseClustersV2Response, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.ListKubernetesEnginesV2(ctx, client.config.ProjectId, request)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetKubeConfig(ctx context.Context, id string, kubeconfigType string) (string, int, error) {
	result, _, err := client.sdk.K8sEngineV2Api.DownloadKubernetesEngineConfigV2(ctx, client.config.ProjectId, id, kubeconfigType)
	if err != nil {
		return "", -1, err
	}
	var statusCode int
	return string(result), statusCode, err
}

func (client *Client) CheckDuplicatedKubernetesEngineName(ctx context.Context, kubernetesEngineName string) (bool, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.CheckKubernetesEngineNameV2(ctx, client.config.ProjectId, kubernetesEngineName)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CheckUsableSubnet(ctx context.Context, subnetId string, vpcId string) (kubernetesengine2.CheckResponse, int, error) {
	result, response, err := client.sdk.K8sEngineV2Api.CheckUsableSubnetV2(ctx, client.config.ProjectId, subnetId, vpcId)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetEngineVersionList(ctx context.Context, request *kubernetesengine2.K8sTemplateV2ApiListKubernetesVersionV21Opts) (kubernetesengine2.PageResponseK8sVersionWithProjectIdResponse, int, error) {
	result, response, err := client.sdk.K8sTemplateV2Api.ListKubernetesVersionV21(ctx, client.config.ProjectId, request)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateNodePool(ctx context.Context, engineId string, request CreateNodePoolRequest) (kubernetesengine2.AsyncResponse, int, error) {

	labels := make([]kubernetesengine2.NodePoolLabelVo, 0)
	for _, label := range request.Labels {
		labels = append(labels, kubernetesengine2.NodePoolLabelVo{
			Key:   label.Key,
			Value: label.Value,
		})
	}

	taints := make([]kubernetesengine2.NodePoolTaintVo, 0)
	for _, taint := range request.Taints {
		taints = append(taints, kubernetesengine2.NodePoolTaintVo{
			Effect: taint.Effect,
			Key:    taint.Key,
			Value:  taint.Value,
		})
	}

	result, response, err := client.sdk.NodePoolV4Api.CreateNodePoolV4(
		ctx,
		client.config.ProjectId,
		engineId,
		kubernetesengine2.NodePoolCreateV4Request{
			AutoRecovery:         &request.AutoRecovery,
			AutoScale:            &request.AutoScale,
			AvailabilityZoneName: request.AvailabilityZoneName,
			ContractName:         request.ContractName,
			DesiredNodeCount:     request.DesiredNodeCount,
			EncryptEnabled:       &request.EncryptEnabled,
			ImageId:              request.ImageId,
			MaxNodeCount:         request.MaxNodeCount,
			MinNodeCount:         request.MinNodeCount,
			NodePoolName:         request.NodePoolName,
			ScaleName:            request.ScaleName,
			ServerType:           request.ServerType,
			ServiceLevelName:     request.ServiceLevelName,
			StorageName:          request.StorageName,
			StorageSize:          request.StorageSize,
			Labels:               labels,
			Taints:               taints,
		})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadNodePool(ctx context.Context, engineId string, nodePoolId string) (kubernetesengine2.NodePoolV2Response, int, error) {
	result, response, err := client.sdk.NodePoolV2Api.DetailNodePoolV2(ctx, client.config.ProjectId, engineId, nodePoolId)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateNodePool(ctx context.Context, engineId string, nodePoolId string, request NodePoolUpdateRequest) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.NodePoolV4Api.UpdateNodePoolV4(ctx, client.config.ProjectId, engineId, nodePoolId, kubernetesengine2.NodePoolUpdateV4Request{
		AutoRecovery:     &request.AutoRecovery,
		AutoScale:        &request.AutoScale,
		DesiredNodeCount: request.DesiredNodeCount,
		MaxNodeCount:     request.MaxNodeCount,
		MinNodeCount:     request.MinNodeCount,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpgradeNodePool(ctx context.Context, engineId string, nodePoolId string) (kubernetesengine2.AsyncResponse, int, error) {
	result, response, err := client.sdk.NodePoolV2Api.UpgradeNodePoolV2(ctx, client.config.ProjectId, engineId, nodePoolId)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateNodePoolLabels(ctx context.Context, engineId string, nodePoolId string, request UpdateNodePoolLabelsRequest) (kubernetesengine2.AsyncResponse, int, error) {

	labels := make([]kubernetesengine2.NodePoolLabelVo, 0)
	for _, label := range request.LabelRequestToUpdate {
		labels = append(labels, kubernetesengine2.NodePoolLabelVo{
			Key:   label.Key,
			Value: label.Value,
		})
	}

	result, response, err := client.sdk.NodePoolV2Api.UpdateNodePoolLabelsV2(ctx, client.config.ProjectId, engineId, nodePoolId, kubernetesengine2.NodePoolLabelsUpdateV2Request{
		Labels: labels,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateNodePoolTaints(ctx context.Context, engineId string, nodePoolId string, request UpdateNodePoolTatintsRequest) (kubernetesengine2.AsyncResponse, int, error) {

	taints := make([]kubernetesengine2.NodePoolTaintVo, 0)
	for _, taint := range request.TaintRequestToUpdate {
		taints = append(taints, kubernetesengine2.NodePoolTaintVo{
			Effect: taint.Effect,
			Key:    taint.Key,
			Value:  taint.Value,
		})
	}

	result, response, err := client.sdk.NodePoolV2Api.UpdateNodePoolTaintsV2(ctx, client.config.ProjectId, engineId, nodePoolId, kubernetesengine2.NodePoolTaintsUpdateV2Request{
		Taints: taints,
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

func (client *Client) GetNodePoolList(ctx context.Context, kubernetesEngineId string, request *kubernetesengine2.NodePoolV2ApiListNodePoolsV2Opts) (kubernetesengine2.PageResponseNodePoolsV2Response, int, error) {
	result, response, err := client.sdk.NodePoolV2Api.ListNodePoolsV2(ctx, client.config.ProjectId, kubernetesEngineId, request)

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}
