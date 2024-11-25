package lifecycle

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/safeblock-dev/wr/taskgroup"
)

const httpServerShutdownTTL = 10 * time.Second

func HTTPServer(server *http.Server) (taskgroup.ExecuteFn, taskgroup.InterruptFn) {
	execute := func() error {
		log.Debug().Msgf("HTTP server starting at addr: %s...", server.Addr)
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		log.Debug().Msg("HTTP server finished")

		return err
	}
	interrupt := func(_ error) {
		ctx, cancel := context.WithTimeout(context.Background(), httpServerShutdownTTL)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Error().Msgf("shutdown server err: %v", err)
		}
		log.Debug().Msg("HTTP server shutdown complete")
	}

	return execute, interrupt
}
