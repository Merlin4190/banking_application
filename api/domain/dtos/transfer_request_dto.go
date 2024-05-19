package dtos

type TransferRequestDto struct {
	Amount                   float32 `json:"amount" validate:"required"`
	SourceAccountNumber      string  `json:"source_account_number" validate:"required"`
	DestinationAccountNumber string  `json:"destination_account_number" validate:"required"`
	TransactionReference     string  `json:"transaction_reference" validate:"required"`
}
