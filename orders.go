package esimgo

import (
	"context"
	"fmt"
	"time"
)

// OrdersService handles order-related operations
type OrdersService struct {
	client *Client
}

// NewOrdersService creates a new orders service
func NewOrdersService(client *Client) *OrdersService {
	return &OrdersService{client: client}
}

// OrderItem represents an order item for creating orders
type OrderItem struct {
	Type          string   `json:"type"`
	Quantity      int      `json:"quantity"`
	Item          string   `json:"item"`
	ICCIDs        []string `json:"iccids,omitempty"`
	AllowReassign bool     `json:"allowReassign,omitempty"`
}

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	Type   string      `json:"type"`   // "validate" or "transaction"
	Assign bool        `json:"assign"` // auto-assign bundles to eSIMs
	Order  []OrderItem `json:"order"`
}

// CreateOrderResponse represents the response from creating an order
type CreateOrderResponse struct {
	Total       float64   `json:"total"`
	Valid       bool      `json:"valid"`
	Currency    string    `json:"currency"`
	CreatedDate time.Time `json:"createdDate"`
	Assigned    bool      `json:"assigned"`
}

// Create creates a new order
func (s *OrdersService) Create(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	var resp CreateOrderResponse
	err := s.client.makeRequest(ctx, "POST", "/orders", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return &resp, nil
}

// Validate validates an order without creating it
func (s *OrdersService) Validate(ctx context.Context, orderItems []OrderItem, assign bool) (*CreateOrderResponse, error) {
	req := &CreateOrderRequest{
		Type:   "validate",
		Assign: assign,
		Order:  orderItems,
	}

	var resp CreateOrderResponse
	err := s.client.makeRequest(ctx, "POST", "/orders", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to validate order: %w", err)
	}
	return &resp, nil
}
