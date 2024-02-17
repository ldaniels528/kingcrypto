package kingcrypto

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

///////////////////////////////////////////////////////
//			QUOTES API
///////////////////////////////////////////////////////

type Quote struct {
	Symbol       string    `json:"symbol"`
	Exchange     string    `json:"exchange"`
	LastSale     float64   `json:"lastSale"`
	LastSaleTime time.Time `json:"lastSaleTime"`
}

var quotes = make(map[string]Quote)

func handleQuoteCRUD(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleQuoteReadAction(w, r, findQuote)
	case "POST":
		handleQuoteWriteAction(w, r, createQuote)
	case "PUT":
		handleQuoteWriteAction(w, r, findAndUpdateQuote)
	default:
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Method "+r.Method+" is unsupported for /quote", http.StatusBadRequest)
	}
}

func handleQuoteWriteAction(w http.ResponseWriter, r *http.Request, op func(quote Quote) error) {
	// decode the quote form
	var form Quote
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// add the new quote
	if err := op(form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// respond with the quote or error
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(form); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func handleQuoteReadAction(w http.ResponseWriter, r *http.Request, find func(quote Quote) (*Quote, error)) {
	// decode the order form
	searchForm, err := toQuoteForm(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// lookup the quote
	quote, err := find(*searchForm)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// encode the searchQuote and send it
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(quote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleQuotes(w http.ResponseWriter, r *http.Request) {
	handleObjectList(w, r, quotes)
}

///////////////////////////////////////////////////////
//			QUOTES FUNCTIONS
///////////////////////////////////////////////////////

func toQuoteForm(v url.Values) (*Quote, error) {
	// extract the symbol
	symbol, err := extractValue(v, "symbol")
	if err != nil {
		return nil, err
	}

	// extract the exchange
	exchange, err := extractValue(v, "exchange")
	if err != nil {
		return nil, err
	}

	// return the order form
	return &Quote{
		Symbol:   *symbol,
		Exchange: *exchange,
	}, nil
}

func createQuote(form Quote) error {
	quotes[quoteKey(form)] = form
	return nil
}

func findQuote(form Quote) (*Quote, error) {
	if quote, ok := quotes[quoteKey(form)]; ok {
		return &quote, nil
	} else {
		return nil, symbolNotFound(form)
	}
}

func findAndUpdateQuote(form Quote) error {
	searchKey := quoteKey(form)
	if _, ok := quotes[searchKey]; ok {
		quotes[searchKey] = form
		return nil
	} else {
		return symbolNotFound(form)
	}
}

func symbolNotFound(quote Quote) error {
	return errors.New("Symbol '" + quote.Symbol + "' was not found")
}

func quoteKey(quote Quote) string {
	return quote.Symbol + "|" + quote.Exchange
}
