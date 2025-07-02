package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/matiasgualino/esimgo-client"
)

// CallbackData represents callback data structure for usage notifications
type CallbackData struct {
	ICCID     string        `json:"iccid"`
	AlertType string        `json:"alertType"`
	Bundle    esimgo.Bundle `json:"bundle"`
}

func main() {
	http.HandleFunc("/webhook", handleWebhook)
	http.HandleFunc("/health", handleHealth)

	fmt.Println("🚀 Webhook server iniciado en :8080")
	fmt.Println("📡 Endpoint: http://localhost:8080/webhook")
	fmt.Println("❤️  Health check: http://localhost:8080/health")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var callback CallbackData
	if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
		log.Printf("Error decoding webhook: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Procesar el callback
	fmt.Printf("📞 Callback recibido:\n")
	fmt.Printf("   ICCID: %s\n", callback.ICCID)
	fmt.Printf("   Tipo: %s\n", callback.AlertType)
	fmt.Printf("   Bundle: %s\n", callback.Bundle.Name)
	fmt.Printf("   Restante: %d/%d\n",
		callback.Bundle.RemainingQuantity,
		callback.Bundle.InitialQuantity)

	// Lógica de negocio basada en el tipo de alerta
	switch callback.AlertType {
	case "usage_warning":
		fmt.Println("⚠️  Advertencia de uso - notificar al cliente")
	case "usage_depleted":
		fmt.Println("🔴 Datos agotados - cliente necesita recarga")
	case "bundle_expired":
		fmt.Println("⏰ Bundle expirado - ofrecer renovación")
	default:
		fmt.Printf("ℹ️  Tipo de alerta no manejado: %s\n", callback.AlertType)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "esimgo-webhook-server",
	})
}
