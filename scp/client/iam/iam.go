package iam

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/iam"
	"github.com/antihax/optional"
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
	result, _, err := client.sdkClient.AccessKeyControllerApi.CreateAccessKey(ctx, iam.AccessKeySecretCreateRequest{
		ProjectId:    projectId,
		DurationDays: durationDays,
	})

	return result, err
}

func (client *Client) ActivateAccessKey(ctx context.Context, accessKeyId string) (bool, error) {
	result, _, err := client.sdkClient.AccessKeyControllerApi.ActivateAccessKey(ctx, accessKeyId)
	return result, err
}

func (client *Client) DeactivateAccessKey(ctx context.Context, accessKeyId string) (bool, error) {
	result, _, err := client.sdkClient.AccessKeyControllerApi.DeactivateAccessKey(ctx, accessKeyId)
	return result, err
}

func (client *Client) DeleteAccessKey(ctx context.Context, accessKeyId string) (int32, error) {
	result, _, err := client.sdkClient.AccessKeyControllerApi.DeleteAccessKeys(ctx, accessKeyId)
	return result, err
}

func (client *Client) ListAccessKeys(ctx context.Context, projectId string, accessKeyProjectType string, accessKeyState string, active optional.Bool, projectName string) (iam.PageResponseV2OfAccessKeysResponse, error) {
	result, _, err := client.sdkClient.AccessKeyControllerApi.ListAccessKeys(ctx, &iam.AccessKeyControllerApiListAccessKeysOpts{
		ProjectId:            optional.NewString(projectId),
		AccessKeyProjectType: optional.NewString(accessKeyProjectType),
		ActiveYn:             active,
		ProjectName:          optional.NewString(projectName),
		AccessKeyState:       optional.NewString(accessKeyState),
		Page:                 optional.NewInt32(0),
		Size:                 optional.NewInt32(10000),
	})

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

func (client *Client) CreateGroup(ctx context.Context, groupName string, description string) (iam.GroupResponse, int, error) {
	result, c, err := client.sdkClient.GroupControllerApi.CreateGroup(ctx, client.config.ProjectId, iam.GroupCreateRequest{
		GroupName:   groupName,
		Description: description,
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
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

func (client *Client) DeleteGroups(ctx context.Context, groupIds string) (int, error) {
	result, err := client.sdkClient.GroupControllerApi.DeleteGroups(ctx, client.config.ProjectId, groupIds)
	var statusCode int
	if result != nil {
		statusCode = result.StatusCode
	}
	return statusCode, err
}

func (client *Client) UpdateSecurityInfo(ctx context.Context, ipAclActivated bool, ipAddresses []string, mfaActivated bool) (iam.SecurityInfoResponse, error) {
	result, _, err := client.sdkClient.SecurityControllerApi.UpdateSecurityInfo(
		ctx,
		iam.SecurityInfoUpdateRequest{
			IpAclActivated: &ipAclActivated,
			IpAddresses:    ipAddresses,
			MfaActivated:   &mfaActivated,
		})

	return result, err
}

/*
func (client *Client) GetSecurityInfo(ctx context.Context) (iam.SecurityInfoResponse, error) {
	result, _, err := client.sdkClient.SecurityControllerApi.GetSecurityInfo(ctx)
	return result, err
}
*/

func (client *Client) DetailMember(ctx context.Context, memberId string) (iam.MemberResponse, error) {
	result, _, err := client.sdkClient.MemberControllerApi.DetailMember(ctx, client.config.ProjectId, memberId)
	return result, err
}

func (client *Client) CreatePolicy(ctx context.Context, policyName string, policyJson string, principals []iam.PolicyPrincipalRequest, tags []interface{}, description string) (iam.PolicyResponse, error) {
	result, _, err := client.sdkClient.PolicyControllerApi.CreatePolicy(ctx, client.config.ProjectId, iam.PolicyCreateRequest{
		PolicyJson:  policyJson,
		PolicyName:  policyName,
		Principals:  principals,
		Tags:        toTagRequestList(tags),
		Description: description,
	})

	return result, err
}

func (client *Client) UpdatePolicy(ctx context.Context, policyId string, policyJson string, policyName string, principals []iam.PolicyPrincipalRequest, description string) (iam.PolicyResponse, error) {
	result, _, err := client.sdkClient.PolicyControllerApi.UpdatePolicy(ctx, client.config.ProjectId, policyId, iam.PolicyUpdateRequest{
		PolicyJson:  policyJson,
		PolicyName:  policyName,
		Principals:  principals,
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

func (client *Client) ValidatePolicyJson(ctx context.Context, policyJson string) (iam.PolicyValidityResponse, int, error) {

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

func (client *Client) CreateRole(ctx context.Context, roleName string, projectIds []string, userSrns []string, tags []interface{}, description string) (iam.RoleResponse, int, error) {
	var principal = iam.TrustPrincipalsResponse{
		ProjectIds: projectIds,
		UserSrns:   userSrns,
	}

	result, c, err := client.sdkClient.RoleControllerApi.CreateRole(ctx, client.config.ProjectId, iam.RoleCreateRequest{
		RoleName:        roleName,
		TrustPrincipals: &principal,
		Description:     description,
		Tags:            toTagRequestList(tags),
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DetailRole(ctx context.Context, roleId string) (iam.RoleResponse, error) {
	result, _, err := client.sdkClient.RoleControllerApi.DetailRole(ctx, client.config.ProjectId, roleId)
	return result, err
}

func (client *Client) UpdateRole(ctx context.Context, roleId string, roleName string, projectIds []string, userSrns []string, description string) (iam.RoleResponse, error) {
	result, _, err := client.sdkClient.RoleControllerApi.UpdateRole(ctx, client.config.ProjectId, roleId, iam.RoleUpdateRequest{
		RoleName: roleName,
		TrustPrincipals: &iam.TrustPrincipalsResponse{
			ProjectIds: projectIds,
			UserSrns:   userSrns,
		},
		Description: description,
	})
	return result, err
}

func (client *Client) DeleteRole(ctx context.Context, roleId string) (*http.Response, error) {
	result, err := client.sdkClient.RoleControllerApi.DeleteRoles(ctx, client.config.ProjectId, roleId)
	return result, err
}

func (client *Client) AddRolePolicies(ctx context.Context, roleId string, policyIds []string) (int, error) {
	c, err := client.sdkClient.RoleControllerApi.AddRolePolicys(ctx, client.config.ProjectId, roleId, iam.RolePolicysAddRequest{PolicyIds: policyIds})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return statusCode, err
}

func (client *Client) RemoveRolePolicies(ctx context.Context, roleId string, policyIds []string) (int, error) {
	c, err := client.sdkClient.RoleControllerApi.RemoveRolePolicys(ctx, client.config.ProjectId, roleId, iam.RolePolicysRemoveRequest{PolicyIds: policyIds})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return statusCode, err
}

func (client *Client) ListMembers(ctx context.Context, companyName string, email string, userName string) (iam.PageResponseV2OfMembersResponse, error) {
	result, _, err := client.sdkClient.MemberControllerApi.ListMembers(ctx, client.config.ProjectId, &iam.MemberControllerApiListMembersOpts{
		CompanyName: optional.NewString(companyName),
		Email:       optional.NewString(email),
		UserName:    optional.NewString(userName),
		Page:        optional.NewInt32(0),
		Size:        optional.NewInt32(10000),
	})
	return result, err
}

func (client *Client) ListGroupMembers(ctx context.Context, groupId string, request ListMemberRequest) (iam.PageResponseV2OfGroupMembersResponse, error) {
	result, _, err := client.sdkClient.GroupControllerApi.ListGroupMembers(ctx, client.config.ProjectId, groupId,
		&iam.GroupControllerApiListGroupMembersOpts{
			CompanyName: optional.NewString(request.CompanyName),
			Email:       optional.NewString(request.Email),
			UserName:    optional.NewString(request.UserName),
			Page:        optional.NewInt32(0),
			Size:        optional.NewInt32(10000),
		})
	return result, err
}

func (client *Client) AddGroupMembers(ctx context.Context, groupId string, userIds []string) (*http.Response, error) {
	result, err := client.sdkClient.GroupControllerApi.AddGroupMembers(ctx, client.config.ProjectId, groupId, iam.GroupMembersAddRequest{
		UserIds: userIds,
	})
	return result, err
}

func (client *Client) GetUserGroupIds(ctx context.Context, groupId string, userIds []string) ([]string, error) {
	if len(userIds) == 0 {
		return []string{}, nil
	}

	result, _, err := client.sdkClient.GroupControllerApi.ListGroupMembers(ctx, client.config.ProjectId, groupId, &iam.GroupControllerApiListGroupMembersOpts{
		Page: optional.NewInt32(0),
		Size: optional.NewInt32(10000),
	})

	if err != nil {
		return []string{}, err
	}

	userIdSet := make(map[string]struct{})
	for _, item := range userIds {
		userIdSet[item] = struct{}{}
	}

	userGroupIds := make([]string, 0)
	for _, c := range result.Contents {
		if _, ok := userIdSet[c.UserId]; ok {
			userGroupIds = append(userGroupIds, c.UserGroupId)
		}
	}

	if len(userIds) != len(userGroupIds) {
		return []string{}, fmt.Errorf("the number of user group ids(%d) and user ids(%d) are different", len(userGroupIds), len(userIds))
	}
	return userGroupIds, nil
}

func (client *Client) GetPrincipalPolicyIds(ctx context.Context, groupId string, policyIds []string) ([]string, error) {
	if len(policyIds) == 0 {
		return []string{}, nil
	}

	result, _, err := client.sdkClient.GroupControllerApi.ListGroupPolicys(ctx, client.config.ProjectId, groupId, &iam.GroupControllerApiListGroupPolicysOpts{
		Page: optional.NewInt32(0),
		Size: optional.NewInt32(10000),
	})

	if err != nil {
		return []string{}, err
	}

	policyIdSet := make(map[string]struct{})
	for _, item := range policyIds {
		policyIdSet[item] = struct{}{}
	}

	principalPolicyIds := make([]string, 0)
	for _, c := range result.Contents {
		if _, ok := policyIdSet[c.PolicyId]; ok {
			principalPolicyIds = append(principalPolicyIds, c.PrincipalPolicyId)
		}
	}

	if len(policyIds) != len(principalPolicyIds) {
		return []string{}, fmt.Errorf("the number of principal policy ids(%d) and policy ids(%d) are different", len(principalPolicyIds), len(policyIds))
	}
	return principalPolicyIds, nil
}

func (client *Client) RemoveGroupMembers(ctx context.Context, groupId string, userIds []string) (*http.Response, error) {
	userGroupIds, err := client.GetUserGroupIds(ctx, groupId, userIds)
	if err != nil {
		return nil, err
	}

	result, err := client.sdkClient.GroupControllerApi.RemoveGroupMembers(ctx, client.config.ProjectId, groupId, iam.GroupMembersRemoveRequest{
		UserGroupIds: userGroupIds,
	})

	return result, err
}

func (client *Client) ListPolicies(ctx context.Context, request ListMemberRequest) (iam.PageResponseV2OfPolicysResponse, error) {
	result, _, err := client.sdkClient.PolicyControllerApi.ListPolicys(ctx, client.config.ProjectId,
		&iam.PolicyControllerApiListPolicysOpts{
			ModifiedByEmail: optional.NewString(request.CompanyName),
			PolicyName:      optional.NewString(request.Email),
			Page:            optional.NewInt32(0),
			Size:            optional.NewInt32(10000),
		})
	return result, err
}

func (client *Client) ListGroupPolicies(ctx context.Context, groupId string, request ListPolicyRequest) (iam.PageResponseV2OfGroupPolicysResponse, error) {
	result, _, err := client.sdkClient.GroupControllerApi.ListGroupPolicys(ctx, client.config.ProjectId, groupId, &iam.GroupControllerApiListGroupPolicysOpts{
		PolicyName: optional.NewString(request.PolicyName),
		PolicyType: optional.NewString(request.PolicyType),
		Page:       optional.NewInt32(0),
		Size:       optional.NewInt32(10000),
	})
	return result, err
}

func (client *Client) ListPolicyGroups(ctx context.Context, policyId string, groupName string) (iam.PageResponseV2OfPolicyGroupsResponse, int, error) {
	result, c, err := client.sdkClient.PolicyControllerApi.ListPolicyGroups(ctx, client.config.ProjectId, policyId, &iam.PolicyControllerApiListPolicyGroupsOpts{
		GroupName: optional.NewString(groupName),
		Page:      optional.NewInt32(0),
		Size:      optional.NewInt32(10000),
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListPolicyRoles(ctx context.Context, policyId string, roleName string) (iam.PageResponseV2OfPolicyRolesResponse, int, error) {
	result, c, err := client.sdkClient.PolicyControllerApi.ListPolicyRoles(ctx, client.config.ProjectId, policyId, &iam.PolicyControllerApiListPolicyRolesOpts{
		RoleName: optional.NewString(roleName),
		Page:     optional.NewInt32(0),
		Size:     optional.NewInt32(10000),
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) AddGroupPolicies(ctx context.Context, groupId string, policyIds []string) (*http.Response, error) {
	result, err := client.sdkClient.GroupControllerApi.AddGroupPolicys(ctx, client.config.ProjectId, groupId, iam.GroupPolicysAddRequest{
		PolicyIds: policyIds,
	})
	return result, err
}

func (client *Client) RemoveGroupPolicies(ctx context.Context, groupId string, policyIds []string) (*http.Response, error) {
	principalPolicyIds, err := client.GetPrincipalPolicyIds(ctx, groupId, policyIds)
	if err != nil {
		return nil, err
	}

	result, err := client.sdkClient.GroupControllerApi.RemoveGroupPolicys(ctx, client.config.ProjectId, groupId, iam.GroupPolicysRemoveRequest{
		PrincipalPolicyIds: principalPolicyIds,
	})

	return result, err
}

func (client *Client) ListMemberGroups(ctx context.Context, memberId string, groupName string) (iam.PageResponseV2OfMemberGroupsResponse, error) {
	result, _, err := client.sdkClient.MemberControllerApi.ListMemberGroups(ctx, client.config.ProjectId, memberId, &iam.MemberControllerApiListMemberGroupsOpts{
		GroupName: optional.NewString(groupName),
		Page:      optional.NewInt32(0),
		Size:      optional.NewInt32(10000),
	})
	return result, err
}

func (client *Client) ListMemberSystemGroups(ctx context.Context, memberId string) (iam.ListResponseV2OfGroupLinkResponse, error) {
	result, _, err := client.sdkClient.MemberControllerApi.ListMemberDefaultSystemGroups(ctx, client.config.ProjectId, memberId)
	return result, err
}

func (client *Client) RemoveMemberGroups(ctx context.Context, memberId string, userGroupIds []string) (int, error) {
	result, err := client.sdkClient.MemberControllerApi.RemoveMemberGroups(ctx, client.config.ProjectId, memberId, iam.MemberGroupsRemoveRequest{
		UserGroupIds: userGroupIds,
	})

	var statusCode int
	if result != nil {
		statusCode = result.StatusCode
	}
	return statusCode, err
}

func (client *Client) AddMemberGroups(ctx context.Context, memberId string, groupIds []string) (int, error) {
	result, err := client.sdkClient.MemberControllerApi.AddMemberGroups(ctx, client.config.ProjectId, memberId, iam.MemberGroupsAddRequest{
		GroupIds: groupIds,
	})

	var statusCode int
	if result != nil {
		statusCode = result.StatusCode
	}
	return statusCode, err
}

func (client *Client) ListGroups(ctx context.Context, groupName string, email string) (iam.PageResponseV2OfGroupsResponse, error) {
	var optName optional.String
	var optEmail optional.String

	if len(groupName) > 0 {
		optName = optional.NewString(groupName)
	}
	if len(email) > 0 {
		optEmail = optional.NewString(email)
	}

	result, _, err := client.sdkClient.GroupControllerApi.ListGroups(ctx, client.config.ProjectId, &iam.GroupControllerApiListGroupsOpts{
		GroupName:       optName,
		ModifiedByEmail: optEmail,
		Page:            optional.NewInt32(0),
		Size:            optional.NewInt32(10000),
	})
	return result, err
}

func (client *Client) ListRoles(ctx context.Context, email string, roleName string) (iam.PageResponseV2OfRolesResponse, error) {
	result, _, err := client.sdkClient.RoleControllerApi.ListRoles(ctx, client.config.ProjectId, &iam.RoleControllerApiListRolesOpts{
		ModifiedByEmail: optional.NewString(email),
		RoleName:        optional.NewString(roleName),
		Page:            optional.NewInt32(0),
		Size:            optional.NewInt32(10000),
	})
	return result, err
}

func (client *Client) ListRolePolicies(ctx context.Context, roleId string, policyName string, policyType string) (iam.PageResponseV2OfRolePolicysResponse, error) {
	result, _, err := client.sdkClient.RoleControllerApi.ListRolePolicys(ctx, client.config.ProjectId, roleId, &iam.RoleControllerApiListRolePolicysOpts{
		PolicyName: optional.NewString(policyName),
		PolicyType: optional.NewString(policyType),
		Page:       optional.NewInt32(0),
		Size:       optional.NewInt32(10000),
	})

	return result, err
}

func (client *Client) AddMember(ctx context.Context, groupIds []string, userEmail string, tags []interface{}) (int32, int, error) {
	userEmails := []string{userEmail}
	result, c, err := client.sdkClient.MemberControllerApi.AddMembers(ctx, client.config.ProjectId, iam.MembersAddRequest{
		GroupIds:   groupIds,
		Tags:       toTagRequestList(tags),
		UserEmails: userEmails,
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) RemoveMember(ctx context.Context, id string) (int, error) {
	userIds := []string{id}
	result, err := client.sdkClient.MemberControllerApi.RemoveMembers(ctx, client.config.ProjectId, iam.MembersRemoveRequest{
		UserIds: userIds,
	})

	var statusCode int
	if result != nil {
		statusCode = result.StatusCode
	}
	return statusCode, err
}

func toTagRequestList(list []interface{}) []iam.TagRequest {
	if len(list) == 0 {
		return nil
	}
	var result []iam.TagRequest

	for _, val := range list {
		kv := val.(common.HclKeyValueObject)
		result = append(result, iam.TagRequest{
			TagKey:   kv["tag_key"].(string),
			TagValue: kv["tag_value"].(string),
		})
	}
	return result
}
