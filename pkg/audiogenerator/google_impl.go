package audiogenerator

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/google/uuid"
	"github.com/leomirandadev/clean-architecture-go/pkg/storage"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer"
	"google.golang.org/api/option"
)

func NewGoogle(fileCredentialPath, bucketName string) AudioGeneratorDoer {
	clientSmallAudio, err := texttospeech.NewClient(context.Background(), option.WithCredentialsFile(fileCredentialPath))
	if err != nil {
		panic(err)
	}

	storageGoogle := storage.NewGoogle(fileCredentialPath, bucketName)

	return &googleImpl{clientSmallAudio, storageGoogle, bucketName}
}

type googleImpl struct {
	clientSmallAudio *texttospeech.Client
	storageGoogle    storage.StorageDoer
	bucketName       string
}

func (g googleImpl) Generate(ctx context.Context, language, text string) (*Result, error) {
	ctx, tr := tracer.Span(ctx, "pkg.audiogenerator.generate")
	defer tr.End()

	voice, err := selectVoiceByLanguage(language)
	if err != nil {
		slog.Error("get voice fails", "err", err)
		return nil, err
	}

	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: text,
			},
		},
		Voice: voice,
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := g.clientSmallAudio.SynthesizeSpeech(ctx, req)
	if err != nil {
		slog.Error("SynthesizeLongAudio fails", "err", err)
		return nil, err
	}

	filename := fmt.Sprintf("audios/%s.mp3", uuid.New().String())

	if err := g.storageGoogle.Upload(ctx, filename, resp.AudioContent, "audio/mp3"); err != nil {
		return nil, err
	}

	return &Result{
		URL: fmt.Sprintf("%s/%s", g.storageGoogle.PublicBaseURL(), filename),
	}, nil
}

func selectVoiceByLanguage(language string) (*texttospeechpb.VoiceSelectionParams, error) {
	switch language {
	case "pt-BR":
		return &texttospeechpb.VoiceSelectionParams{
			LanguageCode: language,
			Name:         "pt-BR-Neural2-B",
		}, nil
	case "en-US":
		return &texttospeechpb.VoiceSelectionParams{
			LanguageCode: language,
			Name:         "en-US-Neural2-A",
		}, nil
	}

	return nil, errors.New("language not found")
}
