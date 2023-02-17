package uistrategy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	log "github.com/dnitsch/simplelog"
	"github.com/dnitsch/uistrategy"
	"github.com/dnitsch/uistrategy/internal/util"
)

var (
	testAuth = &uistrategy.Auth{
		Username: uistrategy.Element{
			Value:    util.Str(`test@example.com`),
			Selector: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[2]/input`),
		},
		RequireConfirm: true,
		Password: uistrategy.Element{

			Value:    util.Str(`P4s$w0rd123!`),
			Selector: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[3]/input`),
		},
		ConfirmPassword: uistrategy.Element{
			Value:    util.Str(`P4s$w0rd123!`),
			Selector: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[4]/input`),
		},
		Navigate: `/_/#/login`,
		Submit: uistrategy.Element{
			Selector: util.Str(`#app > div > div > div.page-wrapper.full-page.center-content > main > div > form > button`),
		},
	}
	testActions = []*uistrategy.ViewAction{
		{
			Name:     "create test collection",
			Navigate: `/_/?#/collections`,
			ElementActions: []*uistrategy.ElementAction{{
				Name: "create new collection",
				Element: uistrategy.Element{
					Selector: util.Str(`#app > div > div > div.page-wrapper.center-content > main > div > button`),
				},
			},
				{
					Name: "Name it test",
					Element: uistrategy.Element{
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-header > form > div > input`),
						Value:    util.Str(`test`),
					},
					// InputText: util.Str("test"),
				},
				{
					Name: "Save it",
					Element: uistrategy.Element{
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-header > form > div > input`),
						// Value:       util.Str(`test`),
					},
				},
				{
					Name: "Add New Field",
					Element: uistrategy.Element{
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-content > div > div > button`),
					},
				},
				{
					Name: "Name Field testField1",
					Element: uistrategy.Element{
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-content > div > div > div.accordions > div > div > form > div > div:nth-child(2) > div > input`),
						Value:    util.Str(`testField1`),
					},
				},
				{
					Name: "Click Done",
					Element: uistrategy.Element{
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-content > div > div > div.accordions > div > div > form > div > div.col-sm-4.txt-right > div.inline-flex.flex-gap-sm.flex-nowrap > button.btn.btn-sm.btn-outline.btn-expanded-sm`),
					},
				},
				{
					Name: "Click Create collection",
					Element: uistrategy.Element{
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-footer > button.btn.btn-expanded`),
					},
				},
			},
		},
	}
	testBaseConfig = uistrategy.BaseConfig{BaseUrl: "http://localhost:8090", ContinueOnError: false}
)

// func Test_DoAuth(t *testing.T) {
// 	tests := map[string]struct {
// 		auth *uistrategy.Auth
// 	}{
// 		"register path": {
// 			auth: testAuth,
// 		},
// 		"no auth": {nil},
// 	}
// 	for name, tt := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			ui := uistrategy.New(testBaseConfig).WithLogger(log.New(os.Stderr, log.DebugLvl))
// 			p, e := ui.DoAuth(tt.auth)
// 			if e != nil {
// 				t.Errorf("wanted %v to be <nil>", e)
// 			}
// 			fmt.Println(p)
// 		})
// 	}
// }

func Test_NoAuthSimulate(t *testing.T) {
	l := log.New(os.Stderr, log.DebugLvl)
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
					w.Write(testHtml_style)
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
							Name:   "asset collection is created and present in sidebar",
							Assert: true,
							Element: uistrategy.Element{
								Selector: util.Str(`//*[@class='sidebar-content']/*[contains(., 'test')]/span`),
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
		},
		// test iframe
		// test with auth
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler(t))
			defer ts.Close()
			ui := uistrategy.New(tt.baseConf(t, ts.URL)).WithLogger(l)
			err := ui.Drive(context.TODO(), tt.auth, tt.actions)
			if err != nil {
				if err.Error() != tt.expect {
					t.Errorf("got: %v\n\nwant: %v", err, nil)
				}
				t.Logf("error: %s \n\nmatches the expected \n\noutput: %s", err.Error(), tt.expect)
				return
			}
		})
	}
}
