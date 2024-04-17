package httpjson

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"text/template"
)

type Payload struct {
	Name    string            `json:"name"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    map[string]any    `json:"body"`
}

func (p *Payload) Request(ctx context.Context) (*http.Request, error) {
	var body []byte
	if p.Body != nil {
		b, err := json.Marshal(p.Body)
		if err != nil {
			return nil, err
		}
		body = b
	}

	req, err := http.NewRequestWithContext(ctx, p.Method, p.URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func FromJSON(path string, vars map[string]any) (*Payload, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var tpl bytes.Buffer
	t := template.New("main")
	if err := template.Must(t.Parse(string(b))).Execute(&tpl, vars); err != nil {
		return nil, err
	}

	var payload *Payload
	if err := json.NewDecoder(&tpl).Decode(&payload); err != nil {
		return nil, err
	}

	return payload, nil
}
