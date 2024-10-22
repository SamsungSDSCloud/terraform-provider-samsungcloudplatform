resource "samsungcloudplatform_backup" "my_scp_backup" {
  backup_name = var.name
  backup_policy_type_category = "VM"
  backup_repository = "SD_STORAGE"
  is_backup_dr_enabled = "N"
  object_id = "INSTANCE-XXXXX"
  object_type = "INSTANCE"
  policy_type = "VMsnapshot"
  product_names = [
    "VM Image"
  ]
  retention_period = "4W"
  service_zone_id = "ZONE-XXXXX"

  dynamic "schedules" {
    for_each = var.schedules
    content {
      schedule_frequency = schedules.value["schedule_frequency"]
      schedule_frequency_detail = schedules.value["schedule_frequency_detail"]
      schedule_type = schedules.value["schedule_type"]
      start_time = schedules.value["start_time"]
    }
  }
}
