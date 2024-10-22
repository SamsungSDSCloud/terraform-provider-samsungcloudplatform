package vpc

type ListVpcRequest struct {
	ServiceZoneId string
	VpcId         string
	VpcName       string
	VpcStates     string
	CreatedBy     string
	Page          int32
	Size          int32
	Sort          string
}
