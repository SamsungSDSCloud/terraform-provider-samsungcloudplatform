output "id" {
  value = samsungcloudplatform_redis.demo_db.id
}

output "natIpAddress" {
  value = samsungcloudplatform_redis.demo_db.redis_servers[0].nat_public_ip_address
}
