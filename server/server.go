package server

import (
	"fmt"
	"path"
	"strings"

	"github.com/labstack/echo"
)

const domain = "zombull.xyz"

type Server struct {
	dir   string
	cache string
	store *KeyValueStore
}

func Init(server, cache string) *Server {
	return &Server{
		dir:   server,
		cache: cache,
	}
}

func (s *Server) newMoon(release bool, year string) *echo.Echo {
	moon := echo.New()
	subName := "moonboard" + year

	common := path.Join(s.dir, "common")
	if release {
		common = path.Join(s.dir, subName)
	}

	index := path.Join(s.dir, subName, "index.html")
	for _, r := range []string{"/", "/:problem", "/q/:query", "/p/:grade", "/t/:grade", "/k/:grade", "/j/:grade", "/o/:grade", "/s/:setter", "/st/:setter"} {
		moon.File(r, index)
	}

	moon.Static("/favicon", path.Join(common, "img", "favicon"))
	moon.Static("/static", path.Join(s.dir, subName))
	moon.Static("/common", common)
	moon.GET("/data/:key", s.store.getValue(subName))

	return moon
}

func (s *Server) Run(port string, release bool) {
	s.store = NewStore(s.cache, s.dir)

	echoes := map[string]*echo.Echo{
		"moon2016": s.newMoon(release, "2016"),
		"moon2017": s.newMoon(release, "2017"),
		"moon2019": s.newMoon(release, "2019"),
	}

	subs := map[string]string{
		"dark2016": "moonboard2016",
		"dark2017": "moonboard2017",
		"dark2019": "moonboard2019",
	}

	// Server
	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		h := c.Request().Host
		ss := strings.SplitN(h, ".", 2)
		if len(ss) != 2 || len(ss[1]) == 0 || strings.Split(ss[1], ":")[0] != domain {
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
