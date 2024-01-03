package storage

import "context"

type StorageDoer interface {
	Upload(ctx context.Context, filename string, data []byte, contentType string) error
	PublicBaseURL() string
}
