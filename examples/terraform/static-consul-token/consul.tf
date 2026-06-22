resource "consul_acl_policy" "consul_cleanup" {
  name = "consul-cleanup"

  rules = <<-HCL
    service_prefix "" {
      policy = "write"
    }

    node_prefix "" {
      policy = "write"
    }
  HCL
}

resource "consul_acl_token" "consul_cleanup" {
  description = var.token_description
  policies    = [consul_acl_policy.consul_cleanup.name]
  local       = var.token_local
}

data "consul_acl_token_secret_id" "consul_cleanup" {
  accessor_id = consul_acl_token.consul_cleanup.accessor_id
}

