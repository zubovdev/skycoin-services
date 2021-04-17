package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// addr which will be tested.
var addr string

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0:8888", "Addr to start HTTP server on.")
	flag.Parse()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", pingHandler).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	_ = srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

// pingHandler handles ping request.
func pingHandler(rw http.ResponseWriter, r *http.Request) {
	var req struct {
		ID   int `json:"id"`
		Rand int `json:"rand"`
	}

	defer r.Body.Close()
	_ = json.NewDecoder(r.Body).Decode(&req)

	req.Rand++

	rw.Header().Set("Content-Type", "application/json")

	b, _ := json.Marshal(&req)
	_, _ = rw.Write(b)
}
