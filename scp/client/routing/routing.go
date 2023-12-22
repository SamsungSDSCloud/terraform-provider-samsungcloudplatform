package routing

import (
	"context"
	"fmt"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/routing2"
	"github.com/antihax/optional"
	"strings"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *routing2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: routing2.NewAPIClient(config),
	}
}

func (client *Client) GetVpcRoutingTableList(ctx context.Context) (routing2.ListResponseOfVpcRoutingTableListResponse, error) {
	result, _, err := client.sdkClient.VpcRoutingTableOpenApiControllerApi.ListVpcRoutingTables(ctx, client.config.ProjectId, &routing2.VpcRoutingTableOpenApiControllerApiListVpcRoutingTablesOpts{
		Size: optional.NewInt32(20),
		Page: optional.NewInt32(0),
	})
	return result, err
}

func (client *Client) GetVpcRoutingTableListV2(ctx context.Context, request ListVpcRoutingTableRequest) (routing2.ListResponseOfVpcRoutingTableListResponse, error) {
	result, _, err := client.sdkClient.VpcRoutingTableOpenApiControllerApi.ListVpcRoutingTables(ctx, client.config.ProjectId, &routing2.VpcRoutingTableOpenApiControllerApiListVpcRoutingTablesOpts{
		RoutingTableId:   optional.NewString(request.RoutingTableId),
		RoutingTableName: optional.NewString(request.RoutingTableName),
		VpcId:            optional.NewString(request.VpcId),
		CreatedBy:        optional.NewString(request.CreatedBy),
		Sort:             optional.NewInterface(request.Sort),
		Size:             optional.NewInt32(10000),
		Page:             optional.NewInt32(0),
	})
	return result, err
}

func (client *Client) GetVpcRoutingTableListByVpcId(ctx context.Context, vpcId string) (routing2.ListResponseOfVpcRoutingTableListResponse, error) {
	result, _, err := client.sdkClient.VpcRoutingTableOpenApiControllerApi.ListVpcRoutingTables(ctx, client.config.ProjectId, &routing2.VpcRoutingTableOpenApiControllerApiListVpcRoutingTablesOpts{
		VpcId: optional.NewString(vpcId),
	})
	return result, err
}

func (client *Client) GetVpcRoutingTableDetail(ctx context.Context, routingTableId string) (routing2.VpcRoutingTableDetailResponse, error) {
	result, _, err := client.sdkClient.VpcRoutingTableOpenApiControllerApi.DetailVpcRoutingTables(ctx, client.config.ProjectId, routingTableId)
	return result, err
}

func (client *Client) GetVpcRoutingRulesRoute(ctx context.Context, routingTableId string) (routing2.ListResponseOfRoutingRuleRouteListResponse, error) {
	result, _, err := client.sdkClient.VpcRoutingRuleOpenApiControllerApi.ListVpcRoutingRulesRoute(ctx, client.config.ProjectId, routingTableId)
	return result, err
}

func (client *Client) GetVpcRoutingRulesList(ctx context.Context, routingTableId string, request ListVpcRoutingRulesRequest) (routing2.ListResponseOfVpcRoutingRuleListResponse, error) {
	options := routing2.VpcRoutingRuleOpenApiControllerApiListVpcRoutingRulesOpts{
		DestinationNetworkCidr:   optional.NewString(request.DestinationNetworkCidr),
		RoutingRuleId:            optional.NewString(request.RoutingRuleId),
		SourceServiceInterfaceId: optional.NewString(request.SourceServiceInterfaceId),
		Page:                     optional.NewInt32(0),
		Size:                     optional.NewInt32(10000),
		Sort:                     optional.NewInterface(request.Sort),
	}
	if request.Editable != "" {
		options.Editable = optional.NewBool(request.Editable == "true")
	}

	result, _, err := client.sdkClient.VpcRoutingRuleOpenApiControllerApi.ListVpcRoutingRules(ctx, client.config.ProjectId, routingTableId, &options)
	return result, err
}

