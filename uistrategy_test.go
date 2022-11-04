package uistrategy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	log "github.com/dnitsch/simplelog"
	"github.com/dnitsch/uistrategy/internal/util"
)

var (
	testAuth = &Auth{
		Username: Element{
			Must:     true,
			Value:    util.Str(`test@example.com`),
			Selector: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[2]/input`),
		},
		RequireConfirm: true,
		Password: Element{
			Must:     true,
			Value:    util.Str(`P4s$w0rd123!`),
			Selector: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[3]/input`),
		},
		ConfirmPassword: Element{
			Must:     true,
			Value:    util.Str(`P4s$w0rd123!`),
			Selector: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[4]/input`),
		},
		Navigate: `/_/#/login`,
		Submit: Element{
			Must:     true,
			Selector: util.Str(`#app > div > div > div.page-wrapper.full-page.center-content > main > div > form > button`),
		},
	}
	testActions = []*ViewAction{
		{
			Name:     "create test collection",
			Navigate: `/_/?#/collections`,
			ElementActions: []*ElementAction{{
				Name: "create new collection",
				Element: Element{
					Must:     false,
					Selector: util.Str(`#app > div > div > div.page-wrapper.center-content > main > div > button`),
				},
			},
				{
					Name: "Name it test",
					Element: Element{
						Must:     false,
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-header > form > div > input`),
						Value:    util.Str(`test`),
					},
					// InputText: util.Str("test"),
				},
				{
					Name: "Save it",
					Element: Element{
						Must:     false,
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-header > form > div > input`),
						// Value:       util.Str(`test`),
					},
				},
				{
					Name: "Add New Field",
					Element: Element{
						Must:     false,
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-content > div > div > button`),
					},
				},
				{
					Name: "Name Field testField1",
					Element: Element{
						Must:     false,
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-content > div > div > div.accordions > div > div > form > div > div:nth-child(2) > div > input`),
						Value:    util.Str(`testField1`),
					},
				},
				{
					Name: "Click Done",
					Element: Element{
						Must:     false,
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-content > div > div > div.accordions > div > div > form > div > div.col-sm-4.txt-right > div.inline-flex.flex-gap-sm.flex-nowrap > button.btn.btn-sm.btn-outline.btn-expanded-sm`),
					},
				},
				{
					Name: "Click Create collection",
					Element: Element{
						Must:     false,
						Selector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-footer > button.btn.btn-expanded`),
					},
				},
			},
		},
	}
	testBaseConfig = BaseConfig{BaseUrl: "http://localhost:8090", ContinueOnError: false}
)

func Test_DoAuth(t *testing.T) {
	tests := []struct {
		name string
		auth *Auth
	}{
		{
			name: "happy path",
			auth: testAuth,
		},
		{
			name: "no path",
			auth: &Auth{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := New(testBaseConfig).WithLogger(log.New(os.Stderr, log.DebugLvl))
			p, e := ui.DoAuth(tt.auth)
			if e != nil {
				t.Errorf("wanted %v to be <nil>", e)
			}
			fmt.Println(p)
		})
	}
}

func Test_Drive(t *testing.T) {

	l := log.New(os.Stderr, log.DebugLvl)
	tests := []struct {
		name string
		auth *Auth

		actions []*ViewAction
		web     *Web
	}{
		{
			name:    "happy path",
			auth:    testAuth,
			web:     New(testBaseConfig).WithLogger(l),
			actions: testActions,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.web.Drive(context.TODO(), tt.auth, tt.actions)
			if len(err) > 0 {
				t.Errorf("expected errors to be nil, got %v", err)
			}
		})
	}
}

func getHtmlHandle(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch values, _ := url.ParseQuery(r.URL.RawQuery); values.Get("simulate_resp") {
		case "with_style":
			w.Write(testHtml_style)
		case "no_style":
			w.Write(testHtml_noStyle)
		case "bad_request":
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{}`))
		case "error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{}`))
		default:
			w.Write(testHtml_style)
		}
	}
}

func Test_NoAuthSimulate(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/route", getHtmlHandle(t))

	ts := httptest.NewServer(mux)

	l := log.New(os.Stderr, log.DebugLvl)
	tests := []struct {
		name string
		auth *Auth

		actions []*ViewAction
		web     *Web
	}{
		{
			name: "happy path - no error - stop on error",
			auth: nil,
			web:  New(BaseConfig{BaseUrl: ts.URL, ContinueOnError: false}).WithLogger(l),
			actions: []*ViewAction{
				{
					Name:     "create test collection",
					Navigate: `/route`,
					ElementActions: []*ElementAction{
						{
							Name:   "asset collection is created and present in sidebar",
							Assert: true,
							Element: Element{
								Must:     false,
								Selector: util.Str(`//*[@class='sidebar-content']/*[contains(., 'test')]/span`),
							},
						},
						{
							Name: "click test collection - just in case",
							Element: Element{
								Must:     false,
								Selector: util.Str(`//*[@class='sidebar-content']/*[contains(., 'test')]/span`),
							},
						},
						{
							Name: "assert field testField1 is created",
							Element: Element{
								Must:     false,
								Selector: util.Str(`//*[@class='page-wrapper']//span[contains(., 'testField1')]`),
							},
							Assert: true,
						},
					},
				}},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.web.Drive(context.TODO(), tt.auth, tt.actions)
		})
	}
}
