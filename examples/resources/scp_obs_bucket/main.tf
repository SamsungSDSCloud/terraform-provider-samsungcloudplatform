resource "scp_obs_bucket" "my_scp_obs_bucket" {
  obs_bucket_name = var.name
  obs_id = "S3OBJECTSTORAGE-XXXXX"
  zone_id = "ZONE-XXXXX"

  obs_bucket_file_encryption_enabled = true
  obs_bucket_file_encryption_algorithm = "AES256"
  obs_bucket_file_encryption_type      = "SSE-S3"
  obs_bucket_version_enabled = true

  is_obs_bucket_ip_address_filter_enabled = true
  dynamic "obs_bucket_access_ip_address_ranges" {
    for_each = var.obs_bucket_access_ip_address_ranges
    content {
      obs_bucket_access_ip_address_range = obs_bucket_access_ip_address_ranges.value["obs_bucket_access_ip_address_range"]
      type = obs_bucket_access_ip_address_ranges.value["type"]
    }
  }
}
