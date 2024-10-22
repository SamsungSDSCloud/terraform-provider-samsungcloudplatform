package iam

type TrustPrincipalsResponse struct {
	ProjectIds []string
	UserSrns   []string
}

type ListMemberRequest struct {
	CompanyName string
	Email       string
	UserName    string
}

type ListPolicyRequest struct {
	PolicyName string
	PolicyType string
}

type ListMemberGroupRequest struct {
	GroupName string
	Page      int32
	Size      int32
	Sort      string
}

type ListGroupRequest struct {
	GroupName       string
	ModifiedByEmail string
	Page            int32
	Size            int32
	Sort            string
}

type ListRoleRequest struct {
	RoleName        string
	ModifiedByEmail string
	Page            int32
	Size            int32
	Sort            string
}
