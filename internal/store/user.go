package store

import (
	"github.com/marcos-brito/booklist/internal/models"
	"gorm.io/gorm"
)

type UserStore struct {
	*gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db}
}

func (us *UserStore) FindProfile(uuid string) (*models.Profile, error) {
	profile := &models.Profile{}
	res := us.DB.Debug().Where(models.Profile{UserUUID: uuid}).
		Attrs(models.Profile{Settings: models.Settings{Private: true}}).
		Preload("Settings").Preload("Collection").Preload("Lists").
		FirstOrCreate(profile)

	if err := res.Error; err != nil {
		return nil, err
	}

	return profile, nil
}
