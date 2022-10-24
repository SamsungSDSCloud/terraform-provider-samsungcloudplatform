package scp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	user2 "os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/ScpDevTerra/trf-provider/scp/service/firewall"
	"github.com/ScpDevTerra/trf-provider/scp/service/image"
	"github.com/ScpDevTerra/trf-provider/scp/service/internetgateway"
	"github.com/ScpDevTerra/trf-provider/scp/service/kubernetes"
	"github.com/ScpDevTerra/trf-provider/scp/service/loadbalancer"
	"github.com/ScpDevTerra/trf-provider/scp/service/natgateway"
	"github.com/ScpDevTerra/trf-provider/scp/service/postgresql"
	"github.com/ScpDevTerra/trf-provider/scp/service/product"
	"github.com/ScpDevTerra/trf-provider/scp/service/project"
	"github.com/ScpDevTerra/trf-provider/scp/service/publicip"
	"github.com/ScpDevTerra/trf-provider/scp/service/region"
	"github.com/ScpDevTerra/trf-provider/scp/service/securitygroup"
	"github.com/ScpDevTerra/trf-provider/scp/service/storage/blockstorage"
	"github.com/ScpDevTerra/trf-provider/scp/service/storage/filestorage"
	"github.com/ScpDevTerra/trf-provider/scp/service/subnet"
	"github.com/ScpDevTerra/trf-provider/scp/service/virtualserver"
	"github.com/ScpDevTerra/trf-provider/scp/service/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:               getSchema(),
		DataSourcesMap:       getDataSourcesMap(),
		ResourcesMap:         getResourcesMap(),
		ConfigureContextFunc: configureProvider,
	}
}

type authResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
}

func getAuthToken(host string, clientId string, username string, password string) (string, error) {
	httpClient := &http.Client{}

	requestBody := url.Values{}
	requestBody.Set("grant_type", "password")
	requestBody.Set("client_id", clientId)
	requestBody.Set("username", username)
	requestBody.Set("password", password)
	encodedBody := requestBody.Encode()

	req, err := http.NewRequest("POST", host+"/accounts/oidc/accessToken", strings.NewReader(encodedBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encodedBody)))
	query := req.URL.Query()
	query.Add("api", "true")
	req.URL.RawQuery = query.Encode()

	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	auth := authResponse{}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseBody, &auth)
	if err != nil {
		return "", err
	}

	return auth.IdToken, nil
}

type credentialConfig struct {
	AuthMethod string `json:"auth-method"`
	AccessKey  string `json:"access-key"`
	SecretKey  string `json:"secret-key"`
	Password   string `json:"password"`
}

type serviceConfig struct {
	Target    string `json:"target"`
	UserId    string `json:"user-id"`
	Email     string `json:"email"`
	ProjectId string `json:"project-id"`
}

type serviceInfo struct {
	host         string
	authUrl      string
	authClientId string
}

const serviceConfigFilename = "config.json"
const credentialConfigFilename = "credentials.json"

var serviceInfoMap = map[string]serviceInfo{
	"production": {
		host:         "https://openapi.samsungsdscloud.com",
		authUrl:      "https://openapi.samsungsdscloud.com",
		authClientId: "cmpClientProd",
	},
}

func loadJson(filename string, result interface{}) error {
	data, err := ioutil.ReadFile(filename)

	if err == nil {
		err = json.Unmarshal(data, result)
		if err != nil {
			return fmt.Errorf("failed to load json file %s", filename)
		}
	}

	return nil
}

func getVariable(rd *schema.ResourceData, name string, env string, getConfig func() string) string {
	res := rd.Get(name).(string)

	if res == "" {
		res = os.Getenv(env)
	}

	if res == "" {
		res = getConfig()
	}

	return res
}

func configureService(rd *schema.ResourceData, target string, service *serviceConfig, config *client.Config) error {
	info, exists := serviceInfoMap[target]
	fmt.Println("target:", target)
	if !exists {
		return fmt.Errorf("property \"target\" invalid\n")
	}

	config.ServiceHost = info.host

	config.ProjectId = getVariable(rd, "project_id", "SCP_TF_PROJECT_ID", func() string { return service.ProjectId })
	if config.ProjectId == "" {
		return fmt.Errorf("failed to get project_id configuration")
	}

	config.UserId = getVariable(rd, "user_id", "SCP_TF_USER_ID", func() string { return service.UserId })
	if config.UserId == "" {
		return fmt.Errorf("failed to get user_id configuration")
	}

	config.Email = getVariable(rd, "email", "SCP_TF_EMAIL", func() string { return service.Email })
	config.LoginId = config.Email

	if config.Email == "" {
		return fmt.Errorf("failed to get email configuration")
	}

	return nil
}

func configureCredential(rd *schema.ResourceData, target string, credential *credentialConfig, config *client.Config) error {
	config.AuthMethod = getVariable(rd, "auth_method", "SCP_TF_AUTH_METHOD", func() string { return credential.AuthMethod })
	config.Credentials.AccessKey = getVariable(rd, "access_key", "SCP_TF_ACCESS_KEY", func() string { return credential.AccessKey })
	config.Credentials.SecretKey = getVariable(rd, "secret_key", "SCP_TF_SECRET_KEY", func() string { return credential.SecretKey })

	if config.AuthMethod == "access-key" {
		return nil
	}

	if config.AuthMethod == "id-token" {
		password := getVariable(rd, "password", "SCP_TF_PASSWORD", func() string { return credential.Password })

		info, exists := serviceInfoMap[target]
		if !exists {
			return fmt.Errorf("property \"target\" invalid\n")
		}

		token, err := getAuthToken(info.authUrl, info.authClientId, config.Email, password)
		if err != nil {
			return err
		}

		config.Token = token
	}

	return fmt.Errorf("unsupported auth method")
}

