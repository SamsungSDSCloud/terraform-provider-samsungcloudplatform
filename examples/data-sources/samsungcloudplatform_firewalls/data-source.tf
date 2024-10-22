data "samsungcloudplatform_firewalls" "my_fws1" {
}

# Find all active firewalls
data "samsungcloudplatform_firewalls" "my_fws2" {
  vpc_id = "VPC-xxxxxx"
  filter {
    name   = "state"
    values = ["ACTIVE"]
  }
}

output "output_my_scp_fw1" {
  value = data.samsungcloudplatform_firewalls.my_fws1
}

output "output_my_scp_fw2" {
  value = data.samsungcloudplatform_firewalls.my_fws2
}
