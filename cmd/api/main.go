package main

import (
	"context"
	"log"
	http2 "net/http"

	"github.com/godsent-code/midtools/configs"
	"github.com/godsent-code/midtools/internal/adapters/http"
	"github.com/godsent-code/midtools/internal/adapters/postgres"
	"github.com/godsent-code/midtools/internal/application/brown_card_service"
	"github.com/godsent-code/midtools/internal/application/policy_verification"
	"github.com/godsent-code/midtools/internal/application/product"
	"github.com/godsent-code/midtools/internal/application/risk_type"
	"github.com/godsent-code/midtools/internal/application/sticker"
	"github.com/godsent-code/midtools/internal/application/ussd_check"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	config, err := configs.LoadConfig("./")
	if err != nil {
		log.Fatal(err)

	}
	conn, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	brownCardRepo := postgres.NewBrownCardRepository(config)
	brownCardService := brown_card_service.NewBrownCard(brownCardRepo)

	stickerRepo := postgres.NewStickerRepository(config)
	stickerService := sticker.NewStickerService(stickerRepo)

	ussdRepo := postgres.NewUSSDCheckerRepository(config)
	ussdService := ussd_check.NewUSSDCheckService(ussdRepo)

	policyVerificationRepo := postgres.NewPolicyVerificationRepository(config)
	policyVerificationService := policy_verification.NewPolicyVerificationService(policyVerificationRepo)

	productRepo := postgres.NewProductRepository(conn, config)
	productService := product.NewProductService(productRepo)

	riskRepo := postgres.NewRiskTypeRepository(conn, config)
	riskService := risk_type.NewRiskTypeService(riskRepo)

	router := http.NewRouter(brownCardService, stickerService, ussdService, policyVerificationService, productService, riskService)

	err = http2.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal(err)
		return
	}
}
