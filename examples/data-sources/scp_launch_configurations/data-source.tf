# Find all Launch Configurations
data "scp_launch_configurations" "my_launch_configurations1" {
}

# Find all Launch Configurations
data "scp_launch_configurations" "my_launch_configurations2" {
  # Sort in ascending order of creation date
  sort = "createdDt:asc"

  # Apply filter for 'lc_name' regex value "test"
  filter {
    name = "lc_name"
    values = ["test"]
    use_regex = true
  }
}

output "output_scp_launch_configurations1" {
  value = data.scp_launch_configurations.my_launch_configurations1
}

output "output_scp_launch_configurations2" {
  value = data.scp_launch_configurations.my_launch_configurations2
}
