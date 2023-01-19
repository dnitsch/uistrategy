package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_runActions_integration(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "integration without configmanager",
			path: "../../test/integration.yml",
		},
		// {
		// 	name: "integration with configmanager",
		// 	path: "../../test/integration-with-configmanager.yml",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path = tt.path
			verbose = true
			if e := runActions(&cobra.Command{}, []string{}); e != nil {
				t.Errorf("%v", e)
			}
		})
	}
}
