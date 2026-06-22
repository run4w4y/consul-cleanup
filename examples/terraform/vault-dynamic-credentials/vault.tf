resource "vault_consul_secret_backend" "consul" {
  path    = "consul"
  address = var.consul_address
  token   = var.consul_token
}

resource "vault_consul_secret_backend_role" "consul_cleanup" {
  name            = "consul-cleanup"
  backend         = vault_consul_secret_backend.consul.path
  consul_policies = [consul_acl_policy.consul_cleanup.name]
}

resource "vault_nomad_secret_backend" "nomad" {
  backend = "nomad"
  address = var.nomad_address
  token   = var.nomad_token
}

resource "vault_nomad_secret_role" "consul_cleanup" {
  backend  = vault_nomad_secret_backend.nomad.backend
  role     = "consul-cleanup"
  policies = [nomad_acl_policy.allocation_observer.name]
}

resource "vault_policy" "consul_cleanup" {
  name = "consul-cleanup"

  policy = <<-HCL
    path "nomad/creds/consul-cleanup" {
      capabilities = ["read"]
    }

    path "consul/creds/consul-cleanup" {
      capabilities = ["read"]
    }

    path "${var.vault_kv_mount}/data/consul-cleanup/*" {
      capabilities = ["read"]
    }
  HCL
}

resource "random_password" "consul_cleanup_access_token" {
  length  = var.cleanup_access_token_length
  special = false
}

resource "vault_kv_secret_v2" "consul_cleanup" {
  mount = var.vault_kv_mount
  name  = "consul-cleanup/root-credentials"

  data_json = jsonencode({
    access_token = random_password.consul_cleanup_access_token.result
  })

  lifecycle {
    ignore_changes = [data_json]
  }
}

