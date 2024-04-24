package peering

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/peering2"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *peering2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: peering2.NewAPIClient(config),
	}
}

func (client *Client) GetVpcPeeringList(ctx context.Context, request VpcPeeringListRequest) (peering2.ListResponseVpcPeeringResponse, error) {
	result, _, err := client.sdkClient.VpcPeeringOpenApiControllerApi.ListVpcPeerings(ctx, client.config.ProjectId, &peering2.VpcPeeringOpenApiControllerApiListVpcPeeringsOpts{
		ApproverVpcId:  optional.NewString(request.ApproverVpcId),
		RequesterVpcId: optional.NewString(request.RequesterVpcId),
		VpcPeeringName: optional.NewString(request.VpcPeeringName),
		CreatedBy:      optional.NewString(request.CreatedBy),
		Page:           optional.NewInt32(request.Page),
		Size:           optional.NewInt32(request.Size),
		Sort:           optional.NewInterface(request.Sort),
	})

	return result, err
}

func (client *Client) GetVpcPeeringForDelete(ctx context.Context, peeringId string) (peering2.VpcPeeringResponse, string, error) {
	result, _, err := client.sdkClient.VpcPeeringOpenApiControllerApi.ListVpcPeerings(ctx, client.config.ProjectId, &peering2.VpcPeeringOpenApiControllerApiListVpcPeeringsOpts{
		Size: optional.NewInt32(1000),
		Page: optional.NewInt32(0),
	})
	if err != nil {
		return peering2.VpcPeeringResponse{}, "", err
	}
	for _, peeringInfo := range result.Contents {
		if peeringInfo.VpcPeeringId == peeringId {
			return peeringInfo, peeringInfo.VpcPeeringState, nil
		}
	}
	return peering2.VpcPeeringResponse{}, "DELETED", nil
}

func (client *Client) GetVpcPeeringDetail(ctx context.Context, peeringId string) (peering2.VpcPeeringDetailResponse, string, error) {
	result, _, err := client.sdkClient.VpcPeeringOpenApiControllerApi.DetailVpcPeering(ctx, client.config.ProjectId, peeringId)
	if err != nil {
		return peering2.VpcPeeringDetailResponse{}, "", err
	}
	return result, result.VpcPeeringState, err
}

func (client *Client) CreateVpcPeering(ctx context.Context, request VpcPeeringCreateRequest) (peering2.VpcPeeringApprovalResponse, error) {
	result, _, err := client.sdkClient.VpcPeeringOpenApiV3ControllerApi.CreateVpcPeering1(ctx, client.config.ProjectId, peering2.VpcPeeringCreateV3Request{
		ApproverProjectId:     request.ApproverProjectId,
		ApproverVpcId:         request.ApproverVpcId,
		FirewallEnabled:       &request.FirewallEnabled,
		RequesterProjectId:    request.RequesterProjectId,
		RequesterVpcId:        request.RequesterVpcId,
		VpcPeeringDescription: request.VpcPeeringDescription,
		Tags:                  client.sdkClient.ToTagRequestList(request.Tags),
	})
	return result, err
}

func (client *Client) UpdateVpcPeeringDescription(ctx context.Context, peeringId string, description string) (peering2.VpcPeeringDetailResponse, error) {
	result, _, err := client.sdkClient.VpcPeeringOpenApiControllerApi.UpdateVpcPeeringDescription(ctx, client.config.ProjectId, peeringId, peering2.VpcPeeringUpdateDescriptionRequest{
		VpcPeeringDescription: description,
	})
	return result, err
}

func (client *Client) DeleteVpcPeering(ctx context.Context, peeringId string) error {
	_, _, err := client.sdkClient.VpcPeeringOpenApiControllerApi.DeleteVpcPeering(ctx, client.config.ProjectId, peeringId)
	return err
}

func (client *Client) ApproveVpcPeering(ctx context.Context, peeringId string, isFirewallEnabled bool) (peering2.VpcPeeringApprovalResponse, error) {
	result, _, err := client.sdkClient.VpcPeeringOpenApiControllerApi.ApproveVpcPeering(ctx, client.config.ProjectId, peeringId, peering2.VpcPeeringApproveRequest{
		FirewallEnabled: &isFirewallEnabled,
	})
	return result, err
}

func (client *Client) RejectVpcPeering(ctx context.Context, peeringId string) (peering2.VpcPeeringApprovalResponse, error) {
	result, _, err := client.sdkClient.VpcPeeringOpenApiControllerApi.RejectVpcPeering(ctx, client.config.ProjectId, peeringId)
	return result, err
}

func (client *Client) CancelVpcPeering(ctx context.Context, peeringId string) (peering2.VpcPeeringApprovalResponse, error) {
	result, _, err := client.sdkClient.VpcPeeringOpenApiControllerApi.CancelVpcPeering(ctx, client.config.ProjectId, peeringId)
	return result, err
}

/*
func (client *Client) MergePeeringRuleId(requesterVpcId, peeringId string) string {
	return requesterVpcId + ":" + peeringId
}

func (client *Client) SplitPeeringRuleId(ruleId string) (requesterVpcId, peeringId string) {
	colon := strings.Index(ruleId, ":")
	return ruleId[:colon], ruleId[colon+1:]
}
*/
