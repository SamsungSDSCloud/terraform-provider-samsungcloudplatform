data "samsungcloudplatform_direct_connect_routing_tables" "tables_01" {
}

output "contents" {
  value = data.samsungcloudplatform_direct_connect_routing_tables.tables_01.contents
}
