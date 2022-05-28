package internal

type contextKey string

const (
	CtxCfgKey   contextKey = "cfg"
	CtxWgKey    contextKey = "wg"
	CtxQueueKey contextKey = "queue"
)
