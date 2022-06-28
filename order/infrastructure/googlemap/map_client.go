package googlemap

import (
	"context"
	"errors"
	"log"

	"github.com/imylam/delivery-test/configs"

	"googlemaps.github.io/maps"
)

// MapClient interface
type MapClient interface {
	GetDistance(string, string) (int, error)
}

type mapClient struct {
	client *maps.Client
}

// NewMapClient creates new a mapClient object representation of MapClient interface
func NewMapClient() MapClient {
	c, err := maps.NewClient(maps.WithAPIKey(configs.Get(configs.KeyGoogleMapAPIKey)))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	return &mapClient{client: c}
}

// GetDistance calls Google Map Distance Matrix API and returns the distance between origin and destination
func (mc *mapClient) GetDistance(origin string, destination string) (distance int, err error) {
	r := &maps.DistanceMatrixRequest{
		Origins:      []string{origin},
		Destinations: []string{destination},
		Units:        maps.UnitsMetric,
	}

	resp, err := mc.client.DistanceMatrix(context.Background(), r)
	if err != nil {
		return
	}

	respStatus := resp.Rows[0].Elements[0].Status

	if respStatus != "OK" {
		err = errors.New("Google Map API error: cannot get distance from coordinates")
		return
	}

	distance = resp.Rows[0].Elements[0].Distance.Meters

	return
}
