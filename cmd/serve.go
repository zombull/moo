package cmd

import (
	"github.com/spf13/cobra"

	"github.com/zombull/floating-castle/server"
)

type serveOpts struct {
	server  string
	cache   string
	port    string
	release bool
}

func serveCmd(server, cache string) *cobra.Command {
	opts := serveOpts{
		server: server,
		cache:  cache,
	}

	cmd := &cobra.Command{
		Use:  "serve <FLAGS>",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			serve(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.port, "port", "p", "", "Port to run the server on")
	cmd.Flags().BoolVarP(&opts.release, "release", "r", false, "Use release directories")

	return cmd
}

func serve(opts *serveOpts) {
	s := server.Init(opts.server, opts.cache)

	if len(opts.port) > 0 {
		opts.port = ":" + opts.port
	}
	s.Run(opts.port, opts.release)
}
