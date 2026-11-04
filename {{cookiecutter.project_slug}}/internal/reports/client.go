package reports

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// RegisteredUser is a user as returned by the users service
type RegisteredUser struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

// UsersClient calls the internal API of the users service
type UsersClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// NewUsersClient creates a client for the users service at baseURL
func NewUsersClient(baseURL, apiKey string) *UsersClient {
	return &UsersClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

// CountUsers returns the number of active users
func (c *UsersClient) CountUsers(ctx context.Context) (int64, error) {
	var body struct {
		Count int64 `json:"count"`
	}
	if err := c.get(ctx, "/internal/v1/users/count", nil, &body); err != nil {
		return 0, err
	}
	return body.Count, nil
}

// ListRegisteredUsers returns one page of active users, newest first
func (c *UsersClient) ListRegisteredUsers(ctx context.Context, limit, offset int) ([]RegisteredUser, error) {
	query := url.Values{
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}
	var body struct {
		Users []RegisteredUser `json:"users"`
	}
	if err := c.get(ctx, "/internal/v1/users/registered", query, &body); err != nil {
		return nil, err
	}
	return body.Users, nil
}

func (c *UsersClient) get(ctx context.Context, path string, query url.Values, out any) error {
	target := c.baseURL + path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return err
	}
	req.Header.Set("api-access-key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("users service returned %d for %s", resp.StatusCode, path)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
