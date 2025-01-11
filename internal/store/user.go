package store

import (
	"github.com/google/uuid"
	"github.com/marcos-brito/booklist/internal/models"
	"gorm.io/gorm"
)

type UserStore struct {
	*gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db}
}

func (us *UserStore) FindProfileByUserUuid(uuid uuid.UUID) (*models.Profile, error) {
	profile := &models.Profile{}
	err := us.DB.Where(models.Profile{UUID: uuid}).
		Attrs(models.Profile{Settings: models.Settings{Private: true}}).FirstOrCreate(profile).Error

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (us *UserStore) FindFullProfileById(id uint) (*models.Profile, error) {
	profile := &models.Profile{}
	err := us.DB.Preload("Settings").First(profile, id).Error

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

func (us *UserStore) FindSettingsByUserUuid(uuid uuid.UUID) (*models.Settings, error) {
	profile := &models.Profile{}
	err := us.DB.First(profile, &models.Profile{UUID: uuid}).Error
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

func (us *UserStore) FindItems(userUuid uuid.UUID) ([]*models.CollectionItem, error) {
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

func (us *UserStore) FindLists(userUuid uuid.UUID) ([]*models.List, error) {
	profile, err := us.FindProfileByUserUuid(userUuid)
	if err != nil {
		return nil, err
	}

	lists := []*models.List{}
	err = us.DB.Find(&lists, models.List{ProfileID: profile.ID}).Error
	if err != nil {
		return nil, err
	}

	return lists, nil
}

func (us *UserStore) FindPublicLists(userUuid uuid.UUID) ([]*models.List, error) {
	profile, err := us.FindProfileByUserUuid(userUuid)
	if err != nil {
		return nil, err
	}

	lists := []*models.List{}
	err = us.DB.Find(&lists, models.List{ProfileID: profile.ID, Published: true}).Error
	if err != nil {
		return nil, err
	}

	return lists, nil
}

func (us *UserStore) AddToCollection(userUuid uuid.UUID, bookID uint, status models.Status) (*models.CollectionItem, error) {
	profile := &models.Profile{}
	err := us.DB.First(profile, &models.Profile{UUID: userUuid}).Error
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

func (us *UserStore) UpdateSettings(uuid uuid.UUID, changes models.UpdateSettings) (*models.Settings, error) {
	settings, err := us.FindSettingsByUserUuid(uuid)
	if err != nil {
		return nil, err
	}

	settings.Private = changes.Private
	settings.ShowName = changes.ShowName
	settings.ShowStats = changes.ShowStats
	settings.ShowCollection = changes.ShowCollection
	settings.ShowListsFollows = changes.ShowListsFollows
	settings.ShowAuthorsFollows = changes.ShowAuthorsFollows

	err = us.Save(settings).Error
	if err != nil {
		return nil, err
	}

	return settings, nil
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
