data "samsungcloudplatform_region" "region" {
}

resource "samsungcloudplatform_auto_scaling_group" "my_asg" {
  asg_name = var.name
  availability_zone_name = var.az_name
  desired_server_count       = var.desired
  min_server_count  = var.min
  max_server_count  = var.max
  desired_server_count_editable = true
  lc_id = data.terraform_remote_state.lc.outputs.id
  multi_availability_zone_enabled = false
  security_group_ids = [
    data.terraform_remote_state.security_group.outputs.id
  ]
  server_name_prefix = var.name
  vpc_info {
    vpc_id          = data.terraform_remote_state.vpc.outputs.id
    subnet_id       = data.terraform_remote_state.subnet.outputs.id
    nat_enabled     = false
  }
}