func configureProvider(ctx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
	providerConfig := client.Config{}
	service := serviceConfig{}
	credential := credentialConfig{}

	user, err := user2.Current()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	err = loadJson(filepath.Join(user.HomeDir, ".scp", serviceConfigFilename), &service)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	err = loadJson(filepath.Join(user.HomeDir, ".scp", credentialConfigFilename), &credential)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	target := getVariable(rd, "target", "SCP_TF_TARGET", func() string { return service.Target })
	if target == "" {
		return nil, diag.FromErr(fmt.Errorf("property \"target\" not found in %s\n", serviceConfigFilename))
	}

	configureService(rd, target, &service, &providerConfig)
	configureCredential(rd, target, &credential, &providerConfig)

	scpClient, err := client.NewSCPClient(&providerConfig)
	if err != nil {
		log.Fatalln("Failed to create SCP EngineClient")
		return nil, diag.FromErr(err)
	}

	inst := client.Instance{
		Client: scpClient,
	}

	return &inst, nil
}

func getSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"target": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SCP target environment (development, stage, scp-maint, production)",
		},
		"user_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SCP user ID",
		},
		"email": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SCP user email",
		},
		"project_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SCP target project id",
		},
		"auth_method": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Auth method (access-key or id-token)",
		},
		"access_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SCP account access key",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SCP account secret key",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SCP account password",
		},
	}
}

func getDataSourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"scp_public_ip":                  publicip.DatasourceVpcPublicIp(),
		"scp_project":                    project.DatasourceProjects(),
		"scp_product":                    product.DatasourceProducts(),
		"scp_vpcs":                       vpc.DatasourceVpcs(),
		"scp_subnets":                    subnet.DatasourceSubnets(),
		"scp_subnet_resources":           subnet.DatasourceSubnetResources(),
		"scp_standard_images":            image.DatasourceStandardImages(),
		"scp_standard_image":             image.DatasourceStandardImage(),
		"scp_regions":                    region.DatasourceRegions(),
		"scp_region":                     region.DatasourceRegion(),
		"scp_kubernetes_apps_image":      kubernetes.DatasourceKubernetesAppsImage(),
		"scp_kubernetes_apps_images":     kubernetes.DatasourceKubernetesAppsImages(),
		"scp_firewalls":                  firewall.DatasourceFirewalls(),
		"scp_firewall":                   firewall.DatasourceFirewall(),
		"scp_public_ips":                 publicip.DatasourcePublicIps(),
		"scp_block_storages":             blockstorage.DatasourceBlockStorages(),
		"scp_file_storages":              filestorage.DatasourceFileStorages(),
		"scp_kubernetes_engines":         kubernetes.DatasourceEngines(),
		"scp_kubernetes_node_pools":      kubernetes.DatasourceNodePools(),
		"scp_kubernetes_subnet":          kubernetes.DatasourceSubnet(),
		"scp_kubernetes_engine_versions": kubernetes.DatasourceEngineVersions(),
		"scp_lb_server_groups":           loadbalancer.DatasourceLBServerGroups(),
		"scp_lb_service_ips":             loadbalancer.DatasourceLBServiceIps(),
		"scp_lb_services":                loadbalancer.DatasourceLBServices(),
		"scp_load_balancers":             loadbalancer.DatasourceLoadBalancers(),
	}
}

func getResourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"scp_vpc":                 vpc.ResourceVpc(),
		"scp_public_ip":           publicip.ResourceVpcPublicIp(),
		"scp_subnet":              subnet.ResourceSubnet(),
		"scp_security_group":      securitygroup.ResourceSecurityGroup(),
		"scp_security_group_rule": securitygroup.ResourceSecurityGroupRule(),
		"scp_internet_gateway":    internetgateway.ResourceInternetGateway(),
		"scp_nat_gateway":         natgateway.ResourceNATGateway(),
		"scp_firewall":            firewall.ResourceFirewall(),
		"scp_firewall_rule":       firewall.ResourceFirewallRule(),

		"scp_virtual_server": virtualserver.ResourceVirtualServer(),

		"scp_kubernetes_engine":    kubernetes.ResourceKubernetesEngine(),
		"scp_kubernetes_node_pool": kubernetes.ResourceKubernetesNodePool(),
		"scp_kubernetes_namespace": kubernetes.ResourceKubernetesNamespace(),
		"scp_kubernetes_apps":      kubernetes.ResourceKubernetesApps(),

		"scp_block_storage": blockstorage.ResourceBlockStorage(),
		"scp_file_storage":  filestorage.ResourceFileStorage(),

		"scp_load_balancer":   loadbalancer.ResourceLoadBalancer(),
		"scp_lb_profile":      loadbalancer.ResourceLbProfile(),
		"scp_lb_server_group": loadbalancer.ResourceLbServerGroup(),
		"scp_lb_service":      loadbalancer.ResourceLbService(),

		"scp_postgresql": postgresql.ResourcePostgresql(),
	}
}
