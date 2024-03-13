package main

import (
	"context"
	"currentPrice/config"
	"currentPrice/database"
	"currentPrice/models"
	"currentPrice/server"
	"time"
)

func main() {
	database.DBCon()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	precoAtualTudo := make(chan []models.PriceResponse)

	go updatePrices(ctx, precoAtualTudo)
	go processPrices(ctx, precoAtualTudo)

	<-ctx.Done()
}

func updatePrices(ctx context.Context, precoAtualTudo chan<- []models.PriceResponse) {
	ticker := time.NewTicker(600 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			precoAtual := server.PrecoAtualTodas()
			precoAtualTudo <- precoAtual
		}
	}
}

func processPrices(ctx context.Context, precoAtualTudo <-chan []models.PriceResponse) {
	defaultTimeSleep := 14
	timeSleep := config.TimeSleep
	if timeSleep == 0 {
		timeSleep = int64(defaultTimeSleep)
	}

	ticker := time.NewTicker(time.Duration(timeSleep) * time.Second)
	defer ticker.Stop()

	runTicker := time.NewTicker(600 * time.Millisecond)
	defer runTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case precoAtual := <-precoAtualTudo:
			server.Run(precoAtual)
			<-runTicker.C
		case <-ticker.C:
			precoAtual := <-precoAtualTudo
			server.Run2(precoAtual)
		}
	}
}
