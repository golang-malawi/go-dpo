package dpo

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client struct represents a client and it's configuration for working with the DPO API.
// The client provides functions to initiate, verify, cancel and revoke payment tokens.
// The client uses a basic net/http http.Client.
type Client struct {
	Debug       bool   // Determines whether to use test or live url
	Token       string // Credentials key for the company
	http        *http.Client
	maxAttempts int // Maximum number of attempts per operation
	GenerateRef func() string
}

// defaultCompanyRefGenerator is a function that generates a string that can be used as a Transaction ID
func defaultCompanyRefGenerator() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		// TODO: default to some other random string scheme
		panic(err) // TODO: don't panic in a library
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// xmlMarshallWithHeader marshals dat into XML with the xml header prepended
func xmlMarshalWithHeader(data any) ([]byte, error) {
	xmlstring, err := xml.Marshal(data) // xml.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, err
	}

	xmlstring = []byte(xml.Header + string(xmlstring))
	return xmlstring, nil
}

// xmlMarshalWithHeaderDebug for debugging, pretty prints the marshalled XML
func xmlMarshalWithHeaderDebug(data any) ([]byte, error) {
	xmlstring, err := xml.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, err
	}

	xmlstring = []byte(xml.Header + string(xmlstring))
	return xmlstring, nil
}

// MakePaymentURL creates a URL which should be passed to the User to redirect to the DPO system to complete the payment
// Requires a non-nil token created using client.CreateToken
func (c *Client) MakePaymentURL(token *CreateTokenResponse) string {
	if token == nil {
		return ""
	}

	if c.Debug {
		return fmt.Sprintf("%s?ID=%s", testPayURL, token.TransToken)
	}
	return fmt.Sprintf("%s?ID=%s", livePayURL, token.TransToken)
}

// NewClient creates a new testing/debug client for 3G service
// companyToken the token to use for API calls
// debug whether to enable debug-mode or not - debug mode uses the test URLs instead of live URLs.
func NewClient(companyToken string, debug bool) *Client {
	return &Client{
		Debug:       debug,
		Token:       companyToken,
		maxAttempts: 5, // other DPO libraries use 10: see - TODO: add link
		GenerateRef: defaultCompanyRefGenerator,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewLiveClient creates a new Client that has debug set to false
// companyToken the token to use for API calls
func NewLiveClient(companyToken string) *Client {
	return NewClient(companyToken, false)
}

// CreateToken creates a token that can be used to perform payments. This is the first step in the payment flow with DPO
// Once the token is created it must be verified using client.VerifyToken
func (c *Client) CreateToken(token *CreateTokenRequest) (*CreateTokenResponse, error) {
	if token == nil {
		return nil, fmt.Errorf("token must not be nil")
	}
	url := testAPIURL
	if !c.Debug {
		url = liveAPIURL
	} else {
		// TODO: log that we are using debug
	}
	xmlData, err := xmlMarshalWithHeader(token)
	if err != nil {
		return nil, fmt.Errorf("failed to form XML request: %s got: %v", string(xmlData), err)
	}

	if c.Debug {
		fmt.Printf("using request body: %s\n", string(xmlData))
	}

	r := bytes.NewReader(xmlData)

	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "go-dpo: https://github.com/nndi-oss/go-dpo/v1-beta")
	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Cache-control", "no-cache")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s got: %v", string(bodyData), err)
	}
	if c.Debug {
		fmt.Printf("got response body: %s\n", string(bodyData))
	}
	var tokenResponse CreateTokenResponse
	if resp.StatusCode == http.StatusOK {
		err = xml.Unmarshal(bodyData, &tokenResponse)
		if err != nil {
			return nil, fmt.Errorf("failed unmarshal response: %v", err)
		}
		if tokenResponse.IsError() {
			return nil, fmt.Errorf("failed to charge card: %s", tokenResponse.ResultExplanation)
		}
		return &tokenResponse, nil
	}

	return nil, fmt.Errorf("invalid response code:%d body: %s", resp.StatusCode, string(bodyData))
}

