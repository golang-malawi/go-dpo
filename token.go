package dpo

import (
	"encoding/xml"
	"math/big"
	"time"
)

// <?xml version="1.0" encoding="utf-8"?>
// <API3G>
// <CompanyToken>57466282-EBD7-4ED5-B699-8659330A6996</CompanyToken>
// <Request>createToken</Request>
// <Transaction>
// <PaymentAmount>450.00</PaymentAmount>
// <PaymentCurrency>USD</PaymentCurrency>
// <CompanyRef>49FKEOA</CompanyRef>
// <RedirectURL>http://www.domain.com/payurl.php</RedirectURL>
// <BackURL>http://www.domain.com/backurl.php </BackURL>
// <CompanyRefUnique>0</CompanyRefUnique>
// <PTL>5</PTL>
// </Transaction>
// <Services>
//   <Service>
//     <ServiceType>45</ServiceType>
//     <ServiceDescription>Flight from Nairobi to Diani</ServiceDescription>
//     <ServiceDate>2013/12/20 19:00</ServiceDate>
//   </Service>
// </Services>
// </API3G>
type CreateTokenRequest struct {
	XMLName xml.Name `xml:"API3G"`

	CompanyToken string                 `xml:"CompanyToken"`
	Request      string                 `xml:"Request"`
	Transaction  CreateTokenTransaction `xml:"Transaction"`
	Services     []Service              `xml:"Services>Service"`
}

func NewCreateTokenRequest(companyToken string, paymentCurrency string, amount *big.Float) *CreateTokenRequest {
	return &CreateTokenRequest{
		CompanyToken: companyToken,
		Request:      "createToken",
		Transaction: CreateTokenTransaction{
			PaymentAmount:    amount.String(),
			PaymentCurrency:  paymentCurrency,
			CompanyRef:       "",
			RedirectURL:      "",
			BackURL:          "",
			CompanyRefUnique: "",
			PTL:              "5",
		},
		Services: []Service{},
	}
}

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

func (c *CreateTokenRequest) SetBackURL(backURL string) {
	c.Transaction.BackURL = backURL
}

func (c *CreateTokenRequest) SetRedirectURL(redirectURL string) {
	c.Transaction.RedirectURL = redirectURL
}

type Service struct {
	ServiceType        string `xml:"ServiceType"`
	ServiceDescription string `xml:"ServiceDescription"`
	ServiceDate        string `xml:"ServiceDate"`
}

type CreateTokenTransaction struct {
	PaymentAmount    string `xml:"PaymentAmount"`
	PaymentCurrency  string `xml:"PaymentCurrency"`
	CompanyRef       string `xml:"CompanyRef"`
	RedirectURL      string `xml:"RedirectURL"`
	BackURL          string `xml:"BackURL"`
	CompanyRefUnique string `xml:"CompanyRefUnique"`
	PTL              string `xml:"PTL"`
}

type CreateTokenResponse struct {
	XMLName xml.Name `xml:"API3G"`

	Result            string      `xml:"Result"`
	ResultExplanation string      `xml:"ResultExplanation"`
	TransToken        string      `xml:"TransToken,omitempty"`
	TransRef          string      `xml:"TransRef,omitempty"`
	Allocations       Allocations `xml:"Allocations,omitempty"`
}

func (c *CreateTokenResponse) IsError() bool {
	return c.Result != "000"
}

type Allocations struct {
	Allocation Allocation `xml:"Allocation"`
}

type Allocation struct {
	AllocationID   string `xml:"AllocationID"`
	AllocationCode string `xml:"AllocationCode"`
}
