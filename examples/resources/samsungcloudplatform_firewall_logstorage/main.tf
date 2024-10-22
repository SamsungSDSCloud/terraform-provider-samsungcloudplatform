data "samsungcloudplatform_region" "region" {
  filter {
    name = "location"
    values = ["KR-WEST-1"]
  }
}

data "samsungcloudplatform_obs_storages" "storages" {
  zone_id = data.samsungcloudplatform_region.region.id
}

resource "samsungcloudplatform_obs_bucket" "mybucket" {
  name = var.name
  obs_id = data.samsungcloudplatform_obs_storages.storages.contents[0].obs_id
  zone_id = data.samsungcloudplatform_region.region.id
  ip_address_filter_enabled = false
  file_encryption_enabled = true
  file_encryption_algorithm = "AES256"
  file_encryption_type      = "SSE-S3"
  version_enabled = true
}

resource "samsungcloudplatform_firewall_logstorage" "this" {
  //vpc_id = var.vpc_id
  vpc_id          = data.terraform_remote_state.vpc.outputs.id
  obs_bucket_id = samsungcloudplatform_obs_bucket.mybucket.id
}
