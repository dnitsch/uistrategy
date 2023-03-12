package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnitsch/uistrategy"
	"github.com/dnitsch/uistrategy/internal/util"
	"gopkg.in/yaml.v2"
)

func helperTestSeed(conf *uistrategy.UiStrategyConf) string {
	b, _ := yaml.Marshal(conf)
	dir, _ := os.MkdirTemp("", "uiseeder-test")
	file := filepath.Join(dir, "uiseeder.yml")
	_ = os.WriteFile(file, b, 0777)
	return file
}

func Test_runActions(t *testing.T) {
	tests := map[string]struct {
		path           func(t *testing.T, baseUrl string) string
		handler        func(t *testing.T) http.Handler
		additionalArgs []string
	}{
		"simple test with 1 action": {
			func(t *testing.T, baseUrl string) string {
				conf := &uistrategy.UiStrategyConf{
					Setup: uistrategy.BaseConfig{
						BaseUrl: baseUrl,
						LauncherConfig: &uistrategy.WebConfig{
							Headless: true,
						},
					},
					Actions: []*uistrategy.ViewAction{
						{
							Name:     "test get",
							Navigate: "/get/index.html",
							ElementActions: []*uistrategy.ElementAction{
								{
									Name: "get input id",
									// SkipOnErrorMessage: ,
									CaptureOutput: false,
									Element: uistrategy.Element{
										Selector: util.Str("//*/input"),
									},
								},
							},
						},
					},
				}
				return helperTestSeed(conf)
			},
			func(t *testing.T) http.Handler {
				mux := http.NewServeMux()
				mux.HandleFunc("/get/index.html", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write([]byte(`<html><body><input id="234"/></body></html>`))
				})
				return mux
			},
			[]string{},
		},
		"verbose with 1 action": {
			func(t *testing.T, baseUrl string) string {
				conf := &uistrategy.UiStrategyConf{
					Setup: uistrategy.BaseConfig{
						BaseUrl: baseUrl,
						LauncherConfig: &uistrategy.WebConfig{
							Headless: true,
						},
					},
					Actions: []*uistrategy.ViewAction{
						{
							Name:     "test get",
							Navigate: "/get/index.html",
							ElementActions: []*uistrategy.ElementAction{
								{
									Name: "get input id",
									// SkipOnErrorMessage: ,
									CaptureOutput: false,
									Element: uistrategy.Element{
										Selector: func(input string) *string { return &input }("//*/input"),
									},
								},
							},
						},
					},
				}
				return helperTestSeed(conf)
			},
			func(t *testing.T) http.Handler {
				mux := http.NewServeMux()
				mux.HandleFunc("/get/index.html", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.Write([]byte(`<html><body><input id="234"/></body></html>`))
				})
				return mux
			},
			[]string{"-v"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler(t))
			defer ts.Close()

			path = tt.path(t, ts.URL)
			stdout, errout := &bytes.Buffer{}, &bytes.Buffer{}
			cmd := rootCmd
			args := []string{"-i", path}
			args = append(args, tt.additionalArgs...)
			cmd.SetArgs(args)
			cmd.SetErr(errout)
			cmd.SetOut(stdout)
			if _, err := cmd.ExecuteC(); err != nil {
				t.Errorf("uiseeder cmd failed: %v", err)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	b := new(bytes.Buffer)
	cmd := rootCmd
	cmd.SetArgs([]string{"version"})
	cmd.SetErr(b)
	cmd.SetOut(b)
	_, err := cmd.ExecuteC()
	if err != nil {
		t.Errorf("...")
	}
	out, _ := io.ReadAll(b)
	if !strings.Contains(string(out), "Version:") {
		t.Errorf("version not shown correctly")
	}
	if !strings.Contains(string(out), "Revision:") {
		t.Errorf("revision not shown correctly")
	}
}
