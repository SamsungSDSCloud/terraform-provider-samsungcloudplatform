data "scp_security_groups" "my_sgs" {
}

data "scp_security_groups" "my_sgs2" {
  filter {
    name   = "is_loggable"
    values = ["true"]
  }
}

output "output_my_scp_sg" {
  value = data.scp_security_groups.my_sgs
}

output "output_my_scp_sg2" {
  value = data.scp_security_groups.my_sgs2
}
