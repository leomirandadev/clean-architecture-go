package maps

import "context"

type MapsClient interface {
	GetAddressByGeocode(ctx context.Context, lat, lng float64) (*Address, error)
	GetPlacesNearby(ctx context.Context, lat, lng float64, radius uint) ([]PlaceDetails, error)
}

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type PlaceDetails struct {
	Name     string   `json:"name"`
	PhotoURL string   `json:"photo_url"`
	Geopoint Geopoint `json:"geopoint"`
}
type Geopoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
