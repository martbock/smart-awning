package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"smart-awning/awning"
	"sync"
	"time"
)

func Start(commandChan chan<- awning.Command, wg *sync.WaitGroup) (*context.CancelFunc, error) {
	router := mux.NewRouter()
	handler := http.HandlerFunc(makeHandler(commandChan))
	router.Handle("/", basicAuthMiddleware(handler)).Methods("POST")

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}

		<-ctx.Done()

		_ = srv.Shutdown(ctx)
		wg.Done()
	}()

	return &cancel, nil
}

func makeHandler(commandChan chan<- awning.Command) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		d := json.NewDecoder(req.Body)
		d.DisallowUnknownFields()

		body := struct {
			Command *string `json:"command"`
		}{}

		err := d.Decode(&body)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		if body.Command == nil {
			http.Error(w, "Missing field 'command' in JSON request", http.StatusBadRequest)
			return
		}

		switch *body.Command {
		case "extend":
			commandChan <- awning.Extend
		case "retract":
			commandChan <- awning.Retract
		case "stop":
			commandChan <- awning.Stop
		default:
			http.Error(w, "Unknown command", http.StatusBadRequest)
		}

		_, _ = fmt.Fprintln(w, "OK")
	}
}
