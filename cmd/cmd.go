package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zombull/floating-castle/bug"
	"github.com/zombull/floating-castle/database"
	"github.com/zombull/floating-castle/server"
)

var rootCmd = &cobra.Command{
	Use:   "floating-castle <COMMAND>",
	Short: "TODO",
	Long:  `Even more verbose`,
}

func Run(d *database.Database, s *server.Server, cache string) {
	rootCmd.AddCommand(moonCmd(d, cache))
	rootCmd.AddCommand(serveCmd(d, s))

	err := rootCmd.Execute()
	bug.OnError(err)
}
