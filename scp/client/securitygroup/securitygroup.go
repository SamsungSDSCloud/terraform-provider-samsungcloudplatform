package securitygroup

import (
	"context"
	"fmt"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	securitygroup2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/security-group2"
	"github.com/antihax/optional"
	"strings"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *securitygroup2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: securitygroup2.NewAPIClient(config),
	}
}

func (client *Client) GetSecurityGroup(ctx context.Context, securityGroupId string) (securitygroup2.DetailSecurityGroupResponse, int, error) {
	result, c, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DetailSecurityGroupV2(ctx, client.config.ProjectId, securityGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateSecurityGroup(ctx context.Context, serviceZoneId string, vpcId string, name string, description string, loggable bool) (securitygroup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV3Api.CreateSecurityGroupV3(ctx, client.config.ProjectId, securitygroup2.SecurityGroupCreateV3Request{
		SecurityGroupName:        name,
		ServiceZoneId:            serviceZoneId,
		VpcId:                    vpcId,
		SecurityGroupDescription: description,
		Loggable:                 &loggable,
	})
	return result, err
}

func (client *Client) UpdateSecurityGroupDescription(ctx context.Context, securityGroupId string, description string) (securitygroup2.DetailSecurityGroupResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.ModifySecurityGroupDescriptionV2(ctx, client.config.ProjectId, securityGroupId, securitygroup2.SecurityGroupModifyDescriptionRequest{
		SecurityGroupDescription: description,
	})
	return result, err
}

func (client *Client) UpdateSecurityGroupIsLoggable(ctx context.Context, securityGroupId string, isLoggable bool) (securitygroup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.UpdateSecurityGroupLoggingV2(ctx, client.config.ProjectId, securityGroupId, securitygroup2.SecurityGroupChangeLoggableRequest{
		IsLoggable: &isLoggable,
	})
	return result, err
}

func (client *Client) DeleteSecurityGroup(ctx context.Context, securityGroupId string) error {
	_, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DeleteSecurityGroupV2(ctx, client.config.ProjectId, securityGroupId)
	return err
}

func (client *Client) CheckSecurityGroupName(ctx context.Context, name string, vpcId string) (bool, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.CheckSecurityGroupNameDuplicationV2(ctx, client.config.ProjectId, name, vpcId)
	if result.Result == nil {
		return false, err
	}
	return *result.Result, err
}

type SecurityGroupServiceRule = securitygroup2.ServiceVo

type SecurityGroupRule = securitygroup2.SecurityGroupCreateRuleRequest

func (client *Client) CreateSecurityGroupRule(ctx context.Context, securityGroupId string, direction string, addresses []string, description string, services []SecurityGroupServiceRule) (securitygroup2.AsyncResponse, error) {
	if strings.ToUpper(direction) == "IN" {
		result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.CreateSecurityGroupRuleV2(ctx, client.config.ProjectId, securityGroupId, securitygroup2.SecurityGroupCreateRuleRequest{
			RuleDirection:     "IN",
			Services:          services,
			SourceIpAddresses: addresses,
			RuleDescription:   description,
		})
		return result, err
	} else if strings.ToUpper(direction) == "OUT" {
		result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.CreateSecurityGroupRuleV2(ctx, client.config.ProjectId, securityGroupId, securitygroup2.SecurityGroupCreateRuleRequest{
			RuleDirection:          "OUT",
			Services:               services,
			DestinationIpAddresses: addresses,
			RuleDescription:        description,
		})
		return result, err
	}
	return securitygroup2.AsyncResponse{}, fmt.Errorf("Invalid rule direction : " + direction)
}

func (client *Client) CreateSecurityGroupBulkRule(ctx context.Context, securityGroupId string, rules []SecurityGroupRule) (securitygroup2.AsyncResponse, error) {

	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.CreateSecurityGroupBulkRuleV2(ctx, client.config.ProjectId, securityGroupId, securitygroup2.SecurityGroupCreateBulkRuleRequest{
		Rules: rules,
	})

	return result, err
}

func (client *Client) AttachUserIpToSecurityGroup(ctx context.Context, securityGroupId string, userIpType string, userIpAddress string, userIpDescription string) (securitygroup2.AsyncResponse, error) {

	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.AttachUserIpV2(ctx, client.config.ProjectId, securityGroupId, securitygroup2.SecurityGroupUserIpAttachRequest{
		UserIpType:        userIpType,
		UserIpAddress:     userIpAddress,
		UserIpDescription: userIpDescription,
	})

	return result, err
}

func (client *Client) DetachUserIpFromSecurityGroup(ctx context.Context, securityGroupId string, userIpAddress string) (securitygroup2.AsyncResponse, error) {

	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DetachUserIpV2(ctx, client.config.ProjectId, securityGroupId, securitygroup2.SecurityGroupUserIpDetachRequest{
		UserIpAddress: userIpAddress,
	})

	return result, err
}

func (client *Client) GetSecurityGroupRule(ctx context.Context, ruleId string, securityGroupId string) (securitygroup2.SecurityGroupRuleResponse, int, error) {
	result, c, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DetailSecurityGroupRuleV2(ctx, client.config.ProjectId, ruleId, securityGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListSecurityGroupRules(ctx context.Context, securityGroupId string, opts *securitygroup2.SecurityGroupOpenApiControllerV2ApiListSecurityGroupRuleV2Opts) (securitygroup2.ListResponseOfSecurityGroupRuleResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.ListSecurityGroupRuleV2(ctx, client.config.ProjectId, securityGroupId, opts)
	return result, err
}

func (client *Client) UpdateSecurityGroupRule(ctx context.Context, ruleId string, securityGroupId string, direction string, addresses []string, description string, services []SecurityGroupServiceRule) (securitygroup2.AsyncResponse, error) {
	if strings.ToUpper(direction) == "IN" {
		result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.ModifySecurityGroupRuleV2(ctx, client.config.ProjectId, ruleId, securityGroupId, securitygroup2.SecurityGroupRuleModifyRequest{
			SourceIpAddresses: addresses,
			RuleDirection:     "IN",
			Services:          services,
			RuleDescription:   description,
		})
		return result, err
	} else if strings.ToUpper(direction) == "OUT" {
		result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.ModifySecurityGroupRuleV2(ctx, client.config.ProjectId, ruleId, securityGroupId, securitygroup2.SecurityGroupRuleModifyRequest{
			DestinationIpAddresses: addresses,
			RuleDirection:          "OUT",
			Services:               services,
			RuleDescription:        description,
		})
		return result, err
	}
	return securitygroup2.AsyncResponse{}, fmt.Errorf("Invalid rule direction : " + direction)
}

func (client *Client) DeleteSecurityGroupRule(ctx context.Context, securityGroupId string, ruleIds []string) (securitygroup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DeleteSecurityGroupRuleV2(ctx, client.config.ProjectId, securityGroupId, securitygroup2.SecurityGroupRuleDeleteRequest{
		RuleDeletionType: "PARTIAL",
		RuleIds:          ruleIds,
	})
	return result, err
}

func (client *Client) DeleteSecurityGroupRuleAll(ctx context.Context, securityGroupId string) (securitygroup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DeleteSecurityGroupRuleV2(ctx, client.config.ProjectId, securityGroupId, securitygroup2.SecurityGroupRuleDeleteRequest{
		RuleDeletionType: "ALL",
	})
	return result, err
}

func (client *Client) CreateSecurityGroupLogStorage(ctx context.Context, request securitygroup2.SecurityGroupLogStorageCreatRequest) (securitygroup2.SecurityGroupLogStorageDetailResponse, int, error) {
	result, c, err := client.sdkClient.SecurityGroupLogStorageOpenApiControllerV2Api.CreateSecurityGroupLogStorageV2(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateSecurityGroupLogStorage(ctx context.Context, logStorageId string, obsBucketId string) (securitygroup2.SecurityGroupLogStorageDetailResponse, int, error) {
	result, c, err := client.sdkClient.SecurityGroupLogStorageOpenApiControllerV2Api.UpdateSecurityGroupLogStorageV2(ctx, client.config.ProjectId, logStorageId, securitygroup2.SecurityGroupStorageUpdateRequest{
		ObsBucketId: obsBucketId,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetSecurityGroupLogStorage(ctx context.Context, logStorageId string) (securitygroup2.SecurityGroupLogStorageDetailResponse, int, error) {
	result, c, err := client.sdkClient.SecurityGroupLogStorageOpenApiControllerV2Api.DetailSecurityGroupLogStorageV2(ctx, client.config.ProjectId, logStorageId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteSecurityGroupLogStorage(ctx context.Context, logStorageId string) error {
	_, err := client.sdkClient.SecurityGroupLogStorageOpenApiControllerV2Api.DeleteSecurityGroupLogStorageV2(ctx, client.config.ProjectId, logStorageId)
	return err
}

func (client *Client) ListSecurityGroupLogStorages(ctx context.Context, vpcId string, obsBucketId string) (securitygroup2.ListResponseOfSecurityGroupLogStorageDetailResponse, int, error) {
	result, c, err := client.sdkClient.SecurityGroupLogStorageOpenApiControllerV2Api.ListSecurityGroupLogStoragesV2(ctx, client.config.ProjectId, vpcId, &securitygroup2.SecurityGroupLogStorageOpenApiControllerV2ApiListSecurityGroupLogStoragesV2Opts{
		ObsBucketId: optional.NewString(obsBucketId),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListSecurityGroups(ctx context.Context, opts *securitygroup2.SecurityGroupOpenApiControllerV2ApiListSecurityGroupV2Opts) (securitygroup2.ListResponseOfSecurityGroupResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.ListSecurityGroupV2(ctx, client.config.ProjectId, opts)
	return result, err
}

func (client *Client) ListSecurityGroupsByLoggable(ctx context.Context, vpcId string, isLoggable bool) (securitygroup2.ListResponseOfSecurityGroupResponse, error) {

	opts := &securitygroup2.SecurityGroupOpenApiControllerV2ApiListSecurityGroupV2Opts{
		VpcId:      optional.NewString(vpcId),
		IsLoggable: optional.NewBool(isLoggable),
	}
	result, err := client.ListSecurityGroups(ctx, opts)
	return result, err
}

func (client *Client) ListUserIpsBySecurityGroupId(ctx context.Context, securityGroupId string) (securitygroup2.ListResponseOfSecurityGroupUserIpResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.ListUserIpV2(ctx, client.config.ProjectId, securityGroupId)
	return result, err
}
