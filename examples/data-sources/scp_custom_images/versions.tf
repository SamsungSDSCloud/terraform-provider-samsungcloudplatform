terraform {
  required_providers {
    scp = {
      version = "1.8.6"
      source  = "SamsungSDSCloud/samsungcloudplatform"
    }
  }
  required_version = ">= 0.13"
}

# Provider setup
provider "scp" {
}
