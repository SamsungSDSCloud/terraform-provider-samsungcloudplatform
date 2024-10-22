data "samsungcloudplatform_project" "my_project" {

}

data "samsungcloudplatform_project_products" "my_scp_products" {
  project_id = data.samsungcloudplatform_project.my_project.id
}

output "samsungcloudplatform_project_products" {
  value = data.samsungcloudplatform_project_products.my_scp_products
}


