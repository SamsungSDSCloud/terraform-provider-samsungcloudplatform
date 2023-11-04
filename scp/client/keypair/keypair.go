package keypair

import (
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	keypair "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/key-pair"
	"github.com/antihax/optional"
	"golang.org/x/net/context"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *keypair.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: keypair.NewAPIClient(config),
	}
}

func (client *Client) CreateKeyPair(ctx context.Context, request CreateRequest) (keypair.KeyPairCreateV1Response, error) {
	tags := make([]keypair.TagRequest, 0)
	for _, tag := range request.Tags {
		tags = append(tags, keypair.TagRequest{
			TagKey:   tag.TagKey,
			TagValue: tag.TagValue,
		})
	}
	result, _, err := client.sdkClient.KeyPairV1Api.CreateKeyPair1(ctx, client.config.ProjectId, keypair.KeyPairCreateV1Request{
		KeyPairName: request.KeyPairName,
		Tags:        tags,
	})

	return result, err
}

func (client *Client) DetailKeyPair(ctx context.Context, keyPairId string) (keypair.KeyPairV1Response, error) {
	result, _, err := client.sdkClient.KeyPairV1Api.DetailKeyPair(ctx, client.config.ProjectId, keyPairId)
	return result, err
}

func (client *Client) DeleteKeyPair(ctx context.Context, keyPairId string) error {
	_, err := client.sdkClient.KeyPairV1Api.DeleteKeyPair(ctx, client.config.ProjectId, keyPairId)
	return err
}

func (client *Client) GetKeyPair(ctx context.Context, keyPairId string) (keypair.KeyPairV1Response, int, error) {
	result, c, err := client.sdkClient.KeyPairV1Api.DetailKeyPair(ctx, client.config.ProjectId, keyPairId)

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ListKeyPairs(ctx context.Context, request ListKeyPairsRequestParam) (keypair.ListResponseOfKeyPairV1Response, error) {
	result, _, err := client.sdkClient.KeyPairV1Api.ListKeyPairs(ctx, client.config.ProjectId, &keypair.KeyPairV1ApiListKeyPairsOpts{
		KeyPairName: optional.NewString(request.KeyPairName),
		CreatedBy:   optional.NewString(request.CreatedBy),
		Page:        optional.NewInt32(request.Page),
		Size:        optional.NewInt32(request.Size),
		Sort:        optional.NewInterface(request.Sort),
	})
	return result, err
}
