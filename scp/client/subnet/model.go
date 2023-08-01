package subnet

type ListSubnetRequest struct {
	SubnetCidrBlock string
	SubnetId        string
	SubnetName      string
	SubnetStates    string
	SubnetTypes     string
	VpcId           string
	CreatedBy       string
	Page            int32
	Size            int32
	Sort            int32
}

type ListSubnetResourceRequest struct {
	IpAddress        string
	SubnetId         string
	LinkedObjectType string
	Page             int32
	Size             int32
	Sort             int32
}

type ListSubnetVirtualIpRequest struct {
	SubnetIpAddress string
	SubnetId        string
	VipState        string
	Page            int32
	Size            int32
	Sort            int32
}
