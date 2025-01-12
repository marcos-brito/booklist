package resolvers

import (
	"github.com/google/uuid"
	"github.com/marcos-brito/booklist/internal/conn"
	"github.com/marcos-brito/booklist/internal/store"
	"gorm.io/gorm"
)

// Reports whether the list is owned by ther user with the given
// UUID. If it's not, a error describing the reason is also returned.
func listIsOwned(listId uint, userUuid uuid.UUID) (bool, error) {
	profile, err := store.NewUserStore(conn.DB).FindProfileByUserUuid(userUuid)
	if err != nil {
		return false, ErrInternal
	}

	list, err := store.NewListStore(conn.DB).FindById(listId)
	if err != nil {
		return false, ErrWithOrInternal(gorm.ErrRecordNotFound, err, ErrBadId(listId, "list"))
	}

	if profile.ID != list.ProfileID {
		return false, ErrBadId(listId, "list")
	}

	return true, nil
}
