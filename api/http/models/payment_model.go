package models

type PaymentRequest struct {
	AccountID string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float32 `json:"amount"`
}

type PaymentResponse struct {
	AccountID string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float32 `json:"amount"`
}
