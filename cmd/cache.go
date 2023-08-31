package cmd

import (
	"github.com/spf13/cobra"

	"github.com/zombull/moo/database"
	"github.com/zombull/moo/server"
)

type cacheOpts struct {
	d      *database.Database
	cache  string
	server string
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

	return cmd
}

func cache(opts *cacheOpts) {
	d := opts.d

	store := server.NewStore(opts.cache, opts.server)

	sets := d.GetSets()
	for _, a := range sets {
		store.Update(opts.d, a)
	}
}
