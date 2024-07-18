package hpclitenew

type HpcLiteNewCreateRequest struct {
	// HPC Lite(New) CO Pool ID
	CoServiceZoneId string
	// HPC Lite(New) Contract
	Contract string
	// HPC Lite(New) HT Enabled
	HyperThreadingEnabled string
	// HPC Lite(New) Image ID
	ImageId string
	// HPC Lite(New) Init Script
	InitScript string
	// HPC Lite(New) OS User ID
	OsUserId string
	// HPC Lite(New) OS User PWD
	OsUserPassword string
	// HPC Lite(New) Product Group ID
	ProductGroupId string
	// HPC Lite(New) block Id
	ResourcePoolId string
	ServerDetails  []ServerDetailRequest
	// HPC Lite(New) Server Type
	ServerType string
	// HPC Lite(New) Service Zone ID
	ServiceZoneId string
	Tags          map[string]interface{}
	// HPC Lite(New) Vlan Pool CIDR
	VlanPoolCidr string
}

type ServerDetailRequest struct {
	// HPC Lite(New) Server Detail ipAddress
	IpAddress string
	// HPC Lite(New) Server Detail Name
	ServerName string
}

type HpcLiteNewDeleteRequest struct {
	ServerIds     []string
	ServiceZoneId string
}
