package db

import (
	"context"
	"fmt"
	"github.com/instinctG/statistics_service/internal/statistics"
)

// GetOrderBook получает книгу ордеров для указанной биржи и пары из базы данных.
func (d *Database) GetOrderBook(ctx context.Context, exchangeName string, pair string) ([]*statistics.DepthOrder, []*statistics.DepthOrder, error) {
	var asks, bids []*statistics.DepthOrder
	row := d.Client.QueryRow(
		ctx,
		`SELECT asks,bids FROM OrderBook WHERE exchange = $1 AND pair = $2`,
		exchangeName,
		pair,
	)
	err := row.Scan(&asks, &bids)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting orderbook asks and bids: %w", err)
	}

	return asks, bids, nil

}

// SaveOrderBook сохраняет книгу ордеров для указанной биржи и пары в базе данных.
func (d *Database) SaveOrderBook(ctx context.Context, orderBook *statistics.OrderBook) error {
	_, err := d.Client.Exec(
		ctx,
		`INSERT INTO OrderBook (exchange, pair, asks, bids) VALUES ($1, $2, $3, $4)`,
		orderBook.Exchange,
		orderBook.Pair,
		orderBook.Asks,
		orderBook.Bids,
	)
	if err != nil {
		return fmt.Errorf("error saving orderbook: %w", err)
	}

	return nil
}

// GetOrderHistory получает историю ордеров для указанного клиента из базы данных.
func (d *Database) GetOrderHistory(ctx context.Context, client *statistics.Client) ([]*statistics.HistoryOrder, error) {
	var historyOrders []*statistics.HistoryOrder
	rows, err := d.Client.Query(
		ctx,
		` SELECT client_name, exchange_name, label, pair, side, type, base_qty, price, 
           algorithm_name_placed, lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed
    	FROM HistoryOrder
    	WHERE client_name = $1 AND exchange_name = $2 AND label = $3 AND pair = $4 `,
		client.ClientName,
		client.ExchangeName,
		client.Label,
		client.Pair,
	)
	if err != nil {
		return nil, fmt.Errorf("could not query order history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order statistics.HistoryOrder
		err := rows.Scan(&order.ClientName, &order.ExchangeName, &order.Label, &order.Pair, &order.Side, &order.TypeOrder,
			&order.BaseQty, &order.Price, &order.AlgorithmNamePlaced, &order.LowestSellPrc, &order.HighestBuyPrc,
			&order.CommissionQuoteQty, &order.TimePlaced)
		if err != nil {
			return nil, fmt.Errorf("could not scan order history row: %w", err)
		}
		historyOrders = append(historyOrders, &order)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", rows.Err())
	}

	return historyOrders, nil
}

// SaveOrder сохраняет ордер и информацию о клиенте в базе данных.
func (d *Database) SaveOrder(ctx context.Context, client *statistics.Client, order *statistics.HistoryOrder) error {
	_, err := d.Client.Exec(
		ctx,
		` INSERT INTO Client (
            		client_name, exchange_name, label, pair
        	  ) VALUES ($1, $2, $3, $4)`,
		client.ClientName,
		client.ExchangeName,
		client.Label,
		client.Pair,
	)
	if err != nil {
		return fmt.Errorf("error saving client in db: %w", err)
	}

	_, err = d.Client.Exec(
		ctx,
		` INSERT INTO HistoryOrder (
            		client_name, exchange_name, label, pair, side, type, base_qty, price, 
            		algorithm_name_placed, lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed
        	  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		order.ClientName,
		order.ExchangeName,
		order.Label,
		order.Pair,
		order.Side,
		order.TypeOrder,
		order.BaseQty,
		order.Price,
		order.AlgorithmNamePlaced,
		order.LowestSellPrc,
		order.HighestBuyPrc,
		order.CommissionQuoteQty,
		order.TimePlaced,
	)

	if err != nil {
		return fmt.Errorf("error saving order: %w", err)
	}

	return nil
}
