package model

type Status string

const (
	NEWSTATUS       Status = "new"
	AWAITINGPAYMENT Status = "awaiting payment"
	FAILED          Status = "failed"
	PAYED           Status = "paid"
	CANCELED        Status = "cancelled"
)
