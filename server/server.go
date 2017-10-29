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
		dir:   root,
		store: newStore(root),
	}
}

func (s *Server) commonRoutes(e *echo.Echo, h string) {
	e.Static("/favicon", path.Join(s.dir, "common/img/favicon"))
	e.Static("/static", path.Join(s.dir, h))
	e.Static("/common", path.Join(s.dir, "common"))

	// // Send back the main HTML page when accessing a front facing URL.
	e.File("/", path.Join(s.dir, h, "index.html"))

	e.GET("/data/:key", s.store.getValue(h))
}

func (s *Server) Log(port string) {
	// Server
	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		fmt.Printf("%s: %s%s\n", c.Request().Method, c.Request().Host, c.Request().URL.Path)
		return echo.ErrNotFound
	})
	e.Logger.Fatal(e.StartTLS(port, path.Join(s.dir, "cert.pem"), path.Join(s.dir, "key.pem")))
}

func (s *Server) Run(port string) {
	beta := echo.New()
	s.commonRoutes(beta, "beta")
	beta.File("/:crag", path.Join(s.dir, "beta", "index.html"))
	beta.File("/:crag/a/:area", path.Join(s.dir, "beta", "index.html"))
	beta.File("/:crag/:route", path.Join(s.dir, "beta", "index.html"))

	// beta.GET('/go', s.getGo)
	// beta.GET("/data/crag/:crag", s.store.getCrag)
	// beta.GET("/data/area/:crag/:area", s.store.getArea)
	// beta.GET("/data/route/:crag/:route", s.store.getRoute)

	moon := echo.New()
	s.commonRoutes(moon, "moonboard")
	moon.File("/:problem", path.Join(s.dir, "moonboard", "index.html"))
	for _, r := range []string{"/p/:grade", "/t/:grade", "/k/:grade", "/j/:grade", "/s/:setter", "/st/:setter"} {
		moon.File(r, path.Join(s.dir, "moonboard", "index.html"))
	}

	echoes := map[string]*echo.Echo{
		"beta": beta,
		"moon": moon,
	}

	subs := map[string]string{
		"northwest": "beta",
		"southwest": "beta",
		"southeast": "beta",
		"norcal":    "beta",
		"socal":     "beta",
		"dark":      "moonboard",
		"side":      "moonboard",
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

		if ee, ok := echoes[h]; ok {
			ee.ServeHTTP(c.Response(), c.Request())
			return nil
		} else if hh, ok := subs[h]; ok {
			return c.File(path.Join(s.dir, hh, "substorage.html"))
		}
		fmt.Printf("host = %s", h)
		return echo.ErrNotFound
	})
	e.Logger.Fatal(e.Start(port))
}

func (s *Server) Update(d *database.Database) {
	// s.client.Nuke()
	s.store.update(d)
}

// func getIndex(c echo.Context) error {
// 	// fmt.Printf("INDEX: %s", c.Request().URL.Path)
// 	// return c.File("public/index.html")
// }
