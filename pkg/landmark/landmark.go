package landmark

import "context"

type LandmarkClient interface {
	FindPlace(ctx context.Context, imageURI string) ([]ImageDetails, error)
	Close() error
}

type ImageDetails struct {
	Description string   `json:"description"`
	Score       float32  `json:"score"`
	Geopoint    Geopoint `json:"geopoint"`
}

type Geopoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
