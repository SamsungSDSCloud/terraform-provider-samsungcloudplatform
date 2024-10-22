data "samsungcloudplatform_vpc_routing_tables" "table" {
}

data "samsungcloudplatform_vpc_routing_rules" "rules" {
  routing_table_id = data.samsungcloudplatform_vpc_routing_tables.table.contents[0].routing_table_id
}

output "contents" {
  value = data.samsungcloudplatform_vpc_routing_rules.rules.contents
}
