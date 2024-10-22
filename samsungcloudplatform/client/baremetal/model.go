package baremetal

type BMServerCreateRequest struct {
	// Block ID Block ID is obtained through @[View Project Details]
	BlockId string
	// Contract Product ID Product ID is obtained through @[Get Product List by Zone ID]
	ContractId string
	// is Delete Protection enabled
	DeletionProtectionEnabled bool
	// Image ID Image ID is obtained through @[List Standard Images] or @[]
	ImageId string
	// Initial Script Content
	InitScript string
	// OS User Id
	OsUserId string
	// OS User Password
	OsUserPassword string
	// ProductGroupId Product Group ID is obtained through @[Get Product Group List by Zone ID]
	ProductGroupId string
	// Server Details
	ServerDetails []BMServerDetailsRequest
	// Service Zone ID Service Zone ID is obtained through @[View Project Details]
	ServiceZoneId string
	// Subnet ID Subnet ID is obtained through @[Get List of Subnet(V2)]
	SubnetId string
	// Time Zone
	TimeZone string
	// VPC ID VPC ID is obtained through @[Get List of VPCs]
	VpcId string
}

type BMServerDetailsRequest struct {
	// is Bare Metal Local Subnet enabled
	BareMetalLocalSubnetEnabled bool
	// Bare Metal Local Subnet ID Subnet ID is obtained through @[Get List of Subnet(V2)]
	BareMetalLocalSubnetId string
	// 1.2.10.0/24
	BareMetalLocalSubnetIpAddress string
	// Bare Metal server Name
	BareMetalServerName string
	// is DNS enabled
	DnsEnabled bool
	// IP Address
	IpAddress string
	// is NAT enabled
	NatEnabled bool
	// Public IP Address ID Reserved IP Address ID is obtained through @[Get List of Public IPs(V2)]
	PublicIpAddressId string
	// Bare Metal Server Product Type ID Product ID is obtained through @[Get Product List by Zone ID]
	ServerTypeId string
	// Storage Details
	StorageDetails []BMAdditionalBlockStorageCreateRequest
	// is HyperThreading enabled.
	UseHyperThreading string
}

type BMAdditionalBlockStorageCreateRequest struct {
	// Bare Metal Block Storage Name
	BareMetalBlockStorageName string
	// Bare Metal Block Storage Size(GB)
	BareMetalBlockStorageSize int32
	// Bare Metal Block Storage 유형
	BareMetalBlockStorageType string
	// Bare Metal Block Storage 유형 ID Product ID is obtained through @[Get Product List by Zone ID]
	BareMetalBlockStorageTypeId string
	// Encryption Use / Not Use
	EncryptionEnabled bool
}

type TagRequest struct {
	//  null
	TagKey   string
	TagValue string
}

type BMStartStopRequest struct {
	BareMetalServerIds []string
}
