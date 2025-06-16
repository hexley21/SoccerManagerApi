package shutdown

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/hexley21/soccer-manager/pkg/logger"
)

func NotifyShutdown(
	ctx context.Context,
	ctxClose context.CancelCauseFunc,
	logger logger.Logger,
	closer io.Closer,
) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		signal.Stop(sigCh)
		logger.Info("caught signal", "signal", sig.String())
		ctxClose(closer.Close())
	case <-ctx.Done():
		return
	}
}
