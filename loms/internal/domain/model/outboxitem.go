package model

type OutboxItem struct {
	Id      int64
	Key     string
	Payload []byte
}