// VerifyToken verifies the token with DPO site to prepare it for use for actual payment process
func (c *Client) VerifyToken(token *CreateTokenResponse) (*VerifyTokenResponse, error) {
	verifyRequest := &VerifyTokenRequest{
		Request:          "verifyToken",
		CompanyToken:     c.Token,
		TransactionToken: token.TransToken,
	}

	//TODO: Validate the token
	url := testAPIURL
	if !c.Debug {
		url = liveAPIURL
	} else {
		// TODO: log that we are using debug
	}
	xmlData, err := xmlMarshalWithHeader(verifyRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to form XML request: %s got: %v", string(xmlData), err)
	}

	if c.Debug {
		fmt.Printf("using request body: %s\n", string(xmlData))
	}

	r := bytes.NewReader(xmlData)
	var created = false

	maxAttempts := c.maxAttempts

	for i := 0; !created && i < maxAttempts; i++ {
		req, err := http.NewRequest("POST", url, r)
		if err != nil {
			return nil, err
		}
		req.Header.Add("User-Agent", "go-dpo: https://github.com/nndi-oss/go-dpo/v1-beta")
		req.Header.Add("Content-Type", "application/xml")
		req.Header.Add("Cache-control", "no-cache")

		resp, err := c.http.Do(req)

		if err != nil {
			return nil, err
		}

		bodyData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %s got: %v", string(bodyData), err)
		}
		if c.Debug {
			fmt.Printf("got response body: %s\n", string(bodyData))
		}
		var verifyTokenResponse VerifyTokenResponse
		if resp.StatusCode == http.StatusOK {
			err = xml.Unmarshal(bodyData, &verifyTokenResponse)
			if err != nil {
				return nil, fmt.Errorf("failed unmarshal response: %v", err)
			}
			// if verifyTokenResponse == "900".IsError() {
			// 	return nil, fmt.Errorf("failed to charge card: %s", verifyTokenResponse.ResultExplanation)
			// }
			return &verifyTokenResponse, nil
		}

		return nil, fmt.Errorf("invalid response code:%d body: %s", resp.StatusCode, string(bodyData))
	}

	return nil, fmt.Errorf("failed to process request after %d attempts", c.maxAttempts)
}

// ChargeCreditCard is used for charging a card directly. Do not use this yet.
func (c *Client) ChargeCreditCard(cardHolder, cardNumber, cvv, cardExpiry string, token *CreateTokenResponse) (*ChargeCreditCardResponse, error) {
	if token == nil {
		return nil, fmt.Errorf("failed to get token: nil value passed as 'token'")
	}
	transactionToken := token.TransToken
	if transactionToken == "" {
		return nil, fmt.Errorf("failed to get token")
	}

	cardRequest := &ChargeCreditCardRequest{
		CompanyToken:     c.Token,
		Request:          opChargeTokenCreditCard,
		TransactionToken: token.TransToken,
		CreditCardNumber: cardNumber,
		// The API doesn't accept  an expiry with MM/YY it requires MMYY
		CreditCardExpiry: strings.ReplaceAll(cardExpiry, "/", ""),
		CreditCardCVV:    cvv,
		CardHolderName:   cardHolder,
		ThreeD: ThreeDRequest{
			Enrolled:    "Y",
			Paresstatus: "Y",
			Eci:         "05",
			Xid:         "",
			Cavv:        "",
			Signature:   "_",
			Veres:       "AUTHENTICATION_SUCCESSFUL",
			Pares:       "",
		},
	}

	url := testAPIURL
	if !c.Debug {
		url = liveAPIURL
	} else {
		// TODO: log that we are using debug
	}
	xmlData, err := xmlMarshalWithHeader(cardRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to form XML request: %s got: %v", string(xmlData), err)
	}

	if c.Debug {
		fmt.Printf("using request body: %s\n", string(xmlData))
	}

	r := bytes.NewReader(xmlData)
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "go-dpo: https://github.com/nndi-oss/go-dpo/v1-beta")
	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Cache-control", "no-cache")

	resp, err := c.http.Do(req)

	if err != nil {
		return nil, err
	}
	// got an error response,
	if err != nil {
		return nil, err
	}
	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s got: %v", string(bodyData), err)
	}
	if c.Debug {
		fmt.Printf("got response body: %s\n", string(bodyData))
	}
	var cardResponse ChargeCreditCardResponse
	if resp.StatusCode == http.StatusOK {
		err = xml.Unmarshal(bodyData, &cardResponse)
		if err != nil {
			return nil, fmt.Errorf("failed unmarshal response: %v", err)
		}
		if cardResponse.IsError() {
			return nil, fmt.Errorf("failed to charge card: %s", cardResponse.Explanation)
		}

		return &cardResponse, nil
	}

	return nil, fmt.Errorf("invalid response code:%d body: %s", resp.StatusCode, string(bodyData))
}

// CancelToken initiates token cancellations - NOT YET IMPLEMENTED
func (c *Client) CancelToken() (*CancelTokenResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

// RefundToken initiates token refunds - NOT YET IMPLEMENTED
func (c *Client) RefundToken() (*RefundTokenResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
