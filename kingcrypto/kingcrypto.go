package kingcrypto

// "github.com/avukadin/goapi/internal/middleware"
// "github.com/go-chi/chi"
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

///////////////////////////////////////////////////////
//			WEB ROUTES
///////////////////////////////////////////////////////

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Welcome to the home page!")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "About Us")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

///////////////////////////////////////////////////////
//			INITIALIZATION
///////////////////////////////////////////////////////

func StartServer() {
	// remove trailing slashes
	//r.Use(chimiddle.StripSlashes)

	// Define routes and their respective handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/order", handleOrderCRUD)
	http.HandleFunc("/orders", handleOrders)
	http.HandleFunc("/quote", handleQuoteCRUD)
	http.HandleFunc("/quotes", handleQuotes)

	// Start the server on port 8080
	log.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Error starting server:", err)
	}
}

///////////////////////////////////////////////////////
//			COMMON / MISCELLANEOUS
///////////////////////////////////////////////////////

func compareTime(date1 time.Time, date2 time.Time) int {
	if date1.Before(date2) {
		return -1
	} else if date1.After(date2) {
		return 1
	} else {
		return 0
	}
}

func extractValue(v url.Values, name string) (*string, error) {
	value := strings.Trim(v.Get(name), "\t\n\r ")
	if len(value) == 0 {
		return nil, errors.New("Required query parameter '" + name + "' is missing")
	}
	return &value, nil
}

func handleObjectList(w http.ResponseWriter, r *http.Request, items any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
