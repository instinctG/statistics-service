package http

import (
	"context"
	"encoding/json"
	"github.com/instinctG/statistics_service/internal/statistics"
	"net/http"
)

type StatisticService interface {
	GetOrderBook(ctx context.Context, exchangeName string, pair string) ([]*statistics.DepthOrder, []*statistics.DepthOrder, error)
	SaveOrderBook(ctx context.Context, orderBook *statistics.OrderBook) error
	GetOrderHistory(ctx context.Context, client *statistics.Client) ([]*statistics.HistoryOrder, error)
	SaveOrder(ctx context.Context, client *statistics.Client, order *statistics.HistoryOrder) error
}

type Response struct {
	Message string `json:"message"`
}

// GetOrderBook обрабатывает HTTP GET запросы для получения данных order book
// на основе параметров exchange_name и pair. Функция проверяет входные параметры,
// извлекает данные с помощью StatisticService.GetOrderBook и возвращает JSON
// ответ с данными о bids и asks.
func (h *Handler) GetOrderBook(w http.ResponseWriter, r *http.Request) {

	exchangeName := r.URL.Query().Get("exchange_name")
	pair := r.URL.Query().Get("pair")

	if exchangeName == "" || pair == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	asks, bids, err := h.Service.GetOrderBook(r.Context(), exchangeName, pair)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string][]*statistics.DepthOrder{
		"asks": asks,
		"bids": bids,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		panic(err)
	}
}

// GetOrderHistory обрабатывает HTTP POST запросы для получения истории ордера
// для конкретного клиента. Функция разбирает информацию о клиенте из тела запроса,
// извлекает данные с помощью StatisticService.GetOrderHistory и возвращает JSON
// ответ с данными истории ордеров.
func (h *Handler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	var client statistics.Client

	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, "failed to parse request", http.StatusBadRequest)
		return
	}

	history, err := h.Service.GetOrderHistory(r.Context(), &client)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(history); err != nil {
		panic(err)
	}
}

// SaveOrderBook обрабатывает HTTP POST запросы для сохранения order book. Функция
// декодирует данные order book из тела запроса, сохраняет данные с помощью
// StatisticService.SaveOrderBook и возвращает JSON ответ, указывающий на успешное
// или неудачное выполнение операции.
func (h *Handler) SaveOrderBook(w http.ResponseWriter, r *http.Request) {
	var orderBook statistics.OrderBook

	if err := json.NewDecoder(r.Body).Decode(&orderBook); err != nil {
		http.Error(w, "failed to decode", http.StatusBadRequest)
		return
	}

	err := h.Service.SaveOrderBook(r.Context(), &orderBook)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Response{
		Message: "Order book saved successfully",
	}); err != nil {
		panic(err)
	}
}

// SaveOrder обрабатывает HTTP POST запросы для сохранения истории ордера.
// Функция декодирует данные истории ордера из тела запроса, извлекает информацию
// о клиенте, сохраняет заказ с помощью StatisticService.SaveOrder и возвращает
// JSON ответ, указывающий на успешное или неудачное выполнение операции.
func (h *Handler) SaveOrder(w http.ResponseWriter, r *http.Request) {
	var history statistics.HistoryOrder
	if err := json.NewDecoder(r.Body).Decode(&history); err != nil {
		http.Error(w, "failed to decode history order", http.StatusBadRequest)
		return
	}

	client := statistics.Client{
		ClientName:   history.ClientName,
		ExchangeName: history.ExchangeName,
		Pair:         history.Pair,
		Label:        history.Label,
	}

	if err := h.Service.SaveOrder(r.Context(), &client, &history); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Response{
		Message: "Order saved successfully",
	}); err != nil {
		panic(err)
	}

}
