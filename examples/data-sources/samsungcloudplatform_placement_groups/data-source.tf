data "samsungcloudplatform_region" "region" {
}

data "samsungcloudplatform_placement_groups" "my_group" {
  service_zone_id = data.samsungcloudplatform_region.region.id
  virtual_server_type = "s1"
}
