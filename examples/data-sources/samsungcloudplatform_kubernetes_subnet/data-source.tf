# Check available subnet
data "samsungcloudplatform_kubernetes_subnet" "my_scp_kubernetes_subnet" {
  vpc_id    = "vpc id"
  subnet_id = "subnet id"
}

output "result_scp_kubernetes_subnet" {
  value = data.samsungcloudplatform_kubernetes_subnet.my_scp_kubernetes_subnet
}
