terraform {
  required_providers {
    samsungcloudplatform = {
      version = "3.1.0"
      source  = "SamsungSDSCloud/samsungcloudplatform"
    }
  }
  required_version = ">= 0.13"
}

# Provider setup
provider "samsungcloudplatform" {
}
