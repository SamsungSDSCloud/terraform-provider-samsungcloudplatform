package loadbalancer

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/loadbalancer2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *loadbalancer2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: loadbalancer2.NewAPIClient(config),
	}
}

/*
 Load Balancer
*/
func (client *Client) CreateLoadBalancer(ctx context.Context, blockId string, filewallEnabled bool, loadBalancerSize string, loadBalancerName string, productGroupId string, productId string, seviceIpCidr string, serviceZoneId string, vpcId string, description string) (loadbalancer2.AsyncResponse, error) {
	result, _, err := client.sdkClient.LoadBalancerOpenApiControllerApi.CreateLoadBalancer(ctx, client.config.ProjectId, loadbalancer2.LbRequest{
		BlockId:                 blockId,
		FirewallEnabled:         filewallEnabled,
		LoadBalancerName:        loadBalancerName,
		LoadBalancerSize:        loadBalancerSize,
		ProductGroupId:          productGroupId,
		ProductId:               productId,
		ServiceIpCidr:           seviceIpCidr,
		ServiceZoneId:           serviceZoneId,
		VpcId:                   vpcId,
		LoadBalancerDescription: description,
	})

	return result, err
}

func (client *Client) CheckLoadBalancerName(ctx context.Context, name string) (bool, error) {
	result, _, err := client.sdkClient.LoadBalancerOpenApiControllerApi.CheckLoadBalancerNameDuplication(ctx, client.config.ProjectId, name)
	v, ok := result["result"]
	if ok {
		return v, nil
	} else {
		return false, err
	}
}

func (client *Client) CheckLoadBalancerLimitValue(ctx context.Context, loadBalancerSize string, vpcId string) (bool, error) {
	result, _, err := client.sdkClient.LoadBalancerOpenApiControllerApi.CheckLoadBalancerLimitValue(ctx, client.config.ProjectId, loadBalancerSize, vpcId)
	v, ok := result["result"]
	if ok {
		return v, nil
	} else {
		return false, err
	}
}

