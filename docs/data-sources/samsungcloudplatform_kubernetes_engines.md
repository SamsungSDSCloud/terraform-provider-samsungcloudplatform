---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "samsungcloudplatform_kubernetes_engines Data Source - scp"
subcategory: ""
description: |-
  Provides list of K8s engines
---

# samsungcloudplatform_kubernetes_engines (Data Source)

Provides list of K8s engines

## Example Usage

```terraform
# Find my engines for current project
data "samsungcloudplatform_kubernetes_engines" "my_scp_kubernetes_engines" {
}

output "result_scp_kubernetes_engines" {
  value = data.samsungcloudplatform_kubernetes_engines.my_scp_kubernetes_engines
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `created_by` (String) The person who created the resource
- `k8s_version` (String) K8s cluster version
- `kubernetes_engine_name` (String) K8s engine name
- `kubernetes_engine_status` (String) K8s engine status
- `page` (Number) Page start number from which to get the list
- `region` (String) Region
- `size` (Number) Size to get list

### Read-Only

- `contents` (Block List) K8s engine list (see [below for nested schema](#nestedblock--contents))
- `id` (String) The ID of this resource.
- `total_count` (Number) Content list size

<a id="nestedblock--contents"></a>
### Nested Schema for `contents`

Read-Only:

- `created_by` (String) The person who created the resource
- `created_dt` (String) Creation time
- `k8s_version` (String) K8s version
- `kubernetes_engine_id` (String) K8s engine id
- `kubernetes_engine_name` (String) K8s engine name
- `kubernetes_engine_status` (String) K8s engine status
- `modified_by` (String) The person who modified the resource
- `modified_dt` (String) Modification time
- `node_count` (Number) K8s node count
- `project_id` (String) Project id
- `region` (String) Region name
- `security_group_id` (String) Security group id
- `subnet_id` (String) Subnet id
- `volume_id` (String) File storage volume id
- `vpc_id` (String) Vpc id

