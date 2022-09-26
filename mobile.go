package dpo

import "encoding/xml"

//https://directpayonline.atlassian.net/wiki/spaces/API/pages/2937749565/ChargeTokenMobile+V7
/*
<?xml version="1.0" encoding="UTF-8"?>
<API3G>
   <CompanyToken>90EC1DA4-A7C5-432C-930C-098715D3130E</CompanyToken>
   <Request>ChargeTokenMobile</Request>
   <TransactionToken>F0C9D5A6-D130-44B7-896C-A0FD701FE132</TransactionToken>
   <PhoneNumber>25412345678</PhoneNumber>
   <MNO>SafaricomC2B</MNO>
   <MNOcountry>kenya</MNOcountry>
</API3G>
*/

type ChargeTokenMobileRequest struct {
	XMLName          xml.Name `xml:"API3G"`
	CompanyToken     string   `xml:"CompanyToken"`
	Request          string   `xml:"Request"`
	TransactionToken string   `xml:"TransactionToken"`
	PhoneNumber      string   `xml:"PhoneNumber"`
	MNO              string   `xml:"MNO"`
	MNOcountry       string   `xml:"MNOcountry"`
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<API3G>
    <Code>130</Code>
	<Explanation>New Invoice</Explanation>
	<RedirectUrl>https://redirect.com</RedirectUrl>
	<BackUrl></BackUrl>
	<declinedUrl></declinedUrl>
	<Instructions>1.Go to the M-PESA menu&lt;br&gt;2. Select Lipa na M-PESA&lt;br&gt;3. Select the Paybill
		option&lt;br&gt;4. Enter business number 927633&lt;br&gt;5. Enter your account number 1776F0C9D&lt;br&gt;6.
		Enter the amount 567&lt;br&gt;7. Press OK to send&lt;br&gt;8. You will receive a confirmation SMS with your
		payment reference number.</Instructions>
	<RedirectOption>0</RedirectOption>
</API3G>
*/
type ChargeTokenMobileResponse struct {
	XMLName        xml.Name `xml:"API3G"`
	Code           int      `xml:"Code"`
	Explanation    string   `xml:"Explanation"`
	RedirectUrl    string   `xml:"RedirectUrl"`
	DeclinedUrl    string   `xml:"declinedUrl"`
	Instructions   string   `xml:"Instructions"`
	RedirectOption int      `xml:"RedirectOption"`
}
