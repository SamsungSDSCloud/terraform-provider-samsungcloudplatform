data "scp_obs_bucket" "my_scp_obs_bucket" {
  object_storage_bucket_id = "OBS_XXXXXX"
}

output "output_my_scp_obs_buckets" {
  value = data.scp_obs_bucket.my_scp_obs_bucket
}
