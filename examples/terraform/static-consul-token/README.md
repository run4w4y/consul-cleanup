# Static Consul token

This example creates the Consul ACL policy and a static Consul token for `consul-cleanup`.

Use this only when you are not using Vault dynamic credentials. It does not create the Nomad token; create a Nomad token with the policy shown in `../../../docs/acl.md` or adapt the `vault-dynamic-credentials` example.

## Usage

```bash
terraform init
terraform apply -var consul_token="$CONSUL_HTTP_TOKEN"
terraform output -json consul_cleanup_token
```

With OpenTofu:

```bash
tofu init
tofu apply -var consul_token="$CONSUL_HTTP_TOKEN"
tofu output -json consul_cleanup_token
```
