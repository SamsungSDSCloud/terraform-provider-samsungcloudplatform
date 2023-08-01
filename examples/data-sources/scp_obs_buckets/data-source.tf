data "scp_obs_buckets" "my_scp_obs_buckets" {
}

output "output_my_scp_obs_buckets" {
  value = data.scp_obs_buckets.my_scp_obs_buckets
}
