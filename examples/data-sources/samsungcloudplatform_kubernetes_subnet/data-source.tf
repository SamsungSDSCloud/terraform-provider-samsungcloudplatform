# Check available subnet
data "samsungcloudplatform_kubernetes_subnet" "my_scp_kubernetes_subnet" {
  vpc_id    = "VPC-XXXXXXXXX"
  subnet_id = "SUBNET-XXXXXXXXX"
}

output "result_scp_kubernetes_subnet" {
  value = data.samsungcloudplatform_kubernetes_subnet.my_scp_kubernetes_subnet
}
