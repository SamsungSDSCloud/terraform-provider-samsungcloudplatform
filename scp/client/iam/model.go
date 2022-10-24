package iam

type PolicyPrincipalRequest struct {
	PrincipalId   string
	PrincipalType string
}

type TrustPrincipalsResponse struct {
	ProjectIds []string
	UserSrns   []string
}
