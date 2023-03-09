package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"

	"github.com/golang-malawi/go-dpo"
)

type DPOConfig struct {
	Token       string `env:"DPO_TOKEN" validate:"required"`
	RedirectURL string `env:"DPO_REDIRECT_URL" validate:"required"`
	BackURL     string `env:"DPO_BACK_URL" validate:"required"`

	ServiceCode string `env:"DPO_SERVICE_CODE" validate:"required"`
	ServiceName string `env:"DPO_SERVICE_NAME" validate:"required"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file %v", err)
	}

	dpoConfig := DPOConfig{
		Token:       os.Getenv("DPO_TOKEN"),
		RedirectURL: os.Getenv("DPO_REDIRECT_URL"),
		BackURL:     os.Getenv("DPO_BACK_URL"),
	}
	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Get("/pay", func(c *fiber.Ctx) error {
		return c.Render("initiate_payment", fiber.Map{})
	})

	app.Post("/payment/initialize", func(ctx *fiber.Ctx) error {
		return InitiatePayment(ctx, dpoConfig)
	})

	app.Get("/payment/verify", func(c *fiber.Ctx) error {
		// When payment completes, DPO will make a callback to this handler
		TransID := c.Query("TransID", "")
		CCDapproval := c.Query("CCDapproval", "")
		PnrID := c.Query("PnrID", "")
		TransactionToken := c.Query("TransactionToken", "")
		CompanyRef := c.Query("CompanyRef", "")

		return c.Render("payment_complete", fiber.Map{
			"TransID":          TransID,
			"CCDapproval":      CCDapproval,
			"PnrID":            PnrID,
			"TransactionToken": TransactionToken,
			"CompanyRef":       CompanyRef,
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func InitiatePayment(ctx *fiber.Ctx, dpoConfig DPOConfig) error {
	clientToken := dpoConfig.Token
	client := dpo.NewClient(clientToken, true)

	// MOTE: Load Plan, Currency and Amount from database
	createTokenRequest := client.NewCreateTokenRequest(clientToken, "USD", big.NewFloat(0.30))

	if dpoConfig.RedirectURL != "" {
		createTokenRequest.SetRedirectURL(dpoConfig.RedirectURL)
	}

	if dpoConfig.BackURL != "" {
		createTokenRequest.SetBackURL(dpoConfig.BackURL)
	}

	// NOTE: load service code and service name from db or config for the specific plan
	createTokenRequest.AddService(dpoConfig.ServiceCode, dpoConfig.ServiceName, time.Now())

	token, err := client.CreateToken(createTokenRequest)
	if err != nil {
		return ctx.Render("payment_error", fiber.Map{
			"errorMessage":       fmt.Sprintf("Failed to Create token. Got: %s", err.Error()),
			"createTokenRequest": fmt.Sprintf("%v", createTokenRequest),
			"chargeToken":        fmt.Sprintf("%v", token),
		})
	}

	time.Sleep(3 * time.Second)
	// URL we will redirect to for the user to make a payment on DPOs site...
	DPOPaymentURL := client.MakePaymentURL(token)

	verifyResponse, err := client.VerifyToken(token)
	if err != nil {
		return ctx.Render("payment_error", fiber.Map{
			"errorMessage":       fmt.Sprintf("Failed to Verify Token. Got: %s", err.Error()),
			"createTokenRequest": fmt.Sprintf("%v", createTokenRequest),
			"chargeToken":        fmt.Sprintf("%v", token),
		})
	}

	transactionId := token.TransRef
	transactionToken := token.TransToken
	companyRef := createTokenRequest.Transaction.CompanyRef
	// TODO: Update transaction data here
	// Verify the token
	if verifyResponse.Result == "900" {
		return ctx.Render("dpo_redirect", fiber.Map{
			"transactionRef":         transactionId,
			"transactionToken":       transactionToken,
			"internalTransactionRef": companyRef,
			"DPOPaymentURL":          DPOPaymentURL,
		})
	}

	return ctx.Render("payment_error", fiber.Map{
		"errorMessage":       fmt.Sprintf("Failed to Process request. Got: %s", err.Error()),
		"createTokenRequest": fmt.Sprintf("%v", createTokenRequest),
		"chargeToken":        fmt.Sprintf("%v", token),
	})
}
