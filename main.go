package main

import (
	"encoding/json"
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

// SendMessage returns the JSON notification to hipchat
func SendMessage(w http.ResponseWriter, color string, message string) {
	response := &Message{
		Color:         color,
		MessageFormat: "text",
		Notify:        "false",
		Message:       message,
	}
	json.NewEncoder(w).Encode(response)
}

// StockPrice return the stock price for NET-B
func StockPrice(w http.ResponseWriter, r *http.Request) {
	var color, comment string
	var difference float64

	const strikePrice float64 = 109.70

	stock, err := yquotes.NewStock("NET-B.ST", false)
	if err != nil {
		color = "yellow"
		message := "Stock information cannot be retrieved"
		SendMessage(w, color, message)
	}

	if stock.Price.Last <= stock.Price.PreviousClose {
		color = "red"
	} else {
		color = "green"
	}

	if stock.Price.Last <= strikePrice {
		difference = strikePrice - stock.Price.Last
		comment = "Stock needs to go up " + strconv.FormatFloat(difference, 'f', 2, 64) + " SEK to hit the strike price (109.70 SEK)"
	} else {
		difference = stock.Price.Last - strikePrice
		comment = "Stock is " + strconv.FormatFloat(difference, 'f', 2, 64) + " SEK above the strike price (109.70 SEK)"
	}

	price := strconv.FormatFloat(stock.Price.Last, 'f', 2, 64)
	oldPrice := strconv.FormatFloat(stock.Price.PreviousClose, 'f', 2, 64)
	message := "Current stock price is: " + price + " SEK\nPrevious stock price was: " + oldPrice + " SEK" + "\n" + comment
	SendMessage(w, color, message)
}
