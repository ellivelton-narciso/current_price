package server

import (
	"currentPrice/config"
	"currentPrice/database"
	"currentPrice/models"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Run(precoAtualTudo []models.PriceResponse) {
	var (
		err                 error
		rows                *sqlx.Rows
		moedasFiltradas     []models.PriceResponse
		bots                []models.Bots
		moedasFiltradasJSON []byte
		count               int
	)

	qry := "select * from bots UNION select * from bots_real;"

	rows, err = database.DB.Queryx(qry)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var bot models.Bots
		err := rows.StructScan(&bot)
		if err != nil {
			fmt.Println("\n erro38 - ", err)
			continue
		}
		bots = append(bots, bot)
	}
	defer rows.Close()

	for _, preco := range precoAtualTudo {
		for _, bot := range bots {
			if preco.Symbol == bot.Coin {
				moedasFiltradas = append(moedasFiltradas, preco)
				break
			}
		}
	}
	if len(bots) == 0 {
		moedasFiltradasJSON = []byte("[]")
	} else {
		moedasFiltradasJSON, err = json.Marshal(moedasFiltradas)
		if err != nil {
			fmt.Println("\n erro39 - ", err)
			return
		}
	}
	qry2 := "SELECT COUNT(*) FROM historico"
	err = database.DB.Get(&count, qry2)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.DB.Exec("INSERT INTO historico (value, created_at) VALUES (?, NOW())", moedasFiltradasJSON)
	if err != nil {
		fmt.Println("\n Erro ao inserir valores no banco de dados (Run1) - ", err)
		return
	}

	if count >= config.Leituras {
		err = apagarPrimeiroHistorico()
		if err != nil {
			fmt.Println("\n Erro ao apagar o primeiro historico - ", err)
			return
		}
	}
}

func Run2(precoAtualTudo []models.PriceResponse) {

	var err error
	var value float64

	for _, p := range precoAtualTudo {
		value, err = strconv.ParseFloat(p.Price, 64)
		if err != nil {
			fmt.Println(err)
		}
		if strings.Contains(p.Symbol, "USDT") {
			_, err = database.DB.Exec("INSERT INTO hist_trading_values (hist_date, trading_name, curr_value) VALUES (NOW(), ?, ?)", p.Symbol, value)
			if err != nil {
				fmt.Println("\n Erro ao inserir valores no banco de dados (Run2) - ", err)
				return
			}
		}

	}
}

func PrecoAtualTodas() []models.PriceResponse {
	config.ReadFile()

	url := config.BaseURL + "fapi/v2/ticker/price"
	req, _ := http.NewRequest("GET", url, nil)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao acessar a API para converter: ", err)
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	var priceResp []models.PriceResponse
	err = json.Unmarshal(body, &priceResp)
	if err != nil {
		log.Println("Erro ao decodificar JSON:", err)
		return nil
	}

	return priceResp
}

func apagarPrimeiroHistorico() error {
	qry := "DELETE FROM historico ORDER BY id LIMIT 1"
	_, err := database.DB.Exec(qry)
	return err
}
