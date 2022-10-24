package filestorage

type CreateFileStorageRequest struct {
	FileStorageName      string
	DiskType             string
	FileStorageProtocol  string
	CifsPassword         string
	ProductGroupId       string
	ProductIds           []string
	SnapshotSchedule     *SnapshotSchedule
	RetentionCount       int32
	IsEncrypted          bool
	ServiceZoneId        string
	SnapshotCapacityRate int32
}
type SnapshotSchedule struct {
	DayOfWeek string
	Frequency string
	Hour      int32
}

type ReadFileStorageRequest struct {
	FileStorageId       string
	FileStorageName     string
	FileStorageProtocol string
	ServiceZoneId       string
	DiskType            string
	CreatedBy           string
	Page                int32
	Size                int32
	Sort                string
}

type CheckFileStorageRequest struct {
	ServiceZoneId   string
	FileStorageName string
	CifsId          string
}

type UpdateFileStorageRequest struct {
	FileStorageId         string
	FileStorageCapacityGb int32
}
