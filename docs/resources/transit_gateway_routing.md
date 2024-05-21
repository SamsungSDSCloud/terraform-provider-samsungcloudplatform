---
page_title: "scp_transit_gateway_routing Resource - scp"
subcategory: ""
description: |-
  Provides a Transit Gateway Connection Routing Rule Resources
---

# Resource: scp_transit_gateway_routing

Provides a Transit Gateway Connection Routing Rule Resources


## Example Usage

```terraform
data "scp_transit_gateway_routing_tables" "tables" {
}

data "scp_transit_gateway_routing_routes" "routes" {
  routing_table_id = data.scp_transit_gateway_routing_tables.tables.contents[0].routing_table_id
}

resource "scp_transit_gateway_routing" "tgwRoutingRule01" {
  routing_table_id  = data.scp_transit_gateway_routing_tables.tables.contents[0].routing_table_id
  destination_network_cidr = var.destinationNetworkCidr
  source_service_interface_id   = data.scp_transit_gateway_routing_routes.routes.contents[0].source_service_interface_id
  source_service_interface_name   = data.scp_transit_gateway_routing_routes.routes.contents[0].source_service_interface_name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `destination_network_cidr` (String) Network CIDR
- `routing_table_id` (String) Routing Table ID for Transit Gateway Connection
- `source_service_interface_id` (String) Source Interface ID
- `source_service_interface_name` (String) Source Interface Name

### Read-Only

- `id` (String) The ID of this resource.