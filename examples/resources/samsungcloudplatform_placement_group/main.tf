locals {
  service-zone-id = data.samsungcloudplatform_region.region.id
  virtual-server-ids = tolist([for i, element in data.samsungcloudplatform_virtual_servers.target_vm.contents : element.virtual_server_id])
}

resource "samsungcloudplatform_placement_group" "placement_group_001" {
  placement_group_name = var.placement-group-name
  virtual_server_ids = var.numb-of-virtual-server-ids > 0 ? local.virtual-server-ids : var.virtual-server-ids
  service_zone_id = local.service-zone-id
  virtual_server_type = var.virtual-server-type
  description = var.description
}
