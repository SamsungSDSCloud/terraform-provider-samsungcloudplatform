package project

type ListProjectRequest struct {
	AccessLevel         string
	ActionName          string
	CmpServiceName      string
	IsUserAuthorization bool
}

type ListAccountRequest struct {
	AccessLevel         string
	ActionName          string
	CmpServiceName      string
	IsUserAuthorization bool
	MyProject           bool
}
