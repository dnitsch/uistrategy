package uistrategy_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/dnitsch/simplelog"
	"github.com/dnitsch/uistrategy"
	"github.com/dnitsch/uistrategy/internal/util"
)

func Test_NoAuthSimulate(t *testing.T) {
	l := log.New(&bytes.Buffer{}, log.DebugLvl)
	tests := map[string]struct {
		name string
		auth *uistrategy.Auth

		actions  []*uistrategy.ViewAction
		handler  func(t *testing.T) http.Handler
		baseConf func(t *testing.T, url string) uistrategy.BaseConfig
		expect   string
	}{
		"happy path - stop on error": {
			auth: nil,
			handler: func(t *testing.T) http.Handler {
				mux := http.NewServeMux()
				mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {

					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write(pocketBaseStyle)
				})
				return mux
			},
			baseConf: func(t *testing.T, url string) uistrategy.BaseConfig {
				return uistrategy.BaseConfig{BaseUrl: url, ContinueOnError: false, LauncherConfig: &uistrategy.WebConfig{Headless: true}}
			},
			actions: []*uistrategy.ViewAction{
				{
					Name:     "test route",
					Navigate: `/route`,
					ElementActions: []*uistrategy.ElementAction{
						{
							Name:   "Click Button",
							Assert: true,
							Element: uistrategy.Element{
								Selector: util.Str(`//*[@id="app"]/div/div/aside/footer/button/./span[text() = 'New collection']`),
							},
						},
						{
							Name: "click test collection - just in case",
							Element: uistrategy.Element{
								Selector: util.Str(`//*[@class='sidebar-content']/*[contains(., 'test')]/span`),
							},
						},
						{
							Name: "assert field testField1 is created",
							Element: uistrategy.Element{
								Selector: util.Str(`//*[@class='page-wrapper']//span[contains(., 'testField1')]`),
							},
							Assert: true,
						},
					},
				}},
		},
		"error on network": {
			auth: nil,
			handler: func(t *testing.T) http.Handler {
				mux := http.NewServeMux()
				mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`<html><body><div>Error Occured</div></body></html>`))
				})
				return mux
			},
			baseConf: func(t *testing.T, url string) uistrategy.BaseConfig {
				return uistrategy.BaseConfig{BaseUrl: url, ContinueOnError: false, LauncherConfig: &uistrategy.WebConfig{Headless: true}}
			},
			actions: []*uistrategy.ViewAction{{
				Name:     "navigate to error",
				Navigate: `/error`,
				ElementActions: []*uistrategy.ElementAction{
					{
						Name:   "asset collection is created and present in sidebar",
						Assert: true,
						Element: uistrategy.Element{
							Selector: util.Str(`//*[@class='sidebar-content']/*[contains(., 'test')]/span`),
						},
					},
				}},
			},
			expect: "following errors occured:\n\n\tin view: asset collection is created and present in sidebar, performing action: //*[@class='sidebar-content']/*[contains(., 'test')]/span, failed on: element not found",
		},
		"not found but continueOnError": {
			auth: nil,
			handler: func(t *testing.T) http.Handler {
				mux := http.NewServeMux()
				mux.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {

					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write(pocketBaseStyle)
				})
				return mux
			},
			baseConf: func(t *testing.T, url string) uistrategy.BaseConfig {
				return uistrategy.BaseConfig{BaseUrl: url, ContinueOnError: true, LauncherConfig: &uistrategy.WebConfig{Headless: true}}
			},
			actions: []*uistrategy.ViewAction{
				{
					Name:     "test route",
					Navigate: `/route`,
					ElementActions: []*uistrategy.ElementAction{
						{
							Name:   "Found element",
							Assert: true,
							Element: uistrategy.Element{
								Selector: util.Str(`//*[@id="app"]/div/div/aside/footer/button/./span[text() = 'New collection']`),
							},
						},
						{
							Name:   "not found id",
							Assert: true,
							Element: uistrategy.Element{
								Selector: util.Str(`#notfound`),
							},
						},
						{
							Name: "assert field testField1 is created",
							Element: uistrategy.Element{
								Selector: util.Str(`//*[@class='page-wrapper']//span[contains(., 'testField1')]`),
							},
							Assert: true,
						},
					},
				}},
			expect: "following errors occured:\n\n\tin view: not found id, performing action: #notfound, failed on: element not found",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler(t))
			defer ts.Close()
			ui := uistrategy.New(tt.baseConf(t, ts.URL)).WithLogger(l)
			_, err := ui.Drive(context.TODO(), tt.auth, tt.actions)
			if err != nil {
				if err.Error() != tt.expect {
					t.Errorf("got: %v\n\nwant: %v", err, nil)
				}
				return
			}
		})
	}
}
