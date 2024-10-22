data "samsungcloudplatform_transit_gateway_routing_tables" "table" {
}

data "samsungcloudplatform_transit_gateway_routing_routes" "routes01" {
  routing_table_id = data.samsungcloudplatform_transit_gateway_routing_tables.table.contents[0].routing_table_id
}

output "contents" {
  value = data.samsungcloudplatform_transit_gateway_routing_routes.routes01.contents
}
