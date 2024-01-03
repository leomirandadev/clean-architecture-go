package maps

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"
	"googlemaps.github.io/maps"
)

func NewGoogle(apiKey string) MapsClient {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}

	return &googleImpl{
		client,
		apiKey,
	}
}

type googleImpl struct {
	client *maps.Client
	apiKey string
}

func (g googleImpl) GetAddressByGeocode(ctx context.Context, lat, lng float64) (*Address, error) {
	result, err := g.client.Geocode(ctx, &maps.GeocodingRequest{
		ResultType: []string{"country", "locality"},
		LatLng: &maps.LatLng{
			Lat: lat,
			Lng: lng,
		},
	})

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("address not found")
	}

	var city string
	var country string

	for _, v := range result[0].AddressComponents {
		if city != "" && country != "" {
			break
		}
		if has(v.Types, "locality") {
			city = v.LongName
		}
		if has(v.Types, "country") {
			country = v.LongName
		}
	}

	return &Address{
		City:    city,
		Country: country,
	}, nil
}

func (g googleImpl) GetPlacesNearby(ctx context.Context, lat, lng float64, radius uint) ([]PlaceDetails, error) {
	body, err := json.Marshal(map[string]any{
		"languageCode":   "en",
		"includedTypes":  []string{"historical_landmark", "tourist_attraction"},
		"maxResultCount": 2,
		"rankPreference": "DISTANCE",
		"locationRestriction": map[string]any{
			"circle": map[string]any{
				"center": map[string]any{
					"latitude":  lat,
					"longitude": lng,
				},
				"radius": radius,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", "https://places.googleapis.com/v1/places:searchNearby", bytes.NewBuffer(body))
	if err != nil {
		slog.Debug("create request nearby places fails", "err", err)
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Goog-Api-Key", g.apiKey)
	request.Header.Add("X-Goog-FieldMask", "places.displayName,places.location,places.photos")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		slog.Debug("do request nearby places fails", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	var resultAPI NearbyResult
	if err := json.NewDecoder(resp.Body).Decode(&resultAPI); err != nil {
		slog.Debug("decode request nearby places fails", "err", err)
		return nil, err
	}

	if len(resultAPI.Places) == 0 {
		return nil, errors.New("no places found")
	}

	response := make([]PlaceDetails, len(resultAPI.Places))
	eg, errCtx := errgroup.WithContext(ctx)

	places := resultAPI.Places
	for i := range places {
		i := i
		eg.Go(func() error {
			var (
				photoURL string
				err      error
			)

			if len(places[i].Photos) > 0 {
				photoURL, err = g.getPublicPhotoURLs(errCtx, places[i].Photos[0])
				if err != nil {
					return err
				}
			}

			response[i] = PlaceDetails{
				Name:     places[i].DisplayName.Text,
				PhotoURL: photoURL,
				Geopoint: Geopoint{
					Latitude:  places[i].Location.Latitude,
					Longitude: places[i].Location.Longitude,
				},
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return response, nil
}

func (g googleImpl) getPublicPhotoURLs(ctx context.Context, photo PhotoNewAPI) (string, error) {
	url := fmt.Sprintf(
		"https://places.googleapis.com/v1/%s/media?skipHttpRedirect=true&maxHeightPx=%v&maxWidthPx=%v&key=%s",
		photo.Name,
		photo.HeightPX,
		photo.WidthPX,
		g.apiKey,
	)

	request, err := http.NewRequestWithContext(ctx, "GET", url, bytes.NewBuffer(nil))
	if err != nil {
		slog.Debug("create request get public photos url", "err", err)
		return "", err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		slog.Debug("do request get public photos url", "err", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Debug("status not ok", "status", resp.StatusCode)
		return "", errors.New("status not ok")
	}

	var resultAPI struct {
		PhotoURI string `json:"photoUri"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&resultAPI); err != nil {
		slog.Debug("decode get public photos url", "err", err)
		return "", err
	}

	return resultAPI.PhotoURI, nil
}

type NearbyResult struct {
	Places []PlaceNewAPI `json:"places"`
}

type PlaceNewAPI struct {
	DisplayName DisplayNameNewAPI  `json:"displayName"`
	Location    LocationNameNewAPI `json:"location"`
	Photos      []PhotoNewAPI      `json:"photos"`
}

type PhotoNewAPI struct {
	Name     string `json:"name"`
	WidthPX  int    `json:"widthPx"`
	HeightPX int    `json:"heightPx"`
}
type DisplayNameNewAPI struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}

type LocationNameNewAPI struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func has(slice []string, lookingFor string) bool {
	for _, v := range slice {
		if v == lookingFor {
			return true
		}
	}

	return false
}
