---
page_title: "scp_kubernetes_engine Resource - scp"
subcategory: ""
description: |-
  Provides a K8s Engine resource.
---

# Resource: scp_kubernetes_engine

Provides a K8s Engine resource.


## Example Usage

```terraform
data "scp_region" "region" {
}

resource "scp_kubernetes_engine" "engine" {
  name               = var.name
  kubernetes_version = "v1.21.8"

  vpc_id            = data.terraform_remote_state.vpc.outputs.id
  subnet_id         = data.terraform_remote_state.subnet.outputs.id
  security_group_id = data.terraform_remote_state.security-group.outputs.id
  volume_id         = data.terraform_remote_state.file-storage.outputs.id

  cloud_logging_enabled = false
  load_balancer_id      = data.terraform_remote_state.load_balancer.outputs.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `kubernetes_version` (String) Kubernetes version (Contact administrator to check supported version)
- `name` (String) Kubernetes engine name
- `security_group_id` (String) Security group ID
- `subnet_id` (String) Subnet ID
- `volume_id` (String) File storage volume ID
- `vpc_id` (String) VPC ID

### Optional

- `cloud_logging_enabled` (Boolean) Enable cloud logging
- `load_balancer_id` (String) Load balancer ID
- `public_acl_ip_address` (String) List of comma separated IP addresses (CIDR or Single IP) for access control

### Read-Only

- `id` (String) The ID of this resource.
- `kube_config` (String) Kube config of the kubernetes cluster
- `public_endpoint` (String) Public endpoint URL for the kubernetes cluster
