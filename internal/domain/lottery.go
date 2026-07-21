package domain

type Lottery struct {
	Id      uint   `json:"id"`
	Round   string `json:"round"`
	Time    string `json:"time"`
	Numbers string `json:"numbers"`
}
