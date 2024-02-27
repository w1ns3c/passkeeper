package entities

import "time"

type Password struct {
	ID          string
	Date        time.Time // date for last changing this line
	Resource    string
	Login       string
	Password    string
	Description string
}
