data "samsungcloudplatform_iam_member_systemgroups" "my_member_systemgroups" {
  member_id = "XXXX"
}

output "result_my_member_systemgroups" {
  value = data.samsungcloudplatform_iam_member_systemgroups.my_member_systemgroups
}
