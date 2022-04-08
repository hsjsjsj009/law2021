package schema

import "time"

type InputData struct {
	Key string `json:"key"`
	Data interface{} `json:"data"`
}

type InputToken struct {
	Key string `json:"key"`
	Token string `json:"token"`
}

type RespondToken struct {
	Token string `json:"token"`
}

type RespondData struct {
	Body interface{} `json:"body"`
	Time time.Time `json:"time"`
}

