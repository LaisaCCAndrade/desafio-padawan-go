package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const coinAPIURL = "https://rest.coinapi.io/v1/exchangerate"
const apiKey = "F0EC2B15-FDF8-413F-92D8-9E328E3BE000"

type ExchangeRatesResponse struct {
	Rate float64 `json:"rate"`
}

var moedas = map[string]string{
	"BRL": "Real (BRL)",
	"USD": "Dólar (USD)",
	"EUR": "Euro (EUR)",
	"BTC": "Bitcoin (BTC)",
	// Adicione outras moedas aqui
}

func getExchangeRate(from, to string) (float64, error) {
	url := fmt.Sprintf("%s/%s/%s", coinAPIURL, from, to)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("X-CoinAPI-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("falha na solicitação: %s", resp.Status)
	}

	var ratesResp ExchangeRatesResponse
	err = json.NewDecoder(resp.Body).Decode(&ratesResp)
	if err != nil {
		return 0, err
	}

	return ratesResp.Rate, nil
}

func converterMoeda(valor float64, de, para string) (float64, error) {
	rate, err := getExchangeRate(de, para)
	if err != nil {
		return 0, err
	}

	valorConvertido := valor * rate
	return valorConvertido, nil
}

type ConvertRequest struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
}

type ConvertResponse struct {
	ConvertedValue float64 `json:"convertedValue"`
}

// Endpoint para conversão de moeda
func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica o corpo da solicitação para ConvertRequest
	var req ConvertRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Falha ao ler o corpo da solicitação", http.StatusBadRequest)
		return
	}

	// Realiza a conversão da moeda
	valorConvertido, err := converterMoeda(req.Amount, req.From, req.To)
	if err != nil {
		http.Error(w, "Falha na conversão", http.StatusInternalServerError)
		return
	}

	// Prepara a resposta como ConvertResponse em formato JSON
	res := ConvertResponse{
		ConvertedValue: valorConvertido,
	}

	// Define o cabeçalho Content-Type para JSON e envia a resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/convert", convertHandler)

	// Servir o arquivo styles.css
	fs := http.FileServer(http.Dir("."))
	http.Handle("/styles.css", fs)

	fmt.Println("Servidor rodando em http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
