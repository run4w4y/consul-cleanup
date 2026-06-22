output "consul_policy_name" {
  value = consul_acl_policy.consul_cleanup.name
}

output "nomad_policy_name" {
  value = nomad_acl_policy.allocation_observer.name
}

output "vault_policy_name" {
  value = vault_policy.consul_cleanup.name
}

output "vault_consul_role" {
  value = vault_consul_secret_backend_role.consul_cleanup.name
}

output "vault_nomad_role" {
  value = vault_nomad_secret_role.consul_cleanup.role
}

output "cleanup_access_token_path" {
  value = "${var.vault_kv_mount}/data/consul-cleanup/root-credentials"
}

