package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_runActions_integration(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "integration.yml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path = "../../.ignore-timesheets.yml"
			verbose = true
			if e := runActions(&cobra.Command{}, []string{}); e != nil {
				t.Errorf("%v", e)
			}
		})
	}
}
