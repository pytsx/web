package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Addr         string
	Handler      http.Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (s Server) ListenAndServe(ctx context.Context) error {

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         s.Addr,
		Handler:      s.Handler,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		IdleTimeout:  s.IdleTimeout,
	}

	srverrs := make(chan error, 1)
	go func() {
		srverrs <- srv.ListenAndServe()
	}()

	select {
	case err := <-srverrs:
		return err
	case sig := <-shutdown:
		fmt.Printf("starting graceful shutdown - signal: %v\n", sig)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("Cannot shutdown server gracefully")

			srv.Close()

			return err
		}

		fmt.Printf("gracefully shutdown - signal: %v\n", sig)
	}

	return nil
}
