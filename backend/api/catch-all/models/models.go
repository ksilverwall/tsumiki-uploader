package models

type FileID string
type TransactionID string

type Transaction struct {
	FileID   FileID
	FilePath string
}

type ThumbnailRequest struct {
	TransactionID         TransactionID
	ArchiveFilePath       string
	ThumbnailFilesKeyPath string
	ThumbnailFilesPrefix  string
}

type DynamodbInfo struct {
	Name string
	Key  string
	TTL  string
}

type PlatformParameters struct {
	DataStorage      string
	TransactionTable DynamodbInfo
}

type BatchParameters struct {
	ThumbnailsCreatingStateMachineArn string
}
