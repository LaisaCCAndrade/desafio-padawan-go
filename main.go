package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<link rel="stylesheet" href="/styles.css">
		<title>Conversor de Moedas</title>
	</head>
	<body>
		<h1>Conversor de Moedas</h1>
		<form action="/convert" method="post">
			<label for="amount">Insira o valor a converter:</label>
			<input type="number" name="amount" step="0.01" placeholder="Digite o valor" required>

			<label for="from">Moeda para conversão:</label>
			<select name="from" required>
				{{range $key, $value := .Moedas}}
				<option value="{{$key}}">{{$value}}</option>
				{{end}}
			</select>

			<label for="to">Moeda convertida:</label>
			<select name="to" required>
				{{range $key, $value := .Moedas}}
				<option value="{{$key}}">{{$value}}</option>
				{{end}}
			</select>

			<button type="submit">Converter</button>
		</form>
	</body>
	</html>
	`))

	data := struct {
		Moedas map[string]string
	}{
		Moedas: moedas,
	}

	tmpl.Execute(w, data)
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Valor inválido", http.StatusBadRequest)
		return
	}

	from := r.FormValue("from")
	to := r.FormValue("to")

	valorConvertido, err := converterMoeda(amount, from, to)
	if err != nil {
		http.Error(w, "Falha na conversão", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%.2f %s = %.2f %s\n", amount, from, valorConvertido, to)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/convert", convertHandler)

	fs := http.FileServer(http.Dir("."))
    http.Handle("/styles.css", fs)

	fmt.Println("Servidor rodando em http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
