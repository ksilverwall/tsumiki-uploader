package models

type Transaction struct {
	ID       string
	FilePath string
}

type ThumbnailRequest struct {
	TransactionID string
	FilePath      string
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
