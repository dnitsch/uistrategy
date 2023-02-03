package cmdutil_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/dnitsch/configmanager/pkg/generator"
	"github.com/dnitsch/uistrategy"
	"github.com/dnitsch/uistrategy/internal/cmdutil"
	"github.com/dnitsch/uistrategy/internal/util"
)

type mockConfManger func(input string, config generator.GenVarsConfig) (string, error)

func (m mockConfManger) RetrieveWithInputReplaced(input string, config generator.GenVarsConfig) (string, error) {
	return m(input, config)
}

func TestYamlParse(t *testing.T) {
	tests := map[string]struct {
		name         string
		confContents []byte
		mockConfMgr  func(t *testing.T) mockConfManger
		expect       uistrategy.UiStrategyConf
	}{
		"all config items": {
			confContents: []byte(`
setup: 
  baseUrl: http://127.0.0.1:8090
  continueOnError: true
auth:
  navigate: /_/#/login
  username: 
    value: test@example.com
    selector: //*[@class="app-body"]/div[1]/main/div/form/div[2]/input
  password:
    value: P4s$w0rd123!
    selector: //*[@class="app-body"]/div[1]/main/div/form/div[3]/input
  submit:
    selector: '#app > div > div > div.page-wrapper.full-page.center-content > main > div > form > button'    
actions:
- name: create test collection
  navigate: /_/?#/collections
  elementActions: 
  - name: Ceate new collection
    element: 
      selector: '#app > div > div > div.page-wrapper.center-content > main > div > button'`),
			expect: uistrategy.UiStrategyConf{
				Setup: uistrategy.BaseConfig{
					BaseUrl:         "http://127.0.0.1:8090",
					ContinueOnError: true,
					LauncherConfig:  nil,
				},
				Auth: &uistrategy.Auth{
					Navigate: "/_/#/login",
					Username: uistrategy.Element{
						Value:    util.Str("test@example.com"),
						Timeout:  0,
						Selector: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[2]/input`),
					},
					Password: uistrategy.Element{
						Value:    util.Str("P4s$w0rd123!"),
						Timeout:  0,
						Selector: util.Str(`//*[@class="app-body"]/div[1]/main/div/form/div[3]/input`),
					},
					Submit: uistrategy.Element{
						Value:    nil,
						Timeout:  0,
						Selector: util.Str(`#app > div > div > div.page-wrapper.full-page.center-content > main > div > form > button`),
					},
				},
				Actions: []*uistrategy.ViewAction{{
					Name:     "create test collection",
					Navigate: "/_/?#/collections",
					ElementActions: []*uistrategy.ElementAction{{
						Name: "Ceate new collection",
						Element: uistrategy.Element{
							Selector: util.Str(`#app > div > div > div.page-wrapper.center-content > main > div > button`),
							Must:     false,
							Timeout:  0,
							Value:    nil,
						},
					}},
				},
				},
			},
			mockConfMgr: func(t *testing.T) mockConfManger {
				return mockConfManger(func(input string, config generator.GenVarsConfig) (string, error) {
					return `
setup: 
  baseUrl: http://127.0.0.1:8090
  continueOnError: true
auth:
  navigate: /_/#/login
  username: 
    value: test@example.com
    selector: //*[@class="app-body"]/div[1]/main/div/form/div[2]/input
  password:
    value: P4s$w0rd123!
    selector: //*[@class="app-body"]/div[1]/main/div/form/div[3]/input
  submit:
    selector: '#app > div > div > div.page-wrapper.full-page.center-content > main > div > form > button'    
actions:
- name: create test collection
  navigate: /_/?#/collections
  elementActions: 
  - name: Ceate new collection
    element: 
      selector: '#app > div > div > div.page-wrapper.center-content > main > div > button'`, nil
				})
			},
		},
		"no auth provided": {
			confContents: []byte(`
setup: 
  baseUrl: http://127.0.0.1:8090  
actions:
- name: create new marketplace app
  navigate: /marketplace
  elementActions: 
  - name: Click MarketPlace
    element: 
      selector: 'body > div.application-main > main > div.MarketplaceHeader.pt-6.pt-lg-10.position-relative.color-bg-default > div.container-lg.p-responsive.text-center.text-md-left > div > div > a'
`),
			expect: uistrategy.UiStrategyConf{
				Setup: uistrategy.BaseConfig{
					BaseUrl:         "http://127.0.0.1:8090",
					ContinueOnError: false,
					LauncherConfig:  nil,
				},
				Auth: nil,
				Actions: []*uistrategy.ViewAction{
					{
						Name:     "create new marketplace app",
						Navigate: "/marketplace",
						ElementActions: []*uistrategy.ElementAction{{
							Name: "Click MarketPlace",
							Element: uistrategy.Element{
								Selector: util.Str(`body > div.application-main > main > div.MarketplaceHeader.pt-6.pt-lg-10.position-relative.color-bg-default > div.container-lg.p-responsive.text-center.text-md-left > div > div > a`),
								Timeout:  0,
								Must:     false,
								Value:    nil,
							},
						}},
					},
				},
			},
			mockConfMgr: func(t *testing.T) mockConfManger {
				return mockConfManger(func(input string, config generator.GenVarsConfig) (string, error) {
					return `
setup: 
  baseUrl: http://127.0.0.1:8090  
actions:
  - name: create new marketplace app
    navigate: /marketplace
    elementActions: 
    - name: Click MarketPlace
      element: 
        selector: 'body > div.application-main > main > div.MarketplaceHeader.pt-6.pt-lg-10.position-relative.color-bg-default > div.container-lg.p-responsive.text-center.text-md-left > div > div > a'
`, nil
				})
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			conf := &uistrategy.UiStrategyConf{}

			err := cmdutil.YamlParseInput(conf, bytes.NewReader(tt.confContents), tt.mockConfMgr(t))
			if err != nil {
				t.Fatalf("failed to parse %q with err: %v", tt.confContents, err.Error())
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
