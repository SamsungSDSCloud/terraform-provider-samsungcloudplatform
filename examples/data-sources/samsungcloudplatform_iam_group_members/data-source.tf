data "samsungcloudplatform_iam_groups" "my_own_groups" {
  group_name = "AdministratorGroup"
}

data "samsungcloudplatform_iam_group_members" "my_group_members" {
  group_id = data.samsungcloudplatform_iam_groups.my_own_groups.contents[0].group_id
}

output "result_my_groups" {
  value = data.samsungcloudplatform_iam_group_members.my_group_members
}
