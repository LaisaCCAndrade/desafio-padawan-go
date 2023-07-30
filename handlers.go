package main

import (
	"encoding/json"
	"net/http"
)

type Quote struct {
	Moeda string  `json:"moeda"`
	Valor float64 `json:"valor"`
}

func GetBitcoinQuote(w http.ResponseWriter, r *http.Request) {
	quote := Quote{
		Moeda: "Bitcoin",
		Valor: 45000.00, // Replace with the actual updated value of Bitcoin quotation
	}
	json.NewEncoder(w).Encode(quote)
}

func GetDolarQuote(w http.ResponseWriter, r *http.Request) {
	quote := Quote{
		Moeda: "Dolar",
		Valor: 5.30, // Replace with the actual updated value of Dollar quotation
	}
	json.NewEncoder(w).Encode(quote)
}

func GetEuroQuote(w http.ResponseWriter, r *http.Request) {
	quote := Quote{
		Moeda: "Euro",
		Valor: 6.20, // Replace with the actual updated value of Euro quotation
	}
	json.NewEncoder(w).Encode(quote)
}

func GetRealQuote(w http.ResponseWriter, r *http.Request) {
	quote := Quote{
		Moeda: "Real",
		Valor: 1.00, // The quotation of Real will always be 1, as it is the base currency
	}
	json.NewEncoder(w).Encode(quote)
}
