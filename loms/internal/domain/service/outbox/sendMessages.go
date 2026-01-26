package outbox

import (
	"context"
	"route256/loms/internal/infra/logger"
)

func (r *OutboxService) SendMessages(ctx context.Context) {
	items, err := r.outboxRepository.Get(ctx)
	if err != nil {
		logger.Errorw("Ошибка чтения сообщений из outbox", "Error", err)
	}
	for _, item := range *items {

		status := SENT
		err = r.producer.SendMessage(ctx, item.Key, item.Payload)
		if err != nil {
			status = ERROR
		}

		err := r.outboxRepository.SetStatus(ctx, int(item.Id), string(status))
		if err != nil {
			logger.Errorw("Ошибка записи статуса сообщения в outbox", "itemId", item.Id, "Error", err)
		}
	}
}
