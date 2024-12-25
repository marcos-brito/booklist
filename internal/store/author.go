package store

import (
	"github.com/marcos-brito/booklist/internal/models"
	"gorm.io/gorm"
)

type AuthorStore struct {
    *gorm.DB
}

func NewAuthorStore(db *gorm.DB) *AuthorStore {
    return &AuthorStore {
        db,
    }
}

func (as *AuthorStore) FindById(id uint) (*models.Author ,error) {
    author := &models.Author{}

    err := as.DB.First(author, id).Error
    if err != nil {
        return nil, err
    }

    return author, nil
}

func (as *AuthorStore) FindManyById(ids ...uint) ([]*models.Author, error, *uint) {
	authors := []*models.Author{}

	for _, id := range ids {
		author, err := as.FindById(id)

		if err != nil {
			return nil, err, &id
		}

		authors = append(authors, author)
	}

    return authors, nil, nil
}
