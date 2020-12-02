package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	configFile = "config.json"
)

var (
	ubirchClientURL     string
	ubirchClientHeaders = map[string]string{}
)

type config struct {
	Uuid string `json:"uuid"`
	Auth string `json:"auth"`
}

type soapEnvelope struct {
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

type fault struct {
	Faultcode   string `xml:"faultcode"`
	Faultstring string `xml:"faultstring"`
}

type CertificationResponse struct {
	Hash     string
	Upp      string
	Response string
}

func setConfig() error {
	conf := config{}
	err := conf.load(configFile)
	if err != nil {
		return err
	}

	ubirchClientURL = fmt.Sprintf("http://localhost:8080/%s", conf.Uuid)
	ubirchClientHeaders = map[string]string{
		"Content-Type": "application/json",
		"X-Auth-Token": conf.Auth,
	}

	return nil
}

func (c *config) load(filename string) error {
	contextBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(contextBytes, c)
}

func parseSoapRequest(reqBody []byte) ([]byte, error) {
	var Envelope soapEnvelope
	err := xml.Unmarshal(reqBody, &Envelope)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(Envelope.Body.Document)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func sendJsonRequest(reqBody []byte) (int, []byte, http.Header, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", ubirchClientURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, nil, nil, fmt.Errorf("can't make new post request: %v", err)
	}

	for k, v := range ubirchClientHeaders {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}

	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, respBodyBytes, resp.Header, err
}

func parseJsonResponse(respBody []byte) ([]byte, error) {
	var resp CertificationResponse
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}

	xmlBytes, err := xml.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return xmlBytes, nil
}

func forwardClientResponse(w http.ResponseWriter, respCode int, respBody []byte, respHeader http.Header) {
	for k, v := range respHeader {
		if strings.ToLower(k) != "content-length" {
			w.Header().Set(k, v[0])
		}
	}
	w.WriteHeader(respCode)
	_, err := w.Write(respBody)
	if err != nil {
		log.Fatalf("unable to write response: %s", err) // todo fatal?
	}
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	soapReq, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to read request body: %v", err), http.StatusBadRequest)
		return
	}

	jsonReq, err := parseSoapRequest(soapReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	respCode, respBody, respHeader, err := sendJsonRequest(jsonReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to send request: %v", err), http.StatusInternalServerError)
		return
	}

	xmlResp, err := parseJsonResponse(respBody)
	if err != nil {
		xmlBytes, err := xml.Marshal(fault{Faultcode: "soap:Server", Faultstring: string(respBody)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		xmlResp = xmlBytes // todo
	}

	forwardClientResponse(w, respCode, xmlResp, respHeader)
}

func main() {
	err := setConfig()
	if err != nil {
		log.Fatalf("Could not set config: %s\n", err.Error())
	}

	http.HandleFunc("/", handleRequest)
	s := &http.Server{
		Addr: ":8090",
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
