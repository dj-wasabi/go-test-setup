package utils

import (
	"context"
)

var LogHeader string = "X-APP-LOG-ID"

type ContextWrapper struct {
	ctx context.Context
	m   map[string]any
}

func NewContextWrapper(con context.Context, logId string) ContextWrapper {
	w := ContextWrapper{
		ctx: con,
		m:   make(map[string]any),
	}
	w.m[LogHeader] = logId
	return w
}

func (w ContextWrapper) Build() context.Context {
	for k, v := range w.m {
		w.ctx = context.WithValue(w.ctx, k, v)
	}
	return w.ctx
}

func GetLogId(c context.Context) string {
	return get(LogHeader, c).(string)
}

func get(k string, c context.Context) any {
	v := c.Value(k)
	if v == nil {
		return "UNKNOWN"
	}
	return v.(string)
}
