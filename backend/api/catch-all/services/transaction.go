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

func (s *Transaction) Get(id models.TransactionID) (models.Transaction, error) {
	t, err := s.TransactionRepository.Get(id)
	if err != nil {
		ErrorLog(fmt.Errorf("failed to get transaction id '%s': %w", id, err).Error())
		return models.Transaction{}, ErrUnexpected
	}

	return t, nil
}

func (s *Transaction) Create(fid models.FileID, filePath string) (models.TransactionID, error) {
	id, err := GenerateID()
	tid := models.TransactionID(id)
	if err != nil {
		ErrorLog(fmt.Errorf("failed to create transaction id: %w", err).Error())
		return models.TransactionID(""), ErrUnexpected
	}

	err = s.TransactionRepository.Put(tid, models.Transaction{
		FileID:   fid,
		FilePath: filePath,
	})
	if err != nil {
		ErrorLog(fmt.Errorf("failed to push transaction: %w", err).Error())
		return models.TransactionID(""), ErrUnexpected
	}

	InfoLog(fmt.Sprintf("transaction created with id '%s'", tid))

	return tid, nil
}
