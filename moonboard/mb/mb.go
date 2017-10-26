package mb

type Setter struct {
	Country   string `json:"Country"`
	City      string `json:"City"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Nickname  string `json:"Nickname"`
	Id        string `json:"Id"`
}

type Hold struct {
	IsEnd    bool   `json:"IsEnd"`
	IsStart  bool   `json:"IsStart"`
	Location string `json:"Description"` // e.g. "G2"
	Id       uint   `json:"Id"`          // 1523436 no idea what this is used for
}

type Problem struct {
	Date      string `json:"DateTimeString"` // "21 Jul 2016 16:59",
	ApiId     uint   `json:"ApiId"`          // 28015
	Id        uint   `json:"Id"`             // 28015
	Url       string `json:"NameForUrl"`
	Setter    Setter `json:"Setter"`
	UserGrade string `json:"UserGrade"` // "consensus grade" "6B+"
	Grade     string `json:"Grade"`     // ?? setter's grade "6B+",
	Name      string `json:"Name"`
	Rating    uint   `json:"Rating"` // setter's rating???
	Stars     uint   `json:"UserRating"`
	Ascents   uint   `json:"Repeats"`
	Benchmark bool   `json:"IsBenchmark"`
	Holds     []Hold `json:"Moves"`
}

type Problems struct {
	Problems []Problem `json:"Data"`
}

type Tick struct {
	Problem       Problem `json:"Problem"`
	Attempts      uint    `json:"Attempts"`      // currently not filled in, i.e. it's worthless
	Grade         string  `json:"Grade"`         // my grade, e.g. "6B+",
	NumberOfTries string  `json:"NumberOfTries"` // "Flashed", "2nd try", "3rd try", "more than 3 tries"
	Stars         uint    `json:"Rating"`
	Date          string  `json:"DateClimbedAsString"` // "06 Aug 2017"
	Comment       string  `json:"Comment"`             // my comment -  "Attempts: 8\nStars: ..."
}

type Ticks struct {
	Ticks []Tick `json:"Data"`
}
