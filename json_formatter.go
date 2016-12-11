package logrus

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// HTMLEscapingDisabled sets whether >, < and & inside JSON stings should be escaped.
	HTMLEscapingDisabled bool
}

func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	data["time"] = entry.Time.Format(timestampFormat)
	data["msg"] = entry.Message
	data["level"] = entry.Level.String()

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}

	if f.HTMLEscapingDisabled {
		serialized = bytes.Replace(serialized, []byte("\\u003c"), []byte("<"), -1)
		serialized = bytes.Replace(serialized, []byte("\\u003e"), []byte(">"), -1)
		serialized = bytes.Replace(serialized, []byte("\\u0026"), []byte("&"), -1)
	}

	return append(serialized, '\n'), nil
}
