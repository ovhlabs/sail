package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"os"
)

// FormatOutput autmatically formats json based output based on user choice.
// when selected formatter is "pretty", call prettyFormatter callback.
func FormatOutput(data []byte, prettyFormatter func([]byte)) {
	switch Format {
	case "pretty":
		prettyFormatter(data)
	case "json":
		jsonFormatter(data)
	case "yaml":
		yamlFormatter(data)
	default:
		fmt.Fprintf(os.Stderr, "Invalid formater %s. Use one of 'pretty', 'json', 'yaml'\n", Format)
		return
	}
}

// FormatOutputDef autmatically formats json based output based on user choice.
// uses yamlFormatter as pretty formatter.
func FormatOutputDef(data []byte) {
	FormatOutput(data, yamlFormatter)
}

// FormatOutputError prints the "message" field of an API return or falls back on FormatOutputDef if the field does not exist
func FormatOutputError(data []byte) {
	var err map[string]interface{}
	Check(json.Unmarshal(data, &err))

	message := err["message"]
	if message == nil {
		message = err["error_details"]
	}

	if message != nil {
		fmt.Printf("Error: %s\n", message)
	} else {
		FormatOutputDef(data)
	}
}

func jsonFormatter(data []byte) {
	var out bytes.Buffer
	json.Indent(&out, data, "", "  ")
	fmt.Println(out.String())
}

func yamlFormatter(data []byte) {
	out, err := yaml.JSONToYAML(data)
	Check(err)
	fmt.Print(string(out))
}
