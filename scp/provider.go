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

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var scpResources map[string]*schema.Resource
var scpDataSources map[string]*schema.Resource

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:               getSchema(),
		DataSourcesMap:       scpDataSources,
		ResourcesMap:         scpResources,
		ConfigureContextFunc: configureProvider,
	}
}

// RegisterResource Register resources terraform for SCP
func RegisterResource(name string, resourceSchema *schema.Resource) {
	if scpResources == nil {
		scpResources = make(map[string]*schema.Resource)
	}
	scpResources[name] = resourceSchema
}

// RegisterDatasource Register datasource terraform for SCP
func RegisterDataSource(name string, DataSourceSchema *schema.Resource) {
	if scpDataSources == nil {
		scpDataSources = make(map[string]*schema.Resource)
	}
	scpDataSources[name] = DataSourceSchema
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
	Host      string `json:"host"`
	UserId    string `json:"user-id"`
	Email     string `json:"email"`
	ProjectId string `json:"project-id"`
}

const serviceConfigFilename = "config.json"
const credentialConfigFilename = "credentials.json"

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

func configureService(rd *schema.ResourceData, service *serviceConfig, config *client.Config) error {
	config.ServiceHost = getVariable(rd, "host", "SCP_TF_HOST", func() string { return service.Host })
	if config.ServiceHost == "" {
		config.ServiceHost = "https://openapi.samsungsdscloud.com" // Fallback to default host
	}

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

func configureCredential(rd *schema.ResourceData, credential *credentialConfig, config *client.Config) error {
	config.AuthMethod = getVariable(rd, "auth_method", "SCP_TF_AUTH_METHOD", func() string { return credential.AuthMethod })
	config.Credentials.AccessKey = getVariable(rd, "access_key", "SCP_TF_ACCESS_KEY", func() string { return credential.AccessKey })
	config.Credentials.SecretKey = getVariable(rd, "secret_key", "SCP_TF_SECRET_KEY", func() string { return credential.SecretKey })

	if config.AuthMethod == "access-key" {
		return nil
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

	configureService(rd, &service, &providerConfig)
	configureCredential(rd, &credential, &providerConfig)

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
		"host": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SCP Host",
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

/*
func getDataSourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"scp_vpc_peering_detail": peering.DataSourceVpcPeeringDetail(),
		"scp_vpc_peerings":       peering.DataSourceVpcPeeringList(),
		"scp_vpc_routing_rules":  routing.DataSourceVpcRoutingRule(),
		"scp_vpc_routing_routes": routing.DataSourceVpcRoutingRoute(),
		"scp_vpc_routing_tables": routing.DataSourceVpcRoutingTable(),
		"scp_vpc_dnss":           vpc.DatasourceVpcDns(),
		"scp_public_ip": publicip.DatasourceVpcPublicIp(),
		"scp_project": project.DatasourceProjects(),
		"scp_product": product.DatasourceProducts(),
		"scp_vpcs":                       vpc.DatasourceVpcs(),
		"scp_subnets":                    subnet.DatasourceSubnets(),
		"scp_subnet_resources":           subnet.DatasourceSubnetResources(),
		"scp_standard_images":            image.DatasourceStandardImages(),
		"scp_standard_image":             image.DatasourceStandardImage(),
		"scp_custom_images":              image.DatasourceCustomImages(),
		"scp_custom_image":               image.DatasourceCustomImage(),
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
		"scp_lb_server_groups": loadbalancer.DatasourceLBServerGroups(),
		"scp_lb_service_ips": loadbalancer.DatasourceLBServiceIps(),
		"scp_lb_services":    loadbalancer.DatasourceLBServices(),
		"scp_load_balancers": loadbalancer.DatasourceLoadBalancers(),
		"scp_members":        iam.DatasourceMembers(),
		"scp_policies":       iam.DatasourcePolicies(),
		"scp_groups":         iam.DatasourceGroups(),
		"scp_roles":          iam.DatasourceRoles(),
		"scp_object_storage":             objectstorage.DatasourceObjectStorage(),
	}
}
func getResourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"scp_vpc_peering_cancel":  peering.ResourceVpcPeeringCancel(),
		"scp_vpc_peering_reject":  peering.ResourceVpcPeeringReject(),
		"scp_vpc_peering_approve": peering.ResourceVpcPeeringApprove(),
		"scp_vpc_peering":         peering.ResourceVpcPeering(),
		"scp_vpc_routing": routing.ResourceVpcRouting(),
		"scp_vpc_dns":             vpc.ResourceVpcDns(),
		"scp_vpc":       vpc.ResourceVpc(),
		"scp_public_ip": publicip.ResourceVpcPublicIp(),
		"scp_subnet":              subnet.ResourceSubnet(),
		"scp_security_group":      securitygroup.ResourceSecurityGroup(),
		"scp_security_group_rule": securitygroup.ResourceSecurityGroupRule(),
		"scp_internet_gateway":    internetgateway.ResourceInternetGateway(),
		"scp_nat_gateway":   natgateway.ResourceNATGateway(),
		"scp_firewall":      firewall.ResourceFirewall(),
		"scp_firewall_rule": firewall.ResourceFirewallRule(),
		"scp_virtual_server": virtualserver.ResourceVirtualServer(),
		"scp_kubernetes_engine":    kubernetes.ResourceKubernetesEngine(),
		"scp_kubernetes_node_pool": kubernetes.ResourceKubernetesNodePool(),
		"scp_kubernetes_namespace": kubernetes.ResourceKubernetesNamespace(),
		"scp_kubernetes_apps": kubernetes.ResourceKubernetesApps(),
		"scp_block_storage": blockstorage.ResourceBlockStorage(),
		"scp_file_storage": filestorage.ResourceFileStorage(),
		"scp_object_storage": objectstorage.ResourceObjectStorage(),
		"scp_load_balancer":   loadbalancer.ResourceLoadBalancer(),
		"scp_lb_profile":      loadbalancer.ResourceLbProfile(),
		"scp_lb_server_group": loadbalancer.ResourceLbServerGroup(),
		"scp_lb_service":      loadbalancer.ResourceLbService(),
		"scp_postgresql": postgresql.ResourcePostgresql(),
	}
}
*/
