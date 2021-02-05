package models

type Chat struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	Flow     *Flow  `json:"flow"`
}
