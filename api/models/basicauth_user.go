package models

import (
	"regexp"
	"time"
	"unicode"
)

type BasicAuthUser struct {
	ID        string    `json:"user_id" gorm:"primary_key"`
	Password  string    `json:"password"`
	Nickname  *string   `json:"nickname"`
	Comment   *string   `json:"comment"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

type ResponseUser struct {
	ID       string  `json:"user_id"`
	Nickname string  `json:"nickname"`
	Comment  *string `json:"comment,omitempty"`
}

type ResponseMessage struct {
	Message string       `json:"message"`
	User    ResponseUser `json:"user"`
}

func check_regexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}
func isAlphanumeric(s string) bool {
	return check_regexp(`[0-9A-Za-z]`, s)
}

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII || c < '!' {
			return false
		}
	}
	return true
}

func (u BasicAuthUser) Validate() bool {

	length := len(u.ID)
	if length < 6 || length > 20 {
		return false
	}
	if !isAlphanumeric(u.ID) {
		return false
	}

	length = len(u.Password)

	if length < 6 || length > 20 {
		return false
	}
	if !isASCII(u.ID) {
		return false
	}
	return true
}
