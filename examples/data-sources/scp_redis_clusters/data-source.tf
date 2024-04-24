data "scp_redis_clusters" "my_scp_redis_clusters" {
}

output "output_my_scp_redis_clusters" {
  value = data.scp_redis_clusters.my_scp_redis_clusters
}
