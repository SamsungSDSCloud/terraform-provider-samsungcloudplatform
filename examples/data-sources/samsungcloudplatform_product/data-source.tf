data "samsungcloudplatform_product" "my_product" {
  product_id = "PRODUCT-XXXXXXXX"
}

output "result_my_product" {
  value = data.samsungcloudplatform_product.my_product
}
