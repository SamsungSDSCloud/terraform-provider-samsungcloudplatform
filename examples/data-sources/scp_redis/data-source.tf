data "scp_redis" "my_scp_redis" {
  redis_id = "SERVICE-123456789"
}

output "output_my_scp_redis" {
  value = data.scp_redis.my_scp_redis
}