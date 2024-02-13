package cli

import (
	"context"
	"github.com/mercury-labs/telemetry-server/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var grpcCmd = &cobra.Command{
	Use:   "start",
	Short: "start",
	Long:  `start server`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		var (
			listenWebAddr = viper.GetString("lister_address")
		)
		if listenWebAddr == "" {
			listenWebAddr = "0.0.0.0:8080"
		}

		handlers := server.NewHandlers()

		log.Info().Str("http", listenWebAddr).Msg("listening")
		httpServer := http.Server{
			Addr:    listenWebAddr,
			Handler: handlers.Router(),
		}

		wg, ctx := errgroup.WithContext(ctx)

		wg.Go(func() error {
			return httpServer.ListenAndServe()
		})

		<-ctx.Done()
		log.Info().Msg("graceful shutdown")

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("shutdown")
		}
		log.Info().Msg("wait for all services being stopped")
		if err := wg.Wait(); err != nil {
			log.Error().Err(err).Msg("wg error")
		}
		log.Info().Msg("all services stopped")
	},
}

func init() {
	RootCmd.AddCommand(grpcCmd)
}
