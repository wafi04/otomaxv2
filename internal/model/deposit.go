package model

import "time"

type DepositData struct {
	ID                int       `json:"id"`
	Username          string    `json:"username"`
	Method            string    `json:"method"`
	PaymentReferee    *string   `json:"paymentReferee,omitempty"`
	DestinationNumber string    `json:"destinationNumber"`
	Amount            int       `json:"amount"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type CreateDeposit struct {
	Amount            int     `json:"amount"`
	Method            string  `json:"method"`
	Username          string  `json:"username"`
	PaymentReferee    *string `json:"paymentReferee,omitempty"`
	DestinationNumber string  `json:"destinationNumber"`
}

type RequestFormClient struct {
	Amount int    `json:"amount"`
	Method string `json:"method"`
}
type FilterDeposit struct {
	Search *string `json:"search,omitempty"`
	Status *string `json:"status,omitempty"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}