func (client *Client) GetVpcRoutingRulesByCidr(ctx context.Context, routingTableId string, destinationNetworkCidr string) (routing2.VpcRoutingRuleListResponse, string, error) {
	result, _, err := client.sdkClient.VpcRoutingRuleOpenApiControllerApi.ListVpcRoutingRules(ctx, client.config.ProjectId, routingTableId, &routing2.VpcRoutingRuleOpenApiControllerApiListVpcRoutingRulesOpts{
		DestinationNetworkCidr: optional.NewString(destinationNetworkCidr),
	})
	if err != nil {
		return routing2.VpcRoutingRuleListResponse{}, "", err
	}
	return result.Contents[0], result.Contents[0].RoutingRuleState, err
}

func (client *Client) GetVpcRoutingRulesById(ctx context.Context, ruleId string) (routing2.VpcRoutingRuleListResponse, string, error) {
	routingTableId, routingRuleId := client.SplitRoutingRuleId(ruleId)

	result, _, err := client.sdkClient.VpcRoutingRuleOpenApiControllerApi.ListVpcRoutingRules(ctx, client.config.ProjectId, routingTableId, &routing2.VpcRoutingRuleOpenApiControllerApiListVpcRoutingRulesOpts{})
	if err != nil {
		return routing2.VpcRoutingRuleListResponse{}, "", err
	}

	for _, rule := range result.Contents {
		if rule.RoutingRuleId == routingRuleId {
			return rule, rule.RoutingRuleState, nil
		}
	}
	return routing2.VpcRoutingRuleListResponse{}, "DELETED", nil
}

func (client *Client) CheckDuplicationRoutingRule(ctx context.Context, routingTableId string, destinationNetworkCidr string) (bool, error) {
	result, _, err := client.sdkClient.VpcRoutingRuleOpenApiControllerApi.CheckDuplicationVpcRoutingRule(ctx, client.config.ProjectId, routingTableId, destinationNetworkCidr)
	if result.Result == nil {
		return false, err
	}
	return *result.Result, err
}

func (client *Client) CreateRoutingRules(ctx context.Context, routingTableId string, request CreateRoutingRulesRequest) error {
	var rules routing2.CreateRoutingRulesRequest
	for _, rule := range request.RoutingRules {
		rules.RoutingRules = append(rules.RoutingRules, routing2.RoutingRule{
			DestinationNetworkCidr:     rule.DestinationNetworkCidr,
			SourceServiceInterfaceId:   rule.SourceServiceInterfaceId,
			SourceServiceInterfaceName: rule.SourceServiceInterfaceName,
		})
	}
	_, _, err := client.sdkClient.VpcRoutingRuleOpenApiControllerApi.CreateVpcRoutingRules(ctx, client.config.ProjectId, routingTableId, routing2.CreateRoutingRulesRequest{
		RoutingRules: rules.RoutingRules,
	})

	return err
}

func (client *Client) DeleteRoutingRules(ctx context.Context, routingTableId string, routingRuleId string) error {
	var rules routing2.DeleteRoutingRulesRequest
	rules.RoutingRuleIds = append(rules.RoutingRuleIds, routingRuleId)
	_, _, err := client.sdkClient.VpcRoutingRuleOpenApiControllerApi.DeleteVpcRoutingRules(ctx, client.config.ProjectId, routingTableId, routing2.DeleteRoutingRulesRequest{
		RoutingRuleIds: rules.RoutingRuleIds,
	})

	return err
}

func (client *Client) MergeRoutingRuleId(routingTableId, routingRuleId string) string {
	return routingTableId + ":" + routingRuleId
}

func (client *Client) SplitRoutingRuleId(ruleId string) (routingTableId, routingRuleId string) {
	colon := strings.Index(ruleId, ":")
	return ruleId[:colon], ruleId[colon+1:]
}

