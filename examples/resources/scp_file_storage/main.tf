data "scp_region" "region" {
}

resource "scp_file_storage" "my_nfs_fs" {
  name            = var.name
  disk_type       = "SSD"
  protocol        = "NFS"
  is_encrypted    = false
  retention_count = 5
  region          = data.scp_region.region.location
}

resource "scp_file_storage" "my_cifs_fs" {
  name            = "fs_cifs_test"
  disk_type       = "HDD"
  protocol        = "CIFS"
  is_encrypted    = false
  cifs_password   = var.password

  snapshot_day_of_week = "SUN"
  snapshot_frequency   = "WEEKLY"
  snapshot_hour        = 11
  retention_count      = 1
  region               = data.scp_region.region.location
}
