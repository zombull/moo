package server

import (
	"fmt"
	"path"
	"strings"

	"github.com/labstack/echo"

	"github.com/zombull/floating-castle/database"
)

const domain = "zombull.xyz"

type Server struct {
	dir   string
	store *KeyValueStore
}

func Init(root string) *Server {
	return &Server{
		dir: root,
	}
}

func (s *Server) Run(port string, release bool) {
	s.store = newStore(path.Join(s.dir, "moonboard"))
	common := path.Join(s.dir, "common")
	if release {
		common = path.Join(s.dir, "moonboard")
	}

	moon := echo.New()

	index := path.Join(s.dir, "moonboard", "index.html")
	for _, r := range []string{"/", "/:problem", "/p/:grade", "/t/:grade", "/k/:grade", "/j/:grade", "/s/:setter", "/st/:setter"} {
		moon.File(r, index)
	}

	moon.Static("/favicon", path.Join(common, "img", "favicon"))
	moon.Static("/static", path.Join(s.dir, "moonboard"))
	moon.Static("/common", common)
	moon.GET("/data/:key", s.store.getValue("moonboard"))

	echoes := map[string]*echo.Echo{
		"moon": moon,
	}

	subs := map[string]string{
		"dark": "moonboard",
	}

	// Server
	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		h := c.Request().Host
		ss := strings.SplitN(h, ".", 2)
		if len(ss) != 2 || ss[1] != domain+port {
			fmt.Printf("%v\n", ss)
			return echo.ErrNotFound
		}
		h = ss[0]

		if release {
			c.Response().Header().Set("Cache-Control", "private, max-age=31536000")
		}
		if ee, ok := echoes[h]; ok {
			ee.ServeHTTP(c.Response(), c.Request())
			return nil
		} else if hh, ok := subs[h]; ok {
			return c.File(path.Join(s.dir, hh, "substorage.html"))
		}
		c.Response().Header().Set("Cache-Control", "no-cache")
		return echo.ErrNotFound
	})
	e.Logger.Fatal(e.Start(port))
}

func (s *Server) Update(d *database.Database) {
	s.store = newStore(path.Join(s.dir, "moonboard"))
	s.store.update(d)
}
