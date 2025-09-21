package client

import (
	"encoding/json"
	"time"
)

// User models
type User struct {
	UserID    int       `json:"user_id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

// Server response models
type ServerResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type ServerErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// Credentials models
type Credentials struct {
	CredentialsID int       `json:"credentials_id"`
	Login         string    `json:"login"`
	Password      string    `json:"password"`
	UserID        int       `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	Meta          string    `json:"meta"`
}

type UpdateCredentialsRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Meta     string `json:"meta"`
}

// Text data models
type TextData struct {
	TextDataID int       `json:"text_data_id"`
	Text       string    `json:"text"`
	UserID     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	Meta       string    `json:"meta"`
}

type UpdateTextDataRequest struct {
	Text string `json:"text"`
	Meta string `json:"meta"`
}

// Bank cards models
type BankCards struct {
	BankCardsID    int       `json:"bank_cards_id"`
	CardNumber     string    `json:"card_number"`
	Holder         string    `json:"holder"`
	CVC            string    `json:"cvc"`
	ExpirationDate time.Time `json:"expiration_date"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	Meta           string    `json:"meta"`
}

type UpdateBankCardsRequest struct {
	CardNumber     string    `json:"card_number"`
	Holder         string    `json:"holder"`
	CVC            string    `json:"cvc"`
	ExpirationDate time.Time `json:"expiration_date"`
	Meta           string    `json:"meta"`
}

// Binary data models
type BinariesData struct {
	BinariesDataID int       `json:"binaries_data_id"`
	BinaryData     []byte    `json:"binary_data"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	Meta           string    `json:"meta"`
}

type UpdateBinariesDataRequest struct {
	BinaryData []byte `json:"binary_data"`
	Meta       string `json:"meta"`
}

// Helper function to parse JSON response
func ParseJSONResponse(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Helper function to create JSON request
func CreateJSONRequest(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Helper function to format time for display
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Helper function to parse time from string
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02", timeStr)
}
