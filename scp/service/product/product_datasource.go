package product

import (
	"context"

	"github.com/ScpDevTerra/trf-provider/scp/client"
	"github.com/ScpDevTerra/trf-provider/scp/client/product"
	"github.com/ScpDevTerra/trf-provider/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func DatasourceProducts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"category_id":    {Type: schema.TypeString, Optional: true, Description: "Product category id"},
			"category_state": {Type: schema.TypeString, Optional: true, Description: "Product category status"},
			"exposure_scope": {Type: schema.TypeString, Optional: true, Description: "Exposure scope"},
			"product_id":     {Type: schema.TypeString, Optional: true, Description: "Product id"},
			"product_state":  {Type: schema.TypeString, Optional: true, Description: "Product status"},
			"language_code":  {Type: schema.TypeString, Optional: true, Default: "ko", Description: "Language code (ko, en)"},
			"contents":       {Type: schema.TypeList, Optional: true, Description: "Product list", Elem: datasourceElem()},
			"total_count":    {Type: schema.TypeInt, Computed: true, Description: "Total list size"},
		},
		Description: "Provides list of products.",
	}
}

func dataSourceList(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	rd.SetId(uuid.NewV4().String())
	rd.Set("contents", contents)
	rd.Set("total_count", responses.TotalCount)

	return nil
}

func datasourceElem() *schema.Resource {
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
			"contents":                          {Type: schema.TypeList, Computed: true, Description: "Contents", Elem: datasourceElem()},
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
