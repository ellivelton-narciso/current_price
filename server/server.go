package server

import (
	"currentPrice/config"
	"currentPrice/database"
	"currentPrice/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Run() {
	qry := "SELECT * FROM bots"

	rows, err := database.DB.Queryx(qry)
	if err != nil {
		log.Fatal(err)
	}
	var bots []models.Bots
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

	precoAtualTudo := precoAtualTodas()

	var moedasFiltradas []models.PriceResponse
	for _, preco := range precoAtualTudo {
		for _, bot := range bots {
			if preco.Symbol == bot.Coin {
				moedasFiltradas = append(moedasFiltradas, preco)
				break
			}
		}
	}

	var moedasFiltradasJSON []byte
	if len(bots) == 0 {
		moedasFiltradasJSON = []byte("[]")
	} else {
		var err error
		moedasFiltradasJSON, err = json.Marshal(moedasFiltradas)
		if err != nil {
			fmt.Println("\n erro39 - ", err)
			return
		}
	}

	var count int
	qry2 := "SELECT COUNT(*) FROM historico"
	err = database.DB.Get(&count, qry2)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.DB.Exec("INSERT INTO historico (value, created_at) VALUES (?, NOW())", moedasFiltradasJSON)
	if err != nil {
		fmt.Println("\n login39 - ", err)
		return
	}

	if count >= 60 {
		err := apagarPrimeiroHistorico()
		if err != nil {
			fmt.Println("\n login39 - ", err)
			return
		}
	}
}

func precoAtualTodas() (allPrice []models.PriceResponse) {
	config.ReadFile()

	url := config.BaseURL + "fapi/v2/ticker/price"
	req, _ := http.NewRequest("GET", url, nil)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao acessar a API para converter: ", err)
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var priceResp []models.PriceResponse
	err = json.Unmarshal(body, &priceResp)
	if err != nil {
		log.Println("Erro ao decodificar JSON:", err)
		return
	}

	return priceResp
}

func apagarPrimeiroHistorico() error {
	qry := "DELETE FROM historico ORDER BY id LIMIT 1"
	_, err := database.DB.Exec(qry)
	return err
}
