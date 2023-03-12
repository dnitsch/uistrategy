package cmdutil

import (
	"context"
	"io"

	"github.com/dnitsch/configmanager"
	"github.com/dnitsch/configmanager/pkg/generator"
	"github.com/dnitsch/uistrategy"
)

type ConfManager interface {
	RetrieveWithInputReplaced(input string, config generator.GenVarsConfig) (string, error)
}

func RunActions(ui *uistrategy.Web, conf *uistrategy.UiStrategyConf) error {
	return ui.Drive(context.Background(), conf.Auth, conf.Actions)
}

// YamlParseInput will return a filled pointer with Unmarshalled data
func YamlParseInput[T any](input *T, conf io.Reader, cm ConfManager) error {
	b, err := io.ReadAll(conf)
	if err != nil {
		return err
	}

	// use custom token separator inline with future releases
	config := generator.NewConfig().WithTokenSeparator("://")
	if _, err := configmanager.RetrieveUnmarshalledFromYaml(b, input, cm, *config); err != nil {
		return err
	}
	return nil
}
