variable "name" {
  default = "backupterraformtest"
}

variable "schedules" {
  type = list(object({
    schedule_frequency = string
    schedule_frequency_detail = string
    schedule_type = string
    start_time = string
  }))
  default = [
    {
      schedule_frequency = "DAYS"
      schedule_frequency_detail = "3"
      schedule_type = "FULL"
      start_time = "12:30:00+09:00"
    },
    {
      schedule_frequency = "WEEKLY"
      schedule_frequency_detail = "MON"
      schedule_type = "INCREMENTAL"
      start_time = "15:30:00+09:00"
    },
    {
      schedule_frequency = "MONTHLY"
      schedule_frequency_detail = "3"
      schedule_type = "INCREMENTAL"
      start_time = "05:30:00+09:00"
    },
  ]
}
