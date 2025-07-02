package esimgo

import (
	"context"
	"fmt"
	"net/url"
)

// ApplyBundleRequest represents a request to apply a bundle to an eSIM
type ApplyBundleRequest struct {
	Name          string `json:"name"`
	ICCID         string `json:"iccid,omitempty"`
	Repeat        int    `json:"repeat,omitempty"`
	AllowReassign bool   `json:"allowReassign,omitempty"`
}

// ApplyBundleResponse represents the response from applying a bundle
type ApplyBundleResponse struct {
	ESIMs          []ESIMBundle `json:"esims"`
	ApplyReference string       `json:"applyReference"`
}

// ESIMBundle represents an eSIM with bundle information
type ESIMBundle struct {
	ICCID  string `json:"iccid"`
	Status string `json:"status"`
	Bundle string `json:"bundle"`
}

// ESIMService handles eSIM-related operations
type ESIMService struct {
	client *Client
}

// NewESIMService creates a new eSIM service
func NewESIMService(client *Client) *ESIMService {
	return &ESIMService{client: client}
}

// ApplyBundle applies a bundle to an eSIM
func (s *ESIMService) ApplyBundle(ctx context.Context, req *ApplyBundleRequest) (*ApplyBundleResponse, error) {
	var resp ApplyBundleResponse
	err := s.client.makeRequest(ctx, "POST", "/esims/apply", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to apply bundle: %w", err)
	}
	return &resp, nil
}

// GetDetails retrieves details for a specific eSIM
func (s *ESIMService) GetDetails(ctx context.Context, iccid string, additionalFields string) (*ESIM, error) {
	endpoint := "/esims/" + iccid
	if additionalFields != "" {
		params := url.Values{}
		params.Set("additionalFields", additionalFields)
		endpoint += "?" + params.Encode()
	}

	var resp ESIM
	err := s.client.makeRequest(ctx, "GET", endpoint, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get eSIM details: %w", err)
	}
	return &resp, nil
}

// Additional structures and methods would go here...
// (This is a simplified version for the setup script)
