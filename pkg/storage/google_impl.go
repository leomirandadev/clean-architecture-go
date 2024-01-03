package storage

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer"
	"google.golang.org/api/option"
)

func NewGoogle(fileCredentialPath, bucketName string) StorageDoer {
	clientStorage, err := storage.NewClient(context.Background(), option.WithCredentialsFile(fileCredentialPath))
	if err != nil {
		panic(err)
	}

	return &googleImpl{clientStorage, bucketName}
}

type googleImpl struct {
	clientStorage *storage.Client
	bucketName    string
}

func (g googleImpl) Upload(ctx context.Context, filename string, data []byte, contentType string) error {
	ctx, tr := tracer.Span(ctx, "pkg.storage.upload")
	defer tr.End()

	obj := g.clientStorage.Bucket(g.bucketName).Object(filename)

	wc := obj.NewWriter(ctx)
	if _, err := wc.Write(data); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	_, err := obj.Update(ctx, storage.ObjectAttrsToUpdate{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g googleImpl) PublicBaseURL() string {
	return "https://storage.googleapis.com/" + g.bucketName
}
