package client

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/scp/common"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/product"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/project"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v2/library/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strconv"
	"time"
)

const DefaultTimeout time.Duration = 120 * time.Minute

type Instance struct {
	Client *SCPClient
}

func selectServiceZone(serviceZones []project.ZoneResponseV3, location string) *project.ZoneResponseV3 {
	var targetServiceZone project.ZoneResponseV3
	for _, serviceZone := range serviceZones {
		if serviceZone.ServiceZoneLocation == location {
			targetServiceZone = serviceZone
			break
		}
	}
	return &targetServiceZone
}

func FindServiceZoneId(ctx context.Context, client *SCPClient, location string) (string, error) {
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

	serviceZoneId, err := FindServiceZoneId(ctx, client, location)
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

func FindProductIdByType(ctx context.Context, client *SCPClient, productGroupId string, productType string, productName string) ([]string, error) {
	productGroupInfo, err := client.Product.GetProductGroup(ctx, productGroupId)

	if err != nil {
		return []string{}, err
	}

	var result []string
	for _, s := range productGroupInfo.Products {
		for _, product := range s {
			if product.ProductType == productType && product.ProductName == productName {
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

func FindScaleInfo(ctx context.Context, client *SCPClient, productGroupId string, scaleProductId string) (int, int, error) {

	scale, err := FindProductById(ctx, client, productGroupId, scaleProductId)
	if err != nil {
		return -1, -1, fmt.Errorf("failed to find scale product.")
	}

	cpuCount := -1
	memorySize := -1

	for _, item := range scale.Item {
		if item.ItemType == "cpu" {
			cpuCount, err = strconv.Atoi(item.ItemValue)

			if err != nil {
				return -1, -1, fmt.Errorf("wrong cpu count.")
			}
		} else if item.ItemType == "memory" {
			memorySize, err = strconv.Atoi(item.ItemValue)

			if err != nil {
				return -1, -1, fmt.Errorf("wrong memory size.")
				continue
			}
		}
	}

	return cpuCount, memorySize, nil
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

func UpdateResourceTag(ctx context.Context, client *SCPClient, resourceId string, oldTagsInterface []interface{}, newTagsInterface []interface{}) error {
	oldTags := toTagRequestList(oldTagsInterface)
	newTags := toTagRequestList(newTagsInterface)

	oldSet := make(map[tag.TagRequest]struct{})
	newSet := make(map[tag.TagRequest]struct{})

	for _, item := range oldTags {
		oldSet[item] = struct{}{}
	}

	for _, item := range newTags {
		newSet[item] = struct{}{}
	}

	var ok bool
	var addList, removeList []tag.TagRequest

	for _, item := range oldTags {
		if _, ok = newSet[item]; !ok {
			removeList = append(removeList, item)
		}
	}

	for _, item := range newTags {
		if _, ok = oldSet[item]; !ok {
			addList = append(addList, item)
		}
	}

	if len(removeList) > 0 {
		for _, tag := range removeList {
			_, err := client.Tag.DetachResourceTag(ctx, resourceId, tag.TagKey)
			if err != nil {
				fmt.Printf("failed to remove tag %s in resource %s", resourceId, tag.TagKey)
				return err
			}
		}
	}

	if len(addList) > 0 {
		_, _, err := client.Tag.AttachResourceTag(ctx, resourceId, addList)
		if err != nil {
			fmt.Printf("failed to add or update tags in resource %s", resourceId)
			return err
		}
	}

	return nil
}

func toTagRequestList(list []interface{}) []tag.TagRequest {
	if len(list) == 0 {
		return nil
	}
	var result []tag.TagRequest

	for _, val := range list {
		kv := val.(common.HclKeyValueObject)
		result = append(result, tag.TagRequest{
			TagKey:   kv["tag_key"].(string),
			TagValue: kv["tag_value"].(string),
		})
	}
	return result
}
