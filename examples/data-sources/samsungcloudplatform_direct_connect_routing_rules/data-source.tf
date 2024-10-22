data "samsungcloudplatform_direct_connect_routing_tables" "table" {
}

data "samsungcloudplatform_direct_connect_routing_rules" "rules" {
  routing_table_id = data.samsungcloudplatform_direct_connect_routing_tables.table.contents[0].routing_table_id
}

output "contents" {
  value = data.samsungcloudplatform_direct_connect_routing_rules.rules.contents
}
