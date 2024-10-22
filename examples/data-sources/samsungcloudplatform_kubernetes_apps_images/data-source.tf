# Find all kubernetes app images
data "samsungcloudplatform_kubernetes_apps_images" "my_k8s_apps_images1" {
}

output "result_scp_kubernetes_apps_images1" {
  value = data.samsungcloudplatform_kubernetes_apps_images.my_k8s_apps_images1
}
