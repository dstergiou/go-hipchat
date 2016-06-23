package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"

	"github.com/doneland/yquotes"
	"github.com/gorilla/mux"
)

// Message builds the basic struct for HipChat communication
type Message struct {
	Color         string `json:"color"`
	Notify        string `json:"notify"`
	MessageFormat string `json:"message_format"`
	Message       string `json:"message"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/stockprice", StockPrice)
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Index vanilla response
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// StockPrice return the stock price for NET-B
func StockPrice(w http.ResponseWriter, r *http.Request) {
	stock, err := yquotes.NewStock("NET-B.ST", false)
	if err != nil {
		fmt.Println(err)
	}
	price := strconv.FormatFloat(stock.Price.Last, 'f', 2, 64)
	oldPrice := strconv.FormatFloat(stock.Price.PreviousClose, 'f', 2, 64)

	response := &Message{
		Color:         "green",
		MessageFormat: "text",
		Notify:        "false",
		Message:       "Stock price is:" + price + "\nPrevious close was: " + oldPrice,
	}
	json.NewEncoder(w).Encode(response)
}
