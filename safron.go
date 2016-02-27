package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
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

func withLogging(l *log.Logger, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		l.WithFields(log.Fields{
			"method": r.Method,
			"path": r.URL.Path,
			"duration": time.Since(startTime),
		}).Info()
	}
}

func main() {
	flag.Parse()

	logger := &log.Logger{
		Out:   os.Stderr,
		Level: log.InfoLevel,
	}
	switch *logFormat {
	case "json":
		logger.Formatter = new(log.JSONFormatter)
	default:
		logger.Formatter = new(log.TextFormatter)
	}

	absRoot, err := filepath.Abs(*relRoot)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	if _, err = os.Stat(absRoot); err != nil {
		logger.Fatalf("%v", err)
	}

	if !*quiet {
		logger.WithField("type", "banner").Infof("Safron version %d", Version)
		logger.WithField("type", "banner").Infof("Listening to http://%s:%d", *listenHost, *listenPort)
		logger.WithField("type", "banner").Infof("Serving %s", absRoot)
	}

	listen := *listenHost + ":" + strconv.Itoa(*listenPort)
	handler := withLogging(logger, http.FileServer(http.Dir(absRoot)))

	if err = http.ListenAndServe(listen, handler); err != nil {
		logger.Fatalf("%v", err)
	}
}
