data "scp_obs_storages" "my_scp_obs_storages" {
  zone_id = "ZONE-XXXXX"
}

output "output_my_scp_obs_storages" {
  value = data.scp_obs_storages.my_scp_obs_storages
}
