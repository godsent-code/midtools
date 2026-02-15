package main

import (
	"log"
	http2 "net/http"

	"github.com/godsent-code/midtools/configs"
	"github.com/godsent-code/midtools/internal/adapters/http"
	"github.com/godsent-code/midtools/internal/adapters/postgres"
	"github.com/godsent-code/midtools/internal/application/brown_card_service"
	"github.com/godsent-code/midtools/internal/application/sticker"
)

func main() {
	config, err := configs.LoadConfig("./")
	if err != nil {
		log.Fatal(err)
	}
	brownCardRepo := postgres.NewBrownCardRepository(config)
	brownCardService := brown_card_service.NewBrownCard(brownCardRepo)

	stickerRepo := postgres.NewStickerRepository(config)
	stickerService := sticker.NewStickerService(stickerRepo)

	router := http.NewRouter(brownCardService, stickerService)

	err = http2.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal(err)
		return
	}
}
