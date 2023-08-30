package cmd

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zombull/moo/bug"
	"github.com/zombull/moo/database"
	"github.com/zombull/moo/moonboard"
)

type moonOpts struct {
	d     *database.Database
	cache string
	json  string
	xfer  string
	sql   bool
}

func moonCmd(db func() *database.Database, cache string) *cobra.Command {
	opts := moonOpts{
		d:     db(),
		cache: path.Join(cache, "source"),
	}

	cmd := &cobra.Command{
		Use:  "moon <FLAGS>",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			moon(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.json, "json", "j", "", "Sync Moonboard problems from JSON (specify year)")
	cmd.Flags().BoolVarP(&opts.sql,  "sql",  "s", true, "JSON is from SQL DB (mb2 layout")
	cmd.Flags().StringVarP(&opts.xfer, "xfer", "x", "", "Transfer Moonboard problems a different database (same schema)")
	return cmd
}

func moonV1(opts *moonOpts, dir string) {
	infos, err := ioutil.ReadDir(dir)
	bug.OnError(err)

	for _, fi := range infos {
		if fi.Mode().IsRegular() {
			name := path.Join(dir, fi.Name())

			if strings.HasPrefix(fi.Name(), "problems") && strings.HasSuffix(fi.Name(), ".json") {
				data, err := ioutil.ReadFile(name)
				bug.OnError(err)
				fmt.Printf("Syncing: %s\n", fi.Name())
				moonboard.SyncProblemsJSONv1(opts.d, opts.json, data)
			}
		}
	}
}

func moon(opts *moonOpts) {
	if len(opts.json) > 0 {
		dir := path.Join(opts.cache, "moon" + opts.json)

		if !opts.sql {
			moonV1(opts, dir)
		} else {
			problems, err := ioutil.ReadFile(path.Join(dir, "Problem.json"))
			bug.OnError(err)

			holds, err := ioutil.ReadFile(path.Join(dir, "Move.json"))
			bug.OnError(err)

			moonboard.SyncProblemsJSONv2(opts.d, opts.json, problems, holds)
		}
	}

	if len(opts.xfer) > 0 {
		src := database.Init(opts.xfer)

		moonboard.Transfer(opts.d, src)
	}
}
