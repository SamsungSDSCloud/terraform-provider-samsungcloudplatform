package publicip

type ListPublicIpRequest struct {
	ServiceZoneId   string
	IpAddress       string
	IsBillable      bool
	IsViewable      bool
	PublicIpPurpose string
	PublicIpState   string
	UplinkType      string
	CreatedBy       string
	Page            int32
	Size            int32
	Sort            string
}
