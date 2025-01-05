package bcryptp

import (
	"fmt"
	"github.com/gofrs/uuid"
)
func CreateSession() (uuid.UUID,error) {
	session, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return uuid.UUID{},err
	}
	return session,nil
}