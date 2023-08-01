resource "scp_iam_access_key" "my_access_key1" {
  project_id = var.project_id
  duration_days = var.duration_days
  access_key_activated = true
}
