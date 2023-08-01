
data "scp_migration_image" "my_image1" {
  image_id = "IMAGE-XXXX"
  service_zone_id = "ZONE-XXXX"
}


output "result_scp_my_migration_images1" {
  value = data.scp_migration_image.my_image1
}
