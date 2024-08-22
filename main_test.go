package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServerStartsSuccessfully(t *testing.T) {
	go main()

	time.Sleep(1 * time.Second)

	req, err := http.NewRequest("GET", "http://localhost:8080/ws", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("Expected status code 101, got %v", resp.StatusCode)
	}
}

// Successful WebSocket connection upgrade
func TestSuccessfulWebSocketConnectionUpgrade(t *testing.T) {
	req, err := http.NewRequest("GET", "/ws", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

	w := httptest.NewRecorder()
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	connection(w, req)

	if w.Result().StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("Expected status code %d, got %d", http.StatusSwitchingProtocols, w.Result().StatusCode)
	}
}
