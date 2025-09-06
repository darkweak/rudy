package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Server helper.
type Server struct{}

// GetRequiredFlags returns the server required flags.
func (*Server) GetRequiredFlags() []string {
	return []string{}
}

// ServeHTTP handle any request.
func (*Server) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write([]byte("Hello"))
}

// SetFlags set the available flags.
func (*Server) SetFlags(_ *pflag.FlagSet) {}

// GetArgs return the args.
func (*Server) GetArgs() cobra.PositionalArgs {
	return nil
}

// GetDescription returns the command description.
func (*Server) GetDescription() string {
	return "Run the rudy web server"
}

// GetLongDescription returns the command long description.
func (*Server) GetLongDescription() string {
	return "Run the rudy web server"
}

// Info returns the command name.
func (*Server) Info() string {
	return "server"
}

// Run executes the script associated to the command.
func (s *Server) Run() RunCmd {
	return func(_ *cobra.Command, _ []string) {
		_ = http.ListenAndServe(":8081", s) //nolint:gosec // not relevant because that's for testing purpose.
	}
}

func newServer() command {
	return &Server{}
}

var (
	_ command             = (*Server)(nil)
	_ commandInstanciator = newServer
)
