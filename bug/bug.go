package bug

import (
	"fmt"
)

func Bug(s string) {
	panic(fmt.Errorf("BUG: %s", s))
}

func On(b bool, s string) {
	if b {
		Bug(s)
	}
}
func OnError(e error) {
	if e != nil {
		Bug(e.Error())
	}
}

func UserBug(s string) {
	panic(fmt.Errorf("USER: %s", s))
}

func UserBugOn(b bool, s string) {
	if b {
		UserBug(s)
	}
}

func UserError(e error) {
	if e != nil {
		UserBug(e.Error())
	}
}
