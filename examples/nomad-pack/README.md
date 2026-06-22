# Job: consul-cleanup

Sanitized Nomad Pack example for running `consul-cleanup` as a Nomad service job.

```bash
nomad-pack render examples/nomad-pack \
  --parser-v1 \
  -f examples/nomad-pack/vars/example.hcl \
  --var docker_image_tag=demo
```

See `../../docs/acl.md` for the Consul, Nomad, and Vault policy pieces expected by the template.

