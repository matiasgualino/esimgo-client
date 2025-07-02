package esimgo

import (
	"context"
	"fmt"
)

// InventoryService handles inventory-related operations
type InventoryService struct {
	client *Client
}

// NewInventoryService creates a new inventory service
func NewInventoryService(client *Client) *InventoryService {
	return &InventoryService{client: client}
}

// InventoryResponse represents the response from getting inventory
type InventoryResponse struct {
	Bundles []InventoryBundle `json:"bundles"`
}

// InventoryBundle represents a bundle in inventory
type InventoryBundle struct {
	Name        string   `json:"name"`
	Description string   `json:"desc"`
	Countries   []string `json:"countries"`
	Data        int      `json:"data"`
	Duration    int      `json:"duration"`
	Unlimited   bool     `json:"unlimited"`
}

// Get retrieves the bundle inventory
func (s *InventoryService) Get(ctx context.Context) (*InventoryResponse, error) {
	var resp InventoryResponse
	err := s.client.makeRequest(ctx, "GET", "/inventory", nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	return &resp, nil
}
