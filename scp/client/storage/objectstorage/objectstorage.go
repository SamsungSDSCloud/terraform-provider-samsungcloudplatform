package objectstorage

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/client"
	objectstorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/object-storage"
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

func (client *Client) ReadObjectStorageList(ctx context.Context, zoneId string, request objectstorage.ObjectStorageV3ControllerApiListObjectStorage3Opts) (objectstorage.PageResponseV2OfS3ObjectStoragesResponse, error) {
	result, _, err := client.sdkClient.ObjectStorageV3ControllerApi.ListObjectStorage3(ctx, client.config.ProjectId, zoneId, &request)
	return result, err
}

func (client *Client) CheckBucketName(ctx context.Context, obsId string, obsBucketName string) (bool, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV3ControllerApi.IsBucketNameDuplicated2(ctx, client.config.ProjectId, obsBucketName, obsId)
	if err != nil || result.IsObsBucketNameDuplicated == nil {
		return true, err
	}
	return *result.IsObsBucketNameDuplicated, err
}

func (client *Client) CreateBucket(ctx context.Context, request CreateBucketRequest) (objectstorage.S3BucketGetSyncResponse, error) {
	var accessIpAddressRanges []objectstorage.BucketIpsIpRange
	for _, b := range request.ObsBucketAccessIpAddressRanges {
		accessIpAddressRanges = append(accessIpAddressRanges, objectstorage.BucketIpsIpRange{
			ObsBucketAccessIpAddressRange: b.ObsBucketAccessIpAddressRange,
			Type_:                         b.Type,
		})
	}

	option := objectstorage.S3BucketCreateRequest{
		ObsId:                             request.ObsId,
		ZoneId:                            request.ZoneId,
		ObsBucketName:                     request.ObsBucketName,
		IsObsBucketIpAddressFilterEnabled: &request.IsObsBucketIpAddressFilterEnabled,
		ObsBucketFileEncryptionEnabled:    &request.ObsBucketFileEncryptionEnabled,
		ObsBucketVersionEnabled:           &request.ObsBucketVersionEnabled,
		ObsBucketAccessIpAddressRanges:    accessIpAddressRanges,
		Tags:                              []objectstorage.TagRequest{},
	}

	if request.ObsBucketFileEncryptionType != "" {
		option.ObsBucketFileEncryptionType = request.ObsBucketFileEncryptionType
	}
	if request.ObsBucketFileEncryptionAlgorithm != "" {
		option.ObsBucketFileEncryptionAlgorithm = request.ObsBucketFileEncryptionAlgorithm
	}

	result, _, err := client.sdkClient.ObjectStorageBucketV3ControllerApi.CreateBucket3(
		ctx,
		client.config.ProjectId,
		option)

	return result, err
}

func (client *Client) ReadBucket(ctx context.Context, obsBucketId string) (objectstorage.S3BucketGetResponse, int, error) {
	result, c, err := client.sdkClient.ObjectStorageBucketV3ControllerApi.ReadBucketInfo2(ctx, client.config.ProjectId, obsBucketId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadBucketList(ctx context.Context, request objectstorage.ObjectStorageBucketV3ControllerApiListBucket3Opts) (objectstorage.PageResponseV2OfS3BucketSearchResponse, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV3ControllerApi.ListBucket3(ctx, client.config.ProjectId, &request)
	return result, err
}

func (client *Client) DeleteBucket(ctx context.Context, obsBucketId string) (objectstorage.S3AsyncResponse, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV3ControllerApi.DeleteBucket3(ctx, client.config.ProjectId, obsBucketId)
	return result, err
}

func (client *Client) UpdateVersioning(ctx context.Context, obsBucketId string, versionEnabled bool) (objectstorage.S3BucketGetSyncResponse, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV3ControllerApi.UpdateVersioning2(ctx, client.config.ProjectId, obsBucketId,
		objectstorage.S3BucketUpdateRequest{
			ObsBucketVersionEnabled: &versionEnabled,
		})
	return result, err
}

func (client *Client) UpdateBucketEncryption(ctx context.Context, obsBucketId string, request UpdateBucketRequest) (objectstorage.S3BucketGetSyncResponse, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV3ControllerApi.UpdateBucketEncryption2(ctx, client.config.ProjectId, obsBucketId,
		objectstorage.S3BucketUpdateRequest{
			ObsBucketVersionEnabled:          &request.ObsBucketVersionEnabled,
			ObsBucketFileEncryptionType:      request.ObsBucketFileEncryptionType,
			ObsBucketFileEncryptionEnabled:   &request.ObsBucketFileEncryptionEnabled,
			ObsBucketFileEncryptionAlgorithm: request.ObsBucketFileEncryptionAlgorithm,
		})
	return result, err
}

func (client *Client) UpdateBucketDr(ctx context.Context, obsId string, drEnabled bool, replicaBucketId string) error {
	_, _, err := client.sdkClient.ObjectStorageDrV3ControllerApi.UpdateBucketDr1(ctx, client.config.ProjectId, obsId, objectstorage.S3BucketDrUpdateRequest{
		IsObsBucketDrEnabled: &drEnabled,
		ReplicaObsBucketId:   replicaBucketId,
	})
	return err
}

func (client *Client) CreateBucketIps(ctx context.Context, obsBucketId string, ipAddressFilterEnabled bool, request []ObsBucketAccessIpAddressInfo) (objectstorage.S3BucketGetSyncResponse, error) {

	var accessIpAddressRanges []objectstorage.BucketIpsIpRange
	for _, b := range request {
		accessIpAddressRanges = append(accessIpAddressRanges, objectstorage.BucketIpsIpRange{
			ObsBucketAccessIpAddressRange: b.ObsBucketAccessIpAddressRange,
			Type_:                         b.Type,
		})
	}

	result, _, err := client.sdkClient.ObjectStorageIpsV3ControllerApi.CreateBucketIps1(ctx, client.config.ProjectId, obsBucketId,
		objectstorage.S3BucketIpsRegisterUpdateRequest{
			IsObsBucketIpAddressFilterEnabled: &ipAddressFilterEnabled,
			ObsBucketAccessIpAddressRanges:    accessIpAddressRanges,
		})
	return result, err
}
