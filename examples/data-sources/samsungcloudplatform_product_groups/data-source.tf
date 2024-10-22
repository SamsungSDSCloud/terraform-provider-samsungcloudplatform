data "samsungcloudplatform_product_groups" "my_groups" {

}

output "result_my_groups" {
  value = data.samsungcloudplatform_product_groups.my_groups
}
