package audit

import (
	"encoding/json"
	"os"
	"time"
)

type Writer struct{ Path string }

type Entry struct {
	TS   time.Time   `json:"ts"`
	Op   string      `json:"op"`
	User string      `json:"user"`
	Data interface{} `json:"data"`
}

func (w Writer) Append(e Entry) error {
	f, err := os.OpenFile(w.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	b, _ := json.Marshal(e)
	b = append(b, '\n')
	_, err = f.Write(b)
	return err
}
