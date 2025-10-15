package util

import (
	"bytes"

	"github.com/goccy/go-yaml"
)

func YamlMarshalWithOptions(val any) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := yaml.NewEncoder(buf,
		yaml.IndentSequence(false),
		yaml.Indent(2),
		yaml.OmitEmpty(),
		yaml.OmitZero(),
		yaml.Flow(false),
	)

	err := encoder.Encode(val)
	return buf.Bytes(), err
}
