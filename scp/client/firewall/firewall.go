package firewall

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/firewall2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *firewall2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: firewall2.NewAPIClient(config),
	}
}

func (client *Client) GetFirewallList(ctx context.Context, vpcId string, targetId string, firewallName string) (firewall2.ListResponseOfFirewallListItemResponse, int, error) {
	var optVpcId optional.String
	if len(vpcId) > 0 {
		optVpcId = optional.NewString(vpcId)
	}
	var optTargetId optional.String
	if len(targetId) > 0 {
		optTargetId = optional.NewString(targetId)
	}
	var optFirewallName optional.String
	if len(firewallName) > 0 {
		optFirewallName = optional.NewString(firewallName)
	}
	result, c, err := client.sdkClient.FirewallV2Api.ListFirewallsV2(ctx, client.config.ProjectId, &firewall2.FirewallV2ApiListFirewallsV2Opts{
		FirewallName:   optFirewallName,
		FirewallStates: optional.Interface{},
		IsLoggable:     optional.Bool{},
		ObjectId:       optTargetId,
		ObjectTypes:    optional.Interface{},
		VpcId:          optVpcId,
		CreatedBy:      optional.String{},
		Page:           optional.NewInt32(0),
		Size:           optional.NewInt32(10000),
		Sort:           optional.Interface{}, //NewInterface([]string{"vpcName:asc"}),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetFirewall(ctx context.Context, firewallId string) (firewall2.FirewallDetailResponse, int, error) {
	result, c, err := client.sdkClient.FirewallV2Api.DetailFirewallV2(ctx, client.config.ProjectId, firewallId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateFirewallEnabled(ctx context.Context, firewallId string, enabled bool) (firewall2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.FirewallV2Api.UpdateFirewallEnabledV2(ctx, client.config.ProjectId, firewallId, firewall2.FirewallChangeEnabledRequest{
		IsEnabled: enabled,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateFirewallRule(ctx context.Context, firewallId string, request firewall2.FirewallCreateRuleRequest) (firewall2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.FirewallRuleV2Api.CreateFirewallRuleV2(ctx, client.config.ProjectId, firewallId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetFirewallRule(ctx context.Context, firewallId string, ruleId string) (firewall2.FirewallRuleDetailResponse, int, error) {
	result, c, err := client.sdkClient.FirewallRuleV2Api.DetailFirewallRuleV2(ctx, client.config.ProjectId, firewallId, ruleId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateFirewallRule(ctx context.Context, firewallId string, ruleId string, request firewall2.FirewallRuleUpdateRequest) (firewall2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.FirewallRuleV2Api.UpdateFirewallRuleV2(ctx, client.config.ProjectId, firewallId, ruleId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateFirewallRuleEnable(ctx context.Context, firewallId string, ruleId string, enable bool) (firewall2.AsyncResponse, int, error) {
	var changeType string
	if enable {
		changeType = "ENABLE_PARTIAL"
	} else {
		changeType = "DISABLE_PARTIAL"
	}
	result, c, err := client.sdkClient.FirewallRuleV2Api.UpdateFirewallRuleEnabledV2(ctx, client.config.ProjectId, firewallId, firewall2.FirewallRuleChangeEnabledRequest{
		RuleEnabledChangeType: changeType,
		RuleIds:               []string{ruleId},
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
func (client *Client) DeleteFirewallRule(ctx context.Context, firewallId string, ruleId string) (firewall2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.FirewallRuleV2Api.DeleteFirewallRuleV2(ctx, client.config.ProjectId, firewallId, firewall2.FirewallRuleDeleteRequest{
		RuleDeletionType: "PARTIAL",
		RuleIds:          []string{ruleId},
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
