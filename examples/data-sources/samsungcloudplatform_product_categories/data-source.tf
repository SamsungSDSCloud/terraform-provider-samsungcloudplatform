data "samsungcloudplatform_region" "my_region" {
}

data "samsungcloudplatform_product_categories" "my_scp_product_categories" {
  language_code = "en_US"
}

output "output_my_scp_product" {
  value = data.samsungcloudplatform_product_categories.my_scp_product_categories
}



