package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERROR_MALFORMED_REQ_LINE = fmt.Errorf("invalid request line")

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rl, err := parseRequestLine(string(data))
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: rl,
	}, nil
}

func parseRequestLine(data string) (RequestLine, error) {
	idx := strings.Index(data, "\r\n")
	if idx == -1 {
		return RequestLine{}, ERROR_MALFORMED_REQ_LINE
	}

	parts := strings.Split(data[:idx], " ")
	if len(parts) != 3 {
		return RequestLine{}, ERROR_MALFORMED_REQ_LINE
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return RequestLine{}, ERROR_MALFORMED_REQ_LINE
	}
	return RequestLine{
		HttpVersion:   httpParts[1],
		RequestTarget: parts[1],
		Method:        parts[0],
	}, nil
}
