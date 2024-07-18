variable co_service_zone_id {
  default = "TA_HPC_Lite"
}
variable contract {
  default = "RI"
}
variable hyper_threading_enabled {
  default = "N"
}
variable image_id {
  default = "IMAGE-56Sgh75jqKfNPyWLl6Ekth"
}
variable init_script {
  default = "init_script"
}
variable os_user_id {
  default = "root"
}
variable os_user_password {
  default = "!q2w3e4r"
  sensitive = true
}
variable product_group_id {
  default = "PRODUCTGROUP-msOf1BSmE8bCfGWzBKcDQE"
}
variable resource_pool_id {
  default = "POOL-KZtNRugNqXcSZdCHkhEl3l"
}
variable server_type {
  default = "2023_001"
}
variable service_zone_id {
  default = "ZONE-9jddEZ88tWoRtF8Z5eQabl"
}
variable vlan_pool_cidr {
  default = "172.24.159.0/24"
}

variable "server_details" {
  default = [
    {
      server_name = "terrahpc01"
    },
    {
      server_name = "terrahpc02"
    }
  ]
}
