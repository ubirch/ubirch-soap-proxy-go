package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	configFile = "config.json"
)

var (
	ubirchClientURL     string
	ubirchClientHeaders map[string]string
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
	ubirchClientHeaders = map[string]string{ // todo use header from original request if it exists
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
		return 0, nil, nil, err
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

	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, respBodyBytes, resp.Header, nil
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

func sendResponse(w http.ResponseWriter, respBody []byte, respCode int) {
	if strings.HasPrefix(string(respBody), "<") {
		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	w.WriteHeader(respCode)
	_, err := w.Write(respBody)
	if err != nil {
		log.Errorf("unable to write response: %v", err)
	}
}

func Error(w http.ResponseWriter, error string, code int) {
	log.Error(error)
	xmlError, err := xml.Marshal(fault{Faultcode: "soap:Server", Faultstring: error})
	if err != nil {
		xmlError = []byte(error)
	}
	sendResponse(w, xmlError, code)
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	soapReq, err := ioutil.ReadAll(req.Body)
	if err != nil {
		Error(w, fmt.Sprintf("unable to read request body: %v", err), http.StatusBadRequest)
		return
	}

	jsonReq, err := parseSoapRequest(soapReq)
	if err != nil {
		Error(w, fmt.Sprintf("unable to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	respCode, respBody, respHeader, err := sendJsonRequest(jsonReq)
	if err != nil {
		Error(w, fmt.Sprintf("unable to send request: %v", err), http.StatusInternalServerError)
		return
	}

	log.Infof("response: (%d) %s, %s", respCode, respBody, respHeader)

	xmlResp, err := parseJsonResponse(respBody)
	if err != nil {
		xmlFault, err := xml.Marshal(fault{Faultcode: "soap:Server", Faultstring: string(respBody)})
		if err != nil {
			xmlResp = respBody
		} else {
			xmlResp = xmlFault
		}
	}

	sendResponse(w, xmlResp, respCode)
}

func main() {
	err := setConfig()
	if err != nil {
		log.Fatalf("Could not set config: %v", err)
	}

	http.HandleFunc("/", handleRequest)
	s := &http.Server{
		Addr: ":8090",
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
