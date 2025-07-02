package esimgo

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// CatalogueService handles catalogue-related operations
type CatalogueService struct {
	client *Client
}

// NewCatalogueService creates a new catalogue service
func NewCatalogueService(client *Client) *CatalogueService {
	return &CatalogueService{client: client}
}

// CatalogueBundle represents a bundle in the catalogue
type CatalogueBundle struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Groups         []string  `json:"groups"`
	Countries      []Country `json:"countries"`
	DataAmount     int       `json:"dataAmount"`
	Duration       int       `json:"duration"`
	Speed          []string  `json:"speed"`
	Autostart      bool      `json:"autostart"`
	Unlimited      bool      `json:"unlimited"`
	RoamingEnabled []Country `json:"roamingEnabled"`
	Price          int       `json:"price"`
}

// ListCatalogueRequest represents query parameters for listing catalogue bundles
type ListCatalogueRequest struct {
	Page        int    `json:"page,omitempty"`
	PerPage     int    `json:"perPage,omitempty"`
	Direction   string `json:"direction,omitempty"`
	OrderBy     string `json:"orderBy,omitempty"`
	Description string `json:"description,omitempty"`
	Group       string `json:"group,omitempty"`
	Countries   string `json:"countries,omitempty"`
	Region      string `json:"region,omitempty"`
}

// List retrieves all bundles available in the catalogue
func (s *CatalogueService) List(ctx context.Context, req *ListCatalogueRequest) ([]CatalogueBundle, error) {
	params := url.Values{}

	if req.Page > 0 {
		params.Set("page", strconv.Itoa(req.Page))
	}
	if req.PerPage > 0 {
		params.Set("perPage", strconv.Itoa(req.PerPage))
	}
	if req.Direction != "" {
		params.Set("direction", req.Direction)
	}
	if req.OrderBy != "" {
		params.Set("orderBy", req.OrderBy)
	}

	endpoint := "/catalogue"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	var resp []CatalogueBundle
	err := s.client.makeRequest(ctx, "GET", endpoint, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list catalogue: %w", err)
	}
	return resp, nil
}
