package objectstorage

import (
	"context"
	sdk "github.com/ScpDevTerra/trf-sdk/client"
	objectstorage "github.com/ScpDevTerra/trf-sdk/library/object-storage"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *objectstorage.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: objectstorage.NewAPIClient(config),
	}
}

func (client *Client) CreateObjectStorage(ctx context.Context, request CreateObjectStorageRequest) (objectstorage.S3BucketGetSyncResponse, error) {

	var accessIpAddressRanges []objectstorage.BucketIpsIpRange
	for _, b := range request.ObsBucketAccessIpAddressRanges {
		accessIpAddressRanges = append(accessIpAddressRanges, objectstorage.BucketIpsIpRange{
			ObsBucketAccessIpAddressRange: b.ObsBucketAccessIpAddressRange,
			Type_:                         b.Type,
		})
	}

	result, _, err := client.sdkClient.ObjectStorageBucketV2ControllerApi.CreateBucket(
		ctx,
		client.config.ProjectId,
		objectstorage.S3BucketCreateRequest{
			IsObsBucketIpAddressFilterEnabled: request.IsObsBucketIpAddressFilterEnabled,
			ObsBucketAccessIpAddressRanges:    accessIpAddressRanges,
			ObsBucketFileEncryptionAlgorithm:  request.ObsBucketFileEncryptionAlgorithm,
			ObsBucketFileEncryptionEnabled:    request.ObsBucketFileEncryptionEnabled,
			ObsBucketFileEncryptionType:       request.ObsBucketFileEncryptionType,
			ObsBucketName:                     request.ObsBucketName,
			ObsBucketVersionEnabled:           request.ObsBucketVersionEnabled,
			ObsId:                             request.ObsId,
			ZoneId:                            request.ZoneId,
		})

	return result, err
}

func (client *Client) ReadObjectStorage(ctx context.Context, obsBucketId string) (objectstorage.S3BucketGetResponse, int, error) {
	result, c, err := client.sdkClient.ObjectStorageBucketV2ControllerApi.ReadBucketInfo(ctx, client.config.ProjectId, obsBucketId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteObjectStorage(ctx context.Context, obsBucketId string) (objectstorage.S3AsyncResponse, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV2ControllerApi.DeleteBucket(ctx, client.config.ProjectId, obsBucketId)
	return result, err
}

func (client *Client) UpdateVersioning(ctx context.Context, obsBucketId string, versionEnabled bool) (objectstorage.S3BucketGetSyncResponse, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV2ControllerApi.UpdateVersioning(ctx, client.config.ProjectId, obsBucketId,
		objectstorage.S3BucketUpdateRequest{
			ObsBucketVersionEnabled: versionEnabled,
		})
	return result, err
}

func (client *Client) UpdateEncryption(ctx context.Context, obsBucketId string, request S3BucketUpdateRequest) (objectstorage.S3BucketGetSyncResponse, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV2ControllerApi.UpdateBucketEncryption(ctx, client.config.ProjectId, obsBucketId,
		objectstorage.S3BucketUpdateRequest{
			ObsBucketVersionEnabled:          request.ObsBucketVersionEnabled,
			ObsBucketFileEncryptionType:      request.ObsBucketFileEncryptionType,
			ObsBucketFileEncryptionEnabled:   request.ObsBucketFileEncryptionEnabled,
			ObsBucketFileEncryptionAlgorithm: request.ObsBucketFileEncryptionAlgorithm,
		})
	return result, err
}

func (client *Client) CreateBucketIps(ctx context.Context, obsBucketId string, ipAddressFilterEnabled bool, request []ObsBucketAccessIpAddressInfo) (objectstorage.S3BucketGetSyncResponse, error) {

	var accessIpAddressRanges []objectstorage.BucketIpsIpRange
	for _, b := range request {
		accessIpAddressRanges = append(accessIpAddressRanges, objectstorage.BucketIpsIpRange{
			ObsBucketAccessIpAddressRange: b.ObsBucketAccessIpAddressRange,
			Type_:                         b.Type,
		})
	}

	result, _, err := client.sdkClient.ObjectStorageIpsV2ControllerApi.CreateBucketIps(ctx, client.config.ProjectId, obsBucketId,
		objectstorage.S3BucketIpsRegisterUpdateRequest{
			IsObsBucketIpAddressFilterEnabled: ipAddressFilterEnabled,
			ObsBucketAccessIpAddressRanges:    accessIpAddressRanges,
		})
	return result, err
}
