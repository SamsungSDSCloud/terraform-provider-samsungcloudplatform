# Find first active firewall
data "samsungcloudplatform_firewall" "my_fws" {
  vpc_id = "vpc id"

  # Apply filter "ACTIVE" state
  filter {
    name   = "state"
    values = ["INACTIVE"]
  }
}

output "output_my_scp_fw" {
  value = data.samsungcloudplatform_firewall.my_fws
}
