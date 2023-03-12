package uistrategy_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/dnitsch/simplelog"
	"github.com/dnitsch/uistrategy"
	"github.com/dnitsch/uistrategy/internal/util"
)

func TestDoLoginBasic(t *testing.T) {
	t.Parallel()
	ttests := map[string]struct {
		baseConf func(t *testing.T, url string) uistrategy.BaseConfig
		auth     func(t *testing.T, url string) *uistrategy.Auth
		handler  func(t *testing.T) http.Handler
		// auth     func(t *testing.T, url string) *uistrategy.Auth
	}{
		"local login": {
			func(t *testing.T, url string) uistrategy.BaseConfig {
				return uistrategy.BaseConfig{
					BaseUrl: url,
					LauncherConfig: &uistrategy.WebConfig{
						Headless: true,
					},
				}
			},
			func(t *testing.T, url string) *uistrategy.Auth {
				return &uistrategy.Auth{
					Username: uistrategy.Element{
						Selector: util.Str("#username"),
						Value:    util.Str("test"),
					},
					Password: uistrategy.Element{
						Selector: util.Str("#password"),
						Value:    util.Str("test"),
					},
					Navigate: "/login",
					Submit: uistrategy.Element{
						Selector: util.Str("#submit"),
					},
				}
			},
			func(t *testing.T) http.Handler {
				mux := http.NewServeMux()
				mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write(localLoginHtml)
				})
				mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write([]byte(`<html><body>Hello!</body></html>`))
				})
				return mux
			},
		},
		"idp auth": {
			func(t *testing.T, url string) uistrategy.BaseConfig {
				return uistrategy.BaseConfig{
					BaseUrl: url,
					LauncherConfig: &uistrategy.WebConfig{
						Headless: true,
					},
				}
			},
			func(t *testing.T, url string) *uistrategy.Auth {
				return &uistrategy.Auth{
					Username: uistrategy.Element{
						Selector: util.Str("#username"),
						Value:    util.Str("test"),
					},
					Password: uistrategy.Element{
						Selector: util.Str("#password"),
						Value:    util.Str("test"),
					},
					Navigate:   "/login-idp",
					IdpManaged: true,
					IdpSelector: &uistrategy.Element{
						Selector: util.Str("#idp-login"),
					},
					IdpUrl: url,
					Submit: uistrategy.Element{
						Selector: util.Str("#submit"),
					},
				}
			},
			func(t *testing.T) http.Handler {
				mux := http.NewServeMux()
				mux.HandleFunc("/login-idp", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write(idpLoginHtml)
				})
				mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write([]byte(`<html><body>Hello!</body></html>`))
				})
				return mux
			},
		},
		"mfa Login": {
			func(t *testing.T, url string) uistrategy.BaseConfig {
				return uistrategy.BaseConfig{
					BaseUrl: url,
					LauncherConfig: &uistrategy.WebConfig{
						Headless: true,
					},
				}
			},
			func(t *testing.T, url string) *uistrategy.Auth {
				return &uistrategy.Auth{
					Username: uistrategy.Element{
						Selector: util.Str("#username"),
						Value:    util.Str("test"),
					},
					Password: uistrategy.Element{
						Selector: util.Str("#password"),
						Value:    util.Str("test"),
					},
					Navigate:   "/login-idp",
					IdpManaged: true,
					MfaSelector: &uistrategy.Element{
						Selector: util.Str("#mfa"),
					},
					IdpSelector: &uistrategy.Element{
						Selector: util.Str("#idp-login"),
					},
					IdpUrl: url,
					Submit: uistrategy.Element{
						Selector: util.Str("#submit"),
					},
				}
			},
			func(t *testing.T) http.Handler {
				mux := http.NewServeMux()
				mux.HandleFunc("/login-idp", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write(mfaLogin)
				})
				mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write([]byte(`<html><body>Hello!</body></html>`))
				})
				return mux
			},
		},
	}
	for name, tt := range ttests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler(t))
			defer ts.Close()
			uiWeb := uistrategy.New(tt.baseConf(t, ts.URL)).WithLogger(log.New(&bytes.Buffer{}, log.ErrorLvl))
			lp, err := uiWeb.DoAuth(tt.auth(t, ts.URL))
			if err != nil {
				t.Errorf("failed to do auth: %v", err)
			}
			if lp == nil {
				t.Errorf("logged in page - got nil expected not nil")
			}
		})
	}
}
