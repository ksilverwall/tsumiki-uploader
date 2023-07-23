package services

import (
	"catch-all/models"
	"catch-all/repositories"
	"fmt"
)

type Transaction struct {
	StorageRepository      repositories.Storage
	TransactionRepository  repositories.Transaction
	StateMachineRepository repositories.StateMachine
}

func (s *Transaction) Get(transactionID string) (models.Transaction, error) {
	t, err := s.TransactionRepository.Get(transactionID)
	if err != nil {
		return models.Transaction{}, ErrUnexpected
	}

	return t, nil
}

func (s *Transaction) Create(id string, filePath string) error {
	err := s.TransactionRepository.Put(id, models.Transaction{
		ID:       id,
		FilePath: filePath,
	})
	if err != nil {
		ErrorLog(fmt.Errorf("failed to push transaction: %w", err).Error())
		return ErrUnexpected
	}

	return nil
}
