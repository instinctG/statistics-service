package statistics

import (
	"context"
	"fmt"
)

type Store interface {
	GetOrderBook(ctx context.Context, exchangeName string, pair string) ([]*DepthOrder, []*DepthOrder, error)
	SaveOrderBook(ctx context.Context, orderBook *OrderBook) error
	GetOrderHistory(ctx context.Context, client *Client) ([]*HistoryOrder, error)
	SaveOrder(ctx context.Context, client *Client, order *HistoryOrder) error
}

type Service struct {
	Store Store
}

// NewService создает новый экземпляр Service.
func NewService(store Store) *Service {
	return &Service{Store: store}
}

// GetOrderBook получает книгу ордеров для указанной биржи и пары.
func (s *Service) GetOrderBook(ctx context.Context, exchangeName string, pair string) ([]*DepthOrder, []*DepthOrder, error) {
	fmt.Println("retrieving order book")

	asks, bids, err := s.Store.GetOrderBook(ctx, exchangeName, pair)
	if err != nil {
		fmt.Println(err)
		return []*DepthOrder{}, []*DepthOrder{}, err
	}

	return asks, bids, nil
}

// SaveOrderBook сохраняет книгу ордеров.
func (s *Service) SaveOrderBook(ctx context.Context, orderBook *OrderBook) error {
	err := s.Store.SaveOrderBook(ctx, orderBook)
	if err != nil {
		return fmt.Errorf("error saving order book: %v", err)
	}
	return nil
}

// GetOrderHistory получает историю ордеров для указанного клиента.
func (s *Service) GetOrderHistory(ctx context.Context, client *Client) ([]*HistoryOrder, error) {
	history, err := s.Store.GetOrderHistory(ctx, client)
	if err != nil {
		return []*HistoryOrder{}, err
	}

	return history, nil
}

// SaveOrder сохраняет ордер для указанного клиента.
func (s *Service) SaveOrder(ctx context.Context, client *Client, order *HistoryOrder) error {
	err := s.Store.SaveOrder(ctx, client, order)
	if err != nil {
		return fmt.Errorf("error saving order: %v", err)
	}

	return nil
}
