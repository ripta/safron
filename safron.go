package main

import (
	"flag"
	"log"
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
	quiet      = flag.Bool("quiet", false, "Quiet")
)

const Version = 1

func withLogging(l *log.Logger, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        startTime := time.Now()
		h.ServeHTTP(w, r)
		l.Printf("%s %s (%v)\n", r.Method, r.URL, time.Since(startTime))
	}
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	flag.Parse()

	absRoot, err := filepath.Abs(*relRoot)
	if err != nil {
		logger.Fatalf("%v\n", err)
	}

	if _, err = os.Stat(absRoot); err != nil {
		logger.Fatalf("%v\n", err)
	}

	if !*quiet {
		logger.Printf("Safron version %d\n", Version)
		logger.Printf("Listening to http://%s:%d\n", *listenHost, *listenPort)
		logger.Printf("Serving %s\n", absRoot)
		logger.Printf("^C to exit\n")
	}

	listen := *listenHost + ":" + strconv.Itoa(*listenPort)
	handler := withLogging(logger, http.FileServer(http.Dir(absRoot)))

	if err = http.ListenAndServe(listen, handler); err != nil {
		logger.Fatalf("%v\n", err)
	}
}
