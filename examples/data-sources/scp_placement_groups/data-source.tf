data "scp_region" "region" {
}

data "scp_placement_groups" "my_group" {
  service_zone_id = data.scp_region.region.id
  virtual_server_type = "s1"
}
