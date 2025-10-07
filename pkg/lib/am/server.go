package am

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Name    string
	Addr    string
	Handler *chi.Mux
}

func StartServers(servers []Server, logger Logger) {
	for _, server := range servers {
		go func(srv Server) {
			logger.Info(fmt.Sprintf("Starting %s server on %s", srv.Name, srv.Addr))
			httpSrv := &http.Server{
				Addr:    srv.Addr,
				Handler: srv.Handler,
			}
			if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error(fmt.Sprintf("%s server error: %v", srv.Name, err))
			}
		}(server)
	}
}

func GracefulShutdown(servers []Server, stops []func(context.Context) error, logger Logger) {
	logger.Info("Shutting down gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, server := range servers {
		httpSrv := &http.Server{
			Addr:    server.Addr,
			Handler: server.Handler,
		}
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			logger.Error(fmt.Sprintf("%s server shutdown error: %v", server.Name, err))
		}
	}

	for i := len(stops) - 1; i >= 0; i-- {
		if err := stops[i](shutdownCtx); err != nil {
			logger.Error(fmt.Sprintf("component shutdown error: %v", err))
		}
	}

	logger.Info("Shutdown complete")
}
