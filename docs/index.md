---
page_title: "Provider: Samsung Cloud Platform"
---

# Samsungcloudplatform Provider

The Samsungcloudplatform provider is used to interact with Samsung Cloud Platform services.
The provider needs to be configured with the proper credentials before it can be used.


## Example Usage
```hcl
terraform {
  required_providers {
    samsungcloudplatform = {
      version = "3.16.1"
      source  = "SamsungSDSCloud/samsungcloudplatform"
    }
  }
  required_version = ">= 0.13"
}

provider "samsungcloudplatform" {
}

//Create a new virtual server instance
resource "samsungcloudplatform_virtual_server" "server_001" {
  image_id        = "Image ID"
  name_prefix     = "server01"
  cpu_count       = 4
  memory_size_gb  = 8
  #...
}
```


## Setup credentials

### Create local setting file

Create `.scp` directory in your OS account home

```
cd %USERPROFILE%
mkdir ".scp"
```

Create `config.json` and `credentials.json` in `.scp` directory

### Add Samsungcloudplatform configuration

Insert following parameters in `.scp/config.json` file

```
{
    "host": "https://openapi.samsungsdscloud.com",
    "user-id": "1234",
    "email" : "your.email@samsung.com",
    "project-id": "PROJECT-XXXXXXXXXXXXXXXX"
}
```

### Add your credentials

Insert following parameters in `.scp/credentials.json` file

```
{
    "auth-method": "access-key",
    "access-key": "XXXXXXXXXXXXXXXX,
    "secret-key": "XXXXXXXXXXXXXXXX"
}
```

## Open-source Software Notice

[OSS Notice Link](https://github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/blob/v3.16.1/OpenSourceNotice.docx)
