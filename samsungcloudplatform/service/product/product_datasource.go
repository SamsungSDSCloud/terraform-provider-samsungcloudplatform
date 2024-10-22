package product

import (
	"context"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/client/product"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/samsungcloudplatform/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func init() {
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_product", DatasourceProduct())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_products_by_zone", DatasourceProductsByZoneId())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_products_by_group", DatasourceProductsByGroup())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_product_categories", DatasourceProductCategories())
	samsungcloudplatform.RegisterDataSource("samsungcloudplatform_product_groups", DatasourceProductGroups())
}

func datasourceProductItemElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code":            {Type: schema.TypeString, Computed: true, Description: "Product item code (ITEM-XXXXXX)"},
			"name":            {Type: schema.TypeString, Computed: true, Description: "Product item name"},
			"type":            {Type: schema.TypeString, Computed: true, Description: "Product item type"},
			"state":           {Type: schema.TypeString, Computed: true, Description: "Product item state"},
			"serviced_for":    {Type: schema.TypeString, Computed: true, Description: "Product item serviced_for info"},
			"unit":            {Type: schema.TypeString, Computed: true, Description: "Product item unit"},
			"value":           {Type: schema.TypeString, Computed: true, Description: "Product item value"},
			"cost":            {Type: schema.TypeString, Computed: true, Description: "Product item cost"},
			"description":     {Type: schema.TypeString, Computed: true, Description: "Product item description"},
			"product_id":      {Type: schema.TypeString, Computed: true, Description: "Product id"},
			"product_id_list": {Type: schema.TypeString, Computed: true, Description: "Product id list"},
			"properties":      {Type: schema.TypeMap, Computed: true, Description: "Product item properties"},
			"version":         {Type: schema.TypeString, Computed: true, Description: "Product item version"},
		},
	}
}

func DatasourceProduct() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProductDetailRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"product_id":   {Type: schema.TypeString, Required: true, Description: "Product id"},
			"item_state":   {Type: schema.TypeString, Optional: true, Description: "Product item state"},
			"id":           {Type: schema.TypeString, Computed: true, Description: "Product id"},
			"type":         {Type: schema.TypeString, Computed: true, Description: "Product type"},
			"name":         {Type: schema.TypeString, Computed: true, Description: "Product name"},
			"description":  {Type: schema.TypeString, Computed: true, Description: "Product description"},
			"state":        {Type: schema.TypeString, Computed: true, Description: "Product state"},
			"items":        {Type: schema.TypeList, Computed: true, Description: "Product item list", Elem: datasourceProductItemElem()},
			"properties":   {Type: schema.TypeMap, Computed: true, Description: "Product properties"},
			"items_map":    {Type: schema.TypeList, Computed: true, Description: "Product items map list", Elem: &schema.Schema{Type: schema.TypeMap, Description: "item map"}},
			"items_string": {Type: schema.TypeString, Computed: true, Description: "Product items string"},
			"rate_id":      {Type: schema.TypeString, Computed: true, Description: "Project sap year for billing"},
			"seq":          {Type: schema.TypeString, Computed: true, Description: "Product display sequence"},
			"created_by":   {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_dt":   {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":  {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_dt":  {Type: schema.TypeString, Computed: true, Description: "Modified date"},
		},
	}
}

func dataSourceProductDetailRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	productId := rd.Get("product_id").(string)
	itemState := rd.Get("item_state").(string)

	info, err := inst.Client.Product.GetProductDetail(ctx, productId, itemState)
	if err != nil {
		return diag.FromErr(err)
	}

	items := common.ConvertStructToMaps(info.Items)
	itemsMap := common.ConvertStructToMaps(info.ItemsMap)

	rd.SetId(info.Id)
	rd.Set("id", info.Id)
	rd.Set("items", items)
	rd.Set("items_map", itemsMap)
	rd.Set("items_string", info.ItemsString)
	rd.Set("name", info.Name)
	rd.Set("properties", common.ToMap(info.Properties))
	rd.Set("rate_id", info.RateId)
	rd.Set("seq", info.Seq)
	rd.Set("state", info.State)
	rd.Set("type", info.Type_)
	rd.Set("description", info.Description)
	rd.Set("created_by", info.CreatedBy)
	rd.Set("created_dt", info.CreatedDt.String())
	rd.Set("modified_by", info.ModifiedBy)
	rd.Set("modified_dt", info.ModifiedDt.String())

	return nil
}

func datasourceItemForCalculatorElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"item_name":  {Type: schema.TypeString, Computed: true, Description: "Product item name"},
			"item_type":  {Type: schema.TypeString, Computed: true, Description: "Product item type (cpu, memory, ...)"},
			"item_value": {Type: schema.TypeString, Computed: true, Description: "Product item value"},
		},
	}
}

func datasourceProductForCalculatorElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"item":                {Type: schema.TypeList, Computed: true, Description: "Product item list", Elem: datasourceItemForCalculatorElem()},
			"product_id":          {Type: schema.TypeString, Computed: true, Description: "Product id"},
			"product_name":        {Type: schema.TypeString, Computed: true, Description: "Product name"},
			"product_state":       {Type: schema.TypeString, Computed: true, Description: "Product state (AVAILABLE, UNAVAILABLE)"},
			"product_type":        {Type: schema.TypeString, Computed: true, Description: "Product type (SCALE, DISK, ...)"},
			"product_description": {Type: schema.TypeString, Computed: true, Description: "Product description"},
		},
	}
}

func DatasourceProductsByGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProductsByGroupRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"product_group_id": {Type: schema.TypeString, Required: true, Description: "Product group id"},
			"products":         {Type: schema.TypeList, Computed: true, Elem: datasourceProductForCalculatorElem()},
		},
	}
}

func dataSourceProductsByGroupRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	productGroupId := rd.Get("product_group_id").(string)

	info, err := inst.Client.Product.GetProductsByGroup(ctx, productGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	productResponses := common.HclListObject{}
	for _, productValue := range info.Products {
		for _, productResponse := range productValue {
			kv := common.HclKeyValueObject{}

			kv["product_id"] = productResponse.ProductId
			kv["product_name"] = productResponse.ProductName
			kv["product_state"] = productResponse.ProductState
			kv["product_type"] = productResponse.ProductType
			kv["product_description"] = productResponse.ProductDescription

			items := common.HclListObject{}
			for _, itemResponse := range productResponse.Item {
				item := common.HclKeyValueObject{}
				item["item_name"] = itemResponse.ItemName
				item["item_type"] = itemResponse.ItemType
				item["item_value"] = itemResponse.ItemValue
				items = append(items, item)
			}
			kv["item"] = items

			productResponses = append(productResponses, kv)
		}
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("products", productResponses)

	return nil
}

func datasourceProductsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_id":   {Type: schema.TypeString, Computed: true, Description: "Product id (PRODUCT-XXXXXXX)"},
			"product_name": {Type: schema.TypeString, Computed: true, Description: "Product name"},
			"product_type": {Type: schema.TypeString, Computed: true, Description: "Product type (SCALE, DISK, ...)"},
			"created_by":   {Type: schema.TypeString, Computed: true, Description: "Creator's ID"},
			"created_dt":   {Type: schema.TypeString, Computed: true, Description: "Created date"},
			"modified_by":  {Type: schema.TypeString, Computed: true, Description: "Modifier's ID"},
			"modified_dt":  {Type: schema.TypeString, Computed: true, Description: "Modified date"},
		},
	}
}

func DatasourceProductsByZoneId() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProductsByZoneIdRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"service_zone_id":  {Type: schema.TypeString, Required: true, Description: "Service zone id"},
			"product_group_id": {Type: schema.TypeString, Optional: true, Description: "Product group id"},
			"product_type":     {Type: schema.TypeString, Optional: true, Description: "Product type (SCALE, DISK, ...)"},
			"contents":         {Type: schema.TypeList, Computed: true, Description: "Contents", Elem: datasourceProductsElem()},
			"total_count":      {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceProductsByZoneIdRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	serviceZoneId := rd.Get("service_zone_id").(string)
	productGroupId := rd.Get("product_group_id").(string)
	productType := rd.Get("product_type").(string)

	responses, err := inst.Client.Product.GetProductsByZoneId(ctx, serviceZoneId, productGroupId, productType)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func DatasourceProductCategories() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceProductCategoriesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"filter":         common.DatasourceFilter(),
			"category_id":    {Type: schema.TypeString, Optional: true, Description: "Product category id"},
			"category_state": {Type: schema.TypeString, Optional: true, Description: "Product category status"},
			"exposure_scope": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ADMIN", "CONSOLE", "LANDING"}, false)),
				Description:      "Exposure scope",
			},
			"language_code": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ko_KR", "en_US"}, false)),
				Description:      "Language code (ko_KR, en_US)",
			},
			"product_id":    {Type: schema.TypeString, Optional: true, Description: "Product id"},
			"product_state": {Type: schema.TypeString, Optional: true, Description: "Product status"},
			"contents":      {Type: schema.TypeList, Optional: true, Description: "Product list", Elem: datasourceProductElem()},
			"total_count":   {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of products.",
	}
}

func datasourceProductCategoriesRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	requestParam := product.ListCategoriesRequest{
		CategoryId:    rd.Get("category_id").(string),
		CategoryState: rd.Get("category_state").(string),
		ExposureScope: rd.Get("exposure_scope").(string),
		ProductId:     rd.Get("product_id").(string),
		ProductState:  rd.Get("product_state").(string),
		LanguageCode:  rd.Get("language_code").(string),
	}

	responses, err := inst.Client.Product.GetCategoryList(ctx, requestParam)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)
	for _, content := range contents {
		delete(content, "product")
	}

	if f, ok := rd.GetOk("filter"); ok {
		contents = common.ApplyFilter(DatasourceProductCategories().Schema, f.(*schema.Set), contents)
	}

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", len(contents))

	return nil
}

func datasourceProductElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"icon_file_name": {Type: schema.TypeString, Computed: true, Description: "Icon file name"},
			//"product":                      {Type: schema.TypeList, Computed: true, Elem: productElem()}, //TODO:
			"product_category_id":          {Type: schema.TypeString, Computed: true, Description: "Product category id"},
			"product_category_name":        {Type: schema.TypeString, Computed: true, Description: "Product category name"},
			"product_category_path":        {Type: schema.TypeString, Computed: true, Description: "Product category path"},
			"product_category_state":       {Type: schema.TypeString, Computed: true, Description: "Product category status"},
			"product_set":                  {Type: schema.TypeString, Computed: true, Description: "Product set type (SE, PAAS)"},
			"product_category_description": {Type: schema.TypeString, Computed: true, Description: "Description of product category"},
		},
	}
}

func DatasourceMenus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMenuList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			common.ToSnakeCase("CategoryId"):    {Type: schema.TypeString, Optional: true, Description: "Category id"},
			common.ToSnakeCase("ExposureType"):  {Type: schema.TypeString, Optional: true, Description: "Category status"},
			common.ToSnakeCase("ExposureScope"): {Type: schema.TypeString, Optional: true, Description: "Exposure scope"},
			common.ToSnakeCase("ProductId"):     {Type: schema.TypeString, Optional: true, Description: "Product id"},
			common.ToSnakeCase("ZoneIds"):       {Type: schema.TypeString, Optional: true, Description: "Service Zone id list"},
			"contents":                          {Type: schema.TypeList, Computed: true, Description: "Contents", Elem: datasourceProductElem()},
			"total_count":                       {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceMenuList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	requestParam := product.ListMenusRequest{
		CategoryId:    rd.Get(common.ToSnakeCase("CategoryId")).(string),
		ExposureType:  rd.Get(common.ToSnakeCase("ExposureType")).(string),
		ExposureScope: rd.Get(common.ToSnakeCase("ExposureScope")).(string),
		ProductId:     rd.Get(common.ToSnakeCase("ProductId")).(string),
		ZoneIds:       rd.Get(common.ToSnakeCase("ZoneIds")).(string),
	}

	responses, err := inst.Client.Product.GetMenuList(ctx, requestParam)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)
	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceProductSubGroupElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			//"product_group":          {Type: schema.TypeList, Computed: true, Elem: datasourceProductGroupElem()},
			"product_group_id":       {Type: schema.TypeString, Computed: true, Description: "Product group id"},
			"product_group_name":     {Type: schema.TypeString, Computed: true, Description: "Product group name"},
			"product_group_sequence": {Type: schema.TypeString, Computed: true, Description: "Product category sequence"},
			"product_group_type":     {Type: schema.TypeString, Computed: true, Description: "Product category type"},
			"target_product":         {Type: schema.TypeString, Computed: true, Description: "Target product (Kubernetes Apps, Virtual Server, MySQL, ...)"},
			"target_product_group":   {Type: schema.TypeString, Computed: true, Description: "Target product group name (CONTAINER, DATABASE, STORAGE, ...)"},
		},
	}
}

func datasourceProductGroupElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_group":          {Type: schema.TypeList, Computed: true, Elem: datasourceProductSubGroupElem()},
			"product_group_id":       {Type: schema.TypeString, Computed: true, Description: "Product group id"},
			"product_group_name":     {Type: schema.TypeString, Computed: true, Description: "Product group name"},
			"product_group_sequence": {Type: schema.TypeString, Computed: true, Description: "Product category sequence"},
			"product_group_type":     {Type: schema.TypeString, Computed: true, Description: "Product category type"},
			"target_product":         {Type: schema.TypeString, Computed: true, Description: "Target product (Kubernetes Apps, Virtual Server, MySQL, ...)"},
			"target_product_group":   {Type: schema.TypeString, Computed: true, Description: "Target product group name (CONTAINER, DATABASE, STORAGE, ...)"},
		},
	}
}

func DatasourceProductGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"target_product":       {Type: schema.TypeString, Optional: true, Description: "Target product"},
			"target_product_group": {Type: schema.TypeString, Optional: true, Description: "Target product group name"},
			"contents":             {Type: schema.TypeList, Computed: true, Description: "Contents", Elem: datasourceProductGroupElem()},
			"total_count":          {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSourceGroupList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	inst := meta.(*client.Instance)

	targetProduct := rd.Get("target_product").(string)
	targetProductGroup := rd.Get("target_product_group").(string)

	responses, err := inst.Client.Product.GetProductGroupsList(ctx, targetProduct, targetProductGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	contents := common.ConvertStructToMaps(responses.Contents)

	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

//TODO:
/*
func productElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"menu":                                  {Type: schema.TypeList, Computed: true, Elem: menuElem()},
			"new_badge_end_date":                    {Type: schema.TypeString, Computed: true},
			"product_offering_console_path":         {Type: schema.TypeString, Computed: true},
			"product_offering_console_request_path": {Type: schema.TypeString, Computed: true},
			"product_offering_detail_info":          {Type: schema.TypeString, Computed: true},
			"product_offering_id":                   {Type: schema.TypeString, Computed: true},
			"product_offering_name":                 {Type: schema.TypeString, Computed: true},
			"product_offering_path":                 {Type: schema.TypeString, Computed: true},
			"product_offering_state":                {Type: schema.TypeString, Computed: true},
			"visible":                               {Type: schema.TypeBool, Computed: true},
			"product_offering_description":          {Type: schema.TypeString, Computed: true},
		},
	}
}

func menuElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"available_button":          {Type: schema.TypeBool, Computed: true},
			"button_name":               {Type: schema.TypeString, Computed: true},
			"default_menu_id":           {Type: schema.TypeString, Computed: true},
			"menu_console_path":         {Type: schema.TypeString, Computed: true},
			"menu_console_request_path": {Type: schema.TypeString, Computed: true},
			"menu_details":              {Type: schema.TypeList, Computed: true, Elem: menuDetailsElem()},
			"menu_exposure_for":         {Type: schema.TypeString, Computed: true},
			"menu_id":                   {Type: schema.TypeString, Computed: true},
			"menu_information_path":     {Type: schema.TypeString, Computed: true},
			"menu_name":                 {Type: schema.TypeString, Computed: true},
			"menu_prior_products":       {Type: schema.TypeList, Computed: true, Elem: menuPriorProductsElem()},
			"menu_resource_type":        {Type: schema.TypeString, Computed: true},
			"menu_state":                {Type: schema.TypeString, Computed: true},
			"menu_type":                 {Type: schema.TypeString, Computed: true},
			"parent_id":                 {Type: schema.TypeString, Computed: true},
		},
	}
}

func menuDetailsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"menu_detail_guide_file_path": {Type: schema.TypeString, Computed: true},
			"menu_detail_guide_path":      {Type: schema.TypeString, Computed: true},
			"menu_detail_id":              {Type: schema.TypeString, Computed: true},
			"menu_detail_type":            {Type: schema.TypeString, Computed: true},
			"menu_detail_description":     {Type: schema.TypeString, Computed: true},
		},
	}
}

func menuPriorProductsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"prior_product_console_path":   {Type: schema.TypeString, Computed: true},
			"prior_product_icon_file_path": {Type: schema.TypeString, Computed: true},
			"prior_product_id":             {Type: schema.TypeString, Computed: true},
			"prior_product_name":           {Type: schema.TypeString, Computed: true},
			"prior_product_offering_id":    {Type: schema.TypeString, Computed: true},
			"prior_product_parent_id":      {Type: schema.TypeString, Computed: true},
			"prior_product_request_path":   {Type: schema.TypeString, Computed: true},
			"prior_product_seq":            {Type: schema.TypeString, Computed: true},
			"prior_product_state":          {Type: schema.TypeString, Computed: true},
			"prior_product_description":    {Type: schema.TypeString, Computed: true},
		},
	}
}
*/
