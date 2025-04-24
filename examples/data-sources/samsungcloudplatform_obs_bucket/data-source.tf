data "samsungcloudplatform_obs_bucket" "my_scp_obs_bucket" {
}

output "output_my_scp_obs_buckets" {
  value = data.samsungcloudplatform_obs_bucket.my_scp_obs_bucket
}
