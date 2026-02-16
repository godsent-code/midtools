package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/godsent-code/midtools/internal/application/brown_card_service"
	"github.com/godsent-code/midtools/internal/application/policy_verification"
	"github.com/godsent-code/midtools/internal/application/product"
	"github.com/godsent-code/midtools/internal/application/risk_type"
	"github.com/godsent-code/midtools/internal/application/sticker"
	"github.com/godsent-code/midtools/internal/application/ussd_check"
)

func NewRouter(
	service brown_card_service.BrownCard,
	stickerService sticker.StickerService,
	ussdCheckService ussd_check.USSDCheckService,
	policyVerification policy_verification.PolicyVerificationService,
	productService product.ProductService,
	riskTypeService risk_type.RiskTypeService,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	BrownCard := NewBrownCardHandler(service)
	Sticker := NewStickerHandler(stickerService)
	Ussd := NewUSSDCheckHandler(ussdCheckService)
	policyVerificationService := NewPolicyVerificationHandler(policyVerification)
	productHandler := NewProductHandler(productService)
	riskTypeHandler := NewRiskTypeHandler(riskTypeService)

	r.Post("/browncard", BrownCard.GetBrownCard)
	r.Post("/sticker", Sticker.GetSticker)
	r.Post("/ussd_check", Ussd.GetUSSDCheck)
	r.Post("/policy_verification", policyVerificationService.GetPolicyVerifications)
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products", productHandler.GetProducts)
	r.Post("/risk_type", riskTypeHandler.CreateRiskType)
	r.Get("/risk_type", riskTypeHandler.GetRiskTypes)
	return r

}
