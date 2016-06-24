package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	listenPort := os.Getenv("PORT")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/stockprice", StockPrice)
	log.Fatal(http.ListenAndServe(":"+listenPort, router))
}

// StockPrice return the stock price for NET-B
func StockPrice(w http.ResponseWriter, r *http.Request) {
	var color string
	stock, err := yquotes.NewStock("NET-B.ST", false)
	if err != nil {
		fmt.Println(err)
	}

	if stock.Price.Last <= stock.Price.PreviousClose {
		color = "red"
	} else {
		color = "green"
	}

	price := strconv.FormatFloat(stock.Price.Last, 'f', 2, 64)
	oldPrice := strconv.FormatFloat(stock.Price.PreviousClose, 'f', 2, 64)

	response := &Message{
		Color:         color,
		MessageFormat: "text",
		Notify:        "false",
		Message:       "Current stock price is: " + price + " SEK\nPrevious stock price was: " + oldPrice + " SEK",
	}
	json.NewEncoder(w).Encode(response)
}
