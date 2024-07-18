package client

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/autoscaling"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/baremetal"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/baremetalvdc"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/database/epas"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/database/mariadb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/database/mysql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/database/postgresql"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/database/redis"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/database/rediscluster"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/database/sqlserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/directconnect"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/dns"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/endpoint"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/gslb"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/hpclitenew"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/image"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/image/customimage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/image/migrationimage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/internetgateway"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/keypair"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/kubernetes"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/kubernetesapps"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/kubernetesengine"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/loadbalancer"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/loggingaudit"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/natgateway"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/peering"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/placementgroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/product"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/project"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/publicip"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/resourcegroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/routing"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/securitygroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/servergroup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/backup"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/blockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/bmblockstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/filestorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/storage/objectstorage"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/subnet"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/tag"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/transitgateway"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/virtualserver"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client/vpc"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type SCPClient struct {
	// Networking
	Vpc             *vpc.Client
	Routing         *routing.Client
	Peering         *peering.Client
	Subnet          *subnet.Client
	SecurityGroup   *securitygroup.Client
	LoadBalancer    *loadbalancer.Client
	InternetGateway *internetgateway.Client
	NatGateway      *natgateway.Client
	Firewall        *firewall.Client
	DirectConnect   *directconnect.Client
	TransitGateway  *transitgateway.Client
	Endpoint        *endpoint.Client
	Dns             *dns.Client
	Gslb            *gslb.Client

	// Kubernetes
	Kubernetes       *kubernetes.Client
	KubernetesEngine *kubernetesengine.Client
	KubernetesApps   *kubernetesapps.Client

	// Compute
	Image          *image.Client
	CustomImage    *customimage.Client
	MigrationImage *migrationimage.Client
	VirtualServer  *virtualserver.Client
	ServerGroup    *servergroup.Client
	BareMetal      *baremetal.Client
	BareMetalVdc   *baremetalvdc.Client
	KeyPair        *keypair.Client
	PlacementGroup *placementgroup.Client
	AutoScaling    *autoscaling.Client
	HpcLiteNew     *hpclitenew.Client

	// Storage
	FileStorage           *filestorage.Client
	BlockStorage          *blockstorage.Client
	ObjectStorage         *objectstorage.Client
	BareMetalBlockStorage *bmblockstorage.Client
	Backup                *backup.Client

	// Database
	Epas         *epas.Client
	Postgresql   *postgresql.Client
	Mariadb      *mariadb.Client
	Sqlserver    *sqlserver.Client
	Mysql        *mysql.Client
	Redis        *redis.Client
	RedisCluster *rediscluster.Client

	// Misc.
	Project       *project.Client
	Product       *product.Client
	Iam           *iam.Client
	PublicIp      *publicip.Client
	ResourceGroup *resourcegroup.Client

	Loggingaudit *loggingaudit.Client
	Tag          *tag.Client

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

	var basePath = serviceHost
	if len(servicePath) != 0 {
		basePath = serviceHost + "/" + servicePath
	}

	cfg := &scpsdk.Configuration{
		BasePath:      basePath,
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
		Routing:         routing.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Peering:         peering.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Subnet:          subnet.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		SecurityGroup:   securitygroup.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		LoadBalancer:    loadbalancer.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		InternetGateway: internetgateway.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		NatGateway:      natgateway.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Firewall:        firewall.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		DirectConnect:   directconnect.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		TransitGateway:  transitgateway.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Endpoint:        endpoint.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Dns:             dns.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Gslb:            gslb.NewClient(NewDefaultConfig(providerConfig, "oss2")),

		// Kubernetes
		Kubernetes:       kubernetes.NewClient(NewDefaultConfig(providerConfig, "kubernetes")),
		KubernetesEngine: kubernetesengine.NewClient(NewDefaultConfig(providerConfig, "kubernetes-engine2")),
		KubernetesApps:   kubernetesapps.NewClient(NewDefaultConfig(providerConfig, "kubernetes-apps")),

		// Compute
		Image:          image.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		CustomImage:    customimage.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		MigrationImage: migrationimage.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		VirtualServer:  virtualserver.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		ServerGroup:    servergroup.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		BareMetal:      baremetal.NewClient(NewDefaultConfig(providerConfig, "baremetal")),
		BareMetalVdc:   baremetalvdc.NewClient(NewDefaultConfig(providerConfig, "baremetal")),
		KeyPair:        keypair.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		PlacementGroup: placementgroup.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		AutoScaling:    autoscaling.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		HpcLiteNew:     hpclitenew.NewClient(NewDefaultConfig(providerConfig, "hpc")),

		// Storage
		FileStorage:           filestorage.NewClient(NewDefaultConfig(providerConfig, "")),
		BlockStorage:          blockstorage.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		ObjectStorage:         objectstorage.NewClient(NewDefaultConfig(providerConfig, "")),
		BareMetalBlockStorage: bmblockstorage.NewClient(NewDefaultConfig(providerConfig, "baremetal")),
		Backup:                backup.NewClient(NewDefaultConfig(providerConfig, "")),

		// Database
		Epas:         epas.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Postgresql:   postgresql.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Mariadb:      mariadb.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Sqlserver:    sqlserver.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Mysql:        mysql.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		Redis:        redis.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		RedisCluster: rediscluster.NewClient(NewDefaultConfig(providerConfig, "oss2")),

		// Common.
		Project:       project.NewClient(NewDefaultConfig(providerConfig, "project")),
		Product:       product.NewClient(NewDefaultConfig(providerConfig, "product")),
		Iam:           iam.NewClient(NewDefaultConfig(providerConfig, "iam")),
		PublicIp:      publicip.NewClient(NewDefaultConfig(providerConfig, "oss2")),
		ResourceGroup: resourcegroup.NewClient(NewDefaultConfig(providerConfig, "resource-group")),

		Loggingaudit: loggingaudit.NewClient(NewDefaultConfig(providerConfig, "logging-audit")),
		Tag:          tag.NewClient(NewDefaultConfig(providerConfig, "tag")),

		// Config
		config: providerConfig,
	}

	return client, nil
}

func (client *SCPClient) GetProjectId() string {
	return client.config.ProjectId
}
