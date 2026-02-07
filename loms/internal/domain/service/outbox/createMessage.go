package outbox

import (
	"context"
	"encoding/json"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/infra/logger"
	"strconv"
)

func (s *OutboxService) CreateMessage(ctx context.Context, orderID int64, status model.Status) {
	key := strconv.FormatInt(orderID, 10)

	kafkaMessage := newKafkaMessage(orderID, string(status))

	message, err := json.Marshal(kafkaMessage)
	if err != nil {
		logger.Errorw("Ошибка создания сообщения kafka", "orderID", orderID, "status", status, "Error", err)
		return
	}

	_, err = s.outboxRepository.Create(ctx, key, string(NEWSTATUS), message)
	if err != nil {
		logger.Errorw("Ошибка записи сообщения в outbox", "orderID", orderID, "status", status, "Error", err)
		return
	}
}
