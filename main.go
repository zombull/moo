package main

import (
	"path"

	"github.com/zombull/moo/cmd"
	"github.com/zombull/moo/database"
	"github.com/zombull/moo/moonboard"
)

func main() {
	c := loadConfig()

	db := func() *database.Database {
		d := database.Init(path.Join(c.Database, "moo.db"))
		moonboard.Init(d)
		return d
	}

	cmd.Run(db, c.Cache, c.Server)
}
