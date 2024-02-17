package kingcrypto

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

///////////////////////////////////////////////////////
//			ORDERS API
///////////////////////////////////////////////////////

type Order struct {
	Type           string    `json:"type"`
	Symbol         string    `json:"symbol"`
	Exchange       string    `json:"exchange"`
	LimitPrice     float64   `json:"limitPrice"`
	CreationTime   time.Time `json:"creationTime"`
	ExpirationTime time.Time `json:"expirationTime"`
}

var orders = make(map[string]Order)

func handleOrderCRUD(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		handleOrderWriteAction(w, r, cancelOrder)
	case "GET":
		handleOrderReadAction(w, r, findOrder)
	case "POST":
		handleOrderWriteAction(w, r, createOrder)
	case "PUT":
		handleOrderWriteAction(w, r, findAndUpdateOrder)
	default:
		http.Error(w, "Method "+r.Method+" is unsupported for /order", http.StatusBadRequest)
	}
}

func handleOrderReadAction(w http.ResponseWriter, r *http.Request, find func(order Order) (*Order, error)) {
	// decode the order form
	form, err := toOrderForm(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// search for the order
	order, err := find(*form)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// encode the order and send it
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func toOrderForm(v url.Values) (*Order, error) {
	// extract the symbol
	symbol, err := extractValue(v, "symbol")
	if err != nil {
		return nil, err
	}
	// extract the type
	_type, err := extractValue(v, "type")
	if err != nil {
		return nil, err
	}
	// extract the exchange
	exchange, err := extractValue(v, "exchange")
	if err != nil {
		return nil, err
	}
	// return the order form
	return &Order{
		Symbol:   *symbol,
		Type:     *_type,
		Exchange: *exchange,
	}, nil
}

func handleOrderWriteAction(w http.ResponseWriter, r *http.Request, write func(order Order) error) {
	// decode the order form
	var form Order
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// add the new order
	if err := write(form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// respond with the order or error
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	var activeOrders []Order
	for _, order := range orders {
		if compareTime(order.ExpirationTime, currentTime) >= 0 {
			activeOrders = append(activeOrders, order)
		}
	}
	handleObjectList(w, r, activeOrders)
}

///////////////////////////////////////////////////////
//			ORDER FUNCTIONS
///////////////////////////////////////////////////////

func cancelOrder(form Order) error {
	searchKey := orderKey(form)
	order, ok := orders[searchKey]
	if ok {
		order.ExpirationTime = time.Now()
		return nil
	} else {
		return orderNotFound(form)
	}
}

func createOrder(form Order) error {
	form.CreationTime = time.Now()
	form.ExpirationTime = form.CreationTime.Add(3 * 24 * time.Hour) // + 3 days
	orders[orderKey(form)] = form
	return nil
}

func findOrder(form Order) (*Order, error) {
	order, ok := orders[orderKey(form)]
	if ok {
		return &order, nil
	} else {
		return nil, orderNotFound(form)
	}
}

func findAndUpdateOrder(form Order) error {
	searchKey := orderKey(form)
	_, ok := orders[searchKey]
	if ok {
		orders[searchKey] = form
		return nil
	} else {
		return orderNotFound(form)
	}
}

func orderKey(order Order) string {
	return order.Symbol + "|" + order.Exchange + "|" + order.Type
}

func orderNotFound(order Order) error {
	return errors.New("Order '" + orderKey(order) + "' was not found")
}
