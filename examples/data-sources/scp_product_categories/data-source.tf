data "scp_region" "my_region" {
}

data "scp_product_categories" "my_scp_product_categories" {
  language_code = "en_US"
}

output "output_my_scp_product" {
  value = data.scp_product_categories.my_scp_product_categories
}



