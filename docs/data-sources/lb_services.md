---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "scp_lb_services Data Source - scp"
subcategory: ""
description: |-
  Provides list of Load Balancer services
---

# scp_lb_services (Data Source)

Provides list of Load Balancer services

## Example Usage

```terraform
data "scp_lb_services" "my_scp_lb_services" {
  load_balancer_id = "lb id"
}
output "output_scp_public_ips" {
  value = data.scp_lb_services.my_scp_lb_services
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `load_balancer_id` (String) Load balancer id

### Optional

- `created_by` (String) The person who created the resource
- `layer_type` (String) Protocol layer type (L4, L7)
- `lb_service_name` (String) Load balancer service Name
- `load_balancer_name` (String) Load balancer name
- `page` (Number) Page start number from which to get the list
- `protocol` (String) The file storage protocol type to create (NFS, CIFS)
- `service_ip_address` (String) Service ip address
- `size` (Number) Size to get list
- `sort` (String) Sort
- `status_check` (Boolean) check status

### Read-Only

- `contents` (Block List) Load balancer service list (see [below for nested schema](#nestedblock--contents))
- `id` (String) The ID of this resource.
- `total_count` (Number)

<a id="nestedblock--contents"></a>
### Nested Schema for `contents`

Read-Only:

- `block_id` (String) Block id of this region
- `created_by` (String) The person who created the resource
- `created_dt` (String) Creation date
- `default_forwarding_ports` (String) Default forwarding ports
- `layer_type` (String) Protocol layer type (L4, L7)
- `lb_service_id` (String) Load balancer service id
- `lb_service_ip_id` (String) Load balancer service ip id
- `lb_service_name` (String) Load balancer service name
- `lb_service_state` (String) Load balancer service status
- `modified_by` (String) The person who modified the resource
- `modified_dt` (String) Modification date
- `nat_ip_address` (String) Nat ip address
- `project_id` (String) Load balancer service ip id
- `protocol` (String) Protocol
- `service_ip_address` (String) Service ip address
- `service_ports` (String) Service ports
- `service_zone_id` (String) Service zone id

