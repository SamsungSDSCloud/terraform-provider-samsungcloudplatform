resource "scp_iam_group" "my_group01" {
  group_name = var.name
  description = var.desc
}
