package kubernetes

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/kubernetes"
)

type Client struct {
	config *sdk.Configuration
	sdk    *kubernetes.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config: config,
		sdk:    kubernetes.NewAPIClient(config),
	}
}

func (client *Client) CreateNamespace(ctx context.Context, clusterId string, name string) (kubernetes.K8sObjectResponse, int, error) {
	result, response, err := client.sdk.K8sObjectYamlServiceApi.CreateK8sObjectV2(ctx, client.config.ProjectId, clusterId, kubernetes.K8sObjectCreateRequest{
		Yaml: "apiVersion: v1\nkind: Namespace\nmetadata:\n  name: " + name,
	})

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) ReadNamespace(ctx context.Context, clusterId string, name string) (kubernetes.NamespaceResponse, int, error) {
	result, response, err := client.sdk.NamespaceServiceApi.DetailNamespaceV2(ctx, client.config.ProjectId, clusterId, name)
	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteNamespace(ctx context.Context, clusterId string, name string) (int, error) {
	response, err := client.sdk.NamespaceServiceApi.DeleteNamespaceV2(ctx, client.config.ProjectId, clusterId, kubernetes.K8sObjectDeleteRequest{
		Names: []string{name},
	})
	return response.StatusCode, err
}
