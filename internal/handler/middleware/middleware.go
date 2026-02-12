package middleware

import (
	"github.com/willbrid/easy-api-prom-alert-sms/config"

	"encoding/base64"
	"net/http"
	"strings"
)

type IAuthMiddleware interface {
	Authenticate(next http.Handler, cfg *config.Config) http.Handler
}

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (a *AuthMiddleware) Authenticate(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")

		if cfg.EasyAPIPromAlertSMS.Auth.Enabled && req.URL.Path != "/healthz" {
			if auth == "" {
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(auth, "Basic ") {
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Basic ")
			decodedToken, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}

			credentialParts := strings.SplitN(string(decodedToken), ":", 2)
			username := credentialParts[0]
			password := credentialParts[1]
			if username != cfg.EasyAPIPromAlertSMS.Auth.Username || password != cfg.EasyAPIPromAlertSMS.Auth.Password {
				http.Error(resp, "invalid credential", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(resp, req)
	})
}
