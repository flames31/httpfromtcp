package headers

import (
	"bytes"
	"errors"
)

type Headers map[string]string

var crlf = []byte("\r\n")

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (int, bool, error) {
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

		h[name] = val
		read += n + len(crlf)
	}
	return read, done, nil
}

func parseHeaderLine(data []byte) (string, string, error) {
	parts := bytes.SplitN(data, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid header format")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", errors.New("invalid header field name")
	}

	return string(name), string(value), nil
}
