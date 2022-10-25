package cmdutil

import (
	"bytes"
	"testing"

	"github.com/dnitsch/uistrategy"
)

func Test_YamlParse(t *testing.T) {
	tests := []struct {
		name         string
		confContents []byte
	}{
		// TODO: Add test cases.
		{
			name: "happy path 1",
			confContents: []byte(`
setup: 
  baseUrl: http://127.0.0.1:8090
  timeout: 30
  reuseBrowserCache: false
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
      cssSelector: '#app > div > div > div.page-wrapper.center-content > main > div > button'
      clickSwipe: true`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &uistrategy.UiStrategyConf{}

			err := YamlParseInput(conf, bytes.NewReader(tt.confContents))
			if err != nil {
				t.Error(" error parsing")
			}
			if conf.Setup.BaseUrl == "" {
				t.Error("incorrect URL")
			}
		})
	}

}
