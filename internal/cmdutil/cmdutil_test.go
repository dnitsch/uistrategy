package cmdutil

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/dnitsch/uistrategy"
	"github.com/dnitsch/uistrategy/internal/util"
)

func Test_YamlParse(t *testing.T) {
	tests := []struct {
		name         string
		confContents []byte
		expect       uistrategy.UiStrategyConf
	}{
		{
			name: "all config items",
			confContents: []byte(`
setup: 
  baseUrl: http://127.0.0.1:8090
  continueOnError: true
auth:
  navigate: /_/#/login
  username: 
    must: true
    value: test@example.com
    xPath: //*[@class="app-body"]/div[1]/main/div/form/div[2]/input
  password:
    must: true
    value: P4s$w0rd123!
    xPath: //*[@class="app-body"]/div[1]/main/div/form/div[3]/input
  submit:
    must: true
    cssSelector: '#app > div > div > div.page-wrapper.full-page.center-content > main > div > form > button'    
actions:
- name: create test collection
  navigate: /_/?#/collections
  elementActions: 
  - name: Ceate new collection
    element: 
      cssSelector: '#app > div > div > div.page-wrapper.center-content > main > div > button'`),
			expect: uistrategy.UiStrategyConf{
				Setup: uistrategy.BaseConfig{
					BaseUrl:         "http://127.0.0.1:8090",
					ContinueOnError: true,
					WebConfig:       nil,
				},
				Auth: &uistrategy.Auth{
					Navigate: "/_/#/login",
					Username: uistrategy.Element{
						Value:    util.Str("test@example.com"),
						Timeout:  0,
						Must:     true,
						XPath:    util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[2]/input`),
						Selector: nil,
					},
					Password: uistrategy.Element{
						Value:    util.Str("P4s$w0rd123!"),
						Timeout:  0,
						Must:     true,
						XPath:    util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[3]/input`),
						Selector: nil,
					},
					Submit: uistrategy.Element{
						Must:     true,
						Value:    nil,
						Timeout:  0,
						Selector: util.Str(`#app > div > div > div.page-wrapper.full-page.center-content > main > div > form > button`),
						XPath:    nil,
					},
				},
				Actions: []uistrategy.ViewAction{
					{
						Name:     "create test collection",
						Navigate: "/_/?#/collections",
						ElementActions: []uistrategy.ElementAction{{
							Name: "Ceate new collection",
							Element: uistrategy.Element{
								Selector: util.Str(`#app > div > div > div.page-wrapper.center-content > main > div > button`),
								Must:     false,
								XPath:    nil,
								Timeout:  0,
								Value:    nil,
							},
						}},
					},
				},
			},
		},
		{
			name: "no auth provided",
			confContents: []byte(`
setup: 
  baseUrl: http://127.0.0.1:8090  
actions:
- name: create new marketplace app
  navigate: /marketplace
  elementActions: 
  - name: Click MarketPlace
    element: 
      cssSelector: 'body > div.application-main > main > div.MarketplaceHeader.pt-6.pt-lg-10.position-relative.color-bg-default > div.container-lg.p-responsive.text-center.text-md-left > div > div > a'
`),
			expect: uistrategy.UiStrategyConf{
				Setup: uistrategy.BaseConfig{
					BaseUrl:         "http://127.0.0.1:8090",
					ContinueOnError: false,
					WebConfig:       nil,
				},
				Auth: nil,
				Actions: []uistrategy.ViewAction{
					{
						Name:     "create new marketplace app",
						Navigate: "/marketplace",
						ElementActions: []uistrategy.ElementAction{{
							Name: "Click MarketPlace",
							Element: uistrategy.Element{
								Selector: util.Str(`body > div.application-main > main > div.MarketplaceHeader.pt-6.pt-lg-10.position-relative.color-bg-default > div.container-lg.p-responsive.text-center.text-md-left > div > div > a`),
								Timeout:  0,
								Must:     false,
								XPath:    nil,
								Value:    nil,
							},
						}},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &uistrategy.UiStrategyConf{}

			err := YamlParseInput(conf, bytes.NewReader(tt.confContents))
			if err != nil {
				t.Error(" error parsing")
			}

			if !reflect.DeepEqual(conf.Setup, tt.expect.Setup) {
				t.Errorf("setup unequal got:\n%+v\nwant:\n%+v", conf.Setup, tt.expect.Setup)
			}
			if !reflect.DeepEqual(conf.Auth, tt.expect.Auth) {
				t.Errorf("auth unequal got:\n%+v\nwant:\n%+v", conf.Auth, tt.expect.Auth)
			}
			if !reflect.DeepEqual(conf.Actions, tt.expect.Actions) {
				t.Errorf("actions unequal got:\n%+v\nwant:\n%+v", conf.Actions, tt.expect.Actions)
			}
		})
	}
}
