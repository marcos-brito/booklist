package store

import (
	"github.com/google/uuid"
	"github.com/marcos-brito/booklist/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ListStore struct {
	*gorm.DB
}

func NewListStore(db *gorm.DB) *ListStore {
	return &ListStore{db}
}

func (ls *ListStore) IsOwner(id uint, userUuid uuid.UUID) (bool, error) {
	list := &models.List{}
	err := ls.DB.First(list, id).Error
	if err != nil {
		return false, err
	}

	profile, err := NewUserStore(ls.DB).FindProfileByUserUuid(userUuid)
	if err != nil {
		return false, err
	}

	if list.ProfileID != profile.ID {
		return false, nil
	}

	return true, nil
}

func (ls *ListStore) FindById(id uint) (*models.List, error) {
	list := &models.List{}
	err := ls.DB.First(list, id).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (ls *ListStore) FindBooks(id uint) ([]*models.Book, error) {
	list := &models.List{}
	list.ID = id
	books := []*models.Book{}
	err := ls.DB.Model(list).Association("Books").Find(&books)

	if err != nil {
		return nil, err
	}

	return books, nil
}

func (ls *ListStore) Create(name string, desc *string, publish bool, userUuid uuid.UUID) (*models.List, error) {
	profile, err := NewUserStore(ls.DB).FindProfileByUserUuid(userUuid)
	if err != nil {
		return nil, err
	}

	list := &models.List{
		Name:        name,
		Description: desc,
		ProfileID:   profile.ID,
		Published:   publish,
	}

	err = ls.DB.Create(list).Error
	if err != nil {
		return nil, err
	}

	list, err = ls.FindById(list.ID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (ls *ListStore) Delete(id uint) (*models.List, error) {
	list := &models.List{}
	err := ls.DB.Delete(list, id).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (ls *ListStore) Clone(id uint, userUuid uuid.UUID) (*models.List, error) {
	original := &models.List{}
	err := ls.DB.Preload("Books").Find(original, id).Error
	if err != nil {
		return nil, err
	}

	profile, err := NewUserStore(ls.DB).FindProfileByUserUuid(userUuid)
	if err != nil {
		return nil, err
	}

	list := &models.List{
		ProfileID:   profile.ID,
		Name:        original.Name,
		Description: original.Description,
		Published:   false,
		Books:       original.Books,
	}

	err = ls.DB.Clauses(clause.Returning{}).Create(list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (ls *ListStore) Publish(id uint) (*models.List, error) {
	list, err := ls.FindById(id)
	if err != nil {
		return nil, err
	}

	list.Published = true
	err = ls.DB.Save(list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (ls *ListStore) Unpublish(id uint) (*models.List, error) {
	list, err := ls.FindById(id)
	if err != nil {
		return nil, err
	}

	list.Published = false
	err = ls.DB.Save(list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (ls *ListStore) Follow(id uint, userUuid uuid.UUID) (*models.List, error) {
	return nil, nil
}

func (ls *ListStore) Unfollow(id uint, userUuid uuid.UUID) (*models.List, error) {
	return nil, nil
}

func (ls *ListStore) AddBook(listId, bookId uint) (*models.List, error) {
	list := &models.List{}
	book := &models.Book{}
	list.ID = listId
	book.ID = bookId

	err := ls.DB.Model(list).Association("Books").Append(book)
	if err != nil {
		return nil, err
	}

	list, err = ls.FindById(listId)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (ls *ListStore) RemoveBook(listId, bookId uint) (*models.List, error) {
	list := &models.List{}
	book := &models.Book{}
	list.ID = listId
	book.ID = bookId

	err := ls.DB.Model(list).Association("Books").Delete(book)
	if err != nil {
		return nil, err
	}

	list, err = ls.FindById(listId)
	if err != nil {
		return nil, err
	}

	return list, nil
}
