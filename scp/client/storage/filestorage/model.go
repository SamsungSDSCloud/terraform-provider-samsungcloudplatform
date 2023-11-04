package filestorage

type CreateFileStorageRequest struct {
	CifsPassword           string
	DiskType               string
	FileStorageName        string
	FileStorageProtocol    string
	MultiAvailabilityZone  *bool
	ProductNames           []string
	ServiceZoneId          string
	SnapshotRetentionCount *int32
	SnapshotSchedule       SnapshotSchedule
	Tags                   []TagRequest
}

type SnapshotSchedule struct {
	DayOfWeek string
	Frequency string
	Hour      *int32
}

type TagRequest struct {
	TagKey   string
	TagValue string
}

type ReadFileStorageRequest struct {
	BlockId             string
	FileStorageId       string
	FileStorageName     string
	FileStorageNameUuid string
	FileStorageProtocol string
	FileStorageState    string
	FileStorageStates   []string
	ServiceZoneId       string
	CreatedBy           string
	Page                int32
	Size                int32
	Sort                []string
}

type CheckFileStorageRequest struct {
	ServiceZoneId   string
	FileStorageName string
}

type UpdateFileStorageRequest struct {
	FileStorageId string
}

type LinkFileStorageObjectRequest struct {
	LinkObjects   []LinkObjectRequest
	UnlinkObjects []LinkObjectRequest
}

type LinkObjectRequest struct {
	LinkObjectId string
	Type         string
}
