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

