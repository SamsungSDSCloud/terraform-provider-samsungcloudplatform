package iam

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/client"
	"github.com/SamsungSDSCloud/terraform-sdk-SamsungCloudPlatform/library/iam"
	"net/http"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *iam.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: iam.NewAPIClient(config),
	}
}

func (client *Client) CreateAccessKey(ctx context.Context, projectId string, durationDays int32) (iam.AccessKeyResponse, error) {
	result, _, err := client.sdkClient.AccessKeyControllerApi.CreateAccessKey(
		ctx,
		iam.AccessKeySecretCreateRequest{
			ProjectId:    projectId,
			DurationDays: durationDays,
		})

	return result, err
}

func (client *Client) ActivateAccessKey(ctx context.Context, accessKeyId string) (bool, error) {
	result, _, err := client.sdkClient.AccessKeyControllerApi.ActivateAccessKey(
		ctx,
		accessKeyId)

	return result, err
}

func (client *Client) DactivateAccessKey(ctx context.Context, accessKeyId string) (bool, error) {
	result, _, err := client.sdkClient.AccessKeyControllerApi.DeactivateAccessKey(
		ctx,
		accessKeyId)

	return result, err
}

func (client *Client) DeleteAccessKey(ctx context.Context, accessKeyId string) (int32, error) {

	result, _, err := client.sdkClient.AccessKeyControllerApi.DeleteAccessKeys(ctx, accessKeyId)
	return result, err
}

func (client *Client) CreateTemporaryAccessKey(ctx context.Context, otp string, durationMinutes int32) (iam.AccessKeyResponse, error) {
	result, _, err := client.sdkClient.TemporaryAccessKeyControllerApi.CreateTemporaryAccessKey(
		ctx,
		iam.TemporaryAccessKeySecretCreateRequest{
			DurationMinutes: durationMinutes,
			Otp:             otp,
		})

	return result, err
}

func (client *Client) ActivateTmpAccessKey(ctx context.Context, accessKeyId string) (bool, error) {
	result, _, err := client.sdkClient.TemporaryAccessKeyControllerApi.ActivateTemporaryAccessKey(ctx, accessKeyId)
	return result, err
}

func (client *Client) DactivateTmpAccessKey(ctx context.Context, accessKeyId string) (bool, error) {
	result, _, err := client.sdkClient.TemporaryAccessKeyControllerApi.DeactivateTemporaryAccessKey(ctx, accessKeyId)
	return result, err
}

func (client *Client) DeleteTmpAccessKey(ctx context.Context, accessKeyId string) (int32, error) {

	result, _, err := client.sdkClient.TemporaryAccessKeyControllerApi.DeleteTemporaryAccessKeys(ctx, accessKeyId)
	return result, err
}

func (client *Client) CreatGroup(ctx context.Context, groupName string, description string) (iam.GroupResponse, error) {
	result, _, err := client.sdkClient.GroupControllerApi.CreateGroup(
		ctx,
		client.config.ProjectId,
		iam.GroupCreateRequest{
			GroupName:   groupName,
			Description: description,
		})

	return result, err
}

func (client *Client) DetailGroup(ctx context.Context, groupId string) (iam.GroupResponse, error) {
	result, _, err := client.sdkClient.GroupControllerApi.DetailGroup(ctx, client.config.ProjectId, groupId)
	return result, err
}

func (client *Client) UpdateGroup(ctx context.Context, groupId string, groupName string, description string) (iam.GroupResponse, error) {
	result, _, err := client.sdkClient.GroupControllerApi.UpdateGroup(
		ctx,
		client.config.ProjectId,
		groupId,
		iam.GroupUpdateRequest{
			GroupName:   groupName,
			Description: description,
		})

	return result, err
}

func (client *Client) UpdateSecurityInfo(ctx context.Context, ipAclActivated bool, ipAddresses []string, mfaActivated bool) (iam.SecurityInfoResponse, error) {
	result, _, err := client.sdkClient.SecurityControllerApi.UpdateSecurityInfo(
		ctx,
		iam.SecurityInfoUpdateRequest{
			IpAclActivated: ipAclActivated,
			IpAddresses:    ipAddresses,
			MfaActivated:   mfaActivated,
		})

	return result, err
}

func (client *Client) GetSecurityInfo(ctx context.Context) (iam.SecurityInfoResponse, error) {
	result, _, err := client.sdkClient.SecurityControllerApi.GetSecurityInfo(ctx)
	return result, err
}

