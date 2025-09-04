package request

import (
	"bytes"
	"fmt"
	"io"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

const SEPERATOR = "\r\n"

var ErrMalformedReqLine = fmt.Errorf("invalid request line")

type Request struct {
	RequestLine RequestLine
	ParserState parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func newRequest() *Request {
	return &Request{
		ParserState: StateInit,
	}
}

func (r *Request) done() bool {
	return r.ParserState == StateDone
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.ParserState {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = rl
			r.ParserState = StateDone
			read += n
		case StateDone:
			break outer
		}
	}
	return read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		parseN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[parseN:bufLen])
		bufLen -= parseN
	}

	return request, nil
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	idx := bytes.Index(data, []byte(SEPERATOR))
	if idx == -1 {
		return RequestLine{}, 0, nil
	}
	startLine := data[:idx]
	n := idx + len(SEPERATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return RequestLine{}, 0, ErrMalformedReqLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return RequestLine{}, 0, ErrMalformedReqLine
	}
	return RequestLine{
		Method:        string(parts[0]),
		HttpVersion:   string(httpParts[1]),
		RequestTarget: string(parts[1]),
	}, n, nil
}
