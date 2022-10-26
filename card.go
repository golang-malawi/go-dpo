package dpo

import "encoding/xml"

type chargeTokenResponseCode string

const (
	TransactionCharged     chargeTokenResponseCode = "000" // TransactionCharged Transaction charged
	TransactionAlreadyPaid chargeTokenResponseCode = "200" // TransactionAlreadyPaid Transaction alreadyp aid
	TokenMissing           chargeTokenResponseCode = "801" // TokenMissing Token missing
	InvalidToken           chargeTokenResponseCode = "802" // InvalidToken Invalid token
	MissingRequestOrName   chargeTokenResponseCode = "803" // MissingRequestOrName Missing request or name
	XMLError               chargeTokenResponseCode = "804" // XMLError Xml error
	DataMismatch           chargeTokenResponseCode = "902" // DataMismatch Data mismatch
	MissingMandatoryFields chargeTokenResponseCode = "950" // MissingMandatoryFields Missing mandatory fields
	TransactionDenied      chargeTokenResponseCode = "999" // TransactionDenied Transaction denied
)

const (
	opChargeTokenCreditCard = "chargeTokenCreditCard"
)

func (v chargeTokenResponseCode) Description() string {
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

// ChargeCreditCardRequest is a request to charge a users card directly.
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

// ThreeDRequest request data for 3D systems
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

// ChargeCreditCardResponse response returned from after processing a credit card charge directly.
type ChargeCreditCardResponse struct {
	XMLName xml.Name `xml:"API3G"`

	Result      string `xml:"Result"`
	Explanation string `xml:"ResultExplanation"`
	RedirectURL string `xml:"RedirectUrl,omitempty"`
	BackURL     string `xml:"BackUrl,omitempty"`
	DeclinedURL string `xml:"declinedUrl,omitempty"`
}

// IsError determines whether the card response is an error or not.
func (c *ChargeCreditCardResponse) IsError() bool {
	return c.Result != "000"
}
