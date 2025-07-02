package esimgo

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// NetworksService handles network-related operations
type NetworksService struct {
	client *Client
}

// NewNetworksService creates a new networks service
func NewNetworksService(client *Client) *NetworksService {
	return &NetworksService{client: client}
}

// Network represents a mobile network
type Network struct {
	Name      string   `json:"name"`
	BrandName string   `json:"brandName"`
	MCC       string   `json:"mcc"`
	MNC       string   `json:"mnc"`
	TagID     string   `json:"tagid"`
	Speed     []string `json:"speed"`
}

// CountryNetwork represents networks available in a country
type CountryNetwork struct {
	Name     string    `json:"name"`
	Networks []Network `json:"networks"`
}

// NetworksResponse represents the response from getting country networks
type NetworksResponse struct {
	CountryNetworks []CountryNetwork `json:"countryNetworks"`
}

// GetNetworksRequest represents query parameters for getting networks
type GetNetworksRequest struct {
	Countries []string `json:"countries,omitempty"`
	ISOs      []string `json:"isos,omitempty"`
	ReturnAll bool     `json:"returnAll,omitempty"`
}

// GetCountryNetworks retrieves network data for specified countries
func (s *NetworksService) GetCountryNetworks(ctx context.Context, req *GetNetworksRequest) (*NetworksResponse, error) {
	params := url.Values{}

	if len(req.Countries) > 0 {
		params.Set("countries", strings.Join(req.Countries, ","))
	}
	if len(req.ISOs) > 0 {
		params.Set("isos", strings.Join(req.ISOs, ","))
	}
	if req.ReturnAll {
		params.Set("returnAll", "true")
	}

	endpoint := "/networks"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	var resp NetworksResponse
	err := s.client.makeRequest(ctx, "GET", endpoint, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get country networks: %w", err)
	}
	return &resp, nil
}
