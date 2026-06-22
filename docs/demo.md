# Demo deployment

The easiest demonstration path is to run `consul-cleanup` as a Nomad service job with short-lived Consul and Nomad credentials rendered from Vault.

## Shape

1. Build and publish an image:

   ```bash
   docker build -t ghcr.io/run4w4y/consul-cleanup:demo .
   ```

2. Create the Consul, Nomad, and Vault policies described in `docs/acl.md`.

   For a Terraform/OpenTofu version of that setup, start with:

   ```bash
   cd examples/terraform/vault-dynamic-credentials
   tofu init
   tofu plan \
     -var consul_token="$CONSUL_HTTP_TOKEN" \
     -var nomad_token="$NOMAD_TOKEN" \
     -var vault_token="$VAULT_TOKEN"
   ```

3. Render the Nomad Pack example:

   ```bash
   nomad-pack render examples/nomad-pack \
     --parser-v1 \
     -f examples/nomad-pack/vars/example.hcl \
     --var docker_image_tag=demo
   ```

4. Run the rendered job in Nomad.

5. Check the HTTP server:

   ```bash
   curl http://<alloc-ip>:8989/api/v1/health
   ```

6. Trigger a one-shot cleanup:

   ```bash
   curl -X POST \
     -H "Authorization: Bearer $CLEANUP_ACCESS_TOKEN" \
     http://<alloc-ip>:8989/api/v1/oneshot
   ```

## Notes

By default, the `serve` command starts the HTTP server, the Nomad event listener, and periodic cleanup. For a quieter demo, add `-disable-events` or `-disable-periodic` to the Nomad task command arguments.

The example job uses host networking so it can address Consul and Nomad on the node IP. If your cluster exposes those APIs differently, adjust `CONSUL_HTTP_ADDR` and `NOMAD_ADDR` in the template.
