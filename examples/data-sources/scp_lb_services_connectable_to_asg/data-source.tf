data "scp_lb_services_connectable_to_asg" "my_scp_lb_services_connectable_to_asg" {
  vpc_id = "VPC-XXXXXXXXXXXXXXXXXXXXXX"
}

#Connectable List
output "result_scp_lb_services_connectable_to_asg" {
  value = data.scp_lb_services_connectable_to_asg.my_scp_lb_services_connectable_to_asg
}
