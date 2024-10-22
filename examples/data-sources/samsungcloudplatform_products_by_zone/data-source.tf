data "samsungcloudplatform_products_by_zone" "my_products" {
  service_zone_id = "ZONE-XXXXXXXXXXX"
}

output "result_my_products" {
  value = data.samsungcloudplatform_products_by_zone.my_products
}
