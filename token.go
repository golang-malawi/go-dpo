package dpo

import (
	"encoding/xml"
	"math/big"
	"time"
)

// CreateTokenRequest is a request to create a token that will be used to process (i.e. initiate, complete, cancel, revoke) payments.
type CreateTokenRequest struct {
	XMLName xml.Name `xml:"API3G"`

	CompanyToken string                 `xml:"CompanyToken"`
	Request      string                 `xml:"Request"`
	Transaction  CreateTokenTransaction `xml:"Transaction"`
	Services     []Service              `xml:"Services>Service"`
}

// NewCreateTokenRequest creates a new token that can be used in client.VerifyToken calls
func (c *Client) NewCreateTokenRequest(companyToken string, paymentCurrency string, amount *big.Float) *CreateTokenRequest {
	// TODO: add validation before creating token
	return &CreateTokenRequest{
		CompanyToken: companyToken,
		Request:      "createToken",
		Transaction: CreateTokenTransaction{
			PaymentAmount:    amount.String(),
			PaymentCurrency:  paymentCurrency,
			CompanyRef:       c.GenerateRef(),
			RedirectURL:      "",
			BackURL:          "",
			CompanyRefUnique: 0, // 0 - not unique, 1 - duplicate request
			PTL:              "5",
		},
		Services: []Service{},
	}
}

// AddService adds a service to the CreateTokenRequests slice of services which indicates which services the payment will be made for.
func (c *CreateTokenRequest) AddService(typeCode, description string, serviceDate time.Time) {
	service := &Service{
		ServiceType:        typeCode,
		ServiceDescription: description,
		ServiceDate:        serviceDate.Format("2006/01/02 15:04"),
	}
	if c.Services == nil || len(c.Services) < 1 {
		c.Services = make([]Service, 0)
		c.Services = append(c.Services, *service)
		return
	}

	c.Services = append(c.Services, *service)
}

// SetBackURL sets the URL that DPO will redirect to when user cancels the payment flow or an error occurs
func (c *CreateTokenRequest) SetBackURL(backURL string) {
	c.Transaction.BackURL = backURL
}

// SetRedirectURL sets the URL that DPO will redirect to when user completes the payment flow
func (c *CreateTokenRequest) SetRedirectURL(redirectURL string) {
	c.Transaction.RedirectURL = redirectURL
}

// Service is a product or service that users can pay for through DPO
type Service struct {
	ServiceType        string `xml:"ServiceType"`
	ServiceDescription string `xml:"ServiceDescription"`
	ServiceDate        string `xml:"ServiceDate"`
}

// CreateTokenTransaction TODO: add docs
type CreateTokenTransaction struct {
	PaymentAmount    string `xml:"PaymentAmount"`
	PaymentCurrency  string `xml:"PaymentCurrency"`
	CompanyRef       string `xml:"CompanyRef"`
	RedirectURL      string `xml:"RedirectURL"`
	BackURL          string `xml:"BackURL"`
	CompanyRefUnique int    `xml:"CompanyRefUnique"`
	PTL              string `xml:"PTL"`
}

// CreateTokenResponse is returned after processing a CreateTokenRequest and depending on the Result may be an error response or not
type CreateTokenResponse struct {
	XMLName xml.Name `xml:"API3G"`

	Result            string      `xml:"Result"`
	ResultExplanation string      `xml:"ResultExplanation"`
	TransToken        string      `xml:"TransToken,omitempty"`
	TransRef          string      `xml:"TransRef,omitempty"`
	Allocations       Allocations `xml:"Allocations,omitempty"`
}

// IsError determines whether the CreateTokenResponse is an error or not.
func (c *CreateTokenResponse) IsError() bool {
	return c.Result != "000"
}

// Allocations collection of allocations
type Allocations struct {
	Allocation Allocation `xml:"Allocation"`
}

// Allocation an allocation as defined by DPO
type Allocation struct {
	AllocationID   string `xml:"AllocationID"`
	AllocationCode string `xml:"AllocationCode"`
}

// VerifyTokenRequest is a request to verify a token that was requested as a CreateTokenRequest
type VerifyTokenRequest struct {
	XMLName xml.Name `xml:"API3G"`

	CompanyToken     string `xml:"CompanyToken"`
	TransactionToken string `xml:"TransactionToken"`
	Request          string `xml:"Request"`
}

// VerifyTokenResponse is returned after processing a VerifyTokenRequet and depending on the .Result may be an error response or not
type VerifyTokenResponse struct {
	XMLName xml.Name `xml:"API3G"`

	Result            string `xml:"Result"`
	ResultExplanation string `xml:"ResultExplanation"`
}

// CancelTokenRequest represents a request to cancel a previously created token.
type CancelTokenRequest struct {
	XMLName xml.Name `xml:"API3G"`
	// TODO: implement me!
}

// CancelTokenResponse is the result of requesting a cancel token and depending on .Result may be an error or not.
type CancelTokenResponse struct {
	XMLName xml.Name `xml:"API3G"`

	Result            string `xml:"Result"`
	ResultExplanation string `xml:"ResultExplanation"`
}

// RefundTokenRequest represents a request to initiate a refund
type RefundTokenRequest struct {
	XMLName xml.Name `xml:"API3G"`
	// TODO: implement me!
}

// RefundTokenResponse represents response from initiating a refund request.
type RefundTokenResponse struct {
	XMLName xml.Name `xml:"API3G"`

	Result            string `xml:"Result"`
	ResultExplanation string `xml:"ResultExplanation"`
}
