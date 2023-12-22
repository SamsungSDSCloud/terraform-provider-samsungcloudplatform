
data "terraform_remote_state" "engine" {
  backend = "local"

  config = {
    path = "../scp_kubernetes_engine/terraform.tfstate"
  }
}

variable "name" {
  default = "nodepooltest"
}

variable "scale_name" {
  default = "s1v2m4"
}

variable "storage_name" {
  default = "SSD"
}

variable "storage_size_gb" {
  default = "100"
}

variable "availability_zone_name" {
  default = "AZ1"
}
