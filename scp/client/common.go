package client

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/product"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/library/project"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"time"
)

const DefaultTimeout time.Duration = 120 * time.Minute

type Instance struct {
	Client *SCPClient
}

func selectServiceZone(serviceZones []project.ProjectZoneV2, location string) *project.ProjectZoneV2 {
	var targetServiceZone project.ProjectZoneV2
	for _, serviceZone := range serviceZones {
		if serviceZone.ServiceZoneLocation == location {
			targetServiceZone = serviceZone
			break
		}
	}
	return &targetServiceZone
}

func findServiceZoneId(ctx context.Context, client *SCPClient, location string) (string, error) {
	projectInfo, err := client.Project.GetProjectInfo(ctx)
	if err != nil {
		return "", err
	}

	targetServiceZone := selectServiceZone(projectInfo.ServiceZones, location)
	if targetServiceZone == nil {
		return "", fmt.Errorf("Target service region not found")
	}

	return targetServiceZone.ServiceZoneId, nil
}

func FindProductGroupId(ctx context.Context, client *SCPClient, serviceZoneId string, productGroup string, product string) (string, error) {
	productGroups, err := client.Product.GetProductGroups(ctx, serviceZoneId, productGroup, product)
	if err != nil {
		return "", err
	}

	var productGroupId string
	for _, c := range productGroups.Contents {
		if c.TargetProduct == product {
			productGroupId = c.ProductGroupId
		}
	}

	if productGroupId == "" {
		return "", fmt.Errorf("failed to find Product group id")
	}

	return productGroupId, nil
}

func FindServiceZoneIdAndProductGroupId(ctx context.Context, client *SCPClient, location string, productGroup string, product string) (string, string, error) {

	serviceZoneId, err := findServiceZoneId(ctx, client, location)
	if err != nil {
		return "", "", err
	}

	productGroupId, err := FindProductGroupId(ctx, client, serviceZoneId, productGroup, product)
	if err != nil {
		return serviceZoneId, "", err
	}

	return serviceZoneId, productGroupId, nil
}

func FindLocationName(ctx context.Context, client *SCPClient, serviceZoneId string) (string, error) {
	projectInfo, err := client.Project.GetProjectInfo(ctx)
	if err != nil {
		return "", nil
	}
	for _, serviceZone := range projectInfo.ServiceZones {
		if serviceZone.ServiceZoneId == serviceZoneId {
			return serviceZone.ServiceZoneLocation, nil
		}
	}
	return "", nil
}

func FindProductIdByType(ctx context.Context, client *SCPClient, productGroupId string, productType string) ([]string, error) {
	productGroupInfo, err := client.Product.GetProductGroup(ctx, productGroupId)

	if err != nil {
		return []string{}, err
	}

	var result []string
	for _, s := range productGroupInfo.Products {
		for _, product := range s {
			if product.ProductType == productType {
				result = append(result, product.ProductId)
			}
		}
	}

	return result, nil
}

func FindProductGroupIdFromServiceZone(ctx context.Context, client *SCPClient, serviceZoneId string, targetProduct string, targetProductGroup string) (string, error) {
	productGroupResult, err := client.Product.GetProductGroupsByZone(ctx, serviceZoneId, targetProduct, targetProductGroup)
	if err != nil {
		return "", err
	}

	for _, pgr := range productGroupResult.Contents {
		if pgr.TargetProduct == targetProduct && pgr.TargetProductGroup == targetProductGroup {
			return pgr.ProductGroupId, nil
		}
	}

	return "", fmt.Errorf("Failed to find product group")
}

func FindProductId(ctx context.Context, client *SCPClient, productGroupId string, productType string, productName string) (string, error) {
	productGroupInfo, err := client.Product.GetProductGroup(ctx, productGroupId)

	if err != nil {
		return "", err
	}

	for _, s := range productGroupInfo.Products {
		for _, product := range s {
			if product.ProductType == productType && product.ProductName == productName {
				return product.ProductId, nil
			}
		}
	}

	return "", fmt.Errorf("Failed to find product")
}

func FindProductById(ctx context.Context, client *SCPClient, productGroupId string, productId string) (*product.ProductForCalculatorResponse, error) {
	productGroupInfo, err := client.Product.GetProductGroup(ctx, productGroupId)

	if err != nil {
		return nil, nil
	}

	for _, s := range productGroupInfo.Products {
		for _, product := range s {
			if product.ProductId == productId {
				return &product, nil
			}
		}
	}

	return nil, nil
}

func FindScaleProductByProducts(products []product.ProductForCalculatorResponse, cpuCount int, memorySizeGB int) string {
	cpuString := strconv.Itoa(cpuCount)
	memoryString := strconv.Itoa(memorySizeGB)

	for _, product := range products {
		if product.ProductType != common.ProductScale {
			continue
		}

		foundCpu := false
		foundMemory := false

		for _, item := range product.Item {
			if item.ItemType == "cpu" && item.ItemValue == cpuString {
				foundCpu = true
			} else if item.ItemType == "memory" && item.ItemValue == memoryString {
				foundMemory = true
			}
		}

		if foundCpu && foundMemory {
			return product.ProductId
		}
	}

	return ""
}

func FindScaleProduct(ctx context.Context, client *SCPClient, productGroupId string, numCpus int, memorySizeGB int) (string, error) {
	productGroupInfo, err := client.Product.GetProductGroup(ctx, productGroupId)

	if err != nil {
		return "", nil
	}

	for _, s := range productGroupInfo.Products {
		productId := FindScaleProductByProducts(s, numCpus, memorySizeGB)
		if len(productId) == 0 {
			continue
		}

		for _, product := range s {
			if product.ProductType != "SCALE" {
				continue
			}
		}

		return productId, nil
	}

	return "", nil
}

func WaitForStatus(ctx context.Context, client *SCPClient, pendingStates []string, targetStates []string, refreshFunc resource.StateRefreshFunc) error {
	stateConf := &resource.StateChangeConf{
		Pending:    pendingStates,
		Target:     targetStates,
		Refresh:    refreshFunc,
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting : %s", err)
	}

	return nil
}
