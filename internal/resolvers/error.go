package resolvers

import (
	"errors"
	"fmt"
)

var (
	ErrInternal     = errors.New("InternalServerError")
	ErrUnauthorized = errors.New("Unauthorized")
)

type BadId struct {
	entity string
	id     uint
}

func ErrBadId(id uint, entity string) *BadId {
	return &BadId{entity, id}
}

func (b *BadId) Error() string {
	return fmt.Sprintf("No %s was found with ID \"%d\"", b.entity, b.id)
}

// Returns defaultErr when err == target. Else returns ErrInternal.
func ErrWithOrInternal(err, target, defaultErr error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, target) {
		return defaultErr
	}

	return ErrInternal
}
