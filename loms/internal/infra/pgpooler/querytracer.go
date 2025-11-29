package pgpooler

import (
	"strings"
	"time"

	"route256/loms/internal/infra/metrics"

	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

type PgxQueryTracer struct {
}

func NewPgxQueryTracer() *PgxQueryTracer {
	return &PgxQueryTracer{}
}

func (t *PgxQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, _ pgx.TraceQueryStartData) context.Context {
	// Здесь можно начать отслеживание (например, создать span OpenTelemetry)
	// и добавить его в контекст.
	startTime := time.Now()

	// Возвращаем новый контекст с временной меткой начала для использования в TraceQueryEnd
	return context.WithValue(ctx, "startTime", startTime)
}

// TraceQueryEnd вызывается по завершении выполнения запроса.
func (t *PgxQueryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	// Извлекаем время начала из контекста
	duration := -1 * time.Second
	startTime, ok := ctx.Value("startTime").(time.Time)
	if !ok {
		duration = -1
	} else {
		duration = time.Since(startTime)
	}

	errorCode := ""
	if data.Err != nil {
		errorCode = data.Err.Error()
	}

	if data.CommandTag.Insert() || data.CommandTag.Delete() || data.CommandTag.Update() || data.CommandTag.Select() {
		queryType := data.CommandTag.String()
		index := strings.Index(data.CommandTag.String(), " ")
		if index != -1 {
			queryType = strings.TrimSpace(queryType[:index])

		}
		metrics.StoreQueryDuration(queryType, errorCode, duration)
		metrics.IncQueryCount(queryType)
	}

}
