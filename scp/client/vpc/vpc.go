package vpc

import (
	"context"
	"errors"
	"strings"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/vpc2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *vpc2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: vpc2.NewAPIClient(config),
	}
}

func (client *Client) GetVpcInfo(ctx context.Context, vpcId string) (vpc2.DetailVpcResponse, int, error) {
	result, c, err := client.sdkClient.VpcOpenApiControllerApi.DetailVpcV2(ctx, client.config.ProjectId, vpcId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateVpc(ctx context.Context, vpcId string, description string) (vpc2.DetailVpcResponse, error) {
	result, _, err := client.sdkClient.VpcOpenApiControllerApi.UpdateVpcDescriptionV2(ctx, client.config.ProjectId, vpcId, vpc2.ModifyVpcDescriptionRequest{
		VpcDescription: description,
	})
	return result, err
}

func (client *Client) GetVpcList(ctx context.Context) (vpc2.ListResponseVpcResponse, error) {
	result, _, err := client.sdkClient.VpcOpenApiControllerApi.ListVpcV2(ctx, client.config.ProjectId, &vpc2.VpcOpenApiControllerApiListVpcV2Opts{
		Size: optional.NewInt32(20),
		Page: optional.NewInt32(0),
	})
	return result, err
}

func (client *Client) GetVpcListV2(ctx context.Context, request ListVpcRequest) (vpc2.ListResponseVpcResponse, error) {
	result, _, err := client.sdkClient.VpcOpenApiControllerApi.ListVpcV2(ctx, client.config.ProjectId, &vpc2.VpcOpenApiControllerApiListVpcV2Opts{
		ServiceZoneId: optional.NewString(request.ServiceZoneId),
		VpcId:         optional.NewString(request.VpcId),
		VpcName:       optional.NewString(request.VpcName),
		VpcStates:     optional.NewInterface(request.VpcStates),
		CreatedBy:     optional.NewString(request.CreatedBy),
		Size:          optional.NewInt32(request.Size),
		Page:          optional.NewInt32(request.Page),
	})
	return result, err
}

func (client *Client) CreateVpc(ctx context.Context, vpcName string, vpcDescription string, serviceZoneId string, tags map[string]interface{}) (vpc2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VpcOpenApiV3ControllerApi.CreateVpcV3(ctx, client.config.ProjectId, vpc2.VpcCreateV3Request{
		VpcName:        vpcName,
		VpcDescription: vpcDescription,
		ServiceZoneId:  serviceZoneId,
		Tags:           client.sdkClient.ToTagRequestList(tags),
	})
	return result, err
}

func (client *Client) DeleteVpc(ctx context.Context, vpcId string) error {
	_, _, err := client.sdkClient.VpcOpenApiControllerApi.DeleteVpcV2(ctx, client.config.ProjectId, vpcId)
	return err
}

func (client *Client) GetVpcDnsList(ctx context.Context, vpcId string) (vpc2.ListResponseDnsUserZoneResponse, error) {
	result, _, err := client.sdkClient.VpcOpenApiControllerApi.ListDnsUserZoneV2(ctx, client.config.ProjectId, vpcId)
	return result, err
}

func (client *Client) CreateVpcDns(ctx context.Context, vpcId string, dnsUserZoneDomain string, dnsUserZoneServerIp string, dnsUserZoneSourceIp string, subnetId string) (vpc2.AsyncResponse, error) {
	result, _, err := client.sdkClient.VpcOpenApiControllerApi.CreateDnsUserZoneV2(ctx, client.config.ProjectId, vpcId, vpc2.CreateDnsUserZoneRequest{
		DnsUserZoneDomain:   dnsUserZoneDomain,
		DnsUserZoneServerIp: dnsUserZoneServerIp,
		DnsUserZoneSourceIp: dnsUserZoneSourceIp,
		SubnetId:            subnetId,
	})
	return result, err
}

func (client *Client) DeleteVpcDns(ctx context.Context, dnsId string) error {
	vpcId, dnsUserZoneId := client.SplitVpcDnsId(dnsId)
	_, _, err := client.sdkClient.VpcOpenApiControllerApi.DeleteDnsUserZoneV2(ctx, client.config.ProjectId, dnsUserZoneId, vpcId)
	return err
}

func (client *Client) GetVpcDnsInfoByDomain(ctx context.Context, vpcId string, dnsUserZoneDomain string) (vpc2.DnsUserZoneResponse, string, error) {
	result, err := client.GetVpcDnsList(ctx, vpcId)
	if err != nil {
		return vpc2.DnsUserZoneResponse{}, "", err
	}
	for _, dnsInfo := range result.Contents {
		if dnsInfo.DnsUserZoneDomain == dnsUserZoneDomain {
			return dnsInfo, dnsInfo.DnsUserZoneState, nil
		}
	}
	return vpc2.DnsUserZoneResponse{}, "", errors.New("domain query failed")
}

func (client *Client) GetVpcDnsInfoById(ctx context.Context, dnsId string) (vpc2.DnsUserZoneResponse, string, error) {
	vpcId, dnsRecordId := client.SplitVpcDnsId(dnsId)

	result, err := client.GetVpcDnsList(ctx, vpcId)
	if err != nil {
		return vpc2.DnsUserZoneResponse{}, "", err
	}
	for _, dnsInfo := range result.Contents {
		if dnsInfo.DnsUserZoneId == dnsRecordId {
			return dnsInfo, dnsInfo.DnsUserZoneState, nil
		}
	}
	return vpc2.DnsUserZoneResponse{}, "DELETED", nil
}

func (client *Client) MergeVpcDnsId(vpcId, dnsRecordId string) string {
	return vpcId + ":" + dnsRecordId
}

func (client *Client) SplitVpcDnsId(dnsId string) (vpcId, dnsRecordId string) {
	colon := strings.Index(dnsId, ":")
	return dnsId[:colon], dnsId[colon+1:]
}
