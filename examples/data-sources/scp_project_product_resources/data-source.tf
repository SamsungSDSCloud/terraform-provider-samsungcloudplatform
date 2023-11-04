data "scp_project" "my_project"{

}

data "scp_project_product_resources" "my_scp_product_resources" {
  project_id = data.scp_project.my_project.id
}

output "output_my_scp_products_resources" {
  value = data.scp_project_product_resources.my_scp_product_resources
}


