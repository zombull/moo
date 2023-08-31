package mb

import (
	"time"
)

type Setter struct {
	Name      string
	Firstname string
	Lastname  string
	Nickname  string
	Id        string
}

type Hold struct {
	IsEnd    bool
	IsStart  bool
	Location string
	Id       uint
}

type Problem struct {
	Date      time.Time
	ApiId     uint
	Id        uint
	Setter    Setter
	UserGrade string
	Grade     string
	Name      string
	Rating    uint
	Stars     uint
	Ascents   uint
	Benchmark bool
	Holds     []Hold
}

type Problems struct {
	Problems []Problem
}
