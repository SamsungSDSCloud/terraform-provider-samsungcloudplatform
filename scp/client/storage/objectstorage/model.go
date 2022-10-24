package objectstorage

type ObsBucketAccessIpAddressInfo struct {
	ObsBucketAccessIpAddressRange string
	Type                          string
}

type CreateObjectStorageRequest struct {
	IsObsBucketIpAddressFilterEnabled bool
	ObsBucketAccessIpAddressRanges    []ObsBucketAccessIpAddressInfo
	ObsBucketFileEncryptionAlgorithm  string
	ObsBucketFileEncryptionEnabled    bool
	ObsBucketFileEncryptionType       string
	ObsBucketName                     string
	ObsBucketVersionEnabled           bool
	ObsId                             string
	ZoneId                            string
}

type S3BucketUpdateRequest struct {
	ObsBucketFileEncryptionAlgorithm string
	ObsBucketFileEncryptionEnabled   bool
	ObsBucketFileEncryptionType      string
	ObsBucketVersionEnabled          bool
}
