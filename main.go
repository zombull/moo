package main

import (
	"os"

	"github.com/zombull/floating-castle/cmd"
	"github.com/zombull/floating-castle/database"
	"github.com/zombull/floating-castle/interactive"
	"github.com/zombull/floating-castle/moonboard"
	"github.com/zombull/floating-castle/server"
)

func main() {
	c := loadConfig()

	d := database.Init(c.Database)

	moonboard.Init(d, c.MoonboardSet)

	s := server.Init(c.Server)

	if len(os.Args) == 1 {
		interactive.Run(d, s)
	} else {
		cmd.Run(d, s, c.Cache)
	}
}
