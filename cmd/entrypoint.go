package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/run4w4y/consul-cleanup/common"
	"github.com/run4w4y/consul-cleanup/server"

	"github.com/urfave/cli/v2"
)

func run(cCtx *cli.Context) error {
	config, err := initConfig(cCtx)

	if err != nil {
		return err
	}

	ctx := bgContextWithLogger()

	return common.OneshotCleanup(ctx, common.OneshotCleanupConfig{
		CleanupConfig: *config,
		ServiceName:   cCtx.String("service"),
	})
}

func serve(cCtx *cli.Context) error {
	config, err := initConfig(cCtx)

	if err != nil {
		return err
	}

	serverConfig := server.ServerCleanupConfig{
		CleanupConfig: *config,
		Port:          cCtx.Uint("port"),
		AccessToken:   cCtx.String("access-token"),
	}

	logger := initLogger()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	ctx = logger.WithContext(ctx)
	defer stop()

	e, srv := server.CreateEchoWithServer(
		logger.With().Str("component", "server").Logger().WithContext(ctx),
		serverConfig,
	)

	// start the http server
	go func() {
		if err := e.StartServer(srv); err != nil && err != http.ErrServerClosed {
			logger.Fatal().
				AnErr("error", err).
				Msg("Error starting HTTP listener")
		}
	}()

	// start the nomad event stream listener
	if !cCtx.Bool("disable-events") {
		go func() {
			ctx := logger.With().Str("component", "events").Logger().WithContext(ctx)
			if err := common.CleanupListenerTask(ctx, *config); err != nil {
				logger.Fatal().
					AnErr("error", err).
					Msg("Error while running the Nomad EventStream listener")
			}
		}()
	}

	// start the periodic cleanups task
	if !cCtx.Bool("disable-periodic") {
		go func() {
			ctx := logger.With().Str("component", "periodic").Logger().WithContext(ctx)
			config := common.PeriodicCleanupConfig{
				CleanupConfig: serverConfig.CleanupConfig,
				Interval:      cCtx.Uint("interval"),
			}

			if err := common.PeriodicCleaningTask(ctx, config); err != nil {
				logger.Fatal().
					AnErr("error", err).
					Msg("Error while running periodic cleanup task")
			}
		}()
	}

	<-ctx.Done()

	logger.Info().Msg("Attempting graceful shutdown, Ctrl+C to force")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx = logger.WithContext(ctx)
	defer cancel()

	// trigger echo graceful shutdown
	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal().
			AnErr("error", err).
			Msg("Error while shutting down")
	}

	return nil
}

func events(cCtx *cli.Context) error {
	config, err := initConfig(cCtx)

	if err != nil {
		return err
	}

	ctx := bgContextWithLogger()

	return common.CleanupListenerTask(ctx, *config)
}

func periodic(cCtx *cli.Context) error {
	config, err := initConfig(cCtx)

	if err != nil {
		return err
	}

	ctx := bgContextWithLogger()

	return common.PeriodicCleaningTask(ctx, common.PeriodicCleanupConfig{
		CleanupConfig: *config,
		Interval:      cCtx.Uint("interval"),
	})
}

func Entrypoint() {
	intervalFlag := &cli.UintFlag{
		Name:  "interval",
		Usage: "Interval with which cleanup jobs should be run in seconds",
		Value: 30,
	}

	app := cli.App{
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Run the cleanup routine once",
				Action: run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "service",
						Usage: "Name of the service to query Consul specifically for",
					},
				},
			},
			{
				Name:   "serve",
				Usage:  "Run the HTTP server",
				Action: serve,
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:    "port",
						Usage:   "Port to bind the HTTP listener to",
						Value:   uint(8080),
						EnvVars: []string{"PORT"},
					},
					&cli.StringFlag{
						Name:    "access-token",
						Usage:   "The token which the server will be using for authenticating incoming requests",
						EnvVars: []string{"CLEANUP_ACCESS_TOKEN"},
					},
					&cli.BoolFlag{
						Name:  "disable-events",
						Usage: "If the flag is present the Nomad EventStream listener won't start alongside the HTTP server",
					},
					&cli.BoolFlag{
						Name:  "disable-periodic",
						Usage: "If the flag is present the periodic cleanups will be disabled",
					},
					intervalFlag,
				},
			},
			{
				Name:   "events",
				Usage:  "Listen to the EventStream from Nomad to look for allocations",
				Action: events,
			},
			{
				Name:   "periodic",
				Usage:  "Run cleanups in the background with a set time interval between each other",
				Action: periodic,
				Flags: []cli.Flag{
					intervalFlag,
				},
			},
		},
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:  "max-concurrent",
				Usage: "Maximum number of concurrent requests to the Nomad and Consul",
				Value: int64(5),
			},
			&cli.StringFlag{
				Name:    "consul-addr",
				Usage:   "Consul address",
				Value:   "127.0.0.1:8500",
				EnvVars: []string{"CONSUL_HTTP_ADDR"},
			},
			&cli.StringFlag{
				Name:    "consul-token",
				Usage:   "Consul token",
				EnvVars: []string{"CONSUL_HTTP_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "nomad-addr",
				Usage:   "Nomad address",
				Value:   "http://127.0.0.1:4646/",
				EnvVars: []string{"NOMAD_ADDR"},
			},
			&cli.StringFlag{
				Name:    "nomad-token",
				Usage:   "Nomad token's SecretID",
				EnvVars: []string{"NOMAD_TOKEN"},
			},
			&cli.StringFlag{
				Name:  "nomad-service-prefix",
				Usage: "ServiceID prefix for services registered within Consul by Nomad (without the trailing dash)",
				Value: "_nomad-task",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
