package blockstorage

type CreateBlockStorageRequest struct {
	BlockStorageName string
	BlockStorageSize int32
	EncryptEnabled   bool
	ProductId        string
	VirtualServerId  string
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
	ProductId        string
}
