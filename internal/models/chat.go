package models

type Chat struct {
	ID   int64 `json:"id"`
	Flow *Flow `json:"flow"`
}
