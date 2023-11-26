package domain

type Message struct {
	ID      int64  `json:"id"`
	Author  string `json:"author"`
	Message string `json:"message"`
}
