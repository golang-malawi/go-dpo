package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/nndi-oss/go-dpo"
)

func main() {

	clientToken := os.Getenv("DPO_TOKEN")

	client := dpo.NewClient(clientToken, true)

	createTokenRequest := dpo.NewCreateTokenRequest(clientToken, "USD", big.NewFloat(0.003))

	createTokenRequest.AddService("3854", "Ecommerce", time.Now())

	token, err := client.CreateToken(createTokenRequest)
	if err != nil {
		log.Fatalf("failed to create token", err)
	}

	time.Sleep(30 * time.Second)

	chargeResponse, err := client.ChargeCreditCard(os.Getenv("CARD_HOLDER"), os.Getenv("CARD_NUMBER"), os.Getenv("CARD_CVV"), os.Getenv("CARD_EXPIRY"), token)
	if err != nil {
		log.Fatalf("failed to charge client :%v", err)
	}

	fmt.Println(chargeResponse)
}
