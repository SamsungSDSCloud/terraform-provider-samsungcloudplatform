---
page_title: "samsungcloudplatform_subnet_public_ip Resource - samsungcloudplatform"
subcategory: ""
description: |-
  Provides a Subnet Vip reserve resource.
---

# Resource: samsungcloudplatform_subnet_public_ip

Provides a Subnet Vip reserve resource.


## Example Usage

```terraform
resource "samsungcloudplatform_subnet_public_ip" "my_subnet_public_ip" {
  subnet_id            = "SUBNET-xxxxxxxxxxxxxxx"
  vip_id               = "SUBNET_VIRTUAL_IP-xxxxxxxxxxxxxxx"
  //public_ip_address_id=""  /*자동할당*/
  public_ip_address_id = "PUBLIC_IP-xxxxxxxxxxxxxxx"   /*지정할당*/
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `public_ip_address_id` (String) Public Ip Address Id (Reserved Public ip id)
- `subnet_id` (String) Target Subnet id
- `vip_id` (String) subnet Virtual ip id. (Reserved Virtual ip id)

### Read-Only

- `id` (String) The ID of this resource.


