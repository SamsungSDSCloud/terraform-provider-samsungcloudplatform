resource "scp_iam_role" "my_role01" {
  role_name = var.name

  trust_principals {
    project_ids = [var.proj_id]
    user_srns = [var.srn]
  }

  tags {
    tag_key = "tk01"
    tag_value = "tv01"
  }
  tags {
    tag_key = "tk02"
    tag_value = "tv02"
  }

  description = ""
}
