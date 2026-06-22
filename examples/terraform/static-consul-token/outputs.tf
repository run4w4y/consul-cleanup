output "consul_cleanup_token" {
  sensitive = true
  value     = data.consul_acl_token_secret_id.consul_cleanup
}

