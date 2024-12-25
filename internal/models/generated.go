// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type CreateBook struct {
	Title       string     `json:"title"`
	Isbn        string     `json:"isbn"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	PageCount   *int       `json:"pageCount,omitempty"`
	Edition     *int       `json:"edition,omitempty"`
	Authors     []uint     `json:"authors"`
	Publisher   *uint      `json:"publisher,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type Status string

const (
	StatusToRead  Status = "TO_READ"
	StatusOnHold  Status = "ON_HOLD"
	StatusDropped Status = "DROPPED"
	StatusReading Status = "READING"
	StatusRead    Status = "READ"
)

var AllStatus = []Status{
	StatusToRead,
	StatusOnHold,
	StatusDropped,
	StatusReading,
	StatusRead,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusToRead, StatusOnHold, StatusDropped, StatusReading, StatusRead:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
