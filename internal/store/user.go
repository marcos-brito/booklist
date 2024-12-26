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

func (us *UserStore) FindProfileByUserUuid(uuid string) (*models.Profile, error) {
	profile := &models.Profile{}
	err := us.DB.Where(models.Profile{UserUUID: uuid}).
		Attrs(models.Profile{Settings: models.Settings{Private: true}}).FirstOrCreate(profile).Error

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (us *UserStore) FindItemById(id uint) (*models.CollectionItem, error) {
	item := &models.CollectionItem{}
	err := us.DB.First(item, id).Error

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (us *UserStore) FindSettingsByUserUuid(uuid string) (*models.Settings, error) {
	profile := &models.Profile{}
    err := us.DB.First(profile, &models.Profile{UserUUID: uuid}).Error
    if err != nil {
        return nil, err
    }

	settings := &models.Settings{}
    err = us.DB.First(settings, &models.Settings{ProfileID: profile.ID}).Error
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (us *UserStore) FindItems(userUuid string) ([]*models.CollectionItem, error) {
    profile, err := us.FindProfileByUserUuid(userUuid)
	if err != nil {
		return nil, err
	}

	items := []*models.CollectionItem{}
    err = us.DB.Find(&items, models.CollectionItem{ProfileID: profile.ID}).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (us *UserStore) AddToCollection(userUuid string, bookID uint, status models.Status) (*models.CollectionItem, error) {
	profile := &models.Profile{}
	err := us.DB.First(profile, &models.Profile{UserUUID: userUuid}).Error
	if err != nil {
		return nil, err
	}

	item := &models.CollectionItem{
		ProfileID: profile.ID,
		BookID:    bookID,
		Status:    status,
	}

	err = us.DB.Create(item).Error
	if err != nil {
		return nil, err
	}

	err = us.DB.First(item, item.ID).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (us *UserStore) DeleteFromCollection(id uint) (*models.CollectionItem, error) {
	item := &models.CollectionItem{}
	err := us.DB.Delete(&item, id).Error

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (us *UserStore) ChangeItemStatus(id uint, status models.Status) (*models.CollectionItem, error) {
	item := &models.CollectionItem{}
	err := us.DB.First(item, id).Error
	if err != nil {
		return nil, err
	}

	item.Status = status
	err = us.DB.Save(item).Error
	if err != nil {
		return nil, err
	}

	err = us.DB.First(item, item.ID).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}
