package logger

import (
	"encoding/json"
	"fmt"
	"os"
)

type Writer interface {
	Write(msg Message)
}

type JSONWriter struct{}

func (w *JSONWriter) Write(msg Message) {
	data, _ := json.MarshalIndent(msg, "", "   ")
	fmt.Fprintln(os.Stdout, string(data))
}
