data "scp_project_user_products_resources" "my_scp_products_resources" {
}

output "output_my_scp_products_resources" {
  value = data.scp_project_user_products_resources.my_scp_products_resources
}


