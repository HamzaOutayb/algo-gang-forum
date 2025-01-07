package bcryptp

import (
	"github.com/gofrs/uuid"
)

func CreateSession() (uuid.UUID, error) {
	session, err := uuid.NewV4()
	if err != nil {
		return uuid.UUID{}, err
	}
	return session, nil
}
