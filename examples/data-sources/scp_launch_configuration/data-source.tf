# Find details of Launch Configuration
data "scp_launch_configuration" "my_launch_configuration1" {
  lc_id = "LAUNCH_CONFIGURATION-XXXXX"
}

output "output_scp_launch_configuration1" {
  value = data.scp_launch_configuration.my_launch_configuration1
}
