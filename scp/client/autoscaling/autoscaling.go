package autoscaling

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	autoscaling2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/autoscaling2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *autoscaling2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: autoscaling2.NewAPIClient(config),
	}
}

func (client *Client) CreateAutoScalingGroup(ctx context.Context, request autoscaling2.AutoScalingGroupCreateV4Request, tags map[string]interface{}) (autoscaling2.AutoScalingGroupResponse, int, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, c, err := client.sdkClient.AutoScalingGroupV4Api.CreateAsgV4(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetAutoScalingGroupDetail(ctx context.Context, asgId string) (autoscaling2.AutoScalingGroupResponse, int, error) {
	result, c, err := client.sdkClient.AutoScalingGroupV3Api.GetAsgDetailV3(ctx, client.config.ProjectId, asgId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetAutoScalingGroupList(ctx context.Context, request *autoscaling2.AutoScalingGroupV2ApiGetAsgListV2Opts) (autoscaling2.ListResponseOfAutoScalingGroupResponse, int, error) {
	result, c, err := client.sdkClient.AutoScalingGroupV2Api.GetAsgListV2(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateAutoScalingGroup(ctx context.Context, asgId string, request autoscaling2.AutoScalingGroupUpdateRequest) (autoscaling2.AutoScalingGroupResponse, int, error) {
	result, c, err := client.sdkClient.AutoScalingGroupV3Api.UpdateAsgV3(ctx, client.config.ProjectId, asgId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateAutoScalingGroupServerCount(ctx context.Context, asgId string, request autoscaling2.AsgServerCountUpdateRequest) (autoscaling2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.AutoScalingGroupV2Api.UpdateAsgServerCountV2(ctx, client.config.ProjectId, asgId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateAutoScalingGroupLoadBalancer(ctx context.Context, asgId string, request autoscaling2.AsgLoadBalancersUpdateRequest) (autoscaling2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.AsgLoadBalancerV2Api.UpdateAsgLoadBalancersV2(ctx, client.config.ProjectId, asgId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteAutoScalingGroup(ctx context.Context, asgId string) (int, error) {
	c, err := client.sdkClient.AutoScalingGroupV2Api.DeleteAsgV2(ctx, client.config.ProjectId, asgId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return statusCode, err
}

func (client *Client) GetLaunchConfigurationList(ctx context.Context, request *autoscaling2.AsgLaunchConfigurationV2ApiGetLaunchConfigListV2Opts) (autoscaling2.ListResponseOfLaunchConfigListItemResponse, int, error) {
	result, c, err := client.sdkClient.AsgLaunchConfigurationV2Api.GetLaunchConfigListV2(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetLaunchConfigurationDetail(ctx context.Context, lcId string) (autoscaling2.LaunchConfigDetailV4Response, int, error) {
	result, c, err := client.sdkClient.AsgLaunchConfigurationV4Api.GetLaunchConfigV4(ctx, client.config.ProjectId, lcId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateLaunchConfigurationGroup(ctx context.Context, request autoscaling2.LaunchConfigCreateV6Request, tags map[string]interface{}) (autoscaling2.LaunchConfigDetailV4Response, int, error) {
	request.Tags = client.sdkClient.ToTagRequestList(tags)
	result, c, err := client.sdkClient.AsgLaunchConfigurationV6Api.CreateLaunchConfigV6(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteLaunchConfigurationGroup(ctx context.Context, lcId string) (int, error) {
	c, err := client.sdkClient.AsgLaunchConfigurationV2Api.DeleteLaunchConfigV2(ctx, client.config.ProjectId, lcId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return statusCode, err
}

func (client *Client) CreateAutoScalingGroupPolicy(ctx context.Context, asgId string, request autoscaling2.AsgPolicyCreateRequest) (autoscaling2.AsgPolicyResponse, int, error) {
	result, c, err := client.sdkClient.AsgPolicyV2Api.CreateAsgPolicyV2(ctx, client.config.ProjectId, asgId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetAutoScalingGroupPolicyDetail(ctx context.Context, asgId string, policyId string) (autoscaling2.AsgPolicyResponse, int, error) {
	result, c, err := client.sdkClient.AsgPolicyV2Api.GetAsgPolicyDetailV2(ctx, client.config.ProjectId, asgId, policyId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetAutoScalingGroupPolicyList(ctx context.Context, asgId string, request *autoscaling2.AsgPolicyV2ApiGetAsgPolicyListV2Opts) (autoscaling2.ListResponseOfAsgPolicyResponse, int, error) {
	result, c, err := client.sdkClient.AsgPolicyV2Api.GetAsgPolicyListV2(ctx, client.config.ProjectId, asgId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateAutoScalingGroupPolicy(ctx context.Context, asgId string, policyId string, request autoscaling2.AsgPolicyUpdateRequest) (autoscaling2.AsgPolicyResponse, int, error) {
	result, c, err := client.sdkClient.AsgPolicyV2Api.UpdateAsgPolicyV2(ctx, client.config.ProjectId, asgId, policyId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteAutoScalingGroupPolicy(ctx context.Context, asgId string, policyId string) (int, error) {
	c, err := client.sdkClient.AsgPolicyV2Api.DeleteAsgPolicyV2(ctx, client.config.ProjectId, asgId, policyId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return statusCode, err
}
