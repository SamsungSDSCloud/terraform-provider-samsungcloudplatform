data "scp_products_by_group" "my_products" {
  product_group_id = "PRODUCTGROUP-XXXXXXXXXX"
}

output "result_my_products" {
  value = data.scp_products_by_group.my_products
}
