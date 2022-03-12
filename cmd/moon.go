package cmd

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zombull/floating-castle/bug"
	"github.com/zombull/floating-castle/database"
	"github.com/zombull/floating-castle/moonboard"
)

type moonOpts struct {
	d     *database.Database
	cache string
	index bool
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

	cmd.Flags().BoolVarP(&opts.index, "index", "i", false, "Sync the Moonboard's index")
	return cmd
}

func moon(opts *moonOpts) {
	onFiles := func(prefix string, f func(data []byte)) {
		infos, err := ioutil.ReadDir(opts.cache)
		bug.OnError(err)

		for _, fi := range infos {
			if fi.Mode().IsRegular() {
				name := path.Join(opts.cache, fi.Name())

				if strings.HasPrefix(fi.Name(), prefix) && strings.HasSuffix(fi.Name(), ".json") {
					data, err := ioutil.ReadFile(name)
					bug.OnError(err)
					fmt.Printf("Syncing: %s\n", fi.Name())
					f(data)
				}
			}
		}
	}

	if opts.index {
		onFiles("problems", func(data []byte) {
			moonboard.SyncProblems(opts.d, data)
		})
	}
}
