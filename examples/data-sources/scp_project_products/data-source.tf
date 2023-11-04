data "scp_project" "my_project" {

}

data "scp_project_products" "my_scp_products" {
  project_id = data.scp_project.my_project.id
}

output "scp_project_products" {
  value = data.scp_project_products.my_scp_products
}


