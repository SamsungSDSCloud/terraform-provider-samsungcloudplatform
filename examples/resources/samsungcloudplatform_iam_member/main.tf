resource "samsungcloudplatform_iam_member" "my_member01" {
  user_email = var.email
  group_ids = [var.group_id_1, var.group_id_2]

  tags = {
    tk01 = "tv01"
    tk02 = "tv02"
  }
}
