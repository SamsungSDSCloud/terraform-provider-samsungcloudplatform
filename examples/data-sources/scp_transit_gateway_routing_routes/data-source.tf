data "scp_transit_gateway_routing_tables" "table" {
}

data "scp_transit_gateway_routing_routes" "routes01" {
  routing_table_id = data.scp_transit_gateway_routing_tables.table.contents[0].routing_table_id
}

output "contents" {
  value = data.scp_transit_gateway_routing_routes.routes01.contents
}
