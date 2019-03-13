package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/zombull/floating-castle/bug"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Config is the location of the user's config file.  This value is not
	// saved/exposed as the location is defined via an environment variable
	// to avoid a Catch-22.
	Config string `yaml:"-"`

	// Cache is the path to the directory where semi-permanent, non-database
	// files are stored, e.g. the JSON representation of the database.
	// Env Var: FC_CACHE
	// Default: $HOME/go/src/github.com/zombull/cache-floating-castle
	Cache string `yaml:"cache"`

	// Database is the path to the directory where the SQLite database
	// exists (or is created).
	// Env Var: FC_DATABASE
	// Default: $HOME/go/src/github.com/zombull/db-floating-castle/sqlite3
	Database string `yaml:"database"`

	// Server is the path to the root directory of the web server.
	// Env Var: FC_SERVER
	// Default: $HOME/go/src/github.com/zombull/floating-castle/server
	Server string `yaml:"server"`

	// Specify the Moonboard set, i.e. year and holds combination.
	// Env Var: FC_MOONBOARD_SET
	// Default: "MoonBoard 2017"
	MoonboardSet string `yaml:"moonboard_set"`
}

func loadEnvVar(name, def string) string {
	name = "FC_" + name
	if os.Getenv(name) != "" {
		def = os.Getenv(name)
	}
	return os.ExpandEnv(def)
}

// LoadConfig reads the configuration from the config path; if the path does
// not exist, it returns a default configuration.
func loadConfig() *Config {
	// Use a default config if a user-defined file does not exist.
	// Basic Windows (not MinGW or MSysGit) may not have $HOME set,
	// look for HOMEDRIVE and HOMEPATH.
	dir := "$HOME"
	if os.Getenv("HOME") == "" && os.Getenv("HOMEDRIVE") != "" && os.Getenv("HOMEPATH") != "" {
		dir = path.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"))
	}
	c := Config{
		Config:       path.Join(dir, ".config", "floating-castle", "config.yml"),
		Cache:        path.Join(dir, "go", "src", "github.com", "zombull", "cache-floating-castle"),
		Database:     path.Join(dir, "go", "src", "github.com", "zombull", "db-floating-castle", "sqlite3"),
		Server:       path.Join(dir, "go", "src", "github.com", "zombull", "floating-castle", "server"),
		MoonboardSet: "MoonBoard 2017",
	}

	path := loadEnvVar("CONFIG", c.Config)
	data, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(data, &c)
		bug.OnError(err)
	} else {
		bug.On(!os.IsNotExist(err), fmt.Sprintf("cannot read config file: %v", err))
	}

	c.Cache = loadEnvVar("CACHE", c.Cache)
	c.Database = loadEnvVar("DATABASE", c.Database)
	c.Server = loadEnvVar("SERVER", c.Server)
	c.MoonboardSet = loadEnvVar("MOONBOARD_SET", c.MoonboardSet)

	bug.On(len(c.Cache) == 0, "CACHE must be a non-empty string")
	bug.On(len(c.Database) == 0, "DATABASE must be a non-empty string")
	bug.On(len(c.Server) == 0, "SERVER must be a non-empty string")

	err = os.MkdirAll(c.Database, 0770)
	bug.OnError(err)

	return &c
}
