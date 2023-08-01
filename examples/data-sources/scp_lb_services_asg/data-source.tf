data "scp_lb_services_connectable_to_asg" "my_scp_lb_services_connectable_to_asg" {
  vpc_id = "VPC-XXXXXXXXXXXXXXXXXXXXXX"
}

data "scp_lb_services_connected_to_asg" "my_scp_lb_services_connected_to_asg" {
  auto_scaling_group_id = "AUTO_SCALING_GROUP-XXXXXXXXXXXXXXXXXXXXXX"
}

#Connectable List
output "result_scp_lb_services_connectable_to_asg" {
  value = data.scp_lb_services_connectable_to_asg.my_scp_lb_services_connectable_to_asg
}

#Connected List
output "result_scp_lb_services_connected_to_asg" {
  value = data.scp_lb_services_connected_to_asg.my_scp_lb_services_connected_to_asg
}
