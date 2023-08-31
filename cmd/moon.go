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
	d      *database.Database
	source string
}

func moonCmd(db func() *database.Database, cache string) *cobra.Command {
	opts := moonOpts{
		d:      db(),
		source: path.Join(cache, "source"),
	}

	cmd := &cobra.Command{
		Use:  "moon (Sync JSON from SQL DB)",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			moon(&opts)
		},
	}
	return cmd
}

func moon(opts *moonOpts) {
	problems, err := ioutil.ReadFile(path.Join(opts.source, "Problem.json"))
	bug.OnError(err)

	holds, err := ioutil.ReadFile(path.Join(opts.source, "Move.json"))
	bug.OnError(err)

	moonboard.SyncProblemsJSONv2(opts.d, problems, holds)
}
