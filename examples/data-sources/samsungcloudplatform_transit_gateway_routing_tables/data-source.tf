data "samsungcloudplatform_transit_gateway_routing_tables" "tables_01" {
}

output "contents" {
  value = data.samsungcloudplatform_transit_gateway_routing_tables.tables_01.contents
}
