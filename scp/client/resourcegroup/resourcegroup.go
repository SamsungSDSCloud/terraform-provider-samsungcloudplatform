package resourcegroup

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/resource-group"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *resourcegroup.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: resourcegroup.NewAPIClient(config),
	}
}

func (client *Client) CreateResourceGroup(ctx context.Context, request ResourceGroupRequest) (resourcegroup.ResourceGroupResponse, int, error) {
	result, c, err := client.sdkClient.ResourceGroupControllerApi.CreateResourceGroup(ctx, client.config.ProjectId, resourcegroup.ResourceGroupCreateRequest{
		ResourceGroupName:        request.ResourceGroupName,
		TargetResourceTags:       toTagRequestList(request.TargetResourceTags),
		TargetResourceTypes:      common.ToStringList(request.TargetResourceTypes),
		ResourceGroupDescription: request.ResourceGroupDescription,
	})

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}

	return result, statusCode, err
}

func (client *Client) GetResourceGroupList(ctx context.Context, request ListResourceGroupRequest) (resourcegroup.PageResponseV2OfResourceGroupsResponse, error) {
	result, _, err := client.sdkClient.ResourceGroupControllerApi.ListResourceGroups(ctx, client.config.ProjectId, &resourcegroup.ResourceGroupControllerApiListResourceGroupsOpts{
		CreatedById:       optional.NewString(request.CreatedById),
		ModifiedByEmail:   optional.NewString(request.ModifiedByEmail),
		ModifiedById:      optional.NewString(request.ModifiedById),
		ResourceGroupName: optional.NewString(request.ResourceGroupName),
		Page:              optional.NewInt32(0),
		Size:              optional.NewInt32(10000),
	})
	return result, err
}

func (client *Client) GetResourceGroupResourcesList(ctx context.Context, resourceGroupId string, request ListResourceGroupResourcesRequest) (resourcegroup.PageResponseV2OfResourcesResponse, error) {
	result, _, err := client.sdkClient.ResourceGroupControllerApi.ListResourceGroupResources(ctx, client.config.ProjectId, resourceGroupId, &resourcegroup.ResourceGroupControllerApiListResourceGroupResourcesOpts{
		CreatedById:  optional.NewString(request.CreatedById),
		ModifiedById: optional.NewString(request.ModifiedById),
		ResourceId:   optional.NewString(request.ResourceId),
		ResourceName: optional.NewString(request.ResourceName),
		Page:         optional.NewInt32(0),
		Size:         optional.NewInt32(10000),
	})
	return result, err
}

func (client *Client) GetResourceGroup(ctx context.Context, resourceGroupId string) (resourcegroup.ResourceGroupResponse, error) {
	result, _, err := client.sdkClient.ResourceGroupControllerApi.DetailResourceGroup(ctx, client.config.ProjectId, resourceGroupId)
	return result, err
}

func (client *Client) UpdateResourceGroup(ctx context.Context, resourceGroupId string, request ResourceGroupRequest) (resourcegroup.ResourceGroupResponse, error) {
	result, _, err := client.sdkClient.ResourceGroupControllerApi.UpdateResourceGroup(ctx, client.config.ProjectId, resourceGroupId, resourcegroup.ResourceGroupUpdateRequest{
		ResourceGroupName:        request.ResourceGroupName,
		TargetResourceTags:       toTagRequestList(request.TargetResourceTags),
		TargetResourceTypes:      common.ToStringList(request.TargetResourceTypes),
		ResourceGroupDescription: request.ResourceGroupDescription,
	})
	return result, err
}

func (client *Client) DeleteResourceGroup(ctx context.Context, resourceGroupId string) error {
	_, err := client.sdkClient.ResourceGroupControllerApi.DeleteResourceGroup(ctx, client.config.ProjectId, resourceGroupId)
	return err
}

func (client *Client) DeleteResourceGroups(ctx context.Context, resourceGroupIds []string) error {
	_, err := client.sdkClient.ResourceGroupControllerApi.DeleteResourceGroups(ctx, client.config.ProjectId, resourceGroupIds)
	return err
}

