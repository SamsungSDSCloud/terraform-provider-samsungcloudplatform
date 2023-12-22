data "scp_transit_gateway_routing_tables" "table" {
}

data "scp_transit_gateway_routing_rules" "rules" {
  routing_table_id = data.scp_transit_gateway_routing_tables.table.contents[0].routing_table_id
}

output "contents" {
  value = data.scp_transit_gateway_routing_rules.rules.contents
}
