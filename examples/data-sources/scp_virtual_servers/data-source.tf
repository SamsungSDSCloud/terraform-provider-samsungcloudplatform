data "scp_virtual_servers" "my_virtual_servers" {
  size = 10
  page = 0
  virtual_server_name = "test"
}

data "scp_virtual_servers" "my_asg_virtual_servers" {
  auto_scaling_group_id = "AUTO_SCALING_GROUP-KMbKbQ3gsVkPZna408fm4g"
}

output "output_my_virtual_servers" {
  value = data.scp_virtual_servers.my_virtual_servers
}

output "output_my_asg_virtual_servers" {
  value = data.scp_virtual_servers.my_asg_virtual_servers
}
