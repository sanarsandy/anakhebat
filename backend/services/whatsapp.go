package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// WhatsAppService handles sending messages via WhatsApp
type WhatsAppService struct {
	gatewayURL   string
	apiKey       string
	senderNumber string
	client       *http.Client
}

// WhatsAppRequest represents the request payload for WhatsApp gateway
type WhatsAppRequest struct {
	Number   string `json:"number"`
	Message  string `json:"message"`
	Filename string `json:"filename,omitempty"`
}

// WhatsAppResponse represents the response from WhatsApp gateway
type WhatsAppResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewWhatsAppService creates a new WhatsApp service instance
func NewWhatsAppService() *WhatsAppService {
	gatewayURL := os.Getenv("WHATSAPP_GATEWAY_URL")
	if gatewayURL == "" {
		// Default gateway URL - menggunakan send-file endpoint sesuai dengan kode PHP
		gatewayURL = "https://anakhebat.web.id/services/wa-gateway/api/send-file"
	}

	return &WhatsAppService{
		gatewayURL:   gatewayURL,
		apiKey:       os.Getenv("WHATSAPP_API_KEY"),
		senderNumber: os.Getenv("WHATSAPP_SENDER_NUMBER"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendOTP sends OTP code via WhatsApp
func (s *WhatsAppService) SendOTP(phoneNumber string, otpCode string) error {
	message := fmt.Sprintf(`🔐 *Kode OTP Anda*

Halo,

Kode OTP Anda untuk masuk ke aplikasi Tukem:

*%s*

Kode ini berlaku selama *5 menit*.

Jangan bagikan kode ini kepada siapapun.

Jika Anda tidak meminta kode ini, abaikan pesan ini.

Terima kasih,
Tim Tukem`, otpCode)

	// Log for development/testing
	fmt.Printf("[WhatsApp OTP] Sending to %s: %s\n", phoneNumber, otpCode)

	// If gateway is configured, send via API
	if s.gatewayURL != "" && s.apiKey != "" {
		// Gunakan send-file endpoint (tanpa attachment untuk text message)
		err := s.sendViaGateway(phoneNumber, message)
		if err != nil {
			// Log error with details
			fmt.Printf("[WhatsApp Error] Failed to send OTP to %s: %v\n", phoneNumber, err)
			return err
		}
		return nil
	}

	// For development: log and return error if gateway not configured
	fmt.Printf("[WhatsApp Message]: %s\n", message)
	fmt.Printf("[Warning] WhatsApp gateway not configured. Set WHATSAPP_GATEWAY_URL and WHATSAPP_API_KEY environment variables.\n")
	return fmt.Errorf("WhatsApp gateway not configured: WHATSAPP_API_KEY is required")
}

// sendViaGateway sends message via WhatsApp gateway API
// Menggunakan multipart/form-data dengan hanya number dan message (tanpa file)
func (s *WhatsAppService) sendViaGateway(phoneNumber string, message string) error {
	// Gateway hanya perlu form-data dengan number dan message
	// Nomor yang berawalan 0 akan otomatis di-format jadi 62 oleh gateway
	// Coba format nomor: hapus + jika ada, karena gateway akan handle format sendiri
	
	// Format nomor: hapus + jika ada
	formattedNumber := phoneNumber
	if len(phoneNumber) > 0 && phoneNumber[0] == '+' {
		formattedNumber = phoneNumber[1:] // Hapus +
	}
	
	fmt.Printf("[WhatsApp Debug] Original number: %s, Formatted: %s\n", phoneNumber, formattedNumber)
	
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add form fields: number dan message saja
	writer.WriteField("number", formattedNumber)
	writer.WriteField("message", message)
	
	err := writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Log request body size untuk debugging
	fmt.Printf("[WhatsApp Debug] Request body size: %d bytes\n", requestBody.Len())

	// Create HTTP request
	req, err := http.NewRequest("POST", s.gatewayURL, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers - multipart/form-data dengan boundary
	contentType := writer.FormDataContentType()
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-api-key", s.apiKey)
	
	fmt.Printf("[WhatsApp Debug] URL: %s, Content-Type: %s\n", s.gatewayURL, contentType)

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Log response untuk debugging
	fmt.Printf("[WhatsApp Debug] Response status: %d, body: %s\n", resp.StatusCode, string(body))

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("WhatsApp gateway returned status %d: %s", resp.StatusCode, string(body))
	}

	// Try to parse response
	var waResp WhatsAppResponse
	if err := json.Unmarshal(body, &waResp); err == nil {
		if !waResp.Success {
			if waResp.Error != "" {
				return fmt.Errorf("WhatsApp gateway error: %s", waResp.Error)
			}
			return fmt.Errorf("WhatsApp gateway returned success=false: %s", string(body))
		}
	}

	fmt.Printf("[WhatsApp] Message sent successfully to %s. Response: %s\n", phoneNumber, string(body))
	return nil
}

