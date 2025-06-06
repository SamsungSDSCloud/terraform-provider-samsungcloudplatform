---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_project_product_resources Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Provides product resources for given project
---

# samsungcloudplatform_project_product_resources (Data Source)

Provides product resources for given project

## Example Usage

```terraform
data "samsungcloudplatform_project" "my_project"{

}

data "samsungcloudplatform_project_product_resources" "my_scp_product_resources" {
  project_id = data.samsungcloudplatform_project.my_project.id
}

output "output_my_scp_products_resources" {
  value = data.samsungcloudplatform_project_product_resources.my_scp_product_resources
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project_id` (String) Project ID

### Optional

- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))
- `product_category_id` (String) Product category ID

### Read-Only

- `contents` (Block List) List of product resources in project (see [below for nested schema](#nestedblock--contents))
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

- `product_category_name` (String) Product category name
- `product_resources` (Block List) List of product resources (see [below for nested schema](#nestedblock--contents--product_resources))

<a id="nestedblock--contents--product_resources"></a>
### Nested Schema for `contents.product_resources`

Read-Only:

- `product_offering_name` (String) Product offering names
- `product_offering_resource_count` (Number) Number of resources provided by product


