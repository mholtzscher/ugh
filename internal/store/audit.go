package store

import "context"

type auditContextKey string

const (
	auditOriginContextKey         auditContextKey = "origin"
	auditShellHistoryIDContextKey auditContextKey = "shell_history_id"
)

func WithAuditOrigin(ctx context.Context, origin string) context.Context {
	return context.WithValue(ctx, auditOriginContextKey, origin)
}

func WithAuditShellHistoryID(ctx context.Context, shellHistoryID int64) context.Context {
	return context.WithValue(ctx, auditShellHistoryIDContextKey, shellHistoryID)
}

func auditOriginFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	origin, _ := ctx.Value(auditOriginContextKey).(string)
	return origin
}

func auditShellHistoryIDFromContext(ctx context.Context) *int64 {
	if ctx == nil {
		return nil
	}
	id, ok := ctx.Value(auditShellHistoryIDContextKey).(int64)
	if !ok {
		return nil
	}
	return &id
}
