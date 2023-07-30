<h1 align="center"> Conversor de Moedas com Gráfico de Linhas </h1>
Este é um pequeno projeto em Go (Golang) que implementa um conversor de moedas com gráfico de linhas para visualizar as flutuações das quatro moedas (USD, BRL, EUR e BTC) em relação ao Dólar Americano nos últimos 12 meses. O conversor utiliza a API CoinAPI para obter as taxas de câmbio das moedas em relação ao Dólar.

## Pré-requisitos
- Antes de executar o projeto, é necessário ter o Go (Golang) instalado em sua máquina. Além disso, você precisará de uma chave de API da CoinAPI para acessar as taxas de câmbio.

- Instalação do Go (Golang)
- Chave de API da CoinAPI
- Como executar o projeto
- Clone este repositório ou faça o download do código.

- No arquivo main.go, substitua o valor da variável apiKey pela sua chave de API da CoinAPI:

- const apiKey = "SUA_CHAVE_DE_API_AQUI"
Abra o terminal e navegue até a pasta do projeto.

Execute o seguinte comando para executar o servidor:

```bash
go run main.go
```

O servidor estará rodando em http://localhost:8000. Abra o navegador e acesse essa URL para utilizar o conversor de moedas com o gráfico de linhas.
Funcionalidades
O conversor de moedas permite que você insira um valor em qualquer uma das moedas (USD, BRL, EUR ou BTC) e selecione a moeda para a qual deseja converter. Após clicar no botão "Converter", o valor convertido será exibido abaixo do formulário.

O gráfico de linhas mostra as flutuações das quatro moedas (USD, BRL, EUR e BTC) em relação ao Dólar Americano nos últimos 12 meses. Cada moeda possui uma linha colorida diferente no gráfico para facilitar a visualização.

## Agradecimentos
Este projeto utiliza a API CoinAPI (https://www.coinapi.io/) para obter as taxas de câmbio das moedas em relação ao Dólar Americano. Muito obrigado à equipe da CoinAPI por fornecer esse serviço!

## Autora
- **Laisa Andrade** - [Github](https://github.com/LaisaCCAndrade)