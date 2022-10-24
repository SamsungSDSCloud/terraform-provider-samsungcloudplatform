data "scp_firewalls" "my_fws1" {
}

# Find all active firewalls
data "scp_firewalls" "my_fws2" {
  vpc_id = "VPC-xxxxxx"
state
  filter {
    name   = "state"
    values = ["ACTIVE"]
  }
}

output "output_my_scp_fw1" {
  value = data.scp_firewalls.my_fws1
}

output "output_my_scp_fw2" {
  value = data.scp_firewalls.my_fws2
}
