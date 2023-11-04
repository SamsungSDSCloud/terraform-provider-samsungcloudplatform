data "scp_product" "my_product" {
  product_id = "PRODUCT-XXXXXXXX"
}

output "result_my_product" {
  value = data.scp_product.my_product
}