func (client *Client) DetailMember(ctx context.Context, memberId string) (iam.MemberResponse, error) {
	result, _, err := client.sdkClient.MemberControllerApi.DetailMember(ctx, client.config.ProjectId, memberId)
	return result, err
}

func (client *Client) CreatePolicy(ctx context.Context, policyName string, policyJson string, policyPrincipal []PolicyPrincipalRequest, description string) (iam.PolicyResponse, error) {

	var principal []iam.PolicyPrincipalRequest
	/*for _, b := range policyPrincipal {
		principal = append(principal, iam.PolicyPrincipalRequest{
			PrincipalId:   b.PrincipalId,
			PrincipalType: b.PrincipalType,
		})
	}*/

	result, _, err := client.sdkClient.PolicyControllerApi.CreatePolicy(
		ctx,
		client.config.ProjectId,
		iam.PolicyCreateRequest{
			PolicyJson:  policyJson,
			PolicyName:  policyName,
			Principals:  principal,
			Description: description,
		})

	return result, err
}

func (client *Client) UpdatePolciy(ctx context.Context, policyId string, policyJson string, policyName string, policyPrincipal []PolicyPrincipalRequest, description string) (iam.PolicyResponse, error) {

	var principal []iam.PolicyPrincipalRequest
	for _, b := range policyPrincipal {
		principal = append(principal, iam.PolicyPrincipalRequest{
			PrincipalId:   b.PrincipalId,
			PrincipalType: b.PrincipalType,
		})
	}

	result, _, err := client.sdkClient.PolicyControllerApi.UpdatePolicy(
		ctx,
		client.config.ProjectId,
		policyId,
		iam.PolicyUpdateRequest{
			PolicyJson:  policyJson,
			PolicyName:  policyName,
			Principals:  principal,
			Description: description,
		})

	return result, err
}

func (client *Client) DetailPolicy(ctx context.Context, policyId string) (iam.PolicyResponse, error) {
	result, _, err := client.sdkClient.PolicyControllerApi.DetailPolicy(ctx, client.config.ProjectId, policyId)
	return result, err
}

func (client *Client) DeletePolicy(ctx context.Context, policyId string) (*http.Response, error) {

	result, err := client.sdkClient.PolicyControllerApi.DeletePolicy(ctx, client.config.ProjectId, policyId)
	return result, err
}

func (client *Client) ValidPolicyJson(ctx context.Context, policyJson string) (iam.PolicyValidityResponse, int, error) {

	result, c, err := client.sdkClient.PolicyControllerApi.ValidatePolicyJson(ctx, client.config.ProjectId,
		iam.PolicyValidationRequest{
			PolicyJson: policyJson,
		})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateRole(ctx context.Context, roleName string, priciple TrustPrincipalsResponse, description string) (iam.RoleResponse, error) {

	var principal = iam.TrustPrincipalsResponse{
		ProjectIds: priciple.ProjectIds,
		UserSrns:   priciple.UserSrns,
	}

	result, _, err := client.sdkClient.RoleControllerApi.CreateRole(
		ctx,
		client.config.ProjectId,
		iam.RoleCreateRequest{
			RoleName:        roleName,
			TrustPrincipals: &principal,
			Description:     description,
		})

	return result, err
}

func (client *Client) DetailRole(ctx context.Context, roleId string) (iam.RoleResponse, error) {

	result, _, err := client.sdkClient.RoleControllerApi.DetailRole(ctx, client.config.ProjectId, roleId)
	return result, err
}

func (client *Client) UpdateRole(ctx context.Context, roleId string, roleName string, priciple TrustPrincipalsResponse, description string) (iam.RoleResponse, error) {

	var principal = iam.TrustPrincipalsResponse{
		ProjectIds: priciple.ProjectIds,
		UserSrns:   priciple.UserSrns,
	}

	result, _, err := client.sdkClient.RoleControllerApi.UpdateRole(
		ctx,
		client.config.ProjectId,
		roleId,
		iam.RoleUpdateRequest{
			RoleName:        roleName,
			TrustPrincipals: &principal,
			Description:     description,
		})
	return result, err
}

func (client *Client) DeleteRole(ctx context.Context, roleId string) (*http.Response, error) {

	result, err := client.sdkClient.RoleControllerApi.DeleteRoles(ctx, client.config.ProjectId, roleId)
	return result, err
}
