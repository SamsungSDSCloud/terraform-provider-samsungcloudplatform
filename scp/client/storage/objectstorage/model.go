package objectstorage

type CreateBucketRequest struct {
	IsObsBucketIpAddressFilterEnabled bool
	ObsBucketAccessIpAddressRanges    []ObsBucketAccessIpAddressInfo
	ObsBucketFileEncryptionAlgorithm  string
	ObsBucketFileEncryptionEnabled    bool
	ObsBucketFileEncryptionType       string
	ObsBucketName                     string
	ObsBucketVersionEnabled           bool
	ObsId                             string
	ZoneId                            string
	Tags                              []TagRequest
}

type ObsBucketAccessIpAddressInfo struct {
	ObsBucketAccessIpAddressRange string
	Type                          string
}

type TagRequest struct {
	TagKey   string
	TagValue string
}

type UpdateBucketRequest struct {
	ObsBucketFileEncryptionAlgorithm string
	ObsBucketFileEncryptionEnabled   bool
	ObsBucketFileEncryptionType      string
	ObsBucketVersionEnabled          bool
}

type ListObjectStorageOpts struct {
	MultiAzYn string
	Page      int32
	Size      int32
	Sort      []string
}

type ApiListObjectOpts struct {
	ObsObjectPath string
	Page          int32
	Size          int32
	Sort          []string
}

type ApiListBucketOpts struct {
	IsObsBucketSync          string
	IsObsSystemBucketEnabled string
	ObsBucketName            string
	ObsBucketNameExact       string
	ObsBucketQueryEndDt      string
	ObsBucketQueryStartDt    string
	ObsBucketState           string
	ObsBucketUsedType        []string
	ObsQuotaId               string
	PoolRegion               string
	ObsBucketIdList          []string
	CreatedBy                string
	Page                     int32
	Size                     int32
	Sort                     []string
}
