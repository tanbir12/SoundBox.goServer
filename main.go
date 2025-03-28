package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
)

var mqttClient mqtt.Client

func main() {
	// Broker, credentials, and topic
	broker := "ssl://3c07990e59c44be0afbc193bdbf9af31.s1.eu.hivemq.cloud:8883"
	username := "pradhantanbir"
	password := "Tanbir123"
	caCertPath := "ca.crt" // Path to your CA certificate file

	// Load CA certificate
	certpool := x509.NewCertPool()
	caCert, err := os.ReadFile(caCertPath) // Use os.ReadFile
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	certpool.AppendCertsFromPEM(caCert)

	// Set up TLS configuration
	tlsConfig := &tls.Config{
		RootCAs: certpool,
	}

	// Configure MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("GoServerClient")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetTLSConfig(tlsConfig)

	// Create and connect the MQTT client
	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to broker: %v", token.Error())
	}
	fmt.Println("Connected to HiveMQ Cloud")

	// Set up Gorilla Mux router
	r := mux.NewRouter()
	r.HandleFunc("/api/payment", controlHandler).Methods("GET")

	// Start the server
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Handler function for the API endpoint
func controlHandler(w http.ResponseWriter, r *http.Request) {
	// Get the "led" query parameter
	led := r.URL.Query().Get("amount")
	if led == "" {
		http.Error(w, "Missing 'led' parameter", http.StatusBadRequest)
		return
	}

	// Define the topic
	topic := "soundbox@tanbir@1001"
	// Publish the message to the MQTT topic
	token := mqttClient.Publish(topic, 1, false, led)
	token.Wait()
	if token.Error() != nil {
		http.Error(w, "Failed to publish message to MQTT broker", http.StatusInternalServerError)
		return
	}

	// Send response to the client
	fmt.Fprintf(w, "Message '%s' published to topic '%s'\n", led, topic)
}
