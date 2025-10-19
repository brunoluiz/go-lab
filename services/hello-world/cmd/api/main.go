package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/brunoluiz/go-lab/services/hello-world/internal/service/greet"
)

func main() {
	greeter := greet.New()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lang := r.URL.Query().Get("lang")
		if lang == "" {
			lang = "en"
		}

		helloMsg, err := greeter.Hello(lang)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.ErrorContext(r.Context(), "unable to greet", slog.String("error", err.Error()))
		}

		if _, writeErr := w.Write([]byte(helloMsg)); writeErr != nil {
			logger.ErrorContext(r.Context(), "unable to write response", slog.String("error", writeErr.Error()))
		}
	})

	log.Println("running server at :3000")
	server := &http.Server{
		Addr:         ":3000",
		Handler:      nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
