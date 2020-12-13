// Copyright (c) 2019-2020 ubirch GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const configFile = "config.json"

var conf config

type config struct {
	Uuid                string `json:"uuid"`
	Auth                string `json:"auth"`
	VerificationBaseURL string `json:"verificationBaseURL"`
	UbirchClientURL     string `json:"ubirchClientURL"`
}

type soapEnvelope struct {
	Body soapBody `xml:"Body"`
}
type soapBody struct {
	Document soapDocument `xml:"Document"`
}

type soapDocument struct {
	ActionReferenceNumber  string `json:"ActionReferenceNumber,omitempty"`
	ActionID               string `json:"ActionID,omitempty"`
	SpecialUseDesc         string `json:"SpecialUseDesc,omitempty"`
	PeriodBeginDate        string `json:"PeriodBeginDate,omitempty"`
	PeriodBeginTime        string `json:"PeriodBeginTime,omitempty"`
	PeriodEndDate          string `json:"PeriodEndDate,omitempty"`
	PeriodEndTime          string `json:"PeriodEndTime,omitempty"`
	PostCode               string `json:"PostCode,omitempty"`
	City                   string `json:"City,omitempty"`
	District               string `json:"District,omitempty"`
	Street                 string `json:"Street,omitempty"`
	FromHouseNumber        string `json:"FromHouseNumber,omitempty"`
	ToHouseNumber          string `json:"ToHouseNumber,omitempty"`
	FromCrossroad          string `json:"FromCrossroad,omitempty"`
	ToCrossroad            string `json:"ToCrossroad,omitempty"`
	LicensePlate           string `json:"LicensePlate,omitempty"`
	GeoAreaCoordinates     string `json:"GeoAreaCoordinates,omitempty"`
	GeoOverviewCoordinates string `json:"GeoOverviewCoordinates,omitempty"`
}

type fault struct {
	Faultcode   string `xml:"faultcode"`
	Faultstring string `xml:"faultstring"`
}

type CertificationResponse struct {
	Hash            string
	Upp             string
	Response        string
	VerificationURL string
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

func getAuth(r *http.Request) string {
	auth := r.Header.Get("X-Auth-Token")
	if auth == "" {
		auth = conf.Auth
	}
	return auth
}

func getUuid(r *http.Request) string {
	uuid := r.Header.Get("X-UUID")
	if uuid == "" {
		uuid = conf.Uuid
	}
	return uuid
}

func sendJsonRequest(reqBody []byte, uuid string, auth string) (int, []byte, http.Header, error) {
	client := &http.Client{}

	requURL, err := url.Parse(conf.UbirchClientURL)
	if err != nil {
		return 0, nil, nil, err
	}
	requURL.Path = path.Join(requURL.Path, uuid)

	req, err := http.NewRequest("POST", requURL.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", auth)

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}

	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, respBodyBytes, resp.Header, nil
}

func createSoapResponse(respBody []byte, reqBody []byte) ([]byte, error) {
	var resp CertificationResponse
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}

	resp.VerificationURL = getVerificationURL(reqBody)

	xmlBytes, err := xml.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return xmlBytes, nil
}

func getVerificationURL(reqBody []byte) string {
	var reqMap map[string]string
	err := json.Unmarshal(reqBody, &reqMap)
	if err != nil {
		log.Errorf("unable to create verification URL: %v", err)
	}

	verificationURL := conf.VerificationBaseURL + "#"
	for k, v := range reqMap {
		verificationURL += fmt.Sprintf("%s=%s;", k, v)
	}
	return strings.TrimSuffix(verificationURL, ";")
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

	uuid := getUuid(req)
	if uuid == "" {
		Error(w, "missing UUID", http.StatusUnauthorized)
		return
	}

	auth := getAuth(req)
	if auth == "" {
		Error(w, "missing auth token", http.StatusUnauthorized)
		return
	}

	respCode, respBody, respHeader, err := sendJsonRequest(jsonReq, uuid, auth)
	if err != nil {
		Error(w, fmt.Sprintf("unable to send request: %v", err), http.StatusInternalServerError)
		return
	}

	log.Infof("response: (%d) %s, %s", respCode, respBody, respHeader)

	xmlResp, err := createSoapResponse(respBody, jsonReq)
	if err != nil {
		log.Error(err)
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
	err := conf.load(configFile)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
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
