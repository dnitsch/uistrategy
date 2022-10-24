package cmdutil

import (
	"os"

	"github.com/dnitsch/uistrategy"
	"gopkg.in/yaml.v2"
)

func RunActions(uistrategy *uistrategy.Web, conf *uistrategy.UiStrategyConf) error {
	return nil
}

// YamlParseInput will return a filled pointer with Unmarshalled data
func YamlParseInput[T any](input *T, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, input)
}
