package store

import (
	"github.com/marcos-brito/booklist/internal/models"
	"gorm.io/gorm"
)

type BookStore struct {
	*gorm.DB
}

func NewBookStore(db *gorm.DB) *BookStore {
	return &BookStore{db}
}

func (bs *BookStore) FindById(id uint) (*models.Book, error) {
	book := &models.Book{}
	err := bs.DB.First(book, id).Error

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (bs *BookStore) Create(input *models.CreateBook, userUuid string) (*models.Book, error) {
	authors := []*models.Author{}
	err := bs.DB.Find(&authors, input.Authors).Error
	if err != nil {
		return nil, err
	}

	profile := &models.Profile{}
	err = bs.DB.First(profile, &models.Profile{UserUUID: userUuid}).Error
	if err != nil {
		return nil, err
	}

	book := models.NewBookFromInput(input, authors, profile.ID)
	err = bs.DB.Create(book).Error
	if err != nil {
		return nil, err
	}

	err = bs.DB.First(book, book.ID).Error
	if err != nil {
		return nil, err
	}

	return book, nil
}