func (client *Client) GetResources(ctx context.Context, request ListResourceRequest) (resourcegroup.PageResponseV2OfResourcesResponse, error) {
	result, _, err := client.sdkClient.ResourceControllerApi.ListResources(ctx, client.config.ProjectId, &resourcegroup.ResourceControllerApiListResourcesOpts{
		CreatedById:         optional.NewString(request.CreatedById),
		DisplayServiceNames: optional.NewInterface(common.ToStringList(request.DisplayServiceNames)),
		FromCreatedAt:       optional.NewString(request.FromCreatedAt),
		IncludeDeleted:      optional.NewString(request.IncludeDeleted),
		Location:            optional.NewString(request.Location),
		ModifiedById:        optional.NewString(request.ModifiedById),
		MyCreate:            optional.NewString(request.MyCreate),
		Partitions:          optional.NewInterface(common.ToStringList(request.Partitions)),
		Regions:             optional.NewInterface(common.ToStringList(request.Regions)),
		ResourceId:          optional.NewString(request.ResourceId),
		ResourceName:        optional.NewString(request.ResourceName),
		ResourceTypes:       optional.NewInterface(common.ToStringList(request.ResourceTypes)),
		ServiceTypes:        optional.NewInterface(common.ToStringList(request.ServiceTypes)),
		ServiceZones:        optional.NewInterface(common.ToStringList(request.ServiceZones)),
		Tags:                optional.NewInterface(common.ToStringList(request.Tags)),
		ToCreatedAt:         optional.NewString(request.ToCreatedAt),
		Page:                optional.NewInt32(0),
		Size:                optional.NewInt32(10000),
	})
	return result, err
}

func (client *Client) GetResource(ctx context.Context, resourceId string, includeDeleted optional.String) (resourcegroup.ResourceResponse, error) {
	result, _, err := client.sdkClient.ResourceControllerApi.DetailResource(ctx, client.config.ProjectId, resourceId, &resourcegroup.ResourceControllerApiDetailResourceOpts{
		IncludeDeleted: includeDeleted,
	})
	return result, err
}

func (client *Client) GetResourceSrn(ctx context.Context, resourceId string, includeDeleted string) (string, error) {
	result, _, err := client.sdkClient.ResourceControllerApi.ResourceSrn(ctx, client.config.ProjectId, resourceId, &resourcegroup.ResourceControllerApiResourceSrnOpts{
		IncludeDeleted: optional.NewString(includeDeleted),
	})
	return result, err
}

func (client *Client) GetResourceTypes(ctx context.Context, resourceType string, serviceType string) ([]resourcegroup.ResourceTypeResponse, error) {
	result, _, err := client.sdkClient.TypeControllerApi.ListResourceTypes(ctx, &resourcegroup.TypeControllerApiListResourceTypesOpts{
		ResourceType: optional.NewString(resourceType),
		ServiceType:  optional.NewString(serviceType),
	})
	return result, err
}

func (client *Client) GetServiceTypes(ctx context.Context, serviceType string) ([]string, error) {
	result, _, err := client.sdkClient.TypeControllerApi.ListServiceTypes(ctx, &resourcegroup.TypeControllerApiListServiceTypesOpts{
		ServiceType: optional.NewString(serviceType),
	})
	return result, err
}

func (client *Client) GetResourceGroupListInMyProjects(ctx context.Context, projectIds []interface{}, request ListResourceGroupRequest) (resourcegroup.PageResponseV2OfMyProjectsResourceGroupsResponse, error) {
	result, _, err := client.sdkClient.MyProjectResourceGroupControllerApi.ListMyProjectsResourceGroups(ctx, &resourcegroup.MyProjectResourceGroupControllerApiListMyProjectsResourceGroupsOpts{
		CreatedById:       optional.NewString(request.CreatedById),
		ModifiedByEmail:   optional.NewString(request.ModifiedByEmail),
		ModifiedById:      optional.NewString(request.ModifiedById),
		ResourceGroupName: optional.NewString(request.ResourceGroupName),
		ProjectIds:        optional.NewInterface(common.ToStringList(projectIds)),
		Page:              optional.NewInt32(0),
		Size:              optional.NewInt32(10000),
	})
	return result, err
}

