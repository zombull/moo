package gd

type Tick struct {
	Date     string `json:"d"`
	Grade    string `json:"g"`
	Stars    uint   `json:"s"`
	Attempts uint   `json:"a"`
	Sessions uint   `json:"e"`
}

type Project struct {
	Attempts uint `json:"a"`
	Sessions uint `json:"s"`
}

type UserData struct {
	Projects map[string]Project `json:"projects"`
	Ticks    map[string]Tick    `json:"ticks"`
}
