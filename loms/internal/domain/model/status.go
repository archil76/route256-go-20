package model

type Status string

const (
	NEW_STATUS       Status = "new"
	AWAITING_PAYMENT Status = "awaiting payment"
	FAILED           Status = "failed"
	PAYED            Status = "payed"
	CANCELED         Status = "canceled"
)
