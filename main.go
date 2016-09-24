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

func floatToString(number float64) (str string) {
	return strconv.FormatFloat(number, 'f', 2, 64)
}

// StockPrice return the stock price for NET-B
func StockPrice(w http.ResponseWriter, r *http.Request) {
	var color, comment string
	var difference, percentage float64

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
		percentage = (difference / stock.Price.Last) * 100
		comment = "Stock needs to go up " + floatToString(difference) + " SEK to hit the strike price (109.70 SEK)"
		comment += "\n This is a " + floatToString(percentage) + "% increase that we need"
	} else {
		difference = stock.Price.Last - strikePrice
		comment = "Stock is " + floatToString(difference) + " SEK above the strike price (109.70 SEK)"
	}

	message := "Price on program start (May 23 2016) was 87.00 SEK"
	message += "Current stock price is: " + floatToString(stock.Price.Last) + " SEK\n"
	message += "Previous stock price was: " + floatToString(stock.Price.PreviousClose) + " SEK" + "\n"
	message += comment
	SendMessage(w, color, message)
}
