package lifecycle

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/VadimOcLock/gophermart/pkg/httpmix"
	"github.com/go-chi/chi/v5"
	"github.com/hellofresh/health-go/v5"
	"github.com/matchsystems/werr"
	"github.com/rs/zerolog/log"
	"github.com/safeblock-dev/wr/taskgroup"
)

const httpServerShutdownTTL = 10 * time.Second

type MuxHandler struct {
	Handler http.Handler
	Name    string
	Path    string
}

func HTTPServer(mux *chi.Mux, port int, handlers ...MuxHandler) (taskgroup.ExecuteFn, taskgroup.InterruptFn) {
	for _, handler := range handlers {
		mux.Handle(handler.Path, handler.Handler)
	}
	server := httpmix.NewSimpleServer(mux, port)

	shutdown := atomic.Bool{}
	execute := func() error {
		l := log.Info().Int("port", port)
		for _, handler := range handlers {
			l.Str(handler.Name, handler.Path)
		}
		l.Msg("HTTP Server starting...")
		err := server.ListenAndServe()
		if shutdown.Load() && errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		log.Err(err).Int("port", port).Msg("HTTP Server finished")

		return werr.Wrap(err)
	}

	interrupt := func(_ error) {
		ctx, cancel := context.WithTimeout(context.Background(), httpServerShutdownTTL)
		defer cancel()

		shutdown.Store(true)
		err := server.Shutdown(ctx)
		log.Err(err).Int("port", port).Msg("HTTP Server shutdown complete")
	}

	return execute, interrupt
}

func HTTPHealthHandler(
	name, version string,
	drivers ...health.Config,
) http.Handler {
	hub, err := health.New(health.WithComponent(health.Component{
		Name:    name,
		Version: version,
	}))
	if err != nil {
		panic(err)
	}

	for _, driver := range drivers {
		err = hub.Register(driver)
		if err != nil {
			panic(err)
		}
	}

	return hub.Handler()
}
