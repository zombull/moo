package cmd

import (
	"github.com/spf13/cobra"

	"github.com/zombull/floating-castle/database"
	"github.com/zombull/floating-castle/server"
)

type serveOpts struct {
	d      *database.Database
	s      *server.Server
	update bool
	port   string
}

func serveCmd(d *database.Database, s *server.Server) *cobra.Command {
	opts := serveOpts{
		d: d,
		s: s,
	}

	cmd := &cobra.Command{
		Use:  "serve <FLAGS>",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			serve(&opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.update, "update", "u", false, "Update storage from database")
	cmd.Flags().StringVarP(&opts.port, "port", "p", "", "Port to run the server on")

	return cmd
}

func serve(opts *serveOpts) {
	if opts.update {
		opts.s.Update(opts.d)
	} else {
		if len(opts.port) > 0 {
			opts.port = ":" + opts.port
		}
		opts.s.Run(opts.port)
	}
}
