package cmdutil

import (
	"context"
	"fmt"
	"io"

	"github.com/dnitsch/uistrategy"
	"gopkg.in/yaml.v2"
)

func RunActions(ui *uistrategy.Web, conf *uistrategy.UiStrategyConf) error {
	if err := ui.Drive(context.Background(), conf.Auth, conf.Actions); len(err) > 0 {
		return fmt.Errorf("%#v", err)
	}
	return nil
}

// YamlParseInput will return a filled pointer with Unmarshalled data
func YamlParseInput[T any](input *T, conf io.Reader) error {
	b, err := io.ReadAll(conf)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, input)
}
