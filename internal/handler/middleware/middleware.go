package middleware

import (
	"github.com/rs/zerolog"
	"github.com/willbrid/easy-api-prom-alert-sms/config"

	"encoding/base64"
	"net/http"
	"strings"
)

type IAuthMiddleware interface {
	Authenticate(next http.Handler, cfg *config.Config) http.Handler
}

type AuthMiddleware struct {
	logger zerolog.Logger
}

func NewAuthMiddleware(logger zerolog.Logger) *AuthMiddleware {
	return &AuthMiddleware{logger}
}

func (a *AuthMiddleware) Authenticate(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")

		if cfg.EasyAPIPromAlertSMS.Enabled && req.URL.Path != "/healthz" {
			if auth == "" {
				a.logger.Error().Str("config_auth", "enabled").Msg("missing authorization header")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(auth, "Basic ") {
				a.logger.Error().Str("config_auth", "enabled").Msg("malformed authorization header")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Basic ")
			decodedToken, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				a.logger.Error().Str("config_auth", "enabled").Msg("unabled to decode authorization header token")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			credentialParts := strings.SplitN(string(decodedToken), ":", 2)
			username := credentialParts[0]
			password := credentialParts[1]
			if username != cfg.EasyAPIPromAlertSMS.Username || password != cfg.EasyAPIPromAlertSMS.Password {
				a.logger.Error().Str("config_auth", "enabled").Msg("invalid login or password")
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(resp, req)
	})
}
