package product

import (
	"context"

	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/product"
	"github.com/antihax/optional"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *product.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: product.NewAPIClient(config),
	}
}

func (client *Client) GetProductGroups(ctx context.Context, serviceZoneId string, targetProductGroup string, targetProduct string) (product.ListResponseV2OfProductGroupsResponse, error) {
	var param product.ProductV2ControllerApiListProductGroupsByZoneIdOpts
	if len(targetProductGroup) != 0 {
		param.TargetProductGroup = optional.NewString(targetProductGroup)
	}
	if len(targetProduct) != 0 {
		param.TargetProduct = optional.NewString(targetProduct)
	}
	result, _, err := client.sdkClient.ProductV2ControllerApi.ListProductGroupsByZoneId(ctx, client.config.ProjectId, serviceZoneId, &param)
	return result, err
}

func (client *Client) GetProductGroup(ctx context.Context, productGroupId string) (product.ProductGroupDetailResponse, error) {
	result, _, err := client.sdkClient.ProductV2ControllerApi.DetailProductGroup(ctx, client.config.ProjectId, productGroupId,
		&product.ProductV2ControllerApiDetailProductGroupOpts{})
	return result, err
}

func (client *Client) GetProducesList(ctx context.Context, serviceZoneId string, productGroupId string, productType string) (product.ListResponseV2OfProductsResponse, error) {
	result, _, err := client.sdkClient.ProductV2ControllerApi.ListProducsByZoneId(ctx, client.config.ProjectId, serviceZoneId, &product.ProductV2ControllerApiListProducsByZoneIdOpts{
		ProductGroupId: optional.NewString(productGroupId),
		ProductType:    optional.NewString(productType),
	})
	return result, err
}

func (client *Client) GetProductGroupsByZone(ctx context.Context, serviceZoneId string, targetProduct string, targetProductGroup string) (product.ListResponseV2OfProductGroupsResponse, error) {
	result, _, err := client.sdkClient.ProductV2ControllerApi.ListProductGroupsByZoneId(ctx, client.config.ProjectId, serviceZoneId, &product.ProductV2ControllerApiListProductGroupsByZoneIdOpts{
		TargetProduct:      optional.NewString(targetProduct),
		TargetProductGroup: optional.NewString(targetProductGroup),
	})
	return result, err
}

func (client *Client) GetProductsByZoneId(ctx context.Context, serviceZoneId string, productGroupId string, productType string) (product.ListResponseV2OfProductsResponse, error) {
	result, _, err := client.sdkClient.ProductV2ControllerApi.ListProducsByZoneId(ctx, client.config.ProjectId, serviceZoneId, &product.ProductV2ControllerApiListProducsByZoneIdOpts{
		ProductGroupId: optional.NewString(productGroupId),
		ProductType:    optional.NewString(productType),
	})
	return result, err
}

func (client *Client) GetProductsByGroup(ctx context.Context, productGroupId string) (product.ProductGroupDetailResponse, error) {
	result, _, err := client.sdkClient.ProductV2ControllerApi.DetailProductGroup(ctx, client.config.ProjectId, productGroupId, &product.ProductV2ControllerApiDetailProductGroupOpts{
		ForCalculator: optional.NewString("false"),
	})
	return result, err
}

func (client *Client) GetProductGroupsList(ctx context.Context, targetProduct string, targetProductGroup string) (product.ListResponseV2OfProductGroupForCalculatorResponse, error) {
	result, _, err := client.sdkClient.ProductV2ControllerApi.ListProductGroups(ctx, client.config.ProjectId, &product.ProductV2ControllerApiListProductGroupsOpts{
		ForCalculator:      optional.NewString("false"),
		ProductGroupType:   optional.NewString("SSC"),
		TargetProduct:      optional.NewString(targetProduct),
		TargetProductGroup: optional.NewString(targetProductGroup),
	})
	return result, err
}

func (client *Client) GetProductDetail(ctx context.Context, productId string, itemState string) (product.ProductResponse, error) {
	result, _, err := client.sdkClient.ProductV2ControllerApi.DetailProduct(ctx, productId, &product.ProductV2ControllerApiDetailProductOpts{
		ItemState: optional.NewString(itemState),
	})
	return result, err
}

func (client *Client) GetCategoryList(ctx context.Context, request ListCategoriesRequest) (product.ListResponseV2OfProductCategoryResponse, error) {
	result, _, err := client.sdkClient.ProductV2ControllerApi.ListCategories(ctx, client.config.ProjectId, request.LanguageCode, &product.ProductV2ControllerApiListCategoriesOpts{
		CategoryId:    optional.NewString(request.CategoryId),
		CategoryState: optional.NewString(request.CategoryState),
		ExposureScope: optional.NewString(request.ExposureScope),
		ProductId:     optional.NewString(request.ProductId),
		ProductState:  optional.NewString(request.ProductState),
	})
	return result, err
}

func (client *Client) GetMenuList(ctx context.Context, request ListMenusRequest) (product.ListResponseV2OfProductCategoryResponse, error) {
	result, _, err := client.sdkClient.ProductV2ControllerApi.ListMenus(ctx, &product.ProductV2ControllerApiListMenusOpts{
		CategoryId:   optional.NewString(request.CategoryId),
		ExposureType: optional.NewString(request.ExposureType),
		ProductId:    optional.NewString(request.ProductId),
		ZoneIds:      optional.NewString(request.ZoneIds),
	})
	return result, err
}

func (client *Client) GetProductsUsingGET(ctx context.Context, productId string) (product.ProductResponse, error) {
	result, _, err := client.sdkClient.ProductV2ControllerApi.DetailProduct(ctx, productId, nil)
	return result, err
}
