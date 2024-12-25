package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	// Comes from Ory
	UserUUID   string
	Name       string `gorm:"-:all"`
	Email      string `gorm:"-:all"`
	Settings   Settings
	Lists      []List
	Collection []CollectionItem
}

func (p *Profile) UUID() string {
	return p.UserUUID
}

type Settings struct {
	gorm.Model
	ProfileID          uint
	Private            bool
	ShowName           bool
	ShowStats          bool
	ShowCollection     bool
	ShowListsFollows   bool
	ShowAuthorsFollows bool
}

type List struct {
	gorm.Model
	ProfileID   uint
	Name        string
	Description string
	Books       []Book `gorm:"many2many:list_books;"`
}

type CollectionItem struct {
	gorm.Model
	ProfileID  uint
	BookID     uint
	Book       Book
	Status     Status
	StartedAt  *time.Time
	FinishedAt *time.Time
}

type Book struct {
	gorm.Model
	Title         string
	ISBN          string
	PublishedAt   *time.Time
	PageCount     *int
	Edition       *int
	NeedsApproval bool
	Authors       []*Author `gorm:"many2many:book_authors;"`
	PublisherID   *uint
	Publisher     Publisher
	ProfileID     uint
	Profile       Profile
}
}

type Author struct {
	gorm.Model
	Name     string
	BirthDay time.Time
	Books    []*Book `gorm:"many2many:book_authors;"`
}

type Publisher struct {
	gorm.Model
	Name string
}
