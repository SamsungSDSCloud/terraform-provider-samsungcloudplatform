data "scp_redis_list" "my_scp_redis_list" {
}

output "output_my_scp_redis_list" {
  value = data.scp_redis_list.my_scp_redis_list
}
