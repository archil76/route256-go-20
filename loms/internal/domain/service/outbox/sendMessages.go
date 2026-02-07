package outbox

import (
	"context"
	"route256/loms/internal/infra/logger"
)

func (s *OutboxService) SendMessages(ctx context.Context) {

	items, err := s.outboxRepository.Get(ctx)
	if err != nil {
		logger.Errorw("Ошибка чтения сообщений из outbox", "Error", err)
	}

	err = s.pooler.InTx(ctx, func(ctx context.Context) error {
		for _, item := range *items {

			status := SENT
			err = s.producer.SendMessage(ctx, item.Key, item.Payload)
			if err != nil {
				status = ERROR
			}

			err = s.outboxRepository.SetStatus(ctx, int(item.Id), string(status))
			if err != nil {
				logger.Errorw("Ошибка записи статуса сообщения в outbox", "itemId", item.Id, "Error", err)
			}
		}

		return nil
	})
}
