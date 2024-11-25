terraform {
  required_providers {
    samsungcloudplatform = {
      version = "3.8.1"
      source  = "SamsungSDSCloud/samsungcloudplatform"
    }
  }
  required_version = ">= 0.13"
}

provider "samsungcloudplatform" {
}