func (client *Client) GetResourceGroupInMyProjects(ctx context.Context, resourceGroupId string) (resourcegroup.MyProjectsResourceGroupResponse, error) {
	result, _, err := client.sdkClient.MyProjectResourceGroupControllerApi.DetailMyProjectsResourceGroup(ctx, resourceGroupId)
	return result, err
}

func (client *Client) GetResourceGroupResourcesInMyProjects(ctx context.Context, resourceGroupId string, request ListResourceGroupResourcesRequest) (resourcegroup.PageResponseV2OfResourcesResponse, error) {
	result, _, err := client.sdkClient.MyProjectResourceGroupControllerApi.ListMyProjectsResourceGroupResources(ctx, resourceGroupId, &resourcegroup.MyProjectResourceGroupControllerApiListMyProjectsResourceGroupResourcesOpts{
		CreatedById:  optional.NewString(request.CreatedById),
		ModifiedById: optional.NewString(request.ModifiedById),
		ResourceId:   optional.NewString(request.ResourceId),
		ResourceName: optional.NewString(request.ResourceName),
		Page:         optional.NewInt32(0),
		Size:         optional.NewInt32(10000),
	})
	return result, err
}

func toTagRequestList(list []interface{}) []resourcegroup.Tag {
	if len(list) == 0 {
		return nil
	}
	var result []resourcegroup.Tag

	for _, val := range list {
		kv := val.(common.HclKeyValueObject)
		result = append(result, resourcegroup.Tag{
			TagKey:   kv["tag_key"].(string),
			TagValue: kv["tag_value"].(string),
		})
	}
	return result
}

func (client *Client) GetMyProjectResources(ctx context.Context, projectIds []string, request ListResourceRequest) (resourcegroup.PageResponseV2OfMyProjectsResourcesResponse, error) {
	result, _, err := client.sdkClient.MyProjectResourceControllerApi.ListMyProjectsResources(ctx, &resourcegroup.MyProjectResourceControllerApiListMyProjectsResourcesOpts{
		CreatedById:         optional.NewString(request.CreatedById),
		DisplayServiceNames: optional.NewInterface(common.ToStringList(request.DisplayServiceNames)),
		FromCreatedAt:       optional.NewString(request.FromCreatedAt),
		IncludeDeleted:      optional.NewString(request.IncludeDeleted),
		Location:            optional.NewString(request.Location),
		ModifiedById:        optional.NewString(request.ModifiedById),
		MyCreate:            optional.NewString(request.MyCreate),
		Partitions:          optional.NewInterface(common.ToStringList(request.Partitions)),
		ProjectIds:          optional.NewInterface(projectIds),
		Regions:             optional.NewInterface(common.ToStringList(request.Regions)),
		ResourceId:          optional.NewString(request.ResourceId),
		ResourceName:        optional.NewString(request.ResourceName),
		ResourceTypes:       optional.NewInterface(common.ToStringList(request.ResourceTypes)),
		ServiceTypes:        optional.NewInterface(common.ToStringList(request.ServiceTypes)),
		ServiceZones:        optional.NewInterface(common.ToStringList(request.ServiceZones)),
		Tags:                optional.NewInterface(common.ToStringList(request.Tags)),
		ToCreatedAt:         optional.NewString(request.ToCreatedAt),
	})
	return result, err
}

func (client *Client) GetMyProjectResource(ctx context.Context, resourceId string, projectId string, includeDeleted optional.String) (resourcegroup.MyProjectsResourceResponse, error) {
	result, _, err := client.sdkClient.MyProjectResourceControllerApi.DetailMyProjectsResource(ctx, resourceId, projectId, &resourcegroup.MyProjectResourceControllerApiDetailMyProjectsResourceOpts{
		IncludeDeleted: includeDeleted,
	})
	return result, err
}