// DirectConnect
func (client *Client) GetDCRoutingTableList(ctx context.Context, routingTableId string, routingTableName string, directConnectConnectionId string, createdBy string) (routing2.ListResponseOfDcRoutingTableListResponse, error) {
	result, _, err := client.sdkClient.DirectConnectRoutingTableOpenApiControllerApi.ListDcRoutingTables(ctx, client.config.ProjectId, &routing2.DirectConnectRoutingTableOpenApiControllerApiListDcRoutingTablesOpts{
		RoutingTableId:            optional.NewString(routingTableId),
		RoutingTableName:          optional.NewString(routingTableName),
		DirectConnectConnectionId: optional.NewString(directConnectConnectionId),
		CreatedBy:                 optional.NewString(createdBy),
		Page:                      optional.NewInt32(0),
		Size:                      optional.NewInt32(10000),
	})
	return result, err
}

func (client *Client) GetDCRoutingRulesRoute(ctx context.Context, routingTableId string) (routing2.ListResponseOfRoutingRuleRouteListResponse, error) {
	result, _, err := client.sdkClient.DirectConnectRoutingRuleOpenApiControllerApi.ListDcRoutingRulesRoute(ctx, client.config.ProjectId, routingTableId)
	return result, err
}

func (client *Client) CreateDCRoutingRules(ctx context.Context, routingTableId string, request CreateRoutingRulesRequest) error {
	var rules routing2.CreateRoutingRulesRequest
	for _, rule := range request.RoutingRules {
		rules.RoutingRules = append(rules.RoutingRules, routing2.RoutingRule{
			DestinationNetworkCidr:     rule.DestinationNetworkCidr,
			SourceServiceInterfaceId:   rule.SourceServiceInterfaceId,
			SourceServiceInterfaceName: rule.SourceServiceInterfaceName,
		})
	}
	_, _, err := client.sdkClient.DirectConnectRoutingRuleOpenApiControllerApi.CreateDcRoutingRules(ctx, client.config.ProjectId, routingTableId, routing2.CreateRoutingRulesRequest{
		RoutingRules: rules.RoutingRules,
	})

	return err
}

func (client *Client) GetDCRoutingRulesByCidr(ctx context.Context, routingTableId string, destinationNetworkCidr string) (routing2.DcRoutingRuleListResponse, string, error) {
	result, _, err := client.sdkClient.DirectConnectRoutingRuleOpenApiControllerApi.ListDcRoutingRules(ctx, client.config.ProjectId, routingTableId, &routing2.DirectConnectRoutingRuleOpenApiControllerApiListDcRoutingRulesOpts{
		DestinationNetworkCidr: optional.NewString(destinationNetworkCidr),
	})
	if err != nil {
		return routing2.DcRoutingRuleListResponse{}, "", err
	}
	if result.TotalCount > 1 {
		return routing2.DcRoutingRuleListResponse{}, "", fmt.Errorf("Duplicate routing rules - %s", routingTableId)
	}
	return result.Contents[0], result.Contents[0].RoutingRuleState, err
}

func (client *Client) GetDCRoutingRulesById(ctx context.Context, ruleId string) (routing2.DcRoutingRuleListResponse, string, error) {
	routingTableId, routingRuleId := client.SplitRoutingRuleId(ruleId)

	result, _, err := client.sdkClient.DirectConnectRoutingRuleOpenApiControllerApi.ListDcRoutingRules(ctx, client.config.ProjectId, routingTableId, &routing2.DirectConnectRoutingRuleOpenApiControllerApiListDcRoutingRulesOpts{})
	if err != nil {
		return routing2.DcRoutingRuleListResponse{}, "", err
	}

	for _, rule := range result.Contents {
		if rule.RoutingRuleId == routingRuleId {
			return rule, rule.RoutingRuleState, nil
		}
	}
	return routing2.DcRoutingRuleListResponse{}, "DELETED", nil
}

