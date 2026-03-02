package alert

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/easy-api-prom-alert-sms/config"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpclient"
	"github.com/willbrid/easy-api-prom-alert-sms/pkg/httpparam"

	"fmt"
	"strings"
)

type AlertMicroservice struct {
	Provider      *config.Provider
	iHTTPClient   httpclient.IHTTPClient
	iParamFactory httpparam.IParamFactory
	logger        zerolog.Logger
}

func NewAlertMicroservice(provider *config.Provider, iHTTPClient httpclient.IHTTPClient, iParamFactory httpparam.IParamFactory, logger zerolog.Logger) *AlertMicroservice {
	return &AlertMicroservice{provider, iHTTPClient, iParamFactory, logger}
}

func (ams *AlertMicroservice) Consume(url string, body string) error {
	provider := ams.Provider
	headers := map[string]string{
		"Content-Type": provider.ContentType,
	}

	if provider.Authentication.Enabled {
		headers["Authorization"] = fmt.Sprintf("%s %s", provider.Authentication.AuthorizationType, provider.Authentication.AuthorizationCredential)
	}

	_, err := ams.iHTTPClient.Post(url, strings.NewReader(body), httpclient.Options{
		Headers:            headers,
		Timeout:            provider.Timeout,
		InsecureSkipVerify: provider.InsecureSkipVerify,
	})

	if err != nil {
		ams.logger.Error().Err(err).Msgf("failed to send alert to %s", url)
		return err
	}

	return nil
}

func (ams *AlertMicroservice) GetUrlAndBody(dest string, message string) (string, string, error) {
	httpParam := ams.iParamFactory.GetNewParam()
	providerParams := ams.Provider.Parameters

	httpParam.AddPostParam(providerParams.Message.ParamName, strings.ReplaceAll(providerParams.Message.ParamValue, config.AlertMessageTemplate, message))
	httpParam.AddParam(providerParams.From.ParamMethod, providerParams.From.ParamName, providerParams.From.ParamValue)
	httpParam.AddParam(providerParams.To.ParamMethod, providerParams.To.ParamName, dest)
	if len(providerParams.ExtraParams) > 0 {
		for _, extraParam := range providerParams.ExtraParams {
			httpParam.AddParam(extraParam.ParamMethod, extraParam.ParamName, extraParam.ParamValue)
		}
	}

	encodedURL := ams.Provider.Url
	if len(httpParam.Query) > 0 {
		encodedURL = fmt.Sprintf("%s?%s", encodedURL, httpParam.EncodeQuery())
	}

	body, err := httpParam.EncodePost(ams.Provider.ContentType)
	if err != nil {
		ams.logger.Error().Err(err).Msg("failed to encode post parameters")
		return "", "", err
	}

	return encodedURL, body, nil
}
