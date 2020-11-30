package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"http_multiplexer/internal/config"
	"http_multiplexer/internal/urlsclient"
)

const shutdownTimeout = 20 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	requestInProgress := make(chan struct{}, config.MaxRequests())

	mux := http.NewServeMux()

	mux.Handle("/urls", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				sendError(w, http.StatusMethodNotAllowed, "")
				return
			}

			// if len(requestInProgress) == config.MaxRequests() {
			// 	sendError(w, http.StatusTooManyRequests, "")
			// 	return
			// }

			requestInProgress <- struct{}{}

			defer func() {
				<-requestInProgress
			}()

			next.ServeHTTP(w, r)
		})
	}(http.HandlerFunc(urlsHandler)))

	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%s", config.Port()),
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	// Run server
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	log.Printf("server started at :%s port\n", config.Port())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	<-signalChan
	log.Println("server shutting down...")

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()

	if err := httpServer.Shutdown(gracefullCtx); err != nil {
		log.Printf("shutdown error: %v\n", err)
		// nolint:gocritic
		os.Exit(1)
	}

	log.Println("gracefully stopped")

	cancel()
}

func urlsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		result map[string]interface{}
		urls   []string
	)

	err := json.NewDecoder(r.Body).Decode(&urls)

	switch {
	case err != nil:
		sendError(w, http.StatusBadRequest, err.Error())
		return
	case len(urls) > config.MaxUrls():
		sendError(w, http.StatusBadRequest, fmt.Sprintf("the amount of URLs cannot be more than %d", config.MaxUrls()))
		return
	case len(urls) == 0:
		sendResponse(w, http.StatusOK, map[string]interface{}{})
		return
	}

	result, err = urlsclient.Get(r.Context(), urls, config.MaxOutRequests())
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, result)
}

func sendResponse(w http.ResponseWriter, statusCode int, body map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if message == "" {
		message = http.StatusText(statusCode)
	}

	_, _ = w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, message)))
}
