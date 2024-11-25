package lifecycle

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/safeblock-dev/werr"
	"github.com/safeblock-dev/wr/taskgroup"
)

type WorkerRunner interface {
	Run(ctx context.Context) error
}

func Worker(
	worker WorkerRunner,
) (taskgroup.ExecuteFn, taskgroup.InterruptFn) {
	ctx, cancel := context.WithCancel(context.Background())

	execute := func() error {
		log.Debug().Msg("worker starting...")
		err := worker.Run(ctx)
		if errors.Is(err, context.Canceled) {
			err = nil
		}
		log.Debug().Msg("worker finished")

		return werr.Wrap(err)
	}

	interrupt := func(_ error) {
		cancel()

		log.Debug().Msg("worker shutdown complete")
	}

	return execute, interrupt
}
