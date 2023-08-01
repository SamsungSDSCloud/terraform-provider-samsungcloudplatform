data "scp_direct_connect_routing_tables" "tables" {
}

data "scp_direct_connect_routing_routes" "routes" {
  routing_table_id = data.scp_direct_connect_routing_tables.tables.contents[0].routing_table_id
}

resource "scp_direct_connect_routing" "routing01" {
  routing_table_id = data.scp_direct_connect_routing_tables.tables.contents[0].routing_table_id
  destination_network_cidr = "192.168.158.0/24"
  source_service_interface_id = data.scp_direct_connect_routing_routes.routes.contents[0].source_service_interface_id
  source_service_interface_name = data.scp_direct_connect_routing_routes.routes.contents[0].source_service_interface_name
}
