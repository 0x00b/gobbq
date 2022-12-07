package testdata

import (
	a "fmt"
	"time"
	_ "underscore" // underscore TODO
)

type foo struct {
	Time time.Time `json:"text"`
}

func fn() {
	a.Println("hello")
}

// Message TODO
type Message struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}
