package blockstorage

type CreateBlockStorageRequest struct {
	BlockStorageName string
	BlockStorageSize int32
	EncryptEnabled   bool
	DiskType         string
	SharedType       string
	Tags             []TagRequest
	VirtualServerId  string
}

type TagRequest struct {
	TagKey   string
	TagValue string
}

type ReadBlockStorageRequest struct {
	BlockStorageName  string
	VirtualServerId   string
	VirtualServerName string
	CreatedBy         string
	Page              int32
	Size              int32
	Sort              string
}

type UpdateBlockStorageRequest struct {
	BlockStorageId   string
	BlockStorageSize int32
}

type BlockStorageAttachRequest struct {
	VirtualServerId string
}

type BlockStorageDetachRequest struct {
	VirtualServerId string
}

type BlockStorageVirtualServersRequest struct {
	Page int32
	Size int32
	Sort []string
}
