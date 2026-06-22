# consul-cleanup

## Description
A set of strategies to find and deregister service entries in Consul that are associated with a dead Nomad allocation.

Can be run as just a CLI tool, as a continuous task, as an HTTP server.

## Issues to watch
This will become redundant once Nomad can deregister those services on itself. Issues are being tracked here:
- https://github.com/hashicorp/nomad/issues/23534
- https://github.com/hashicorp/nomad/issues/16616
- https://github.com/hashicorp/nomad/issues/23494

## Building
You can build it the same way you would build any other Go application, for the specifics refer to the `Dockerfile`

With the included Nix shell:
```bash
direnv allow
go build ./...
```

## Running with Docker
```bash
docker build -t consul-cleanup .
docker run consul-cleanup # this will display the help message
```

## Local development
This repository includes a Nix flake and `.envrc` for a reproducible development shell.

Useful commands:
```bash
go test ./...
go vet ./...
golangci-lint run ./...
tofu fmt -recursive examples/terraform
act
```

`act` is included in the shell so the GitHub Actions workflow can be exercised locally.

## ACL requirements
- Consul token ACL requirements: `service:write`, `node:write`
- Nomad token ACL requirements: `namespace:read-job`

### Example Nomad ACL policy
```hcl
namespace "*" { // you can limit this to just one namespace
    policy = "read"
    capabilities = ["read-job"]
}
```

### Example Consul ACL policy
```hcl
service_prefix "" { // you can limit that to only the nomad prefix
    policy = "write"
}

node_prefix "" {
    policy = "write"
}
```

## Available strategies/modes

### Default strategy (Consul -> Nomad -> Consul)
- Query Consul for all of the registered services (that are related a Nomad allocation)
- For each service check whether its related allocation is still relevant with Nomad
- Deregister from Consul services that failed to pass the check above

### Reactive strategy (Nomad -> Consul -> Consul)
- Attach to the Nomad Event Stream and listen for the Allocation topic
- Whenever an allocation stops running, find services related to it in Consul
- Deregister the found services from Consul

### CLI subcommands
- `run` - performs a cleanup using the default strategy once
- `serve` - starts an HTTP server that can perform default strategy cleanups on a request
- `events` - a continuous task which implements the reactive strategy described above
- `periodic` - continuously performs cleanups using the default strategy with a set time interval between each

For more information about flags and each subcommand refer to the CLI `--help` (works with a subcommand as well)

## HTTP server
By default, `events` and `periodic` tasks are going to be run alongside the HTTP server. To disable either, refer to the `serve` subcommand `--help`.

Available endpoints are:
- GET `/api/v1/health` - Healthcheck endpoint
- POST `/api/v1/oneshot` - performs default strategy for all of the services
- POST `/api/v1/oneshot/:service` - performs default strategy for the service specified in the path

When `CLEANUP_ACCESS_TOKEN` or `--access-token` is set, the `oneshot` endpoints expect:
```bash
Authorization: Bearer <token>
```

## Demo and deployment examples
The `examples/nomad-pack` directory contains a sanitized Nomad Pack based on the original deployment setup.

The `examples/terraform` directory contains focused Terraform/OpenTofu examples for the required Consul, Nomad, and Vault ACL pieces.

For the Consul, Nomad, and Vault ACL pieces needed to run it with least-privilege tokens, see `docs/acl.md`.

For a short walkthrough of the demo deployment shape, see `docs/demo.md`.
