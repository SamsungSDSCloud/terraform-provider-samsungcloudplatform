data "samsungcloudplatform_vpc_routing_tables" "tables_01" {
}

output "contents" {
  value = data.samsungcloudplatform_vpc_routing_tables.tables_01.contents
}
