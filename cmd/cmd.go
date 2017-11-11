package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zombull/floating-castle/bug"
	"github.com/zombull/floating-castle/database"
)

var rootCmd = &cobra.Command{
	Use:   "floating-castle <COMMAND>",
	Short: "TODO",
	Long:  `Even more verbose`,
}

func Run(db func() *database.Database, cache, server string) {
	rootCmd.AddCommand(moonCmd(db, cache))
	rootCmd.AddCommand(serveCmd(server, cache))
	rootCmd.AddCommand(cacheCmd(db, cache, server))

	err := rootCmd.Execute()
	bug.OnError(err)
}
