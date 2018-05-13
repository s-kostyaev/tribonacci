package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/s-kostyaev/tribonacci/internal/health"

	"github.com/s-kostyaev/tribonacci/pkg/handler"

	"github.com/agalitsyn/flagenv"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/s-kostyaev/tribonacci/pkg/tribonacci"
)

// Version will be set during build process
var Version string

func main() {
	cfg := parseFlags()
	log.Printf("started with config: %+v", *cfg)

	// note: order of middlewares is important
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
	)
	r.Mount("/readiness", health.Routes())
	r.Route("/v1", func(r chi.Router) {
		r.HandleFunc("/tribonacci", tribonacciHandler)
	})
	handler.FileServer(r, "/docs", http.Dir(cfg.DocsPath))
	srv := &http.Server{Addr: cfg.ListenAddr, Handler: r}

	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigquit
		log.Printf("captured %v, exiting...", s)

		health.SetReadinessStatus(http.StatusServiceUnavailable)

		log.Println("gracefully shutdown server")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("could not shutdown server: %v", err)
		}
	}()

	log.Printf("starting http service...")
	log.Printf("listening on %s", cfg.ListenAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("server error: %v", err)
	}
}

type cliFlags struct {
	Version    bool
	ListenAddr string
	DocsPath   string
}

func parseFlags() *cliFlags {
	var cfg cliFlags

	flag.StringVar(&cfg.ListenAddr, "listen-addr", ":8080", "HTTP service address.")
	flag.StringVar(&cfg.DocsPath, "docs-path", "docs", "Path to documentation folder.")
	flag.BoolVar(&cfg.Version, "version", false, "Service version.")

	flagenv.Prefix = "tribonacci_web_"
	flagenv.Parse()
	flag.Parse()

	if cfg.Version {
		log.Println(Version)
		os.Exit(0)
	}

	return &cfg
}

func tribonacciHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	rawN, ok := query["number"]
	if !ok || len(rawN) != 1 {
		render.Render(w, r, handler.ErrBadRequest(errors.New("number is required")))
		return
	}
	n, err := strconv.ParseUint(rawN[0], 10, 64)
	if err != nil {
		render.Render(w, r, handler.ErrBadRequest(fmt.Errorf("can't parse number param %v: %v", rawN[0], err)))
		return
	}
	if n == 0 {
		render.Render(w, r, handler.ErrBadRequest(errors.New("indexing starting with 1")))
		return
	}
	render.JSON(w, r, tribonacci.Number(n))
}