func (client *Client) DeleteDCRoutingRules(ctx context.Context, routingTableId string, routingRuleId string) error {
	var rules routing2.DeleteRoutingRulesRequest
	rules.RoutingRuleIds = append(rules.RoutingRuleIds, routingRuleId)
	_, _, err := client.sdkClient.DirectConnectRoutingRuleOpenApiControllerApi.DeleteDcRoutingRules(ctx, client.config.ProjectId, routingTableId, routing2.DeleteRoutingRulesRequest{
		RoutingRuleIds: rules.RoutingRuleIds,
	})

	return err
}

func (client *Client) CheckDCDuplicationRoutingRule(ctx context.Context, routingTableId string, destinationNetworkCidr string) (bool, error) {
	result, _, err := client.sdkClient.DirectConnectRoutingRuleOpenApiControllerApi.CheckDuplicationDcRoutingRule(ctx, client.config.ProjectId, routingTableId, destinationNetworkCidr)
	if result.Result == nil {
		return false, err
	}
	return *result.Result, err
}

func (client *Client) GetDCRoutingRulesList(ctx context.Context, routingTableId string, request ListVpcRoutingRulesRequest) (routing2.ListResponseOfDcRoutingRuleListResponse, error) {
	options := routing2.DirectConnectRoutingRuleOpenApiControllerApiListDcRoutingRulesOpts{
		DestinationNetworkCidr:   optional.NewString(request.DestinationNetworkCidr),
		RoutingRuleId:            optional.NewString(request.RoutingRuleId),
		SourceServiceInterfaceId: optional.NewString(request.SourceServiceInterfaceId),
		Page:                     optional.NewInt32(0),
		Size:                     optional.NewInt32(10000),
		Sort:                     optional.NewInterface(request.Sort),
	}
	if request.Editable != "" {
		options.Editable = optional.NewBool(request.Editable == "true")
	}

	result, _, err := client.sdkClient.DirectConnectRoutingRuleOpenApiControllerApi.ListDcRoutingRules(ctx, client.config.ProjectId, routingTableId, &options)
	return result, err
}

//Transit Gateway

func (client *Client) GetTgwRoutingTableList(ctx context.Context, request ListTgwRoutingTableRequest) (routing2.ListResponseOfTgwRoutingTableListResponse, error) {
	result, _, err := client.sdkClient.TransitGatewayRoutingTableOpenApiControllerApi.ListTgwRoutingTables(ctx, client.config.ProjectId, &routing2.TransitGatewayRoutingTableOpenApiControllerApiListTgwRoutingTablesOpts{
		RoutingTableId:             optional.NewString(request.RoutingTableId),
		RoutingTableName:           optional.NewString(request.RoutingTableName),
		TransitGatewayConnectionId: optional.NewString(request.TransitGatewayConnectionId),
		CreatedBy:                  optional.NewString(request.CreatedBy),
		Sort:                       optional.NewInterface(request.Sort),
		Page:                       optional.NewInt32(0),
		Size:                       optional.NewInt32(20),
	})

	return result, err
}

func (client *Client) GetTgwRoutingTableDetail(ctx context.Context, routingTableId string) (routing2.TgwRoutingTableDetailResponse, error) {
	result, _, err := client.sdkClient.TransitGatewayRoutingTableOpenApiControllerApi.DetailTgwRoutingTables(ctx, client.config.ProjectId, routingTableId)
	return result, err
}

func (client *Client) GetTgwRoutingRuleList(ctx context.Context, routingTableId string, request ListTgwRoutingRuleRequest) (routing2.ListResponseOfTgwRoutingRuleListResponse, error) {

	result, _, err := client.sdkClient.TransitGatewayRoutingRuleOpenApiControllerApi.ListTgwRoutingRules(ctx, client.config.ProjectId, routingTableId, &routing2.TransitGatewayRoutingRuleOpenApiControllerApiListTgwRoutingRulesOpts{
		DestinationNetworkCidr:   optional.NewString(request.DestinationNetworkCidr),
		RoutingRuleId:            optional.NewString(request.RoutingRuleId),
		SourceServiceInterfaceId: optional.NewString(request.SourceServiceInterfaceId),
		Page:                     optional.NewInt32(0),
		Size:                     optional.NewInt32(20),
		Sort:                     optional.NewInterface(request.Sort),
	})
	return result, err
}

