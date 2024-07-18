data "scp_obs_buckets" "buckets" {

}

resource "scp_trail" "my_trail" {
  name          = var.name
  obs_bucket_id = data.scp_obs_buckets.buckets.contents[0].object_storage_bucket_id
  save_type     = var.save_type

  use_verification = true

  is_logging_target_all_user = true
  is_logging_target_all_resource = true
  is_logging_target_all_region = true

  description = "Trail from Terraform"
}

