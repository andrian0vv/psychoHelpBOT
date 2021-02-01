package models

type Step struct {
	Name     string   `json:"name"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}

func (s Step) IsAnswered() bool {
	return len(s.Answer) != 0
}
