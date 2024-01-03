package landmark

import (
	"context"
	"errors"
	"log/slog"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/leomirandadev/clean-architecture-go/pkg/tracer"
	"google.golang.org/api/option"
)

func NewGoogle(fileCredentialPath string) LandmarkClient {
	ctx := context.Background()
	landmarkClient, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(fileCredentialPath))
	if err != nil {
		panic(err)
	}

	return &googleImpl{landmarkClient}
}

type googleImpl struct {
	landmarkClient *vision.ImageAnnotatorClient
}

func (g googleImpl) FindPlace(ctx context.Context, imageURI string) ([]ImageDetails, error) {
	ctx, tr := tracer.Span(ctx, "pkg.landmark.find_place")
	defer tr.End()

	image := vision.NewImageFromURI(imageURI)
	placesFound, err := g.landmarkClient.DetectLandmarks(ctx, image, nil, 10)
	if err != nil {
		slog.Warn("client detect landmark fails", "err", err)
		return nil, err
	}

	if len(placesFound) == 0 {
		slog.Warn("no landmarks found", "err", err)
		return nil, errors.New("no landmarks found")
	}

	places := make([]ImageDetails, 0, len(placesFound))
	for _, placeFound := range placesFound {

		var geopoint Geopoint
		if len(placeFound.Locations) > 0 {
			geopoint.Latitude = placeFound.Locations[0].LatLng.Latitude
			geopoint.Longitude = placeFound.Locations[0].LatLng.Longitude
		}

		places = append(places, ImageDetails{
			Description: placeFound.Description,
			Score:       placeFound.Score,
			Geopoint:    geopoint,
		})

	}

	return places, nil
}

func (g googleImpl) Close() error {
	return g.landmarkClient.Close()
}
