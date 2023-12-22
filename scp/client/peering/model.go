package peering

type VpcPeeringListRequest struct {
	ApproverVpcId  string
	RequesterVpcId string
	VpcPeeringName string
	CreatedBy      string
	Page           int32
	Size           int32
	Sort           string
}

type VpcPeeringCreateRequest struct {
	// 수락 VPC 의 Project ID Project ID 는 @[View Project List] 를 통해 획득
	ApproverProjectId string
	// 수락 VPC ID vpcId is obtained through @[Get List of VPCs]
	ApproverVpcId string
	// 방화벽 사용 여부
	FirewallEnabled bool
	// Product group Id productGroupId is obtained through @[Get Product Groups List By Zone ID]
	ProductGroupId string
	// 신청 VPC 의 Project ID Project ID 는 @[View Project List] 를 통해 획득
	RequesterProjectId string
	// 신청 VPC ID vpcId is obtained through @[Get List of VPCs]
	RequesterVpcId string
	// VPC Peering Type
	VpcPeeringType string
	// VPC Peering Description
	VpcPeeringDescription string
	Tags                  map[string]interface{}
}

type VpcPeeringApprovalResponse struct {
	// 요청 등록 성공 여부
	Success bool
	// VPC Peering ID
	VpcPeeringId string
}

type VpcPeeringDetailResponse struct {
	// Project ID
	ProjectId string
	// 수락자
	ApprovedBy string
	// 수락일시
	ApprovedDt string
	// 수락 VPC 의 Project ID
	ApproverProjectId string
	// 수락 VPC 의 방화벽 사용 여부(승인시 firewallEnabled 값과 매치)
	ApproverVpcFirewallEnabled bool
	// 수락 VPC ID
	ApproverVpcId string
	// Block Id
	BlockId string
	// 연결완료일
	CompletedDt string
	// Product group Id
	ProductGroupId string
	// 신청자
	RequestedBy string
	// 신청일자
	RequestedDt string
	// 신청 VPC 의 Project ID
	RequesterProjectId string
	// 신청 VPC 의 방화벽 사용 여부(연결 신청시 firewallEnabled 값과 매치)
	RequesterVpcFirewallEnabled bool
	// 신청 VPC ID
	RequesterVpcId string
	// Service Zone ID
	ServiceZoneId string
	// VPC Peering ID
	VpcPeeringId string
	// VPC Peering Name
	VpcPeeringName string
	// VPC Peering 상태
	VpcPeeringState string
	// VPC Peering Type
	VpcPeeringType string
	// VPC Peering Description
	VpcPeeringDescription string
	// creator
	CreatedBy string
	// created datetime
	CreatedDt string
	// last modified user
	ModifiedBy string
	// Resource modified datetime
	ModifiedDt string
}

type VpcPeeringResponse struct {
	// 수락 VPC 의 Project ID
	ApproverProjectId string
	// 수락 VPC ID
	ApproverVpcId string
	Automated     bool
	// 연결완료일
	CompletedDt string
	// 신청 VPC 의 Project ID
	RequesterProjectId string
	// 신청 VPC ID
	RequesterVpcId string
	// VPC Peering ID
	VpcPeeringId string
	// VPC Peering Name
	VpcPeeringName string
	// VPC Peering 상태
	VpcPeeringState string
	// VPC Peering Type
	VpcPeeringType string
	// VPC Peering Description
	VpcPeeringDescription string
	// 생성자
	CreatedBy string
	// 생성일
	CreatedDt string
	// 수정자
	ModifiedBy string
	// 수정일
	ModifiedDt string
}
