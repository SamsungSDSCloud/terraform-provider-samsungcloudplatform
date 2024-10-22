package servergroup

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	servergroup2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/server-group2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *servergroup2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: servergroup2.NewAPIClient(config),
	}
}

//func (client *Client) CreateServerGroup(ctx context.Context, affinityPolicyType string,) {
//	client.sdkClient.ServerGroupOperationControllerApi.CreateServerGroup(ctx, client.config.ProjectId, servergroup2.ServerGroupCreateRequest{
//		AffinityPolicyType: "",
//		DeploymentEnvType:  "",
//		ServerGroupName:    "",
//		ServerGroupType:    "",
//		ServiceZoneId:      "",
//		ServicedFor:        "",
//		ServicedGroupFor:   "",
//	})
//}

func (client *Client) GetServerGroup(ctx context.Context, serverGroupName string, servicedFor []string) (servergroup2.ListResponseServerGroupsResponse, error) {
	var optServerGroupName optional.String
	if len(serverGroupName) > 0 {
		optServerGroupName = optional.NewString(serverGroupName)
	}
	var optServicedForList optional.Interface
	if len(servicedFor) > 0 {
		optServicedForList = optional.NewInterface(servicedFor)
	}

	result, _, err := client.sdkClient.ServerGroupSearchControllerApi.ListServerGroup(ctx, client.config.ProjectId, &servergroup2.ServerGroupSearchControllerApiListServerGroupOpts{
		ServerGroupName: optServerGroupName,
		ServicedForList: optServicedForList,
		Page:            optional.Int32{},
		Size:            optional.Int32{},
		Sort:            optional.Interface{},
	})

	return result, err
}

func (client *Client) GetServerGroupByServicedForCondition(ctx context.Context, serverGroupName string, servicedFor string) (servergroup2.ListResponseServerGroupsResponse, error) {
	var optServerGroupName optional.String
	if len(serverGroupName) > 0 {
		optServerGroupName = optional.NewString(serverGroupName)
	}
	var optServicedForList optional.Interface
	if len(servicedFor) > 0 {
		optServicedForList = optional.NewInterface(servicedFor)
	}

	result, _, err := client.sdkClient.ServerGroupSearchControllerApi.ListServerGroup(ctx, client.config.ProjectId, &servergroup2.ServerGroupSearchControllerApiListServerGroupOpts{
		ServerGroupName: optServerGroupName,
		ServicedForList: optServicedForList,
		Page:            optional.Int32{},
		Size:            optional.Int32{},
		Sort:            optional.Interface{},
	})

	return result, err
}
