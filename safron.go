package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	relRoot    = flag.String("path", ".", "Document root")
	listenPort = flag.Int("port", 8080, "Port")
	listenHost = flag.String("host", "0.0.0.0", "Host")
	logFormat  = flag.String("log-format", "text", "Log format: text, json")
	quiet      = flag.Bool("quiet", false, "No banner on startup")
)

const Version = 2

func withLogging(l *slog.Logger, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		l.Info("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", time.Since(startTime)),
		)
	}
}

func main() {
	flag.Parse()

	var handler slog.Handler
	switch *logFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, nil)
	default:
		handler = slog.NewTextHandler(os.Stderr, nil)
	}
	logger := slog.New(handler)

	if args := flag.Args(); len(args) > 0 {
		logger.Error("unexpected positional arguments", slog.Any("args", args))
		os.Exit(1)
	}

	absRoot, err := filepath.Abs(*relRoot)
	if err != nil {
		logger.Error("failed to get absolute path", slog.Any("error", err))
		os.Exit(1)
	}

	if _, err = os.Stat(absRoot); err != nil {
		logger.Error("failed to stat root directory", slog.Any("error", err))
		os.Exit(1)
	}

	if !*quiet {
		logger.Info("Safron version", slog.Int("version", Version), slog.String("type", "banner"))
		logger.Info("Listening", slog.String("address", "http://"+*listenHost+":"+strconv.Itoa(*listenPort)), slog.String("type", "banner"))
		logger.Info("Serving", slog.String("path", absRoot), slog.String("type", "banner"))
	}

	listen := *listenHost + ":" + strconv.Itoa(*listenPort)
	handlerFunc := withLogging(logger, http.FileServer(http.Dir(absRoot)))

	if err = http.ListenAndServe(listen, handlerFunc); err != nil {
		logger.Error("server error", slog.Any("error", err))
		os.Exit(1)
	}
}
