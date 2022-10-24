output "id" {
  value = scp_lb_profile.my_lb_profile_app_l4.id
}

output "persistence_id" {
  value = scp_lb_profile.my_lb_profile_persistence.id
}

output "l7_id" {
  value = scp_lb_profile.my_lb_profile_app_l7.id
}
