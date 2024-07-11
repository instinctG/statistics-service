package http

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Handler struct {
	Router  *mux.Router
	Service StatisticService
	Server  *http.Server
}

// NewHandler создает новый экземпляр Handler.
func NewHandler(service StatisticService) *Handler {
	h := &Handler{Service: service}

	h.Router = mux.NewRouter()
	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: h.Router,
	}

	return h
}

// mapRoutes задает маршруты для API.
func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/api/get-order-book", h.GetOrderBook).Methods("GET")
	h.Router.HandleFunc("/api/save-order-book", h.SaveOrderBook).Methods("POST")
	h.Router.HandleFunc("/api/get-history", h.GetOrderHistory).Methods("GET")
	h.Router.HandleFunc("/api/save-history", h.SaveOrder).Methods("POST")
}

// Serve запускает HTTP-сервер и обрабатывает остановку сервера.
func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)
	log.Println("shut down gracefully")
	return nil
}
