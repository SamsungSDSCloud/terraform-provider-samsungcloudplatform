package objectstorage

type CreateBucketRequest struct {
	ObjectStorageBucketAccessControlEnabled  bool
	AccessControlRules                       []AccessControlRule
	ObjectStorageBucketFileEncryptionEnabled bool
	ObjectStorageBucketName                  string
	ObjectStorageBucketVersionEnabled        bool
	ObjectStorageId                          string
	ServiceZoneId                            string
	ProductNames                             []string
	Tags                                     []TagRequest
}

type AccessControlRule struct {
	RuleType  string
	RuleValue string
}

type TagRequest struct {
	TagKey   string
	TagValue string
}
