package server

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/zombull/floating-castle/database"
	"github.com/zombull/floating-castle/moonboard"
)

const domain = "zombull.xyz"

type Server struct {
	dir      string
	cache    string
	password string
	db       *database.Database
	store    *KeyValueStore
}

func Init(d *database.Database, server, cache, password string) *Server {
	return &Server{
		db:       d,
		dir:      server,
		cache:    cache,
		password: password,
	}
}

func (s *Server) Run(port string, release bool) {
	s.store = newStore(path.Join(s.dir, "moonboard"), path.Join(s.cache, "moonboard"))
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
	if len(s.password) > 16 {
		moon.POST("/data/tocks", s.PostTocks)
	}

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

func (s *Server) Update() {
	s.store = newStore(path.Join(s.dir, "moonboard"), path.Join(s.cache, "moonboard"))
	s.store.update(s.db)
}

type tock struct {
	Problem  string `json:"p"`
	Date     string `json:"d"`
	Grade    string `json:"g"`
	Stars    uint   `json:"s"`
	Attempts uint   `json:"a"`
	Sessions uint   `json:"e"`
}

type tocks struct {
	Password string `json:"password"`
	Tocks    []tock `json:"tocks"`
}

func (s *Server) PostTocks(c echo.Context) error {
	t := &tocks{}

	var err error
	if err = c.Bind(&t); err != nil {
		return err
	}
	if t.Password != s.password {
		return echo.ErrUnauthorized
	}

	ticks := make([]*database.Tick, 0, len(t.Tocks))
	cragId := moonboard.Id(s.db)
	setId := moonboard.SetId(s.db)
	for _, o := range t.Tocks {
		route := s.db.FindRoute(setId, o.Problem)
		if route == nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("'%s' not in database, Moonboard index needs to be synced", o.Problem))
		}

		if e := s.db.GetTicks(route.Id); len(e) > 0 {
			continue
		}

		if _, ok := database.HuecoToFontainebleau[strings.ToUpper(o.Grade)]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid grade '%s' for problem '%s'", o.Grade, o.Problem))
		}

		tick := &database.Tick{
			RouteId:  route.Id,
			AreaId:   setId,
			CragId:   cragId,
			Grade:    o.Grade,
			Stars:    o.Stars,
			Attempts: o.Attempts,
			Flash:    o.Attempts == 1,
			Redpoint: true,
		}
		tick.Date, err = time.Parse("January 02, 2006", o.Date)
		if err != nil {
			return err
		}
		if o.Sessions > 0 {
			tick.Sessions = o.Sessions
		}
		ticks = append(ticks, tick)
	}
	for _, tick := range ticks {
		s.db.Insert(tick)
	}
	return c.String(http.StatusOK, "High Five!")
}
