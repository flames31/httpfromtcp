package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/flames31/httpfromtcp/internal/headers"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
)

const SEPERATOR = "\r\n"

var ErrMalformedReqLine = fmt.Errorf("invalid request line")
var ErrRequestBodyLenMismtach = fmt.Errorf("invalid body length")

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        []byte
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
		Headers:     headers.NewHeaders(),
		Body:        make([]byte, 0, 1024),
	}
}

func (r *Request) done() bool {
	return r.ParserState == StateDone
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]

		switch r.ParserState {
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			r.ParserState = StateHeaders
			read += n

		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				return 0, err
			}

			if done {
				read += len(SEPERATOR)
				r.ParserState = StateBody
				continue
			}

			if n == 0 {
				break outer
			}
			read += n

		case StateBody:
			contentLength := r.Headers.Get("content-length")
			if contentLength == "" {
				r.ParserState = StateDone
				break outer
			}

			length, err := strconv.Atoi(contentLength)
			if err != nil {
				return 0, err
			}

			r.Body = append(r.Body, currentData...)
			if len(r.Body) > length {
				return 0, ErrRequestBodyLenMismtach
			}

			if len(r.Body) == length {
				r.ParserState = StateDone
			}

			read += len(currentData)
			break outer

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

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(SEPERATOR))
	if idx == -1 {
		return nil, 0, nil
	}
	startLine := data[:idx]
	n := idx + len(SEPERATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrMalformedReqLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrMalformedReqLine
	}
	return &RequestLine{
		Method:        string(parts[0]),
		HttpVersion:   string(httpParts[1]),
		RequestTarget: string(parts[1]),
	}, n, nil
}
