variable "block_storages" {
  default = [
    {
      "block_storage_size": 100,
      "disk_type": "SSD",
      "encryption_enabled": false,
      "is_boot_disk": true
    },
    {
      "block_storage_size": 4,
      "disk_type": "SSD",
      "encryption_enabled": false,
      "is_boot_disk": false
    }
  ]
}

variable "image_id" {
  default = "IMAGE-XXXXX"
}

variable "initial_script" {
  default = "ls"
}

variable "key_pair_id" {
  default = "KEY_PAIR-XXXXX"
}

variable "lc_name" {
  default = "my-lc"
}

variable "server_type" {
  default = "s1v1m2"
}

variable "service_zone_id" {
  default = "ZONE-XXXXX"
}
