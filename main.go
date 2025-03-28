package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Broker, credentials, and topic
	broker := "ssl://3c07990e59c44be0afbc193bdbf9af31.s1.eu.hivemq.cloud:8883"
	username := "pradhantanbir"
	password := "Tanbir123"
	topic := "esp32/led"
	payload := "ON" // Message to publish

	// Load CA certificate
	certpool := x509.NewCertPool()
	caCertPath := "ca.crt"                 // Path to your CA certificate file
	caCert, err := os.ReadFile(caCertPath) // Use os.ReadFile instead of ioutil.ReadFile
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
	opts.SetClientID("GoClient")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetTLSConfig(tlsConfig)

	// Create and start MQTT client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to broker: %v", token.Error())
	}
	fmt.Println("Connected to HiveMQ Cloud")

	// Publish a message to the topic
	token := client.Publish(topic, 1, false, payload)
	token.Wait()
	fmt.Printf("Message '%s' published to topic '%s'\n", payload, topic)

	// Disconnect client
	time.Sleep(2 * time.Second)
	client.Disconnect(250)
	fmt.Println("Disconnected from HiveMQ Cloud")
}
