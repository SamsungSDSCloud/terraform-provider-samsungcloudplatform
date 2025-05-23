---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_postgresqls Data Source - samsungcloudplatform"
subcategory: ""
description: |-
  Provides list of postgresql databases.
---

# samsungcloudplatform_postgresqls (Data Source)

Provides list of postgresql databases.

## Example Usage

```terraform
data "samsungcloudplatform_postgresqls" "my_scp_postgresqls" {
}

output "output_my_scp_postgresqls" {
  value = data.samsungcloudplatform_postgresqls.my_scp_postgresqls
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `page` (Number) Page start number from which to get the list.
- `postgresql_cluster_name` (String) Database name.
- `size` (Number) Size to get list.
- `sort` (String) Sort

### Read-Only

- `contents` (Block List) PostgreSQL list (see [below for nested schema](#nestedblock--contents))
- `id` (String) The ID of this resource.
- `total_count` (Number)

<a id="nestedblock--contents"></a>
### Nested Schema for `contents`

Read-Only:

- `block_id` (String) Block ID.
- `created_by` (String) The person who created the resource
- `created_dt` (String) Creation date
- `modified_by` (String) The person who modified the resource
- `modified_dt` (String) Modification date
- `postgresql_cluster_id` (String) PostgreSQL Cluster ID
- `postgresql_cluster_name` (String) PostgreSQL Cluster Name
- `postgresql_cluster_state` (String) PostgreSQL Cluster State
- `project_id` (String) Project ID.
- `service_zone_id` (String) Service Zone ID


