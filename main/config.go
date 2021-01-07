package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Uuid                string             `json:"uuid"`
	Auth                string             `json:"auth"`
	VerificationBaseURL string             `json:"verificationBaseURL"`
	UbirchClientURL     string             `json:"ubirchClientURL"`
	XmlMapping          *map[string]string `json:"xmlMapping"`
}

var defaultXmlMapping = map[string]string{
	"ActionReferenceNumber":  "arn",
	"ActionID":               "id",
	"SpecialUseDesc":         "ud",
	"PeriodBeginDate":        "bd",
	"PeriodBeginTime":        "bt",
	"PeriodEndDate":          "ed",
	"PeriodEndTime":          "et",
	"PostCode":               "pc",
	"City":                   "c",
	"District":               "d",
	"Street":                 "s",
	"FromHouseNumber":        "fn",
	"ToHouseNumber":          "tn",
	"FromCrossroad":          "fc",
	"ToCrossroad":            "tc",
	"LicensePlate":           "lp",
	"GeoAreaCoordinates":     "gac",
	"GeoOverviewCoordinates": "goc",
}

func (c *Config) Load(configDir string, filename string) error {
	var err error
	if os.Getenv("UBIRCH_UBIRCH_CLIENT_URL") != "" {
		err = c.loadEnv()
	} else {
		err = c.loadFile(filepath.Join(configDir, filename))
	}
	if err == nil && c.XmlMapping == nil {
		c.XmlMapping = &defaultXmlMapping
	}
	return err
}

// loadEnv reads the configuration from environment variables
func (c *Config) loadEnv() error {
	return envconfig.Process("ubirch", c)
}

// LoadFile reads the configuration from a json file
func (c *Config) loadFile(filename string) error {
	contextBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(contextBytes, c)
}
