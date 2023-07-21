package models

type Transaction struct {
	ID       string
	FilePath string
}

type ThumbnailRequest struct {
	TransactionID string
	FilePath      string
}
