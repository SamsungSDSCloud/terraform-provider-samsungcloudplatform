data "scp_products_by_zone" "my_products" {
  service_zone_id = "ZONE-XXXXXXXXXXX"
}

output "result_my_products" {
  value = data.scp_products_by_zone.my_products
}
