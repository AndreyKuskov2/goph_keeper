package client

import (
	"fmt"
	"net/http"
)

// CreateTextData creates a new text data entry
func (c *CLIClient) CreateTextData(text, meta string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	textData := TextData{
		Text: text,
		Meta: meta,
	}

	reqBody, err := CreateJSONRequest(textData)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("POST", "/api/v1/text-data", reqBody, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// ListTextData retrieves all text data for the authenticated user
func (c *CLIClient) ListTextData() error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("GET", "/api/v1/text-data", nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	// Extract text data from data field
	var textData []TextData
	if err := c.parseServerResponseData(resp, &textData); err != nil {
		return fmt.Errorf("failed to parse text data: %v", err)
	}

	fmt.Println("Text Data:")
	for _, td := range textData {
		fmt.Printf("  ID: %d, Meta: %s, Created: %s\n",
			td.TextDataID, td.Meta, FormatTime(td.CreatedAt))
		fmt.Printf("  Text: %s\n", td.Text)
		fmt.Println()
	}

	return nil
}

// GetTextData retrieves a specific text data entry by ID
func (c *CLIClient) GetTextData(id string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("GET", fmt.Sprintf("/api/v1/text-data/%s", id), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	// Extract text data from data field
	var textData TextData
	if err := c.parseServerResponseData(resp, &textData); err != nil {
		return fmt.Errorf("failed to parse text data: %v", err)
	}

	fmt.Printf("Text Data ID: %d\n", textData.TextDataID)
	fmt.Printf("Text: %s\n", textData.Text)
	fmt.Printf("Meta: %s\n", textData.Meta)
	fmt.Printf("Created: %s\n", FormatTime(textData.CreatedAt))

	return nil
}

// UpdateTextData updates an existing text data entry
func (c *CLIClient) UpdateTextData(id, text, meta string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	updateReq := UpdateTextDataRequest{
		Text: text,
		Meta: meta,
	}

	reqBody, err := CreateJSONRequest(updateReq)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("PUT", fmt.Sprintf("/api/v1/text-data/%s", id), reqBody, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// DeleteTextData deletes a text data entry by ID
func (c *CLIClient) DeleteTextData(id string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("DELETE", fmt.Sprintf("/api/v1/text-data/%s", id), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}
