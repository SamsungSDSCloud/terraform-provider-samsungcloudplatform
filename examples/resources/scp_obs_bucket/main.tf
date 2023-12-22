resource "scp_obs_bucket" "my_scp_obs_bucket" {
  object_storage_bucket_name = var.name
  object_storage_id = "S3OBJECTSTORAGE-XXXXXX"
  service_zone_id = "ZONE-XXXXXXXX"

  object_storage_bucket_file_encryption_enabled = true
  object_storage_bucket_version_enabled = true
  object_storage_bucket_user_purpose = "PRIVATE"
  object_storage_bucket_access_control_enabled = true
  product_names = ["Object Storage"]
  dynamic "access_control_rules" {
    for_each = var.access_control_rules
    content {
      rule_value = access_control_rules.value["rule_value"]
      rule_type = access_control_rules.value["rule_type"]
    }
  }
  tags = {test1:"test2",test3:"test4"}
  object_storage_bucket_dr_enabled = false
}
