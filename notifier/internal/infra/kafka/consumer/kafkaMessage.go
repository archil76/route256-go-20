package consumer

type KafkaMessage struct {
	OrderId int64  `json:"order_id"`
	Status  string `json:"status"`
	Moment  string `json:"moment"`
}
