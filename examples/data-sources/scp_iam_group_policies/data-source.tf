data "scp_iam_groups" "my_own_groups" {
  group_name = "AdministratorGroup"
}

data "scp_iam_group_policies" "my_group_policies" {
  group_id = data.scp_iam_groups.my_own_groups.contents[0].group_id
}

output "result_my_groups" {
  value = data.scp_iam_group_policies.my_group_policies
}
