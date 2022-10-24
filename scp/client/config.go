package client

import scpsdk "github.com/ScpDevTerra/trf-sdk/client"

type Config struct {
	ServiceHost     string
	Oss2ServiceHost string
	ProjectId       string
	Email           string
	UserId          string
	LoginId         string
	CertFilePath    string
	AuthMethod      string
	Credentials     scpsdk.Credentials
	Token           string
}
