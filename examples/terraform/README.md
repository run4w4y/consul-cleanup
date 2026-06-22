# Terraform examples

These examples are intentionally small roots that can be copied into an existing Consul/Nomad/Vault setup.

- `vault-dynamic-credentials` matches the deployment pattern used by the Nomad Pack: Vault mints short-lived Consul and Nomad tokens for the job.
- `static-consul-token` is a smaller fallback for clusters that do not use Vault. It only creates the Consul side of the credentials.

Both examples assume the Consul, Nomad, and Vault services already exist.

