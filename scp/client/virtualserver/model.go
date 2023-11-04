package virtualserver

type InitialScriptInfo struct {
	// Initial Script encoding type
	EncodingType string
	// Initial Script content
	InitialScriptContent string
	// Initial Script shell type (Bash or Powershell)
	InitialScriptShell string
	// Initial Script type
	InitialScriptType string
}

type BlockStorageInfo struct {
	// Block Storage name
	BlockStorageName string
	// Block Storage size(GB)
	DiskSize int32
	// 암호화 여부
	EncryptEnabled bool
	//DiskType
	DiskType string
}

type LocalSubnetInfo struct {
	// 로컬 서브넷 IP 지정을 위한 IP 주소
	LocalSubnetIpAddress string
	// Subnet ID subnetId is obtained through @[Get List of VPC Subnet]
	SubnetId string
}

type NicInfo struct {
	// IP address for Internal IP
	InternalIpAddress string
	// Is NAT IP used
	NatEnabled bool
	// IP Public ID for NAT IP publicIpAddressId is obtained through @[Get List of Public IPs]
	PublicIpAddressId string
	// Subnet ID subnetId is obtained through @[Get List of VPC Subnet]
	SubnetId string
}

type OsAdminInfo struct {
	// OS accoount(Linux:root (fixed), Windows: administrator (or other accounts))
	OsUserId string
	// OS account password (length between 8~20, consist of lowercase letter, symbol and number, and at least 2 combination of category (like {lowercase letter, symbol} or {lowercase letter, number}))
	OsUserPassword string
}

type CreateRequest struct {
	// Information of Block Storage for default Volume(OS)
	BlockStorage BlockStorageInfo
	// Contract Prduct ID productId is obtained through @[Get Product List By Zone ID]
	ContractDiscount string
	// is Delete Protection enabled
	DeletionProtectionEnabled bool
	// Is DNS used
	DnsEnabled bool
	// Information of Block Storage for extra Volume
	ExtraBlockStorages []BlockStorageInfo
	// Image ID imageId is obtained through @[List Standard Images]
	ImageId string
	// Initial Script information
	InitialScript InitialScriptInfo
	// 생성할 vm의 로컬 서브넷 정보
	LocalSubnet LocalSubnetInfo
	// NIC(Network Interface Card) Information
	Nic NicInfo
	// OS Admin account information
	OsAdmin OsAdminInfo
	// Security Group ID list securityGroupId is obtained through @[Get List of Security Groups]
	SecurityGroupIds []string
	// Server Group ID serverGroupId is obtained through @[Get List of Server Group]
	ServerGroupId string
	// Server Type ID productId is obtained through @[Get Product List By Zone ID]
	ServerType string
	// Service Level Product ID productId is obtained through @[Get Product List By Zone ID]
	// Zone ID serviceZoneId is obtained through @[View Project Details]
	ServiceZoneId string
	// Time zone
	Timezone string
	// Virtual Server name
	VirtualServerName    string
	AvailabilityZoneName string
	Tags                 []TagRequest
	KeyPairId            string
	PlacementGroupId     string
}

type ListVirtualServersRequestParam struct {
	AutoscalingEnabled   bool
	ServerGroupId        string
	ServicedForList      []string
	ServicedGroupForList []string
	VirtualServerName    string
	Page                 int32
	Size                 int32
	Sort                 string
}

type VirtualServerSubnetIpUpdateRequest struct {
	InternalIpAddress string
	SubnetId          string
}

type VirtualServerContractUpdateRequest struct {
	ContractDiscount string
}

type ExternalStorage struct {
	BlockStorageName string
	BlockStorageId   string
	StorageSizeGb    int32
	SharedType       string
	EncryptEnabled   bool
	ProductName      string
	ProductId        string
	Tags             []TagRequest
}

type TagRequest struct {
	TagKey   string
	TagValue string
}
