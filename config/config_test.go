package config_test

import (
	"github.com/willbrid/easy-api-prom-alert-sms/config"

	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func triggerTest(t *testing.T, yamlConfig []byte) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewBuffer([]byte(yamlConfig))); err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	_, err := config.LoadConfig(v, validate)

	var fieldErr validator.FieldError
	if !errors.As(err, &fieldErr) && !strings.Contains(err.Error(), "unable to unmarshal config struct") {
		t.Errorf("wrong error: %v", err)
	}
}

func TestReadConfigFile_ReturnFileNotFoundError(t *testing.T) {
	t.Parallel()

	var filename string

	_, err := config.ReadConfigFile(filename)

	if err == nil {
		t.Error("no error returned for file not found")
	}
}

func TestReadConfigFile_ReturnFileNotExistError(t *testing.T) {
	t.Parallel()

	filename := "nonexistentfile.yaml"
	_, err := config.ReadConfigFile(filename)

	if err == nil {
		t.Error("no error returned for file no exist")
	}
}

func TestLoadConfig_ReturnErrorWithBadUsernameInputWhenAuthEnabled(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "x"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadPasswordInputWhenAuthEnabled(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxx
`),
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadProviderUrl(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    content_type: ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    content_type: "xxxxx"
`),
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadProviderAuthInput(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: ''
      authorization_credential: ''
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: 'xxxxx'
      authorization_credential: ''
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: ''
      authorization_credential: 'xxxxxxxxxx'
`),
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadProviderParamInput(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_value: ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_value: "xxxxx"
        param_method: "xxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "xxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadRecipientInput(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxx"
  recipients:
  - name: ""
    members:
    - "xxxxx"
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxx"
  recipients:
  - name: "admin"
    members:
    - ""
`),
		[]byte(`---
easy_api_prom_sms_alert:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  provider:
    url: "http://localhost:5797"
    authentication:
      enabled: true
      authorization_type: "xxxxx"
      authorization_credential: "xxxxx"
    parameters:
      from:
        param_name: "xxxxx"
        param_method: "query"
        param_value: "xxxxx"
      to:
        param_name: "xxxxx"
        param_method: "post"
        param_value: "xxxxx"
      message:
        param_name: "xxxxx"
  recipients:
  - name: "admin"
    members:
    - "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig)
		})
	}
}
