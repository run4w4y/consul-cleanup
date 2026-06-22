package common

import (
	capi "github.com/hashicorp/consul/api"
	napi "github.com/hashicorp/nomad/api"
)

type ApiClients struct {
	Nomad  *napi.Client
	Consul *capi.Client // this has locks internally so we pass it by reference
}

type CleanupConfig struct {
	Clients            ApiClients
	NomadServicePrefix string
	MaxConcurrent      int64
}

type OneshotCleanupConfig struct {
	CleanupConfig
	ServiceName string
}

type PeriodicCleanupConfig struct {
	CleanupConfig
	Interval uint // in seconds
}
