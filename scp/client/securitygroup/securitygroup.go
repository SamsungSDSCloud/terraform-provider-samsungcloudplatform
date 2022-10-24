package securitygroup

import (
	"context"
	"fmt"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	securitygroup2 "github.com/ScpDevTerra/trf-sdk/library/security-group2"
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
	result, c, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DetailSecurityGroupV2(ctx, securityGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateSecurityGroup(ctx context.Context, productGroupId string, serviceZoneId string, vpcId string, name string, description string) (securitygroup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.CreateSecurityGroupV2(ctx, securitygroup2.SecurityGroupCreateRequest{
		ProductGroupId:           productGroupId,
		SecurityGroupName:        name,
		ServiceZoneId:            serviceZoneId,
		VpcId:                    vpcId,
		SecurityGroupDescription: description,
	})
	return result, err
}

func (client *Client) UpdateSecurityGroupDescription(ctx context.Context, securityGroupId string, description string) (securitygroup2.DetailSecurityGroupResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.ModifySecurityGroupDescriptionV2(ctx, securityGroupId, securitygroup2.SecurityGroupModifyDescriptionRequest{
		SecurityGroupDescription: description,
	})
	return result, err
}

func (client *Client) UpdateSecurityGroupIsLoggable(ctx context.Context, securityGroupId string, isLoggable bool) (securitygroup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.UpdateSecurityGroupLoggingV2(ctx, securityGroupId, securitygroup2.SecurityGroupChangeLoggableRequest{
		IsLoggable: isLoggable,
	})
	return result, err
}

func (client *Client) DeleteSecurityGroup(ctx context.Context, securityGroupId string) error {
	_, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DeleteSecurityGroupV2(ctx, securityGroupId)
	return err
}

func (client *Client) CheckSecurityGroupName(ctx context.Context, name string, vpcId string) (bool, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.CheckSecurityGroupNameDuplicationV2(ctx, name, vpcId)
	return result.Result, err
}

type SecurityGroupServiceRule = securitygroup2.ServiceVo

func (client *Client) CreateSecurityGroupRule(ctx context.Context, securityGroupId string, direction string, addresses []string, description string, services []SecurityGroupServiceRule) (securitygroup2.AsyncResponse, error) {
	if strings.ToUpper(direction) == "IN" {
		result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.CreateSecurityGroupRuleV2(ctx, securityGroupId, securitygroup2.SecurityGroupCreateRuleRequest{
			RuleDirection:     "IN",
			Services:          services,
			SourceIpAddresses: addresses,
			RuleDescription:   description,
		})
		return result, err
	} else if strings.ToUpper(direction) == "OUT" {
		result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.CreateSecurityGroupRuleV2(ctx, securityGroupId, securitygroup2.SecurityGroupCreateRuleRequest{
			RuleDirection:          "OUT",
			Services:               services,
			DestinationIpAddresses: addresses,
			RuleDescription:        description,
		})
		return result, err
	}
	return securitygroup2.AsyncResponse{}, fmt.Errorf("Invalid rule direction : " + direction)
}

func (client *Client) GetSecurityGroupRule(ctx context.Context, ruleId string, securityGroupId string) (securitygroup2.SecurityGroupRuleResponse, int, error) {
	result, c, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DetailSecurityGroupRuleV2(ctx, ruleId, securityGroupId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateSecurityGroupRule(ctx context.Context, ruleId string, securityGroupId string, direction string, addresses []string, description string, services []SecurityGroupServiceRule) (securitygroup2.AsyncResponse, error) {
	if direction == "IN" {
		result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.ModifySecurityGroupRuleV2(ctx, ruleId, securityGroupId, securitygroup2.SecurityGroupRuleModifyRequest{
			SourceIpAddresses: addresses,
			RuleDirection:     "IN",
			Services:          services,
			RuleDescription:   description,
		})
		return result, err
	} else if direction == "OUT" {
		result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.ModifySecurityGroupRuleV2(ctx, ruleId, securityGroupId, securitygroup2.SecurityGroupRuleModifyRequest{
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
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DeleteSecurityGroupRuleV2(ctx, securityGroupId, securitygroup2.SecurityGroupRuleDeleteRequest{
		RuleDeletionType: "PARTIAL",
		RuleIds:          ruleIds,
	})
	return result, err
}

func (client *Client) DeleteSecurityGroupRuleAll(ctx context.Context, securityGroupId string) (securitygroup2.AsyncResponse, error) {
	result, _, err := client.sdkClient.SecurityGroupOpenApiControllerV2Api.DeleteSecurityGroupRuleV2(ctx, securityGroupId, securitygroup2.SecurityGroupRuleDeleteRequest{
		RuleDeletionType: "ALL",
	})
	return result, err
}
