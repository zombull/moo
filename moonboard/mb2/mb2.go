package mb2

type Hold struct {
	IsEnd		uint   `json:"IsEnd"`
	IsStart  	uint   `json:"IsStart"`
	Location 	string `json:"Description"` // e.g. "G2"
	ProblemId	uint   `json:"Problem_Id"`  // Not ApiId, index into array
}

/*
type Problem struct {
	Active		uint	`json:"Active"`
	ApiId		uint	`json:"ApiId"`
	Comment		string	`json:"Comment"`
	DateDeleted	*int64	`json:"DateDeleted"`
	DateInserted	uint64	`json:""`
	DateUpdated	uint64	`json:""`
	Downgraded	uint	`json:"Downgraded"`
	FirstAscent	string	`json:"FirstAscent"`
	Grade		string	`json:"Grade"`
	HasBetaVideo	uint	`json:"HasBetaVideo"`
	Holdsets	string	`json:"Holdsets"`
	Id		uint	`json:"Id"`
	IsBenchmark	uint	`json:"IsBenchmark"`
	IsMaster	uint	`json:"IsMaster"`
	Method		string	`json:"Method"`
	ConfigId	uint	`json:"MoonBoardConfigurationId"`
	MoonId		uint	`json:"MoonId"`
	Name		string	`json:"Name"`
	Repeats		uint	`json:"Repeats"`
	Setby		string	`json:"Setby"`
	SetbyId		string	`json:"SetbyId"` // UUID
	SetupId		uint	`json:"SetupId"`
	Upgraded	uint	`json:"Upgraded"`
	UserGrade	string	`json:"UserGrade"`
	Stars		uint	`json:"UserRating"`
}
*/

type Problem struct {
	ApiId		uint	`json:"ApiId"`
	Date		int64	`json:"DateInserted"`
	Grade		string	`json:"Grade"`
	Id		uint	`json:"Id"`
	IsBenchmark	uint	`json:"IsBenchmark"`
	Name		string	`json:"Name"`
	Ascents		uint	`json:"Repeats"`
	Setter		string	`json:"Setby"`
	SetbyId		string	`json:"SetbyId"` // UUID
	SetupId		uint	`json:"SetupId"`
	UserGrade	string	`json:"UserGrade"`
	Stars		uint	`json:"UserRating"`
}
