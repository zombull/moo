package cmd

import (
	"path"

	"github.com/spf13/cobra"

	"github.com/zombull/moo/database"
	"github.com/zombull/moo/server"
)

type cacheOpts struct {
	d      *database.Database
	cache  string
	server string
	update bool
}

func cacheCmd(db func() *database.Database, c, s string) *cobra.Command {
	opts := cacheOpts{
		d:      db(),
		cache:  c,
		server: s,
	}

	cmd := &cobra.Command{
		Use:  "cache <FLAGS>",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cache(&opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.update, "update", "u", false, "Update cache from database")

	return cmd
}

func cache(opts *cacheOpts) {
	if opts.update {
		store := server.NewStore(path.Join(opts.cache, "moonboard"))
		store.Update(opts.d, opts.server, "moonboard")
	}
}
