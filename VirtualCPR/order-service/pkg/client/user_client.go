package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserClient interface {
	ValidateUser(userID string) (User, error)
}

type User struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpddatedAt time.Time `json:"updated_at"`
}

type HttpUserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHttpUserClient(baseURL string, timeout int) *HttpUserClient {
	return &HttpUserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (c *HttpUserClient) ValidateUser(userID string) (User, error) {
	url := fmt.Sprintf("%s/api/users/%s", c.baseURL, userID)
	response, err := c.httpClient.Get(url)

	if err != nil {
		return User{}, fmt.Errorf("failed to call user service: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return User{}, fmt.Errorf("user not found")
	}
	if response.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("user service returned status code %d", response.StatusCode)
	}

	var user User
	if err = json.NewDecoder(response.Body).Decode(&user); err != nil {
		return User{}, fmt.Errorf("failed to parse response: %w", err)
	}
	return user, nil
}
