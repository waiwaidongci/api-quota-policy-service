package transport

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type requestIDKey struct{}

func WithRequestID(c context.Context, id string) context.Context {
	return context.WithValue(c, requestIDKey{}, id)
}
func RequestID(c context.Context) string { v, _ := c.Value(requestIDKey{}).(string); return v }
func ID() string                         { b := make([]byte, 8); _, _ = rand.Read(b); return hex.EncodeToString(b) }
