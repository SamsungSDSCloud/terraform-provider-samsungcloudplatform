package bmblockstorage

type BmBlockStorageCreateRequest struct {
	// Bare Metal Block Storage Name
	BareMetalBlockStorageName string
	// Bare Metal Block Storage Size(GB)
	BareMetalBlockStorageSize int32
	// List of Bare Metal Server ID linked Bare Metal Server ID is obtained through @[Get the List of Bare Metal Servers]
	BareMetalServerIds []string
	// Encryption Use / Not Use
	EncryptionEnabled bool
	// snapshot enable policy
	IsSnapshotPolicy bool
	// Disk Type ID Product ID is obtained through @[Get Product List by Zone ID]
	ProductId string
	// Service Zone ID Service Zone ID is obtained through @[View Project Details]
	ServiceZoneId string
	// Snapshot Capacity rate
	SnapshotCapacityRate int32
	// Snapshot Schedule
	SnapshotSchedule SnapshotSchedule
	Tags             map[string]interface{}
}

type SnapshotSchedule struct {
	// \"Schedule Week Name\"
	DayOfWeek string `json:"dayOfWeek,omitempty"`
	// \"Schedule Period\"
	Frequency string `json:"frequency,omitempty"`
	// \"Hour\"
	Hour int32 `json:"hour,omitempty"`
}

type TagRequest struct {
	//  null
	TagKey   string `json:"tagKey"`
	TagValue string `json:"tagValue,omitempty"`
}
