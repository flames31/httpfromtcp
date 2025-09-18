package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string]string
}

var crlf = []byte("\r\n")

var (
	ErrInvalidHeaderFieldName = fmt.Errorf("invalid header field name")
	ErrInvalidHeaderFormat    = fmt.Errorf("invalid header format")
	ErrInvalidHeaderName      = fmt.Errorf("invalid header character in name")
)

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	done := false
	read := 0
	for {
		n := bytes.Index(data[read:], crlf)
		if n == -1 {
			break
		}

		if n == 0 {
			done = true
			break
		}

		name, val, err := parseHeaderLine(data[read : read+n])
		if err != nil {
			return 0, false, err
		}

		h.Set(name, val)
		read += n + len(crlf)
	}
	return read, done, nil
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if val, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s,%s", val, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) Replace(key, val string) {
	h.headers[strings.ToLower(key)] = val
}

func (h *Headers) ForEach(cb func(k, v string)) {
	for k, v := range h.headers {
		cb(k, v)
	}
}
func parseHeaderLine(data []byte) (string, string, error) {
	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ErrInvalidHeaderFormat
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if !headerNameValid(name) {
		return "", "", ErrInvalidHeaderName
	}
	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", ErrInvalidHeaderFieldName
	}

	return string(name), string(value), nil
}

func headerNameValid(name []byte) bool {
	for _, c := range name {
		found := false
		switch c {
		case '!', '#', '$', '%', '&', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}

		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') && (c >= '0' || c <= '9') {
			found = true
		}

		if !found {
			return false
		}
	}

	return true
}
