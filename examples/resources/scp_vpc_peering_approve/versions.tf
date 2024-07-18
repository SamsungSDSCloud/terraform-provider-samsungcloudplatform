terraform {
  required_providers {
    scp = {
      version = "3.7.0"
      source  = "SamsungSDSCloud/samsungcloudplatform"
    }
  }
  required_version = ">= 0.13"
}

provider "scp" {
}
