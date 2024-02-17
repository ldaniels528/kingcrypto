package kingcrypto

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestCreateOrder(t *testing.T) {
	// request quote update for AAPL
	lastSaleTime, err := parseISODate("2024-02-18T22:32:28.230Z")
	err = postQuote(Quote{
		Symbol:       "AAPL",
		Exchange:     "NYSE",
		LastSale:     156.88,
		LastSaleTime: lastSaleTime,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	// request quote update for INTC
	lastSaleTime, err = parseISODate("2024-02-18T22:32:28.230Z")
	err = postQuote(Quote{
		Symbol:       "INTC",
		Exchange:     "NYSE",
		LastSale:     52.22,
		LastSaleTime: lastSaleTime,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	// request quote update for AMD
	lastSaleTime, err = parseISODate("2024-02-18T22:32:28.230Z")
	err = postQuote(Quote{
		Symbol:       "AMD",
		Exchange:     "NASDAQ",
		LastSale:     69.65,
		LastSaleTime: lastSaleTime,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	// request buy order for AAPL
	err = postOrder(Order{
		Symbol:     "AAPL",
		Exchange:   "NYSE",
		Type:       "BUY",
		LimitPrice: 112.78,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	// request buy order for INTC
	err = postOrder(Order{
		Symbol:     "INTC",
		Exchange:   "NYSE",
		LimitPrice: 52.22,
		Type:       "BUY",
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	// request buy order for AMD
	err = postOrder(Order{
		Symbol:     "AMD",
		Exchange:   "NASDAQ",
		LimitPrice: 69.65,
		Type:       "BUY",
	})
	if err != nil {
		t.Errorf(err.Error())
	}
}

func postOrder(order Order) error {
	// encode the quote into JSON (bytes)
	jsonBody, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// send a POST request to create the new quote
	resp, err := http.Post("http://localhost:8080/order", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check the status code inside response
	if resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Error creating order. Status code: %d", resp.StatusCode))
	}

	return nil
}

func postQuote(quote Quote) error {
	// encode the quote into JSON (bytes)
	jsonBody, err := json.Marshal(quote)
	if err != nil {
		return err
	}

	// send a POST request to create the new quote
	resp, err := http.Post("http://localhost:8080/quote", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check the status code inside response
	if resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Error creating order. Status code: %d", resp.StatusCode))
	}

	return nil
}

func parseISODate(inputDate string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, inputDate)
}
