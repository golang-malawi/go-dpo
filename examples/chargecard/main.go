package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/golang-malawi/go-dpo"
)

func main() {
	clientToken := os.Getenv("DPO_TOKEN")

	client := dpo.NewClient(clientToken, true)

	client.SetUserAgent("Example User Agent")

	createTokenRequest := client.NewCreateTokenRequest(clientToken, "USD", big.NewFloat(0.30))

	createTokenRequest.AddService("3854", "Ecommerce", time.Now())

	token, err := client.CreateToken(createTokenRequest)
	if err != nil {
		log.Fatalf("failed to create token %v", err)
	}

	time.Sleep(30 * time.Second)

	verifyResponse, err := client.VerifyToken(token)
	if err != nil {
		log.Fatalf("failed to charge client :%v", err)
	}
	fmt.Println("=== verify token", verifyResponse)

	if verifyResponse.Result == "900" {
		fmt.Printf("Click here to verify payment: %s", client.MakePaymentURL(token))
	}

	// time.Sleep(30 * time.Second)
	// chargeResponse, err := client.ChargeCreditCard(os.Getenv("CARD_HOLDER"), os.Getenv("CARD_NUMBER"), os.Getenv("CARD_CVV"), os.Getenv("CARD_EXPIRY"), token)
	// if err != nil {
	// 	log.Fatalf("failed to charge client :%v", err)
	// }

	// fmt.Println("=== chargeResponse", chargeResponse)
}
