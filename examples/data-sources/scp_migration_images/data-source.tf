
data "scp_migration_images" "my_images1" {
  service_group = "COMPUTE"
  service       = "Virtual Server"
}


output "result_scp_my_migration_images1" {
  value = data.scp_migration_images.my_images1
}
