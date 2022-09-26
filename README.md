# Unofficial Go Library for DPO

## Usage

```go
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

	client := dpo.NewClient(clientToken)

	createTokenRequest := dpo.NewCreateTokenRequest(clientToken, "USD", *big.NewFloat(30.00))

	createTokenRequest.AddService("3854", "Ecommerce", time.Now())

	token, err := client.CreateToken(createTokenRequest)
	if err != nil {
		log.Fatalf("failed to create token", err)
	}

	chargeResponse, err := client.ChargeCreditCard("JOHN DOE", "1234-1234-1234-1234", "123", "12/31", token)
	if err != nil {
		log.Fatalf("failed to charge client :%v", err)
	}

	fmt.Println(chargeResponse)
}
```
