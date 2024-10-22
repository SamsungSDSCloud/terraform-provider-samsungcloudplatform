package project

type ListProjectRequest struct {
	AccountName          string
	BillYearMonth        string
	IsBillingInfoDemand  bool
	IsResourceInfoDemand bool
	IsUserInfoDemand     bool
	ProjectName          string
	CreatedByEmail       string
}

type ListAccountRequest struct {
	AccessLevel         string
	ActionName          string
	CmpServiceName      string
	IsUserAuthorization bool
	MyProject           bool
}
