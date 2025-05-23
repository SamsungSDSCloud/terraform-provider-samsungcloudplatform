---
page_title: "samsungcloudplatform_dns_record Resource - samsungcloudplatform"
subcategory: ""
description: |-
  Provides a DnsDomain resource.
---

# Resource: samsungcloudplatform_dns_record

Provides a DnsDomain resource.


## Example Usage

```terraform
resource "samsungcloudplatform_dns_record" "my_dns_record" {
  dns_domain_id         = data.terraform_remote_state.dns_domain.outputs.id
  dns_record_name       = var.name
  dns_record_type       = "MX"
  ttl = 300
  dns_record_mapping {
    record_destination            = "192.168.0.1"
    preference = 1
  }
  dns_record_mapping {
    record_destination            = "192.168.0.2"
    preference = 2
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `dns_domain_id` (String) DNS Domain Id
- `dns_record_mapping` (Block Set, Min: 1) DNS Record Mappings. Record Type CNAME, SPF, TXT can have only one record mapping. Record Type A, AAAA, MX can have 1 or more record mappings. (see [below for nested schema](#nestedblock--dns_record_mapping))
- `dns_record_name` (String) DNS Record Name (0 to 63, lowercase, number and -_.@)
- `dns_record_type` (String) DNS Record Type. One of A, TXT, CNAME, MX, AAAA, SPF
- `ttl` (Number) DNS TTL. (300 to 86400)

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--dns_record_mapping"></a>
### Nested Schema for `dns_record_mapping`

Required:

- `record_destination` (String) DnsDomain Resource Destination

Optional:

- `preference` (Number) DnsDomain Resource Weight


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)


