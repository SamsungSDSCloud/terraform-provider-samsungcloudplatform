data "samsungcloudplatform_project_user_products_resources" "my_scp_products_resources" {
}

output "output_my_scp_products_resources" {
  value = data.samsungcloudplatform_project_user_products_resources.my_scp_products_resources
}


