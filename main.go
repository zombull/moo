package main

import (
	"os"

	"github.com/zombull/floating-castle/cmd"
	"github.com/zombull/floating-castle/database"
	"github.com/zombull/floating-castle/interactive"
	"github.com/zombull/floating-castle/moonboard"
)

func main() {
	c := loadConfig()

	db := func() *database.Database {
		d := database.Init(c.Database)
		moonboard.Init(d, c.MoonboardSet)
		return d
	}

	if len(os.Args) == 1 {
		interactive.Run(db())
	} else {
		cmd.Run(db, c.Cache, c.Server)
	}
}