func (client *Client) GetTgwRoutingRuleByCidr(ctx context.Context, routingTableId string, destinationNetworkCidr string) (routing2.TgwRoutingRuleListResponse, string, error) {
	result, _, err := client.sdkClient.TransitGatewayRoutingRuleOpenApiControllerApi.ListTgwRoutingRules(ctx, client.config.ProjectId, routingTableId, &routing2.TransitGatewayRoutingRuleOpenApiControllerApiListTgwRoutingRulesOpts{
		DestinationNetworkCidr: optional.NewString(destinationNetworkCidr),
		Page:                   optional.NewInt32(0),
		Size:                   optional.NewInt32(20),
	})

	if err != nil {
		return routing2.TgwRoutingRuleListResponse{}, "", err
	}
	if result.TotalCount > 1 {
		return routing2.TgwRoutingRuleListResponse{}, "", err
	}

	return result.Contents[0], result.Contents[0].RoutingRuleState, err
}

func (client *Client) GetTgwRoutingRuleById(ctx context.Context, routingTableId string, routingRuleId string) (routing2.TgwRoutingRuleListResponse, string, error) {
	result, _, err := client.sdkClient.TransitGatewayRoutingRuleOpenApiControllerApi.ListTgwRoutingRules(ctx, client.config.ProjectId, routingTableId, &routing2.TransitGatewayRoutingRuleOpenApiControllerApiListTgwRoutingRulesOpts{})

	if err != nil {
		return routing2.TgwRoutingRuleListResponse{}, "", err
	}
	for _, rule := range result.Contents {
		if rule.RoutingRuleId == routingRuleId {
			return rule, rule.RoutingRuleState, nil
		}
	}

	return routing2.TgwRoutingRuleListResponse{}, "DELETED", nil
}

func (client *Client) GetTgwRoutingRoutes(ctx context.Context, routingTableId string) (routing2.ListResponseOfRoutingRuleRouteListResponse, error) {
	result, _, err := client.sdkClient.TransitGatewayRoutingRuleOpenApiControllerApi.ListTgwRoutingRulesRoute(ctx, client.config.ProjectId, routingTableId)

	return result, err
}

func (client *Client) CreateTgwRoutingRules(ctx context.Context, routingTableId string, request CreateRoutingRulesRequest) error {
	var rules routing2.CreateRoutingRulesRequest
	for _, rule := range request.RoutingRules {
		rules.RoutingRules = append(rules.RoutingRules, routing2.RoutingRule{
			DestinationNetworkCidr:     rule.DestinationNetworkCidr,
			SourceServiceInterfaceId:   rule.SourceServiceInterfaceId,
			SourceServiceInterfaceName: rule.SourceServiceInterfaceName,
		})
	}
	_, _, err := client.sdkClient.TransitGatewayRoutingRuleOpenApiControllerApi.CreateTgwRoutingRules(ctx, client.config.ProjectId, routingTableId, rules)

	return err
}

func (client *Client) DeleteTgwRoutingRules(ctx context.Context, routingTableId string, routingRuleId string) error {
	var rules routing2.DeleteRoutingRulesRequest
	rules.RoutingRuleIds = append(rules.RoutingRuleIds, routingRuleId)
	_, _, err := client.sdkClient.TransitGatewayRoutingRuleOpenApiControllerApi.DeleteTgwRoutingRules(ctx, client.config.ProjectId, routingTableId, rules)
	return err
}

func (client *Client) CheckDuplicationTgwRoutingRule(ctx context.Context, routingTableId string, destinationNetworkCidr string) (bool, error) {
	result, _, err := client.sdkClient.TransitGatewayRoutingRuleOpenApiControllerApi.CheckDuplicationTgwRoutingRule(ctx, client.config.ProjectId, routingTableId, destinationNetworkCidr)
	if result.Result == nil {
		return false, err
	}

	return *result.Result, err
}
