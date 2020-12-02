package main

import (
	"encoding/xml"
	"fmt"
)

var payload = `<s11:Envelope xmlns:s11='http://schemas.xmlsoap.org/soap/envelope/'>
  <s11:Body>
    <ubirch:Document xmlns:ubirch='http://ubirch.com/'>
      <ActionReferenceNumber>?XXX?</ActionReferenceNumber>
      <ActionID>?XXX?</ActionID>
      <SpecialUseID>?XXX?</SpecialUseID>
      <PeriodBeginDate>?XXX?</PeriodBeginDate>
      <PeriodBeginTime>?XXX?</PeriodBeginTime>
      <PeriodEndDate>?XXX?</PeriodEndDate>
      <PeriodEndTime>?XXX?</PeriodEndTime>
      <PostCode>?XXX?</PostCode>
      <City>?XXX?</City>
      <District>?XXX?</District>
      <Street>?XXX?</Street>
      <HouseNumber>?XXX?</HouseNumber>
      <FromCrossroad>?XXX?</FromCrossroad>
      <ToCrossroad>?XXX?</ToCrossroad>
<!-- optional -->
      <LicensePlate>?XXX?</LicensePlate>
<!-- optional -->
      <GeoAreaCoordinates>?XXX?</GeoAreaCoordinates>
<!-- optional -->
      <GeoOverviewCoordinates>?XXX?</GeoOverviewCoordinates>
    </ubirch:Document>
  </s11:Body>
</s11:Envelope>`

type SoapEnvelope struct {
	Body soapBody `xml:"Body"`
}
type soapBody struct {
	Document soapDocument `xml:"Document"`
}

type soapDocument struct {
	ActionReferenceNumber string `xml:"ActionReferenceNumber,omitempty"`
	ActionID              string `xml:"ActionID,omitempty"`
	SpecialUseID          string `xml:"SpecialUseID"`    // todo: enum?
	PeriodBeginDate       string `xml:"PeriodBeginDate"` // todo: type -> date
	PeriodBeginTime       string `xml:"PeriodBeginTime"` // todo: type -> time
	PeriodEndDate         string `xml:"PeriodEndDate"`   // todo: type -> date
	PeriodEndTime         string `xml:"PeriodEndTime"`   // todo: type -> time
	PostCode              string `xml:"PostCode"`
	City                  string `xml:"City"`
	District              string `xml:"District"`
	Street                string `xml:"Street"`
	HouseNumber           string `xml:"HouseNumber"`
	FromCrossroad         string `xml:"FromCrossroad"`
	ToCrossroad           string `xml:"ToCrossroad"`
}

func main() {

	var Envelope SoapEnvelope
	err := xml.Unmarshal([]byte(payload), &Envelope)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%v\n", Envelope.Body.Document.SpecialUseID)

}
