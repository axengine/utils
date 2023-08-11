package utils

import (
	"encoding/json"
	"os"
)

func JsonPretty(v interface{}) string {
	bz, _ := json.MarshalIndent(v, "", "  ")
	return string(bz)
}

func JsonPrettyToStdout(v interface{}) {
	en := json.NewEncoder(os.Stdout)
	en.SetIndent("", " ")
	_ = en.Encode(v)
}
