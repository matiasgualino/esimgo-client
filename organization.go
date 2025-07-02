package esimgo

import (
	"context"
	"fmt"
	"time"
)

// OrganizationService handles organization-related operations
type OrganizationService struct {
	client *Client
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(client *Client) *OrganizationService {
	return &OrganizationService{client: client}
}

// Organization represents organization details
type Organization struct {
	Name               string    `json:"name"`
	APIKey             string    `json:"apiKey"`
	TaxLiable          string    `json:"taxLiable"`
	Addr1              string    `json:"addr1"`
	Addr2              string    `json:"addr2"`
	City               string    `json:"city"`
	Country            string    `json:"country"`
	Postcode           string    `json:"postcode"`
	CallbackURL        string    `json:"callbackUrl"`
	Notes              string    `json:"notes"`
	Groups             []string  `json:"groups"`
	Currency           string    `json:"currency"`
	Balance            int       `json:"balance"`
	TestCredit         int       `json:"testCredit"`
	TestCreditExpiry   time.Time `json:"testCreditExpiry"`
	BusinessType       string    `json:"businessType"`
	Website            string    `json:"website"`
	ProductDescription string    `json:"productDescription"`
	Users              []User    `json:"users"`
}

// User represents a user in the organization
type User struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Role         string `json:"role"`
	EmailAddress string `json:"emailAddress"`
	PhoneNumber  string `json:"phoneNumber"`
	TimeZone     string `json:"timeZone"`
}

// GetDetails retrieves organization details
func (s *OrganizationService) GetDetails(ctx context.Context) (*Organization, error) {
	var resp Organization
	err := s.client.makeRequest(ctx, "GET", "/organisation", nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization details: %w", err)
	}
	return &resp, nil
}
