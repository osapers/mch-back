package types

type Event struct {
	tableName        struct{} `pg:"mch.event,alias:e"` // nolint:structcheck,unused
	ID               string   `json:"id" pg:"type:uuid"`
	Name             string   `json:"name"`
	Image            string   `json:"image"`
	Category         string   `json:"category"`
	Date             int64    `json:"date"` // unix date
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	Address          Address  `json:"address"`
	Website          string   `json:"website"`
	Email            string   `json:"email"`
}

type Address struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Raw string  `json:"raw"`
}

var EventCategories = map[string]string{
	"webinar":           "вебинар",
	"forum":             "форум",
	"session":           "сессия",
	"exhibition":        "выставка",
	"lecture":           "лекция",
	"online_lecture":    "on-line лекция",
	"demo_day":          "демо-день",
	"round_table":       "круглый стол",
	"strategic_session": "стратегическая сессия",
}
