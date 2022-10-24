package uistrategy

import (
	"context"
	"fmt"
	"os"
	"testing"

	log "github.com/dnitsch/simplelog"
	"github.com/dnitsch/uistrategy/internal/util"
)

var (
	testAuth = Auth{
		Username: Element{
			Must:  true,
			Value: util.Str(`test@example.com`),
			XPath: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[2]/input`),
		},
		RequireConfirm: true,
		Password: Element{
			Must:  true,
			Value: util.Str(`P4s$w0rd123!`),
			XPath: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[3]/input`),
		},
		ConfirmPassword: Element{
			Must:  true,
			Value: util.Str(`P4s$w0rd123!`),
			XPath: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[4]/input`),
		},
		Navigate: `/_/#/login`,
		Submit: Element{
			Must:        true,
			CSSSelector: util.Str(`#app > div > div > div.page-wrapper.full-page.center-content > main > div > form > button`),
		},
	}
	testActions = []ViewAction{
		{
			Name:     "create test collection",
			Navigate: `/_/?#/collections`,
			ElementActions: []ElementAction{{
				Name: "create new collection",
				Element: Element{
					Must:        false,
					CSSSelector: util.Str(`#app > div > div > div.page-wrapper.center-content > main > div > button`),
				},
				ClickSwipe: true,
			},
				{
					Name: "Name it test",
					Element: Element{
						Must:        false,
						CSSSelector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-header > form > div > input`),
						Value:       util.Str(`test`),
					},
					// InputText: util.Str("test"),
				},
				{
					Name: "Save it",
					Element: Element{
						Must:        false,
						CSSSelector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-header > form > div > input`),
						// Value:       util.Str(`test`),
					},
				},
				{
					Name: "Add New Field",
					Element: Element{
						Must:        false,
						CSSSelector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-content > div > div > button`),
					},
					ClickSwipe: true,
				},
				{
					Name: "Name Field testField1",
					Element: Element{
						Must:        false,
						CSSSelector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-content > div > div > div.accordions > div > div > form > div > div:nth-child(2) > div > input`),
						Value:       util.Str(`testField1`),
					},
				},
				{
					Name: "Click Done",
					Element: Element{
						Must:        false,
						CSSSelector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-content > div > div > div.accordions > div > div > form > div > div.col-sm-4.txt-right > div.inline-flex.flex-gap-sm.flex-nowrap > button.btn.btn-sm.btn-outline.btn-expanded-sm`),
					},
					ClickSwipe: true,
				},
				{
					Name: "Click Create collection",
					Element: Element{
						Must:        false,
						CSSSelector: util.Str(`body > div.overlays > div:nth-child(2) > div > div.overlay-panel.overlay-panel-lg.colored-header.compact-header.collection-panel > div.overlay-panel-section.panel-footer > button.btn.btn-expanded`),
					},
					ClickSwipe: true,
				},
			},
		},
	}
	testBaseConfig = BaseConfig{BaseUrl: "http://localhost:8090", Timeout: 30}
)

func Test_DoAuth(t *testing.T) {
	tests := []struct {
		name string
		auth Auth
	}{
		{
			name: "happy path",
			auth: testAuth,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := New(testBaseConfig)
			p, e := ui.DoAuth(tt.auth)
			if e != nil {
				t.Errorf("wanted %v to be <nil>", e)
			}
			fmt.Println(p)
		})
	}
}

func Test_Drive(t *testing.T) {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:8888")
	l := log.New(os.Stderr, log.DebugLvl)
	tests := []struct {
		name string
		auth Auth

		actions []ViewAction
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
