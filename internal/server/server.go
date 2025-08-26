package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maker2413/daycount/internal/compute"
	"github.com/maker2413/daycount/internal/parse"
)

// Response is the JSON structure returned by the API.
type Response struct {
	Date      string `json:"date"`
	Reference string `json:"reference"`
	Days      int    `json:"days"`
	Phrase    string `json:"phrase"`
}

// Handler returns an http.Handler with all routes registered.
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/daycount", daycountHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK); w.Write([]byte("ok")) })
	return mux
}

func daycountHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	dateStr := q.Get("date")
	if dateStr == "" {
		http.Error(w, "missing 'date' query parameter", http.StatusBadRequest)
		return
	}

	target, err := parse.ParseDate(dateStr)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, parse.ErrUnrecognizedFormat) {
			http.Error(w, "unrecognized date format", status)
			return
		}
		http.Error(w, "parse error", status)
		return
	}

	referenceStr := q.Get("reference")
	var now time.Time
	var refOut string
	if referenceStr != "" {
		ref, err := parse.ParseDate(referenceStr)
		if err != nil {
			http.Error(w, "invalid reference date", http.StatusBadRequest)
			return
		}
		now = ref
		refOut = ref.Format(time.DateOnly)
	} else {
		now = time.Now().UTC()
		refOut = now.Format(time.DateOnly)
	}

	days := compute.DaysSince(target, now)
	phrase := "days since"
	if days < 0 {
		phrase = "days until"
	}

	if q.Get("format") == "plain" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(formatPlain(days, phrase, dateStr)))
		return
	}

	resp := Response{
		Date:      dateStr,
		Reference: refOut,
		Days:      days,
		Phrase:    formatPlain(days, phrase, dateStr),
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

func formatPlain(days int, phrase, input string) string {
	if days < 0 {
		return sprintf("%d %s %s", -days, phrase, input)
	}
	return sprintf("%d %s %s", days, phrase, input)
}

// Options allows injecting dependencies for testability.
type Options struct {
	Server  *http.Server     // pre-built server; if nil one is constructed from addr
	Signals <-chan os.Signal // custom signal channel; if nil real OS signals used
}

// Run starts the HTTP server listening on addr (fallback to :8080) and blocks until context cancelled or shutdown signal.
// Deprecated: use RunContext for dependency injection in tests.
func Run(ctx context.Context, addr string) error {
	return RunContext(ctx, addr, nil)
}

// RunContext is like Run but allows dependency injection via opts for testing.
func RunContext(ctx context.Context, addr string, opts *Options) error {
	if addr == "" {
		addr = os.Getenv("DAYCOUNT_ADDR")
		if addr == "" {
			addr = ":8080"
		}
	}

	var srv *http.Server
	if opts != nil && opts.Server != nil {
		srv = opts.Server
	} else {
		srv = &http.Server{Addr: addr, Handler: Handler()}
	}

	// signal handling for graceful shutdown
	var sigCh <-chan os.Signal
	if opts != nil && opts.Signals != nil {
		sigCh = opts.Signals
	} else {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		sigCh = ch
	}
	go func() {
		select {
		case <-ctx.Done():
		case <-sigCh:
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("http shutdown error: %v", err)
		}
	}()

	log.Printf("listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// small wrapper to allow test to stub formatting easily
var sprintf = func(format string, a ...any) string { return fmt.Sprintf(format, a...) }
