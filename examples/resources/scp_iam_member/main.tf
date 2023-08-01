resource "scp_iam_member" "my_member01" {
  user_email = var.email
  group_ids = [var.group_id_1, var.group_id_2]

  tags {
    tag_key = "tk01"
    tag_value = "tv01"
  }
  tags {
    tag_key = "tk02"
    tag_value = "tv02"
  }
}
