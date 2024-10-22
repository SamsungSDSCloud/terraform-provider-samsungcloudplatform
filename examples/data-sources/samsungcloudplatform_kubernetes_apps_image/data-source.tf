# Find all kubernetes app images
data "samsungcloudplatform_kubernetes_apps_image" "my_k8s_apps_image" {
  filter {
    name      = "image_name"
    values    = ["Core *"]
    use_regex = true
  }
}

output "result_scp_kubernetes_apps_image" {
  value = data.samsungcloudplatform_kubernetes_apps_image.my_k8s_apps_image
}


