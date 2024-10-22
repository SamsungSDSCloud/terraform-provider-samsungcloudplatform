data "samsungcloudplatform_transit_gateway_routing_tables" "tables" {
}

data "samsungcloudplatform_transit_gateway_routing_routes" "routes" {
  routing_table_id = data.samsungcloudplatform_transit_gateway_routing_tables.tables.contents[0].routing_table_id
}

resource "samsungcloudplatform_transit_gateway_routing" "tgwRoutingRule01" {
  routing_table_id  = data.samsungcloudplatform_transit_gateway_routing_tables.tables.contents[0].routing_table_id
  destination_network_cidr = var.destinationNetworkCidr
  source_service_interface_id   = data.samsungcloudplatform_transit_gateway_routing_routes.routes.contents[0].source_service_interface_id
  source_service_interface_name   = data.samsungcloudplatform_transit_gateway_routing_routes.routes.contents[0].source_service_interface_name
}

