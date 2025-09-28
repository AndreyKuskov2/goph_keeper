package client

import (
	"fmt"
	"net/http"
)

// CreateCredentials creates a new credentials entry
func (c *CLIClient) CreateCredentials(login, password, meta string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	credentials := Credentials{
		Login:    login,
		Password: password,
		Meta:     meta,
	}

	reqBody, err := CreateJSONRequest(credentials)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("POST", "/api/v1/credentials", reqBody, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// ListCredentials retrieves all credentials for the authenticated user
func (c *CLIClient) ListCredentials() error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("GET", "/api/v1/credentials", nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	// Extract credentials from data field
	var credentials []Credentials
	if err := c.parseServerResponseData(resp, &credentials); err != nil {
		return fmt.Errorf("failed to parse credentials: %v", err)
	}

	fmt.Println("Credentials:")
	for _, cred := range credentials {
		fmt.Printf("  ID: %d, Login: %s, Meta: %s, Created: %s\n",
			cred.CredentialsID, cred.Login, cred.Meta, FormatTime(cred.CreatedAt))
	}

	return nil
}

// GetCredentials retrieves a specific credentials entry by ID
func (c *CLIClient) GetCredentials(id string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("GET", fmt.Sprintf("/api/v1/credentials/%s", id), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	// Extract credentials from data field
	var credentials Credentials
	if err := c.parseServerResponseData(resp, &credentials); err != nil {
		return fmt.Errorf("failed to parse credentials: %v", err)
	}

	fmt.Printf("Credentials ID: %d\n", credentials.CredentialsID)
	fmt.Printf("Login: %s\n", credentials.Login)
	fmt.Printf("Password: %s\n", credentials.Password)
	fmt.Printf("Meta: %s\n", credentials.Meta)
	fmt.Printf("Created: %s\n", FormatTime(credentials.CreatedAt))

	return nil
}

// UpdateCredentials updates an existing credentials entry
func (c *CLIClient) UpdateCredentials(id, login, password, meta string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	updateReq := UpdateCredentialsRequest{
		Login:    login,
		Password: password,
		Meta:     meta,
	}

	reqBody, err := CreateJSONRequest(updateReq)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("PUT", fmt.Sprintf("/api/v1/credentials/%s", id), reqBody, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// DeleteCredentials deletes a credentials entry by ID
func (c *CLIClient) DeleteCredentials(id string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("DELETE", fmt.Sprintf("/api/v1/credentials/%s", id), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}
