data "samsungcloudplatform_virtual_servers" "virtual_server_list" {
  filter {
    name   = "serviced_group_for"
    values = ["COMPUTE"]
  }
  filter {
    name   = "serviced_for"
    values = ["Virtual Server"]
  }
  filter {
    name   = "virtual_server_state"
    values = ["RUNNING"]
  }
}

variable "name" {
  default = "tf_test_image_001"
}

variable "desc" {
  default = "description"
}
