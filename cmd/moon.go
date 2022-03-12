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
	sql   string
}

func moonCmd(db func() *database.Database, cache string) *cobra.Command {
	opts := moonOpts{
		d:     db(),
		cache: path.Join(cache, "www.moonboard.com"),
	}

	cmd := &cobra.Command{
		Use:  "moon <FLAGS>",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			moon(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.json, "json", "j", "", "Sync Moonboard problems from JSON (specify year)")
	cmd.Flags().StringVarP(&opts.xfer, "xfer", "x", "", "Transfer Moonboard problems a different database (same schema)")
	cmd.Flags().StringVarP(&opts.sql,  "import", "i", "", "Import Moonboard problems from their database")
	return cmd
}

func moon(opts *moonOpts) {

	// Sync problems from JSON files.
	if len(opts.json) > 0 {
		dir := path.Join(opts.cache, "moon" + opts.json)

		infos, err := ioutil.ReadDir(dir)
		bug.OnError(err)

		for _, fi := range infos {
			if fi.Mode().IsRegular() {
				name := path.Join(dir, fi.Name())

				if strings.HasPrefix(fi.Name(), "problems") && strings.HasSuffix(fi.Name(), ".json") {
					data, err := ioutil.ReadFile(name)
					bug.OnError(err)
					fmt.Printf("Syncing: %s\n", fi.Name())
					moonboard.SyncProblemsJSON(opts.d, opts.json, data)
				}
			}
		}
	}

	if len(opts.xfer) > 0 {
		src := database.Init(opts.xfer)

		moonboard.Transfer(opts.d, src)
	}
}
