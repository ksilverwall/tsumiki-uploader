package services

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func GenerateID() (string, error) {
	u7, err := uuid.NewV7()
	if err != nil {
		ErrorLog(fmt.Errorf("failed to create transaction id: %w", err).Error())
		return "", ErrUnexpected
	}

	id := u7.String()

	return id, nil
}
