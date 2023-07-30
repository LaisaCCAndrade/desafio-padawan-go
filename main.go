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
		<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    	<title>Conversor de Moedas</title>
	</head>
	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 20px;
			background: #b1b1b11a;
			display: flex;
			flex-direction: column;
			align-items: center;
  		}
  
  		h1 {
			background: #ffffff;
			text-align: center;
			box-shadow: 0 0 1.5em #737875, 0 0 0.5em #595b53;
			color: #1c221e;
			border: none;
  		}
  
  		form {
			display: flex;
			width: 400px;
			height: 271px;
			flex-direction: column;
			align-items: center;
			border-radius: 8px;
			margin-top: 20px;
			box-shadow: 0 0 1.5em #737875, 0 0 0.5em #595b53;
			justify-content: center;
  		}
  
		label {
			color: #000000;
			font-size: 18px;
			margin-top: 10px;
		}

		input, select {
			padding: 5px;
			margin-top: 5px;
			border: 1px solid #000000;
		}
  
		button {
			padding: 5px;
			margin-top: 5px;
			border: none;
		}
  
		button {
			background-color: #4caf50;
			color: white;
			border: none;
			cursor: pointer;
		}
  
		button:hover {
			background-color: #45a049;
		}
		
		.moedas {
			border-collapse: collapse;
			margin: 22px 64px 64px 58px;
			display: flex;
			width: 400px;
			height: 271px;
			flex-direction: column;
			align-items: center;
			border-radius: 8px;
			box-shadow: 0 0 1.5em #737875, 0 0 0.5em #595b53;
			justify-content: center;
		}
		
		th {
			border: 1px solid black;
			padding: 8px;
		}

		td {
			border: 1px solid black;
			padding: 8px;
			color: #ffffff;
		}
	
		th {
			background-color: #f2f2f2;
			font-weight: bold;
		}

		.conversao {
			display: flex;
			flex-direction: column;
			align-items: center;
		}

		canvas {
			box-sizing: border-box;
			height: 400px;
			width: 77%;
			background: white;
			margin: 30px auto;
			box-shadow: 0 0 1.5em #737875, 0 0 0.5em #595b53;
		}

		.usd {
			background: #0903a8;
		}

		.brl {
			background: #0e633e;
		}

		.eur {
			background: #d08542;
		}

		.btc {
			background: #6b0957;
		}

		@media(min-width: 900px) {
			.conversao {
				flex-direction: row
			}

			form {
				margin-bottom: 61px;
			}

			canvas {
				width: 867px;
				height: 700px;
				margin-left: -133px;
			}

			.moedas {
				margin: 22px 64px 64px 0;
			}
		}
  
	</style>
	<body>
    	<h1>Genesis Bank Conversor de Moedas</h1>
   
		<div class="conversao">
			<table class="moedas" style="float: left;">
				<tr>
					<th>Moeda</th>
					<th>Valor em relação ao Dólar Americano</th>
				</tr>
				<tr class="usd">
					<td>USD</td>
					<td>$ 1.00</td>
				</tr>
				<tr class="brl">
					<td>BRL</td>
					<td>R$ 4.72</td>
				</tr>
				<tr class="eur">
					<td>EUR</td>
					<td>€ 0.91</td>
				</tr>
				<tr class="btc">
					<td>BTC</td>
					<td>₿ 0.000034</td>
				</tr>
			</table>

			<form action="/exchange" method="post">
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

		</div>

		<table style="float: left; padding-right: 20px;">
			<!-- Conteúdo da tabela aqui -->
		</table>

		<div style="width: 600px; height: 480px;">
			<canvas id="graficoLinhas"></canvas>
		</div>

		<script>
			const valoresMoedas = {
				USD: [1.00, 1.10, 1.05, 1.08, 1.15, 1.20, 1.18, 1.25, 1.30, 1.35, 1.40, 1.38],
				BRL: [4.72, 5.00, 5.20, 5.40, 5.50, 5.60, 5.55, 5.45, 5.35, 5.30, 5.25, 5.28],
				EUR: [0.91, 0.95, 0.89, 0.88, 0.86, 0.84, 0.80, 0.82, 0.85, 0.88, 0.90, 0.87],
				BTC: [0.000034, 0.000040, 0.000032, 0.000038, 0.000035, 0.000036, 0.000039, 0.000037, 0.000036, 0.000039, 0.000038, 0.000035]
			};

			const ctx = document.getElementById('graficoLinhas').getContext('2d');

			new Chart(ctx, {
				type: 'line',
				data: {
					labels: ['Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun', 'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez'],
					datasets: [
						{
							label: 'USD',
							data: valoresMoedas.USD,
							borderColor: 'blue',
							fill: false
						},
						{
							label: 'BRL',
							data: valoresMoedas.BRL,
							borderColor: 'green',
							fill: false
						},
						{
							label: 'EUR',
							data: valoresMoedas.EUR,
							borderColor: 'orange',
							fill: false
						},
						{
							label: 'BTC',
							data: valoresMoedas.BTC,
							borderColor: 'purple',
							fill: false
						}
					]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false,
					scales: {
						x: {
							display: true,
							title: {
								display: true,
								text: 'Mês'
							}
						},
						y: {
							display: true,
							title: {
								display: true,
								text: 'Valor em relação ao Dólar Americano'
							}
						}
					}
				}
			});
		</script>  
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
	http.HandleFunc("/exchange", convertHandler)

	fs := http.FileServer(http.Dir("."))
	http.Handle("./styles.css", fs)

	fmt.Println("Servidor rodando em http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
