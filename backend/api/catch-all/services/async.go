package services

import (
	"catch-all/models"
	"catch-all/repositories"
	"fmt"
)

type Async struct {
	StateMachineRepository repositories.StateMachine
}

func (s *Async) CreateThumbnails(m models.ThumbnailRequest) error {
	err := s.StateMachineRepository.Execute(m)
	if err != nil {
		ErrorLog(fmt.Errorf("transaction id not found: %w", err).Error())
		return err
	}

	return nil
}