func (client *Client) GetLoadBalancer(ctx context.Context, loadBalancerId string) (loadbalancer2.LbDetailResponse, int, error) {
	result, c, err := client.sdkClient.LoadBalancerOpenApiControllerApi.GetLoadBalancer(ctx, client.config.ProjectId, loadBalancerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteLoadBalancer(ctx context.Context, loadBalancerId string) (loadbalancer2.AsyncResponse, error) {
	result, _, err := client.sdkClient.LoadBalancerOpenApiControllerApi.DeleteLoadBalancer(ctx, client.config.ProjectId, loadBalancerId)
	return result, err
}

func (client *Client) UpdateLoadBalancerDescription(ctx context.Context, loadBalancerId string, description string) (loadbalancer2.LbDetailResponse, error) {
	result, _, err := client.sdkClient.LoadBalancerOpenApiControllerApi.UpdateLoadBalancer(ctx, client.config.ProjectId, loadBalancerId, loadbalancer2.LbChangeRequest{
		LoadBalancerDescription: description,
	})
	return result, err
}

func (client *Client) GetLoadBalancerList(ctx context.Context, request *loadbalancer2.LoadBalancerOpenApiControllerApiGetLoadBalancerListOpts) (loadbalancer2.ListResponseOfLbResponse, int, error) {
	result, c, err := client.sdkClient.LoadBalancerOpenApiControllerApi.GetLoadBalancerList(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

/*
 LB Profile
*/
func (client *Client) CheckLbProfileName(ctx context.Context, loadBalancerId string, name string) (bool, error) {
	result, _, err := client.sdkClient.LbProfileOpenApiControllerApi.CheckLbProfileNameDuplication(ctx, client.config.ProjectId, loadBalancerId, name)
	v, ok := result["result"]
	if ok {
		return v, nil
	} else {
		return false, err
	}
}

func (client *Client) CreateLbProfile(ctx context.Context, loadBalancerId string, layerType string, category string, name string, pfType string, protocol string, requestHeaderSize int, responseHeaderSize int, responseTimeout int, sessionTimeout int, xforwardedFor string) (loadbalancer2.AsyncResponse, error) {
	attr := loadbalancer2.LbProfileAttrCreateRequest{
		RequestHeaderSize:  int32(requestHeaderSize),
		ResponseHeaderSize: int32(responseHeaderSize),
		ResponseTimeout:    int32(responseTimeout),
		SessionTimeout:     int32(sessionTimeout),
		XforwardedFor:      xforwardedFor,
	}

	result, _, err := client.sdkClient.LbProfileOpenApiControllerApi.CreateLoadBalancerProfile(ctx, client.config.ProjectId, loadBalancerId, loadbalancer2.LbProfileCreateRequest{
		LayerType:         layerType,
		LbProfileAttrs:    &attr, // *LbProfileAttrCreateRequest `json:"lbProfileAttrs,omitempty"`
		LbProfileCategory: category,
		LbProfileName:     name,
		LbProfileType:     pfType,
		Protocol:          protocol,
	})

	return result, err
}

func (client *Client) GetLbProfile(ctx context.Context, lbProfileId string, loadBalancerId string) (loadbalancer2.LbProfileDetailResponse, int, error) {
	result, c, err := client.sdkClient.LbProfileOpenApiControllerApi.GetLoadBalancerProfile(ctx, client.config.ProjectId, lbProfileId, loadBalancerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateLbProfile(ctx context.Context, lbProfileId string, loadBalancerId string, requestHeaderSize int, responseHeaderSize int, responseTimeout int, sessionTimeout int, XforwardedFor string) (loadbalancer2.AsyncResponse, error) {
	request := loadbalancer2.LbProfileAttrModifyRequest{
		RequestHeaderSize:  int32(requestHeaderSize),
		ResponseHeaderSize: int32(responseHeaderSize),
		ResponseTimeout:    int32(responseTimeout),
		SessionTimeout:     int32(sessionTimeout),
		XforwardedFor:      XforwardedFor,
	}
	result, _, err := client.sdkClient.LbProfileOpenApiControllerApi.UpdateLoadBalancerProfile(ctx, client.config.ProjectId, lbProfileId, loadBalancerId, loadbalancer2.LbProfileChangeRequest{
		LbProfileAttrs: &request,
	})
	return result, err
}

func (client *Client) DeleteLbProfile(ctx context.Context, lbProfileId string, loadBalancerId string) (loadbalancer2.AsyncResponse, error) {
	result, _, err := client.sdkClient.LbProfileOpenApiControllerApi.DeleteLoadBalancerProfile(ctx, client.config.ProjectId, lbProfileId, loadBalancerId)
	return result, err
}

/*
 LB Server Group
*/
type LbServerGroupMember = loadbalancer2.LbServerGroupMemberRequest
type LbServerGroupMonitor = loadbalancer2.LbMonitorRequest

func (client *Client) CheckServerGroupNameDuplicated(ctx context.Context, loadBalancerId string, lbServerGroupName string) (bool, error) {
	result, _, err := client.sdkClient.LbServerGroupOpenApiControllerApi.CheckLbServerGroupNameDuplication(ctx, client.config.ProjectId, loadBalancerId, lbServerGroupName)

	v, ok := result["result"]
	if ok {
		return v, nil
	} else {
		return false, err
	}
}

func (client *Client) GetLbServerGroup(ctx context.Context, lbServerGroupId string, loadBalancerId string) (loadbalancer2.LbServerGroupDetailResponse, int, error) {
	result, c, err := client.sdkClient.LbServerGroupOpenApiControllerApi.GetLoadBalancerServerGroup(ctx, client.config.ProjectId, lbServerGroupId, loadBalancerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateLbServerGroup(ctx context.Context, loadBalancerId string, algorithm string, name string, monitor *LbServerGroupMonitor, members []LbServerGroupMember, tcpMultiplexingEnabled bool) (loadbalancer2.AsyncResponse, error) {

	result, _, err := client.sdkClient.LbServerGroupOpenApiControllerApi.CreateLoadBalancerServerGroup(ctx, client.config.ProjectId, loadBalancerId, loadbalancer2.LbServerGroupCreateRequest{
		LbAlgorithm:            algorithm,
		LbMonitor:              monitor,
		LbServerGroupMembers:   members,
		LbServerGroupName:      name,
		TcpMultiplexingEnabled: tcpMultiplexingEnabled,
	})
	return result, err
}

func (client *Client) UpdateLbServerGroup(ctx context.Context, tcpMultiplexingEnabled bool, lbAlgorithm string, lbServerGroupId string, loadBalancerId string, monitor *LbServerGroupMonitor, members []LbServerGroupMember) (loadbalancer2.AsyncResponse, error) {
	result, _, err := client.sdkClient.LbServerGroupOpenApiControllerApi.UpdateLoadBalancerServerGroup(ctx, client.config.ProjectId, lbServerGroupId, loadBalancerId, loadbalancer2.LbServerGroupChangeRequest{
		LbAlgorithm:            lbAlgorithm,
		LbMonitor:              monitor,
		TcpMultiplexingEnabled: tcpMultiplexingEnabled,
		LbServerGroupMembers:   members,
	})

	return result, err
}

func (client *Client) DeleteLbServerGroup(ctx context.Context, lbServerGroupId string, loadBalancerId string) (loadbalancer2.AsyncResponse, error) {
	result, _, err := client.sdkClient.LbServerGroupOpenApiControllerApi.DeleteLoadBalancerServerGroup(ctx, client.config.ProjectId, lbServerGroupId, loadBalancerId)
	return result, err
}

func (client *Client) GetLbServerGroupList(ctx context.Context, loadBalancerId string, request *loadbalancer2.LbServerGroupOpenApiControllerApiGetLoadBalancerServerGroupListOpts) (loadbalancer2.ListResponseOfLbServerGroupResponse, int, error) {
	result, c, err := client.sdkClient.LbServerGroupOpenApiControllerApi.GetLoadBalancerServerGroupList(ctx, client.config.ProjectId, loadBalancerId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

/*
 LB Server Group
*/
type LbServiceRule = loadbalancer2.LbRuleRequest

func (client *Client) CheckLbServiceNameDuplicated(ctx context.Context, loadBalancerId string, lbServiceName string) (bool, error) {
	result, _, err := client.sdkClient.LbServiceOpenApiControllerApi.CheckLoadBalancerServiceNameDuplication(ctx, client.config.ProjectId, loadBalancerId, lbServiceName)

	v, ok := result["result"]
	if ok {
		return v, nil
	} else {
		return false, err
	}
}

func (client *Client) CreateLbService(
	ctx context.Context,
	loadBalancerId string,
	appProfileId string,
	defaultForwardingPorts string,
	layerType string,
	lbServiceName string,
	persistence string, persistenceProfileId string,
	protocol string,
	lbRules []LbServiceRule,
	serviceIpAddr string,
	servicePorts string,
	lbServiceIpId string,
	serverCertificateId string,
	clientCertificateId string,
	useAccessLog bool) (loadbalancer2.AsyncResponse, error) {

	result, _, err := client.sdkClient.LbServiceOpenApiControllerApi.CreateLoadBalancerService(ctx, client.config.ProjectId, loadBalancerId, loadbalancer2.LbServiceRequest{
		ApplicationProfileId:   appProfileId,
		DefaultForwardingPorts: defaultForwardingPorts,
		LayerType:              layerType,
		LbRules:                lbRules,
		LbServiceName:          lbServiceName,
		Persistence:            persistence,
		PersistenceProfileId:   persistenceProfileId,
		Protocol:               protocol,
		ServiceIpAddress:       serviceIpAddr,
		ServicePorts:           servicePorts,
		// the belows are not in use for now
		LbServiceIpId:       lbServiceIpId,
		ServerCertificateId: serverCertificateId,
		ClientCertificateId: clientCertificateId,
		UseAccessLog:        useAccessLog,
	})

	return result, err
}

func (client *Client) GetLbService(ctx context.Context, lbServiceId string, loadBalancerId string) (loadbalancer2.LbServiceDetailResponse, int, error) {
	result, c, err := client.sdkClient.LbServiceOpenApiControllerApi.GetLoadBalancerService(ctx, client.config.ProjectId, lbServiceId, loadBalancerId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteLbService(ctx context.Context, lbServiceId string, loadBalancerId string) (loadbalancer2.AsyncResponse, error) {
	result, _, err := client.sdkClient.LbServiceOpenApiControllerApi.DeleteLoadBalancerService(ctx, client.config.ProjectId, lbServiceId, loadBalancerId)
	return result, err
}

func (client *Client) UpdateLbService(ctx context.Context, lbServiceId string, loadBalancerId string,
	appProfileId string, clientCertId string, defaultForwardingPorts string, lbRules []LbServiceRule, persistence string,
	persistenceProfileId string, serverCertId string, servicePorts string, useAccessLog bool) (loadbalancer2.AsyncResponse, error) {
	result, _, err := client.sdkClient.LbServiceOpenApiControllerApi.UpdateLoadBalancerService(ctx, client.config.ProjectId, lbServiceId, loadBalancerId, loadbalancer2.LbServiceChangeRequest{
		ApplicationProfileId:   appProfileId,
		ClientCertificateId:    clientCertId,
		DefaultForwardingPorts: defaultForwardingPorts,
		//LbRules:                lbRules,
		Persistence:          persistence,
		PersistenceProfileId: persistenceProfileId,
		ServerCertificateId:  serverCertId,
		ServicePorts:         servicePorts,
		UseAccessLog:         true,
	})

	return result, err
}

func (client *Client) GetLbServiceList(ctx context.Context, loadBalancerId string, request *loadbalancer2.LbServiceOpenApiControllerApiGetLoadBalancerServiceListOpts) (loadbalancer2.ListResponseOfLbServiceResponse, int, error) {
	result, c, err := client.sdkClient.LbServiceOpenApiControllerApi.GetLoadBalancerServiceList(ctx, client.config.ProjectId, loadBalancerId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetLbServiceIpList(ctx context.Context, loadBalancerId string, request *loadbalancer2.LbServiceOpenApiControllerApiGetLoadBalancerServiceIpListOpts) (loadbalancer2.ListResponseOfLbServiceIpResponse, int, error) {
	result, c, err := client.sdkClient.LbServiceOpenApiControllerApi.GetLoadBalancerServiceIpList(ctx, client.config.ProjectId, loadBalancerId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
