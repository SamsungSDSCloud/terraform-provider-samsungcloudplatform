# Find all Auto Scaling Groups
data "samsungcloudplatform_auto_scaling_groups" "my_auto_scaling_groups1" {
}

# Find all Auto Scaling Groups
data "samsungcloudplatform_auto_scaling_groups" "my_auto_scaling_groups2" {
  # Sort in ascending order of creation date
  sort = "createdDt:asc"

  # Set paging condition
  page = 0
  size = 100

  # Apply filter for 'asg_name' regex value "test"
  filter {
    name = "asg_name"
    values = ["test"]
    use_regex = true
  }
}

output "output_scp_auto_scaling_groups1" {
  value = data.samsungcloudplatform_auto_scaling_groups.my_auto_scaling_groups1
}

output "output_scp_auto_scaling_groups2" {
  value = data.samsungcloudplatform_auto_scaling_groups.my_auto_scaling_groups2
}
