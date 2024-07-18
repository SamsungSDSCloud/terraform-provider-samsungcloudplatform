package objectstorage

import (
	"github.com/antihax/optional"
)

type CreateBucketRequest struct {
	ObjectStorageBucketAccessControlEnabled  bool
	AccessControlRules                       []AccessControlRule
	ObjectStorageBucketFileEncryptionEnabled bool
	ObjectStorageBucketName                  string
	ObjectStorageBucketUserPurpose           string
	ObjectStorageBucketVersionEnabled        bool
	ObjectStorageId                          string
	ServiceZoneId                            string
	ProductNames                             []string
	Tags                                     map[string]interface{}
}

type ReadBucketRequest struct {
	EndModifiedDt                    optional.Time
	ObjectStorageBucketIds           optional.Interface
	ObjectStorageBucketName          optional.String
	ObjectStorageBucketPurposes      optional.Interface
	ObjectStorageBucketState         optional.String
	ObjectStorageBucketStates        optional.Interface
	ObjectStorageBucketUserPurpose   optional.String
	ObjectStorageId                  optional.String
	ObjectStorageQuotaId             optional.String
	ObjectStorageSystemBucketEnabled optional.Bool
	ServiceZoneId                    optional.String
	StartModifiedDt                  optional.Time
	CreatedBy                        optional.String
	Page                             optional.Int32
	Size                             optional.Int32
	Sort                             optional.Interface
}

type ReadBucketListRequest struct {
	EndModifiedDt                    optional.Time
	ObjectStorageBucketIds           optional.Interface
	ObjectStorageBucketName          optional.String
	ObjectStorageBucketPurposes      optional.Interface
	ObjectStorageBucketState         optional.String
	ObjectStorageBucketStates        optional.Interface
	ObjectStorageBucketUserPurpose   optional.String
	ObjectStorageId                  optional.String
	ObjectStorageQuotaId             optional.String
	ObjectStorageSystemBucketEnabled optional.Bool
	ServiceZoneId                    optional.String
	StartModifiedDt                  optional.Time
	CreatedBy                        optional.String
	Page                             optional.Int32
	Size                             optional.Int32
	Sort                             optional.Interface
}

type ReadObjectStorageListRequest struct {
	IsMultiAvailabilityZone optional.Bool
	ObjectStorageName       optional.String
	Page                    optional.Int32
	Size                    optional.Int32
	Sort                    optional.Interface
}

type AccessControlRule struct {
	RuleType  string
	RuleValue string
}

type TagRequest struct {
	TagKey   string
	TagValue string
}
