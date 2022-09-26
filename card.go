package dpo

import "encoding/xml"

type ChargeTokenResponseCode string

const (
	TransactionCharged     ChargeTokenResponseCode = "000"
	TransactionAlreadyPaid ChargeTokenResponseCode = "200"
	TokenMissing           ChargeTokenResponseCode = "801"
	InvalidToken           ChargeTokenResponseCode = "802"
	MissingRequestOrName   ChargeTokenResponseCode = "803"
	XMLError               ChargeTokenResponseCode = "804"
	DataMismatch           ChargeTokenResponseCode = "902"
	MissingMandatoryFields ChargeTokenResponseCode = "950"
	TransactionDenied      ChargeTokenResponseCode = "999"
)

const (
	OpChargeTokenCreditCard = "chargeTokenCreditCard"
)

func (v ChargeTokenResponseCode) Description() string {
	switch v {
	case TransactionCharged:
		return "Transaction charged"
	case TransactionAlreadyPaid:
		return "Transaction already paid"

	case TokenMissing:
		return "Request missing company token"

	case InvalidToken:
		return "Wrong CompanyToken"

	case MissingRequestOrName:
		return "No request or error in Request type name"

	case XMLError:
		return "Error in XML"

	case DataMismatch:
		return "Data mismatch in one of the fields – fieldname"

	case MissingMandatoryFields:
		return "Request missing mandatory fields – fieldname"

	case TransactionDenied:
		return "Transaction Declined - Explanation"
	default:
		return "Unknown"
	}
}

// <?xml version="1.0" encoding="utf-8"?>
// <API3G>
//   <CompanyToken>57466282-EBD7-4ED5-B699-8659330A6996</CompanyToken>
//   <Request>chargeTokenCreditCard</Request>
//   <TransactionToken>72983CAC-5DB1-4C7F-BD88-352066B71592</TransactionToken>
//   <CreditCardNumber>123412341234</CreditCardNumber>
//   <CreditCardExpiry>1214</CreditCardExpiry>
//   <CreditCardCVV>333</CreditCardCVV>
//   <CardHolderName>John Doe</CardHolderName>
//   <ThreeD>
//         <Enrolled>Y</Enrolled>
//         <Paresstatus>Y</Paresstatus>
//         <Eci>05</Eci>
//         <Xid>DYYVcrwnujRMnHDy1wlP1Ggz8w0=</Xid>
//         <Cavv>mHyn+7YFi1EUAREAAAAvNUe6Hv8=</Cavv>
//         <Signature>_</Signature>
//         <Veres>AUTHENTICATION_SUCCESSFUL</Veres>
//         <Pares>eAHNV1mzokgW/isVPY9GFSCL0EEZkeyg7</Pares>
//     </ThreeD>
// </API3G>
type ChargeCreditCardRequest struct {
	XMLName xml.Name `xml:"API3G"`

	CompanyToken     string        `xml:"CompanyToken"`
	Request          string        `xml:"Request"`
	TransactionToken string        `xml:"TransactionToken"`
	CreditCardNumber string        `xml:"CreditCardNumber"`
	CreditCardExpiry string        `xml:"CreditCardExpiry"`
	CreditCardCVV    string        `xml:"CreditCardCVV"`
	CardHolderName   string        `xml:"CardHolderName"`
	ThreeD           ThreeDRequest `xml:"ThreeD"`
}

type ThreeDRequest struct {
	Enrolled    string `xml:"Enrolled"`
	Paresstatus string `xml:"Paresstatus"`
	Eci         string `xml:"Eci"`
	Xid         string `xml:"Xid"`
	Cavv        string `xml:"Cavv"`
	Signature   string `xml:"Signature"`
	Veres       string `xml:"Veres"`
	Pares       string `xml:"Pares"`
}

// <?xml version="1.0" encoding="UTF-8"?>
// <API3G><Code>200</Code>
// 	<Explanation>Transaction already paid</Explanation>
// 	<RedirectUrl>https://redirect.com</RedirectUrl>
// 	<BackUrl></BackUrl>
// 	<declinedUrl></declinedUrl>
// </API3G>
type ChargeCreditCardResponse struct {
	XMLName xml.Name `xml:"API3G"`

	Result      string `xml:"Result"`
	Explanation string `xml:"ResultExplanation"`
	RedirectUrl string `xml:"RedirectUrl,omitempty"`
	BackUrl     string `xml:"BackUrl,omitempty"`
	DeclinedUrl string `xml:"declinedUrl,omitempty"`
}

func (c *ChargeCreditCardResponse) IsError() bool {
	return c.Result != "000"
}
