package models

type TransactionID string

type PlatformParameters struct {
	DataStorage string
}

type ThumbnailRequest struct {
	TransactionID         TransactionID
	ArchiveFilePath       string
	ThumbnailFilesKeyPath string
	ThumbnailFilesPrefix  string
}
