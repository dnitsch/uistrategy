package uistrategy

import (
	"context"
	"fmt"
	"testing"

	"github.com/dnitsch/uistrategy/internal/util"
)

var (
	testAuth = Auth{
		Username: Element{
			Must:  true,
			Value: "test@example.com",
			XPath: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[2]/input`),
		},
		RequireConfirm: true,
		Password: Element{
			Must:  true,
			Value: "P4s$w0rd123!",
			XPath: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[3]/input`),
		},
		ConfirmPassword: Element{
			Must:  true,
			Value: "P4s$w0rd123!",
			XPath: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[4]/input`),
		},
		Navigate: `/_/?installer#`,
		Submit: Element{
			Must:        true,
			Value:       "",
			CSSSelector: util.Str(`#app > div > div > div.page-wrapper.full-page.center-content > main > div > form > button`),
		},
	}
	testActions = []Action{
		{
			Navigate: `/_/?#/collections?collectionId=&filter=&sort=-created`,
			Element: Element{
				Must:        false,
				CSSSelector: util.Str(`#app > div > div > div.page-wrapper.center-content > main > div > button`),
				Value:       "",
			},
			ClickSwipe: true,
			// #app > div > div > div.page-wrapper.center-content > main > div > button
		},
	}
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
			ui := New()
			p, e := ui.DoAuth(tt.auth)
			if e != nil {
				t.Errorf("wanted %v to be <nil>", e)
			}
			fmt.Println(p)
		})
	}
}

func Test_ActionsPerform(t *testing.T) {
	tests := []struct {
		name    string
		auth    Auth
		actions []Action
		web     *Web
	}{
		{
			name:    "happy path",
			auth:    testAuth,
			web:     New(),
			actions: testActions,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.web.ActionsPerform(context.TODO(), tt.auth, tt.actions)
			if len(err) > 0 {
				t.Errorf("expected errors to be nil, got %v", err)
			}
		})
	}
}
