resource "samsungcloudplatform_migration_image" "my_migration_image" {
  image_name = var.name
  original_image_id = var.image_id
  ova_url = var.url
  access_key = var.access_key
  secret_key = var.secret_key
  os_user_id = var.os_id
  os_user_password = var.os_pw
  image_description = var.desc
  az_name =var.az_name
  service_zone_id = var.service_zone_id
}

