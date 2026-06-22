# ACL reference

`consul-cleanup` needs one token for Consul and one token for Nomad.

The examples below are based on the deployment Terraform setup this service was originally used with. Scope the prefixes/namespaces more tightly if your cluster layout allows it.

## Consul

The cleanup task must be able to deregister service catalog entries and write node catalog updates.

```hcl
service_prefix "" {
  policy = "write"
}

node_prefix "" {
  policy = "write"
}
```

The original Terraform shape used a Consul ACL policy named `consul-cleanup` and then exposed it through a Vault Consul secrets role:

```hcl
resource "vault_consul_secret_backend_role" "consul-cleanup" {
  name            = "consul-cleanup"
  backend         = vault_consul_secret_backend.config.path
  consul_policies = ["consul-cleanup"]
}
```

## Nomad

The cleanup task must be able to inspect allocations and jobs. The original setup used an `allocation-observer` policy:

```hcl
namespace "*" {
  policy       = "read"
  capabilities = ["read-job"]
}
```

That policy was exposed through a Vault Nomad secrets role:

```hcl
resource "vault_nomad_secret_role" "consul-cleanup" {
  backend  = vault_nomad_secret_backend.config.backend
  role     = "consul-cleanup"
  policies = ["allocation-observer"]
}
```

## Vault

For a Nomad job that renders credentials from Vault templates, the task policy needs to read the generated Consul and Nomad credentials and the HTTP access token:

```hcl
path "nomad/creds/consul-cleanup" {
  capabilities = ["read"]
}

path "consul/creds/consul-cleanup" {
  capabilities = ["read"]
}

path "secret/data/consul-cleanup/*" {
  capabilities = ["read"]
}
```

The bearer token used by the HTTP server can be stored as KV data:

```hcl
resource "vault_kv_secret_v2" "consul_cleanup_secret" {
  mount = vault_mount.kvv2.path
  name  = "consul-cleanup/root-credentials"

  data_json = jsonencode({
    access_token = random_password.consul_cleanup_password.result
  })
}
```

## Terraform examples

The `examples/terraform` directory contains focused examples extracted from the original deployment pattern.

Use `examples/terraform/vault-dynamic-credentials` when the Nomad job should render short-lived Consul and Nomad credentials from Vault. This matches the `examples/nomad-pack` template.

Use `examples/terraform/static-consul-token` only when Vault is not part of the demo. It creates the Consul policy and a static Consul token, but the Nomad token still needs to be created separately.
