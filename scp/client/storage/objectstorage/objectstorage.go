package objectstorage

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	objectstorage "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/object-storage"
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

func (client *Client) ReadObjectStorageList(ctx context.Context, serviceZoneId string, request objectstorage.ObjectStorageV4ControllerApiListObjectStorage6Opts) (objectstorage.ListResponseOfObjectStorageListV4Response, error) {
	result, _, err := client.sdkClient.ObjectStorageV4ControllerApi.ListObjectStorage6(ctx, client.config.ProjectId, serviceZoneId, &request)
	return result, err
}

func (client *Client) CheckBucketName(ctx context.Context, objectStorageId string, objectStorageBucketName string) (bool, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV4ControllerApi.CheckObjectStorageBucketDuplication1(ctx, client.config.ProjectId, objectStorageBucketName, objectStorageId)
	if err != nil || result.Result == nil {
		return true, err
	}
	return *result.Result, err
}

func (client *Client) CreateBucket(ctx context.Context, request CreateBucketRequest) (objectstorage.ObjectStorageBucketDetailV4Response, error) {
	var accessControlRules []objectstorage.AccessControlRule
	for _, b := range request.AccessControlRules {
		accessControlRules = append(accessControlRules, objectstorage.AccessControlRule{
			RuleValue: b.RuleValue,
			RuleType:  b.RuleType,
		})
	}

	option := objectstorage.ObjectStorageBucketCreateV4Request{
		ObjectStorageId:                          request.ObjectStorageId,
		ServiceZoneId:                            request.ServiceZoneId,
		ObjectStorageBucketName:                  request.ObjectStorageBucketName,
		ObjectStorageBucketAccessControlEnabled:  &request.ObjectStorageBucketAccessControlEnabled,
		ObjectStorageBucketFileEncryptionEnabled: &request.ObjectStorageBucketFileEncryptionEnabled,
		ObjectStorageBucketVersionEnabled:        &request.ObjectStorageBucketVersionEnabled,
		AccessControlRules:                       accessControlRules,
		ProductNames:                             request.ProductNames,
		Tags:                                     []objectstorage.TagRequest{},
	}

	result, _, err := client.sdkClient.ObjectStorageBucketV4ControllerApi.CreateObjectStorageBucket(
		ctx,
		client.config.ProjectId,
		option)

	return result, err
}

func (client *Client) ReadBucket(ctx context.Context, objectStorageBucketId string) (objectstorage.ObjectStorageBucketDetailV4Response, int, error) {
	result, c, err := client.sdkClient.ObjectStorageBucketV4ControllerApi.DetailObjectStorageBucket1(ctx, client.config.ProjectId, objectStorageBucketId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadBucketList(ctx context.Context, request objectstorage.ObjectStorageBucketV4ControllerApiListObjectStorageBuckets2Opts) (objectstorage.ListResponseOfObjectStorageBucketListV4Response, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV4ControllerApi.ListObjectStorageBuckets2(ctx, client.config.ProjectId, &request)
	return result, err
}

func (client *Client) DeleteBucket(ctx context.Context, objectStorageBucketId string) (objectstorage.AsyncResponse, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV4ControllerApi.DeleteObjectStorageBucket(ctx, client.config.ProjectId, objectStorageBucketId)
	return result, err
}

func (client *Client) UpdateVersioning(ctx context.Context, objectStorageBucketId string, versionEnabled bool) (objectstorage.ObjectStorageBucketDetailV4Response, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV4ControllerApi.UpdateObjectStorageBucketVersionEnabled1(ctx, client.config.ProjectId, objectStorageBucketId,
		objectstorage.ObjectStorageBucketVersionUpdateV4Request{
			ObjectStorageBucketVersionEnabled: &versionEnabled,
		})
	return result, err
}

func (client *Client) UpdateBucketEncryption(ctx context.Context, objectStorageBucketId string, fileEncryptionEnabled bool) (objectstorage.ObjectStorageBucketDetailV4Response, error) {
	result, _, err := client.sdkClient.ObjectStorageBucketV4ControllerApi.UpdateObjectStorageBucketFileEncryptionEnabled(ctx, client.config.ProjectId, objectStorageBucketId,
		objectstorage.ObjectStorageBucketFileEncryptionUpdateV4Request{
			ObjectStorageBucketFileEncryptionEnabled: &fileEncryptionEnabled,
		})
	return result, err
}

func (client *Client) UpdateBucketDr(ctx context.Context, objectStorageBucketId string, drEnabled bool, syncBucketId string) error {
	_, _, err := client.sdkClient.ObjectStorageDrV4ControllerApi.UpdateObjectStorageBucketDrEnabled1(ctx, client.config.ProjectId, objectStorageBucketId, objectstorage.ObjectStorageBucketDrUpdateV4Request{
		ObjectStorageBucketDrEnabled: &drEnabled,
		SyncObjectStorageBucketId:    syncBucketId,
	})
	return err
}

func (client *Client) CreateBucketIps(ctx context.Context, objectStorageBucketId string, objectStorageBucketAccessControlEnabled bool, request []AccessControlRule) (objectstorage.ObjectStorageBucketDetailV4Response, error) {

	var accessControlRules []objectstorage.AccessControlRule
	for _, b := range request {
		accessControlRules = append(accessControlRules, objectstorage.AccessControlRule{
			RuleValue: b.RuleValue,
			RuleType:  b.RuleType,
		})
	}

	result, _, err := client.sdkClient.ObjectStorageIpsV4ControllerApi.UpdateObjectStorageBucketAccessControl(ctx, client.config.ProjectId, objectStorageBucketId,
		objectstorage.ObjectStorageBucketAccessControlV4Request{
			ObjectStorageBucketAccessControlEnabled: &objectStorageBucketAccessControlEnabled,
			ObjectStorageBucketAccessControlRules:   accessControlRules,
		})
	return result, err
}
