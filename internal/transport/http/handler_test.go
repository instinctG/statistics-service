package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/instinctG/statistics_service/internal/statistics"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStatisticService struct {
	mock.Mock
}

func (m *MockStatisticService) GetOrderBook(ctx context.Context, exchangeName string, pair string) ([]*statistics.DepthOrder, []*statistics.DepthOrder, error) {
	args := m.Called(ctx, exchangeName, pair)
	return args.Get(0).([]*statistics.DepthOrder), args.Get(1).([]*statistics.DepthOrder), args.Error(2)
}

func (m *MockStatisticService) SaveOrderBook(ctx context.Context, orderBook *statistics.OrderBook) error {
	args := m.Called(ctx, orderBook)
	return args.Error(0)
}

func (m *MockStatisticService) GetOrderHistory(ctx context.Context, client *statistics.Client) ([]*statistics.HistoryOrder, error) {
	args := m.Called(ctx, client)
	return args.Get(0).([]*statistics.HistoryOrder), args.Error(1)
}

func (m *MockStatisticService) SaveOrder(ctx context.Context, client *statistics.Client, order *statistics.HistoryOrder) error {
	args := m.Called(ctx, client, order)
	return args.Error(0)
}

type HandlerSuite struct {
	suite.Suite
	handler *Handler
	service *MockStatisticService
}

func (suite *HandlerSuite) SetupTest() {
	suite.service = new(MockStatisticService)
	suite.handler = &Handler{Service: suite.service}
}

func (suite *HandlerSuite) TestGetOrderBook() {
	req := httptest.NewRequest(http.MethodGet, "/api/get?exchange_name=Binance&pair=BTC_USD", nil)
	w := httptest.NewRecorder()

	asks := []*statistics.DepthOrder{{Price: 10000, BaseQty: 1}}
	bids := []*statistics.DepthOrder{{Price: 9500, BaseQty: 1}}
	suite.service.On("GetOrderBook", mock.Anything, "Binance", "BTC_USD").Return(asks, bids, nil)

	suite.handler.GetOrderBook(w, req)

	require.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)
	var response map[string][]*statistics.DepthOrder
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), asks, response["asks"])
	require.Equal(suite.T(), bids, response["bids"])
}

func (suite *HandlerSuite) TestGetOrderHistory() {
	client := statistics.Client{ClientName: "test_client"}
	clientJSON, _ := json.Marshal(client)
	req := httptest.NewRequest(http.MethodPost, "/api/get_history", bytes.NewReader(clientJSON))
	w := httptest.NewRecorder()

	history := []*statistics.HistoryOrder{{ClientName: "test_client"}}
	suite.service.On("GetOrderHistory", mock.Anything, &client).Return(history, nil)

	suite.handler.GetOrderHistory(w, req)

	require.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)
	var response []*statistics.HistoryOrder
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), history, response)
}

func (suite *HandlerSuite) TestSaveOrderBook() {
	orderBook := statistics.OrderBook{Exchange: "Binance", Pair: "BTC_USD"}
	orderBookJSON, _ := json.Marshal(orderBook)
	req := httptest.NewRequest(http.MethodPost, "/api/save", bytes.NewReader(orderBookJSON))
	w := httptest.NewRecorder()

	suite.service.On("SaveOrderBook", mock.Anything, &orderBook).Return(nil)

	suite.handler.SaveOrderBook(w, req)

	require.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)
	var response Response
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), "Order book saved successfully", response.Message)
}

func (suite *HandlerSuite) TestSaveOrder() {
	historyOrder := statistics.HistoryOrder{ClientName: "test_client"}
	historyOrderJSON, _ := json.Marshal(historyOrder)
	req := httptest.NewRequest(http.MethodPost, "/api/save_order", bytes.NewReader(historyOrderJSON))
	w := httptest.NewRecorder()

	client := statistics.Client{
		ClientName:   historyOrder.ClientName,
		ExchangeName: historyOrder.ExchangeName,
		Pair:         historyOrder.Pair,
		Label:        historyOrder.Label,
	}
	suite.service.On("SaveOrder", mock.Anything, &client, &historyOrder).Return(nil)

	suite.handler.SaveOrder(w, req)

	require.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)
	var response Response
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), "Order saved successfully", response.Message)
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
