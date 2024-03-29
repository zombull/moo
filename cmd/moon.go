package cmd

import (
	"io/ioutil"
	"path"

	"github.com/spf13/cobra"

	"github.com/zombull/moo/bug"
	"github.com/zombull/moo/database"
	"github.com/zombull/moo/moonboard"
)

type moonOpts struct {
	db     func() *database.Database
	source string
	purge  bool
}

func moonCmd(db func() *database.Database, cache string) *cobra.Command {
	opts := moonOpts{
		db:     db,
		source: path.Join(cache, "source"),
	}

	cmd := &cobra.Command{
		Use:  "moon (Sync JSON from SQL DB)",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			moon(&opts)
		},
	}
	cmd.Flags().BoolVarP(&opts.purge, "purge", "p", false, "Purge (temporary)")
	return cmd
}

func moon(opts *moonOpts) {
	d := opts.db()

	if (opts.purge) {
		moonboard.Purge(d)
		return
	}
	problems, err := ioutil.ReadFile(path.Join(opts.source, "Problem.json"))
	bug.OnError(err)

	holds, err := ioutil.ReadFile(path.Join(opts.source, "Move.json"))
	bug.OnError(err)

	moonboard.SyncProblemsJSONv2(d, problems, holds)
}
