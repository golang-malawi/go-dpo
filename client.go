package dpo

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Debug bool   // Determines whether to use test or live url
	Token string // Credentials key for the company
	http  *http.Client
}

// NewClient creates a new testing/debug client for 3G service
// companyToken the token to use for API calls
// debug whether to enable debug-mode or not - debug mode uses the test URLs instead of live URLs.
func NewClient(companyToken string, debug bool) *Client {
	return &Client{
		Debug: debug,
		Token: companyToken,
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

func (c *Client) CreateToken(token *CreateTokenRequest) (*CreateTokenResponse, error) {
	//TODO: Validate the token
	url := TestApiUrl
	if !c.Debug {
		url = LiveApiUrl
	} else {
		// TODO: log that we are using debug
	}
	xmlData, err := xml.Marshal(token)
	if err != nil {
		return nil, fmt.Errorf("failed to form XML request: %s got: %v", string(xmlData), err)
	}

	if c.Debug {
		fmt.Printf("using request body: %s\n", string(xmlData))
	}

	r := bytes.NewReader(xmlData)
	resp, err := c.http.Post(url, "text/xml", r)
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
		Request:          OpChargeTokenCreditCard,
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

	url := TestApiUrl
	if !c.Debug {
		url = LiveApiUrl
	} else {
		// TODO: log that we are using debug
	}
	xmlData, err := xml.Marshal(cardRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to form XML request: %s got: %v", string(xmlData), err)
	}

	if c.Debug {
		fmt.Printf("using request body: %s\n", string(xmlData))
	}

	r := bytes.NewReader(xmlData)
	resp, err := c.http.Post(url, "text/xml", r)
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
