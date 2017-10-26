package customs

import "github.com/zombull/floating-castle/database"

var ImportTypes = map[string]func(d *database.Database, files []string){
	"gym": ImportGym,
	"set": ImportGymSet,
}
