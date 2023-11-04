
terraform {
  required_providers {
    scp = {
      version = "3.0.0"
      source  = "SamsungSDSCloud/samsungcloudplatform"
    }
  }
  required_version = ">= 0.13"
}

# Provider setup
provider "scp" {
}
