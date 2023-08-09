package utils

import (
	"encoding/json"
	"os"
)

func JsonPretty(v interface{}) string {
	bz, _ := json.Marshal(v)
	return string(bz)
}

func JsonPrettyToStdout(v interface{}) {
	_ = json.NewEncoder(os.Stdout).Encode(v)
}
