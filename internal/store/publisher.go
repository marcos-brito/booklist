package store

import (
	"github.com/marcos-brito/booklist/internal/models"
	"gorm.io/gorm"
)

type PublisherStore struct {
    *gorm.DB
}

func NewPublisherStore(db *gorm.DB) *PublisherStore {
    return &PublisherStore {db}
}

func (as *PublisherStore) FindById(id uint) (*models.Publisher ,error) {
    publisher := &models.Publisher{}

    err := as.DB.First(publisher, id).Error
    if err != nil {
        return nil, err
    }

    return publisher, nil
}
