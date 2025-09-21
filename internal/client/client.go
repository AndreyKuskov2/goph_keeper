package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CLIClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewCLIClient(baseURL string) *CLIClient {
	return &CLIClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Token management
func (c *CLIClient) SaveToken(token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	tokenFile := filepath.Join(homeDir, ".goph_keeper_token")
	return os.WriteFile(tokenFile, []byte(token), 0600)
}

func (c *CLIClient) LoadToken() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	tokenFile := filepath.Join(homeDir, ".goph_keeper_token")
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return fmt.Errorf("failed to read token file: %v", err)
	}

	// Remove extra whitespace and newlines
	c.token = strings.TrimSpace(string(data))
	return nil
}

func (c *CLIClient) SetToken(token string) {
	c.token = token
}

// Authentication methods
func (c *CLIClient) Register(login, password string) error {
	user := User{
		Login:    login,
		Password: password,
	}

	reqBody, err := CreateJSONRequest(user)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("POST", "/api/v1/register", reqBody, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return c.handleErrorResponse(resp)
	}

	return nil
}

func (c *CLIClient) Login(login, password string) (string, error) {
	user := UserLoginRequest{
		Login:    login,
		Password: password,
	}

	reqBody, err := CreateJSONRequest(user)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.makeRequest("POST", "/api/v1/login", reqBody, false)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", c.handleErrorResponse(resp)
	}

	// Parse server response structure
	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract token from data field
	var loginResp UserLoginResponse
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &loginResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal login response: %v", err)
	}

	return loginResp.Token, nil
}

// Credentials methods
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

	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract credentials from data field - handle []interface{} case
	var credentials []Credentials
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &credentials); err != nil {
		return fmt.Errorf("failed to unmarshal credentials: %v", err)
	}

	fmt.Println("Credentials:")
	for _, cred := range credentials {
		fmt.Printf("  ID: %d, Login: %s, Meta: %s, Created: %s\n",
			cred.CredentialsID, cred.Login, cred.Meta, FormatTime(cred.CreatedAt))
	}

	return nil
}

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

	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract credentials from data field
	var credentials Credentials
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &credentials); err != nil {
		return fmt.Errorf("failed to unmarshal credentials: %v", err)
	}

	fmt.Printf("Credentials ID: %d\n", credentials.CredentialsID)
	fmt.Printf("Login: %s\n", credentials.Login)
	fmt.Printf("Password: %s\n", credentials.Password)
	fmt.Printf("Meta: %s\n", credentials.Meta)
	fmt.Printf("Created: %s\n", FormatTime(credentials.CreatedAt))

	return nil
}

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

// Text data methods
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

	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract text data from data field
	var textData []TextData
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &textData); err != nil {
		return fmt.Errorf("failed to unmarshal text data: %v", err)
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

	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract text data from data field
	var textData TextData
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &textData); err != nil {
		return fmt.Errorf("failed to unmarshal text data: %v", err)
	}

	fmt.Printf("Text Data ID: %d\n", textData.TextDataID)
	fmt.Printf("Text: %s\n", textData.Text)
	fmt.Printf("Meta: %s\n", textData.Meta)
	fmt.Printf("Created: %s\n", FormatTime(textData.CreatedAt))

	return nil
}

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

// Bank cards methods
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

	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract bank cards from data field
	var bankCards []BankCards
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &bankCards); err != nil {
		return fmt.Errorf("failed to unmarshal bank cards: %v", err)
	}

	fmt.Println("Bank Cards:")
	for _, card := range bankCards {
		fmt.Printf("  ID: %d, Card: %s, Holder: %s, Meta: %s, Created: %s\n",
			card.BankCardsID, card.CardNumber, card.Holder, card.Meta, FormatTime(card.CreatedAt))
	}

	return nil
}

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

	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract bank card from data field
	var bankCard BankCards
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &bankCard); err != nil {
		return fmt.Errorf("failed to unmarshal bank card: %v", err)
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

// Binary data methods
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

	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract binary data from data field
	var binaryData []BinariesData
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &binaryData); err != nil {
		return fmt.Errorf("failed to unmarshal binary data: %v", err)
	}

	fmt.Println("Binary Data:")
	for _, bd := range binaryData {
		fmt.Printf("  ID: %d, Size: %d bytes, Meta: %s, Created: %s\n",
			bd.BinariesDataID, len(bd.BinaryData), bd.Meta, FormatTime(bd.CreatedAt))
	}

	return nil
}

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

	var serverResp ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract binary data from data field
	var binaryData BinariesData
	dataBytes, err := json.Marshal(serverResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, &binaryData); err != nil {
		return fmt.Errorf("failed to unmarshal binary data: %v", err)
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

// Helper methods
func (c *CLIClient) makeRequest(method, path string, body []byte, useAuth bool) (*http.Response, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if useAuth {
		if c.token == "" {
			return nil, fmt.Errorf("authentication token required")
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}

func (c *CLIClient) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response: %v", err)
	}

	var errorResp ServerErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return fmt.Errorf("server error (status %d): %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("server error (status %d): %s", resp.StatusCode, errorResp.Error)
}
