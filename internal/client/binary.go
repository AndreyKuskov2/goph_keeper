package client

import (
	"fmt"
	"net/http"
	"os"
)

// CreateBinaryData creates a new binary data entry from a file
func (c *CLIClient) CreateBinaryData(filePath, meta string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	binaryData := BinariesData{
		BinaryData: data,
		Meta:       meta,
	}

	reqBody, err := CreateJSONRequest(binaryData)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("POST", "/api/v1/binaries", reqBody, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// ListBinaryData retrieves all binary data for the authenticated user
func (c *CLIClient) ListBinaryData() error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("GET", "/api/v1/binaries", nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	// Extract binary data from data field
	var binaryData []BinariesData
	if err := c.parseServerResponseData(resp, &binaryData); err != nil {
		return fmt.Errorf("failed to parse binary data: %v", err)
	}

	fmt.Println("Binary Data:")
	for _, bd := range binaryData {
		fmt.Printf("  ID: %d, Size: %d bytes, Meta: %s, Created: %s\n",
			bd.BinariesDataID, len(bd.BinaryData), bd.Meta, FormatTime(bd.CreatedAt))
	}

	return nil
}

// GetBinaryData retrieves a specific binary data entry by ID and saves it to a file
func (c *CLIClient) GetBinaryData(id, outputFile string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("GET", fmt.Sprintf("/api/v1/binaries/%s", id), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	// Extract binary data from data field
	var binaryData BinariesData
	if err := c.parseServerResponseData(resp, &binaryData); err != nil {
		return fmt.Errorf("failed to parse binary data: %v", err)
	}

	if outputFile == "" {
		outputFile = fmt.Sprintf("binary_%s.bin", id)
	}

	if err := os.WriteFile(outputFile, binaryData.BinaryData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("Binary data saved to: %s\n", outputFile)
	fmt.Printf("Size: %d bytes\n", len(binaryData.BinaryData))
	fmt.Printf("Meta: %s\n", binaryData.Meta)

	return nil
}

// UpdateBinaryData updates an existing binary data entry with a new file
func (c *CLIClient) UpdateBinaryData(id, filePath, meta string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	updateReq := UpdateBinariesDataRequest{
		BinaryData: data,
		Meta:       meta,
	}

	reqBody, err := CreateJSONRequest(updateReq)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("PUT", fmt.Sprintf("/api/v1/binaries/%s", id), reqBody, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// DeleteBinaryData deletes a binary data entry by ID
func (c *CLIClient) DeleteBinaryData(id string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("DELETE", fmt.Sprintf("/api/v1/binaries/%s", id), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}
