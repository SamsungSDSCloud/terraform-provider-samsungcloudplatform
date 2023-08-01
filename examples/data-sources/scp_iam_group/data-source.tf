data "scp_iam_group" "my_group" {

}

output "result_my_group" {
  value = data.scp_iam_group.my_group
}
