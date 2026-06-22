package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/run4w4y/consul-cleanup/common"

	capi "github.com/hashicorp/consul/api"
	napi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

// initialize both consul and nomad api clients
func initClients(cCtx *cli.Context) (*common.ApiClients, error) {
	// init consul client
	consul, err := capi.NewClient(&capi.Config{
		Address: cCtx.String("consul-addr"),
		Token:   cCtx.String("consul-token"),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize Consul API client: %w", err)
	}

	// init nomad client
	nomad, err := napi.NewClient(&napi.Config{
		Address:  cCtx.String("nomad-addr"),
		SecretID: cCtx.String("nomad-token"),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize Nomad API client: %w", err)
	}

	return &common.ApiClients{Nomad: nomad, Consul: consul}, nil
}

// initialize the shared base config from the cli
func initConfig(cCtx *cli.Context) (*common.CleanupConfig, error) {
	clients, err := initClients(cCtx)
	if err != nil {
		return nil, err
	}

	return &common.CleanupConfig{
		Clients:            *clients,
		MaxConcurrent:      cCtx.Int64("max-concurrent"),
		NomadServicePrefix: cCtx.String("nomad-service-prefix"),
	}, nil
}

func initLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	return logger
}

func bgContextWithLogger() context.Context {
	ctx := context.Background()
	logger := initLogger()

	return logger.WithContext(ctx)
}
