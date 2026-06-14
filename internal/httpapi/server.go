package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RunAwareServer struct {
	server *http.Server
}

func NewServer(router *gin.Engine, listenAddr string) *RunAwareServer {
	return &RunAwareServer{
		server: &http.Server{
			Addr:    listenAddr,
			Handler: router,
		},
	}
}

func (server *RunAwareServer) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.server.Shutdown(shutdownCtx)
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
