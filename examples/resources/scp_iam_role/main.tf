resource "scp_iam_role" "my_role01" {
  role_name = var.name

  trust_principals {
    project_ids = [var.proj_id]
    user_srns = [var.srn]
  }

  tags = {
    tk01 = "tv01"
    tk02 = "tv02"
  }

  description = ""
}
