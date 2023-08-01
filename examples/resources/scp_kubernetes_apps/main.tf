data "scp_region" "region" {
}

data "scp_kubernetes_apps_image" "apps_image" {
  filter {
    name   = "category"
    values = ["Base"]
  }

  filter {
    name      = "image_name"
    values    = ["Alpine Community *"]
    use_regex = true
  }

  filter {
    name   = "version"
    values = ["3.13.12"]
  }
}

resource "scp_kubernetes_apps" "apps" {
  name      = var.name
  engine_id = data.terraform_remote_state.engine.outputs.id
  namespace = data.terraform_remote_state.namespace.outputs.id
  image_id  = data.scp_kubernetes_apps_image.apps_image.id
}

