package oakserver

import (
	"context"
	"sync"
)

type contextKey struct{}

type Traceable interface {
	GetTraceID() string
}

type immediateTracing struct {
	id string
}

func (i *immediateTracing) GetTraceID() string {
	return i.id
}

type lazyTracing struct {
	generator func() string

	mu sync.Mutex
	id string
}

func (l *lazyTracing) GetTraceID() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.id == "" {
		l.id = l.generator()
	}
	return l.id
}

func ContextWithTraceIDGenerator(parent context.Context, generator func() string) context.Context {
	if generator == nil {
		generator = func() string {
			return ""
		}
	}
	return ContextWithTracing(parent, &lazyTracing{
		generator: generator,
		mu:        sync.Mutex{},
	})
}

func ContextWithTraceID(parent context.Context, ID string) context.Context {
	return ContextWithTracing(parent, &immediateTracing{id: ID})
}

func ContextWithTracing(parent context.Context, t Traceable) context.Context {
	return context.WithValue(parent, contextKey{}, t)
}

func TraceIDFromContext(ctx context.Context) string {
	t, _ := ctx.Value(contextKey{}).(Traceable)
	if t == nil {
		return ""
	}
	return t.GetTraceID()
}
