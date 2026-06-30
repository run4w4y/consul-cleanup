# Job: consul-cleanup

Example Nomad Pack for running `consul-cleanup` as a Nomad service job.

```bash
nomad-pack render examples/nomad-pack \
  --parser-v1 \
  -f examples/nomad-pack/vars/example.hcl
```

By default the pack uses `ghcr.io/run4w4y/consul-cleanup:latest`. Override `docker_image_registry_url`, `docker_image_name`, or `docker_image_tag` if you publish the image elsewhere or want a specific release tag.

See `../../docs/acl.md` for the Consul, Nomad, and Vault policy pieces expected by the template.
