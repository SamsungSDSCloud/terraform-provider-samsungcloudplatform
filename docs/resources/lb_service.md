---
page_title: "scp_lb_service Resource - scp"
subcategory: ""
description: |-
  Provides a Load Balancer Service resource.
---

# Resource: scp_lb_service

Provides a Load Balancer Service resource.


## Example Usage

```terraform
resource "scp_lb_service" "my_lb_service_l4" {
  lb_id            = data.terraform_remote_state.load_balancer.outputs.id
  name             = var.namel4
  layer_type       = "L4"
  protocol         = "TCP"
  service_ports    = "8090"
  forwarding_ports = "8091"
  service_ipv4     = "192.168.102.20"
  persistence      = "SOURCE_IP"
  app_profile_id   = data.terraform_remote_state.load_balancer_profile.outputs.id
  persistence_profile_id = data.terraform_remote_state.load_balancer_profile.outputs.persistence_id
}

resource "scp_lb_service" "my_lb_service_l7" {
  lb_id            = data.terraform_remote_state.load_balancer.outputs.id
  app_profile_id   = data.terraform_remote_state.load_balancer_profile.outputs.id
  name             = var.namel7
  layer_type       = "L7"
  protocol         = "HTTP"
  service_ports    = "8088"
  forwarding_ports = "8089"
  service_ipv4     = "192.168.102.11"
  persistence      = "DISABLED"

  lb_rules {
    lb_rule_seq = 1
    pattern_url = "/promise"
  }
  lb_rules {
    lb_rule_seq = 2
    pattern_url = "/devotion"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_profile_id` (String) Application Profile ID
- `layer_type` (String) Servicing protocol layer. (L4 for TCP, L7 for HTTP or HTTPS)
- `lb_id` (String) Target Load-Balancer id.
- `name` (String) Name of Load-Balancer Service. (3 to 20 characters without specials)
- `persistence` (String) Persistence option. (DISABLED, SOURCE_IP, COOKIE)
- `protocol` (String) Servicing protocol. (TCP, HTTP, HTTPS)
- `service_ports` (String) Servicing port numbers. Multiple ports can be inserted using comma and dash. (e.g. 8000-8100,8200)

### Optional

- `forwarding_ports` (String) Forwarding port numbers. Multiple ports can be inserted using comma and dash. (e.g. 8000-8100,8200)
- `lb_rules` (Block List) Server-Group rules. (see [below for nested schema](#nestedblock--lb_rules))
- `lb_service_ip_id` (String)
- `persistence_profile_id` (String) Persistence target profile id.
- `service_ipv4` (String) Servicing IP address
- `use_access_log` (Boolean)

### Read-Only

- `client_certificate_id` (String) SSL client certification id.
- `id` (String) The ID of this resource.
- `server_certificate_id` (String) SSL server certification id.

<a id="nestedblock--lb_rules"></a>
### Nested Schema for `lb_rules`

Required:

- `lb_rule_seq` (Number)

Optional:

- `lb_server_group_id` (String) Target server-group id.
- `pattern_url` (String) Pattern URL.

Read-Only:

- `lb_rule_id` (String)