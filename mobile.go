package dpo

import "encoding/xml"

// ChargeTokenMobileRequest is a request to charge a subscriber's mobile money directly.
type ChargeTokenMobileRequest struct {
	XMLName          xml.Name `xml:"API3G"`
	CompanyToken     string   `xml:"CompanyToken"`
	Request          string   `xml:"Request"`
	TransactionToken string   `xml:"TransactionToken"`
	PhoneNumber      string   `xml:"PhoneNumber"`
	MNO              string   `xml:"MNO"`
	MNOcountry       string   `xml:"MNOcountry"`
}

// ChargeTokenMobileResponse is a response from a ChargeTokenMobileRequest.
type ChargeTokenMobileResponse struct {
	XMLName        xml.Name `xml:"API3G"`
	Code           int      `xml:"Code"`
	Explanation    string   `xml:"Explanation"`
	RedirectURL    string   `xml:"RedirectUrl"`
	DeclinedURL    string   `xml:"declinedUrl"`
	Instructions   string   `xml:"Instructions"`
	RedirectOption int      `xml:"RedirectOption"`
}
