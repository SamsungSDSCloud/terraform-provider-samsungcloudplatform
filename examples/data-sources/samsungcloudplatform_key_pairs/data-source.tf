# Find all Key pairs
data "samsungcloudplatform_key_pairs" "my_scp_key_pairs1" {
}

# Find all Key pairs
data "samsungcloudplatform_key_pairs" "my_scp_key_pairs2" {
  # Sort in ascending order of creation date
  sort = "createdDt:asc"
}

output "output_scp_key_pairs1" {
  value = data.samsungcloudplatform_key_pairs.my_scp_key_pairs1
}

output "output_scp_key_pairs2" {
  value = data.samsungcloudplatform_key_pairs.my_scp_key_pairs2
}
