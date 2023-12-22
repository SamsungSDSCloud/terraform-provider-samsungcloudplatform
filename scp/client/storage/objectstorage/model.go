package objectstorage

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

type AccessControlRule struct {
	RuleType  string
	RuleValue string
}

type TagRequest struct {
	TagKey   string
	TagValue string
}
