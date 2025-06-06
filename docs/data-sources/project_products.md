---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_project_products Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Provides list of products in given project
---

# samsungcloudplatform_project_products (Data Source)

Provides list of products in given project

## Example Usage

```terraform
data "samsungcloudplatform_project" "my_project" {

}

data "samsungcloudplatform_project_products" "my_scp_products" {
  project_id = data.samsungcloudplatform_project.my_project.id
}

output "samsungcloudplatform_project_products" {
  value = data.samsungcloudplatform_project_products.my_scp_products
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project_id` (String) Project ID

### Optional

- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))
- `language_code` (String) Language code

### Read-Only

- `contents` (Block List) List of products  in project (see [below for nested schema](#nestedblock--contents))
- `id` (String) The ID of this resource.
- `total_count` (Number) Total list size

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Filtering target name
- `values` (List of String) Filtering values. Each matching value is appended. (OR rule)

Optional:

- `use_regex` (Boolean) Enable regex match for values


<a id="nestedblock--contents"></a>
### Nested Schema for `contents`

Read-Only:

- `product_category_description` (String) Product category description
- `product_category_id` (String) Product category id
- `product_category_name` (String) Product category name
- `product_category_state` (String) Product category state
- `product_set` (String) Product category set
- `products` (Block List) List of product resources (see [below for nested schema](#nestedblock--contents--products))

<a id="nestedblock--contents--products"></a>
### Nested Schema for `contents.products`

Read-Only:

- `is_product_creatable` (String) Product creation availability
- `product_offering_description` (String) Product offering description
- `product_offering_detail_info` (String) Product offering details
- `product_offering_id` (String) Product offering ID
- `product_offering_name` (String) Product offering name
- `product_offering_state` (String) Product offering state


