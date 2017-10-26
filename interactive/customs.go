package interactive

import (
	"github.com/zombull/floating-castle/bug"
	"github.com/zombull/floating-castle/customs"
	"github.com/zombull/floating-castle/database"
)

func import_(d *database.Database) {
	ac := newMapAutocompleter(customs.ImportTypes)
	l := newReader("Import Type: ", ac)
	doReadline(l, true, func(line string) string {
		if f, ok := customs.ImportTypes[line]; ok {
			f(d, getFiles(line))
			return line
		}
		return ""
	})
}

var xx struct{}

func export(d *database.Database) {
	m := database.Set{
		"file": xx,
	}

	t := getSet(m, "Export To")
	if t == "file" {
		bug.On(true, "file not yet implemented")
	}
}
