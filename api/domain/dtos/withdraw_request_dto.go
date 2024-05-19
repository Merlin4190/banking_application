package dtos

type WithdrawRequestDto struct {
	Amount               float32 `json:"amount" validate:"required"`
	AccountNumber        string  `json:"account_number" validate:"required"`
	TransactionReference string  `json:"transaction_reference" validate:"required"`
}
