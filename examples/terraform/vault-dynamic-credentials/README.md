# Vault dynamic credentials

This example supports `examples/nomad-pack`: Terraform creates the Consul/Nomad policies, configures Vault roles for dynamic credentials, and stores the HTTP bearer token used by `consul-cleanup serve`.

It expects existing Consul, Nomad, and Vault clusters. The tokens passed to Terraform must be operator/bootstrap tokens with enough permissions to create ACL policies and configure Vault secrets backends.

## Usage

```bash
terraform init
terraform plan \
  -var consul_token="$CONSUL_HTTP_TOKEN" \
  -var nomad_token="$NOMAD_TOKEN" \
  -var vault_token="$VAULT_TOKEN"
```

With OpenTofu:

```bash
tofu init
tofu plan \
  -var consul_token="$CONSUL_HTTP_TOKEN" \
  -var nomad_token="$NOMAD_TOKEN" \
  -var vault_token="$VAULT_TOKEN"
```

The Nomad job template can then render:

- `CONSUL_HTTP_TOKEN` from `consul/creds/consul-cleanup`
- `NOMAD_TOKEN` from `nomad/creds/consul-cleanup`
- `CLEANUP_ACCESS_TOKEN` from `secret/data/consul-cleanup/root-credentials`
