package esimgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents the eSIM Go API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new eSIM Go API client
func NewClient(apiKey string) *Client {
	return &Client{
		baseURL:    "https://api.esim-go.com/v2.4",
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SetHTTPClient allows setting a custom HTTP client
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// SetBaseURL allows setting a custom base URL (useful for testing)
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = strings.TrimSuffix(baseURL, "/")
}

// makeRequest performs HTTP requests with proper authentication
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
		}
		apiErr.StatusCode = resp.StatusCode
		return &apiErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// APIError represents an API error response
type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// Common structures used across the API

// Country represents a country
type Country struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	ISO    string `json:"iso"`
}

// Bundle represents a data bundle
type Bundle struct {
	ID          string `json:"id,omitempty"`
	Reference   string `json:"reference,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Bundle data details
	InitialQuantity   int64     `json:"initialQuantity,omitempty"`
	RemainingQuantity int64     `json:"remainingQuantity,omitempty"`
	StartTime         time.Time `json:"startTime,omitempty"`
	EndTime           time.Time `json:"endTime,omitempty"`
	Unlimited         bool      `json:"unlimited,omitempty"`
}

// Assignment represents a bundle assignment
type Assignment struct {
	ID                  string    `json:"id"`
	CallTypeGroup       string    `json:"callTypeGroup"`
	InitialQuantity     int64     `json:"initialQuantity"`
	RemainingQuantity   int64     `json:"remainingQuantity"`
	AssignmentDateTime  time.Time `json:"assignmentDateTime"`
	AssignmentReference string    `json:"assignmentReference"`
	BundleState         string    `json:"bundleState"`
	Unlimited           bool      `json:"unlimited"`
}

// ESIM represents an eSIM
type ESIM struct {
	ICCID                  string `json:"iccid"`
	PIN                    string `json:"pin,omitempty"`
	PUK                    string `json:"puk,omitempty"`
	MatchingID             string `json:"matchingId,omitempty"`
	SMDPAddress            string `json:"smdpAddress,omitempty"`
	ProfileStatus          string `json:"profileStatus,omitempty"`
	FirstInstalledDateTime int64  `json:"firstInstalledDateTime,omitempty"`
	CustomerRef            string `json:"customerRef,omitempty"`
	LastAction             string `json:"lastAction,omitempty"`
	ActionDate             string `json:"actionDate,omitempty"`
	Physical               bool   `json:"physical,omitempty"`
	AssignedDate           string `json:"assignedDate,omitempty"`
}

// ESIMGoClient is the main client that provides access to all services
type ESIMGoClient struct {
	*Client
	ESIMs        *ESIMService
	Organization *OrganizationService
	Orders       *OrdersService
	Inventory    *InventoryService
	Catalogue    *CatalogueService
	Networks     *NetworksService
}

// NewESIMGoClient creates a new complete eSIM Go API client
func NewESIMGoClient(apiKey string) *ESIMGoClient {
	baseClient := NewClient(apiKey)

	return &ESIMGoClient{
		Client:       baseClient,
		ESIMs:        NewESIMService(baseClient),
		Organization: NewOrganizationService(baseClient),
		Orders:       NewOrdersService(baseClient),
		Inventory:    NewInventoryService(baseClient),
		Catalogue:    NewCatalogueService(baseClient),
		Networks:     NewNetworksService(baseClient),
	}
}

// Common constants for API usage
const (
	// Bundle states
	BundleStateProcessing = "Processing"
	BundleStateQueued     = "Queued"
	BundleStateActive     = "Active"
	BundleStateDepleted   = "Depleted"
	BundleStateExpired    = "Expired"
	BundleStateRevoked    = "Revoked"
	BundleStateLapsed     = "Lapsed"

	// Order types
	OrderTypeValidate    = "validate"
	OrderTypeTransaction = "transaction"

	// Bundle types
	BundleTypeBundle = "bundle"

	// Direction constants
	DirectionAsc  = "asc"
	DirectionDesc = "desc"

	// Duration units
	DurationUnitDay   = "day"
	DurationUnitMonth = "month"
)
