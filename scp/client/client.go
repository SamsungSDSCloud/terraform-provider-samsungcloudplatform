package client

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/epas"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/image"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/internetgateway"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/kubernetes"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/kubernetesapps"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/kubernetesengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/loggingaudit"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/mariadb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/mysql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/natgateway"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/postgresql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/product"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/project"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/publicip"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/servergroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/sqlserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/storage/blockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/storage/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/storage/objectstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/subnet"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/tibero"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/client/vpc"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/client"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type SCPClient struct {
	// Networking
	Vpc             *vpc.Client
	Subnet          *subnet.Client
	SecurityGroup   *securitygroup.Client
	LoadBalancer    *loadbalancer.Client
	InternetGateway *internetgateway.Client
	NatGateway      *natgateway.Client
	Firewall        *firewall.Client

	// Kubernetes
	Kubernetes       *kubernetes.Client
	KubernetesEngine *kubernetesengine.Client
	KubernetesApps   *kubernetesapps.Client

	// Compute
	Image         *image.Client
	VirtualServer *virtualserver.Client
	ServerGroup   *servergroup.Client

	// Storage
	FileStorage   *filestorage.Client
	BlockStorage  *blockstorage.Client
	ObjectStorage *objectstorage.Client

	// Database
	Postgresql *postgresql.Client
	Mariadb    *mariadb.Client
	MySql      *mysql.Client
	Epas       *epas.Client
	SqlServer  *sqlserver.Client
	Tibero     *tibero.Client

	// Misc.
	Project  *project.Client
	Product  *product.Client
	Iam      *iam.Client
	PublicIp *publicip.Client

	Loggingaudit *loggingaudit.Client

	// Config
	config *Config
}

func createTlsConfig(serverHost string) *tls.Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Failed to get user home directory")
		homeDir = ""
	}
	certPath := homeDir + string(os.PathSeparator) + ".cmp" + string(os.PathSeparator) + "scp.cer"
	crt, err := ioutil.ReadFile(certPath)
	var certPool *x509.CertPool
	if err == nil {
		certPool, err = x509.SystemCertPool()
		//certPool := x509.NewCertPool()
		if err == nil {
			//certPool.AppendCertsFromPEM(crt)
			certPool.AppendCertsFromPEM(crt)
		}
	}

	return &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
		ServerName:         serverHost,
	}
}

func NewDefaultConfig(config *Config, servicePath string) *scpsdk.Configuration {
	serviceHost := config.ServiceHost

	if servicePath == "oss2" && len(config.Oss2ServiceHost) != 0 {
		serviceHost = config.Oss2ServiceHost
	}

	tlsConfig := createTlsConfig(serviceHost)

	cfg := &scpsdk.Configuration{
		BasePath:      serviceHost + "/" + servicePath,
		DefaultHeader: make(map[string]string),
		UserAgent:     "scpclient/0.0.1",
		ProjectId:     config.ProjectId,
		UserId:        config.UserId,
		Email:         config.Email,
		LoginId:       config.LoginId,
		AuthMethod:    config.AuthMethod,
		Credentials:   &config.Credentials,
		Token:         config.Token,
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
				Proxy:           nil, // Ignore host machine proxy
			},
			//Timeout: DefaultTimeout, // Default timeout
		},
	}

	return cfg
}

func NewSCPClient(providerConfig *Config) (*SCPClient, error) {
	client := &SCPClient{
		// Networking
		Vpc:             vpc.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Subnet:          subnet.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		SecurityGroup:   securitygroup.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		LoadBalancer:    loadbalancer.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		InternetGateway: internetgateway.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		NatGateway:      natgateway.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Firewall:        firewall.NewClient(NewDefaultConfig(providerConfig, "oss2")),

		// Kubernetes
		Kubernetes:       kubernetes.NewClient(NewDefaultConfig(providerConfig, "kubernetes")),
		KubernetesEngine: kubernetesengine.NewClient(NewDefaultConfig(providerConfig, "kubernetes-engine2")),
		KubernetesApps:   kubernetesapps.NewClient(NewDefaultConfig(providerConfig, "kubernetes-apps")),

		// Compute
		Image:         image.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		VirtualServer: virtualserver.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		ServerGroup:   servergroup.NewClient(NewDefaultConfig(providerConfig, "oss2")),

		// Storage
		FileStorage:   filestorage.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		BlockStorage:  blockstorage.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		ObjectStorage: objectstorage.NewClient(NewDefaultConfig(providerConfig, "object-storage")),

		// Database
		Postgresql: postgresql.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Mariadb:    mariadb.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		MySql:      mysql.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		SqlServer:  sqlserver.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Tibero:     tibero.NewClient(NewDefaultConfig(providerConfig, "oss2")),

		// Misc.
		Project:  project.NewClient(NewDefaultConfig(providerConfig, "project")),
		Product:  product.NewClient(NewDefaultConfig(providerConfig, "product")),
		Iam:      iam.NewClient(NewDefaultConfig(providerConfig, "iam")),
		PublicIp: publicip.NewClient(NewDefaultConfig(providerConfig, "oss2")),

		Loggingaudit: loggingaudit.NewClient(NewDefaultConfig(providerConfig, "logging-audit")),

		// Config
		config: providerConfig,
	}

	return client, nil
}

func (client *SCPClient) GetProjectId() string {
	return client.config.ProjectId
}
