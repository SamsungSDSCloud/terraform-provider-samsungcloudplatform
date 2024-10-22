package baremetalvdc

import "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client/baremetal"

type BMVDCServerCreateRequest struct {
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
	// Subnet ID
	SubnetId string
	// VDC ID
	VdcId string
}

type BMServerDetailsRequest struct {
	// Bare Metal server Name
	BareMetalServerName string
	// is DNS enabled
	DnsEnabled bool
	// IP Address
	IpAddress string
	// Bare Metal Server Product Type ID Product ID is obtained through @[Get Product List by Zone ID]
	ServerTypeId string
	// Storage Details
	StorageDetails []baremetal.BMAdditionalBlockStorageCreateRequest
	// is HyperThreading enabled.
	UseHyperThreading string
}
