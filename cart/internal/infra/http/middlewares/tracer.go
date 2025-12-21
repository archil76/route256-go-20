package middlewares

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

type TraceMux struct {
	h http.Handler
	t Tracer
}

func NewTraceMux(h http.Handler, t Tracer) http.Handler {
	return &TraceMux{h: h, t: t}
}

func (m *TraceMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := m.t.Start(r.Context(), r.Pattern)
	defer span.End()

	r = r.WithContext(ctx)
	m.h.ServeHTTP(w, r)
}
