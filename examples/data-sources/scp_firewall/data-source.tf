# Find first active firewall
data "scp_firewall" "my_fws" {
  vpc_id = "vpc id"

  # Apply filter "ACTIVE" state
  filter {
    name   = "state"
    values = ["INACTIVE"]
  }
}

output "output_my_scp_fw" {
  value = data.scp_firewall.my_fws
}
