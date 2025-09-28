package client

import (
	"fmt"
	"net/http"
)

// CreateBankCard creates a new bank card entry
func (c *CLIClient) CreateBankCard(cardNumber, holder, cvc, expirationDate, meta string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	expDate, err := ParseTime(expirationDate)
	if err != nil {
		return fmt.Errorf("invalid expiration date format (use YYYY-MM-DD): %v", err)
	}

	bankCard := BankCards{
		CardNumber:     cardNumber,
		Holder:         holder,
		CVC:            cvc,
		ExpirationDate: expDate,
		Meta:           meta,
	}

	reqBody, err := CreateJSONRequest(bankCard)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("POST", "/api/v1/bank-cards", reqBody, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// ListBankCards retrieves all bank cards for the authenticated user
func (c *CLIClient) ListBankCards() error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("GET", "/api/v1/bank-cards", nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	// Extract bank cards from data field
	var bankCards []BankCards
	if err := c.parseServerResponseData(resp, &bankCards); err != nil {
		return fmt.Errorf("failed to parse bank cards: %v", err)
	}

	fmt.Println("Bank Cards:")
	for _, card := range bankCards {
		fmt.Printf("  ID: %d, Card: %s, Holder: %s, Meta: %s, Created: %s\n",
			card.BankCardsID, card.CardNumber, card.Holder, card.Meta, FormatTime(card.CreatedAt))
	}

	return nil
}

// GetBankCard retrieves a specific bank card entry by ID
func (c *CLIClient) GetBankCard(id string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("GET", fmt.Sprintf("/api/v1/bank-cards/%s", id), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	// Extract bank card from data field
	var bankCard BankCards
	if err := c.parseServerResponseData(resp, &bankCard); err != nil {
		return fmt.Errorf("failed to parse bank card: %v", err)
	}

	fmt.Printf("Bank Card ID: %d\n", bankCard.BankCardsID)
	fmt.Printf("Card Number: %s\n", bankCard.CardNumber)
	fmt.Printf("Holder: %s\n", bankCard.Holder)
	fmt.Printf("CVC: %s\n", bankCard.CVC)
	fmt.Printf("Expiration Date: %s\n", FormatTime(bankCard.ExpirationDate))
	fmt.Printf("Meta: %s\n", bankCard.Meta)
	fmt.Printf("Created: %s\n", FormatTime(bankCard.CreatedAt))

	return nil
}

// UpdateBankCard updates an existing bank card entry
func (c *CLIClient) UpdateBankCard(id, cardNumber, holder, cvc, expirationDate, meta string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	expDate, err := ParseTime(expirationDate)
	if err != nil {
		return fmt.Errorf("invalid expiration date format (use YYYY-MM-DD): %v", err)
	}

	updateReq := UpdateBankCardsRequest{
		CardNumber:     cardNumber,
		Holder:         holder,
		CVC:            cvc,
		ExpirationDate: expDate,
		Meta:           meta,
	}

	reqBody, err := CreateJSONRequest(updateReq)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("PUT", fmt.Sprintf("/api/v1/bank-cards/%s", id), reqBody, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// DeleteBankCard deletes a bank card entry by ID
func (c *CLIClient) DeleteBankCard(id string) error {
	if err := c.LoadToken(); err != nil {
		return fmt.Errorf("authentication required: %v", err)
	}

	resp, err := c.makeRequest("DELETE", fmt.Sprintf("/api/v1/bank-cards/%s", id), nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}
