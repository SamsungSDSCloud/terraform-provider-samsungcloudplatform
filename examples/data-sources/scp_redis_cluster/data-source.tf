data "scp_redis_cluster" "my_scp_redis_cluster" {
  redis_cluster_id = "SERVICE-123456789"
}

output "output_my_scp_redis_cluster" {
  value = data.scp_redis_cluster.my_scp_redis_cluster
}
