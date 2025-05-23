---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_product_categories Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Provides list of products.
---

# samsungcloudplatform_product_categories (Data Source)

Provides list of products.

## Example Usage

```terraform
data "samsungcloudplatform_region" "my_region" {
}

data "samsungcloudplatform_product_categories" "my_scp_product_categories" {
  language_code = "en_US"
}

output "output_my_scp_product" {
  value = data.samsungcloudplatform_product_categories.my_scp_product_categories
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `language_code` (String) Language code (ko_KR, en_US)

### Optional

- `category_id` (String) Product category id
- `category_state` (String) Product category status
- `exposure_scope` (String) Exposure scope
- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))
- `product_id` (String) Product id
- `product_state` (String) Product status

### Read-Only

- `contents` (Block List) Product list (see [below for nested schema](#nestedblock--contents))
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

- `icon_file_name` (String) Icon file name
- `product_category_description` (String) Description of product category
- `product_category_id` (String) Product category id
- `product_category_name` (String) Product category name
- `product_category_path` (String) Product category path
- `product_category_state` (String) Product category status
- `product_set` (String) Product set type (SE, PAAS)


