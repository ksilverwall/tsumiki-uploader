// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package openapi

// Defines values for ClientErrorCode.
const (
	ClientErrorCodeFileNotFound      ClientErrorCode = "FileNotFound"
	ClientErrorCodeThumbnailNotFound ClientErrorCode = "ThumbnailNotFound"
	ClientErrorCodeUnknown           ClientErrorCode = "Unknown"
)

// Defines values for ServerErrorCode.
const (
	ServerErrorCodeUnknown ServerErrorCode = "Unknown"
)

// Defines values for TransactionStatus.
const (
	READY    TransactionStatus = "READY"
	UPLOADED TransactionStatus = "UPLOADED"
)

// ClientErrorCode defines model for ClientErrorCode.
type ClientErrorCode string

// DownloadInfo defines model for DownloadInfo.
type DownloadInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// FileThumbnail defines model for FileThumbnail.
type FileThumbnail struct {
	Url string `json:"url"`
}

// FileThumbnails defines model for FileThumbnails.
type FileThumbnails struct {
	Items []FileThumbnail `json:"items"`
}

// ServerErrorCode defines model for ServerErrorCode.
type ServerErrorCode string

// Transaction defines model for Transaction.
type Transaction struct {
	Id     string            `json:"id"`
	Status TransactionStatus `json:"status"`
	Url    string            `json:"url"`
}

// TransactionStatus defines model for TransactionStatus.
type TransactionStatus string

// ClientError defines model for ClientError.
type ClientError struct {
	Code    ClientErrorCode `json:"code"`
	Message string          `json:"message"`
}

// ServerError defines model for ServerError.
type ServerError struct {
	Code    ServerErrorCode `json:"code"`
	Message string          `json:"message"`
}

// UpdateTransactionJSONBody defines parameters for UpdateTransaction.
type UpdateTransactionJSONBody struct {
	Status TransactionStatus `json:"status"`
}

// UpdateTransactionJSONRequestBody defines body for UpdateTransaction for application/json ContentType.
type UpdateTransactionJSONRequestBody UpdateTransactionJSONBody
