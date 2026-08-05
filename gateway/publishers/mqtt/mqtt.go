package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
)

const (
	publisherType = "mqtt"
	topicTemplate = "openindustrial/data/upload/%s/%s" // openindustrial/data/upload/{gatewayId}/{deviceId}
)

func init() {
	gateway.RegisterPublisher(publisherType, func() gateway.Publisher {
		return New()
	})
}

// Publisher implements the gateway.Publisher interface for MQTT.
type Publisher struct {
	id        string
	config    Config
	client    mqtt.Client
	collector *gateway.Collector
	cancel    context.CancelFunc
}

// Config holds the specific configuration for the MQTT publisher.
type Config struct {
	Broker   string `mapstructure:"broker"`
	ClientID string `mapstructure:"clientId"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	QoS      byte   `mapstructure:"qos"`
}

// New creates a new MQTT publisher instance, returning it as a gateway.Publisher interface.
func New() gateway.Publisher {
	return &Publisher{}
}

// Init initializes the publisher.
func (p *Publisher) Init(config gateway.PublisherConfig, collector *gateway.Collector) error {
	if config.Type != publisherType {
		return fmt.Errorf("invalid publisher type: expected '%s', got '%s'", publisherType, config.Type)
	}
	p.id = config.ID
	p.collector = collector

	// Parse settings
	if err := mapstructure.Decode(config.Settings, &p.config); err != nil {
		return fmt.Errorf("failed to decode mqtt settings: %w", err)
	}

	// Configure MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(p.config.Broker)
	opts.SetClientID(p.config.ClientID)
	opts.SetUsername(p.config.Username)
	opts.SetPassword(p.config.Password)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectTimeout(10 * time.Second)
	opts.OnConnect = func(client mqtt.Client) {
		log.Printf("[%s] connected to MQTT broker: %s", p.id, p.config.Broker)
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("[%s] connection lost to MQTT broker: %v", p.id, err)
	}

	p.client = mqtt.NewClient(opts)
	return nil
}

// Start connects to the MQTT broker and begins the publishing loop.
func (p *Publisher) Start(ctx context.Context) error {
	if token := p.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	pubCtx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	go p.publishLoop(pubCtx)

	return nil
}

// Stop disconnects from the MQTT broker.
func (p *Publisher) Stop(ctx context.Context) error {
	if p.cancel != nil {
		p.cancel()
	}
	p.client.Disconnect(250) // 250ms timeout
	log.Printf("[%s] MQTT publisher stopped.", p.id)
	return nil
}

// ID returns the publisher's ID.
func (p *Publisher) ID() string {
	return p.id
}

// publishLoop subscribes to the collector and publishes events to MQTT.
func (p *Publisher) publishLoop(ctx context.Context) {
	subscriber := p.collector.Subscribe()
	// No need to explicitly unsubscribe. The collector will handle closed channels.

	log.Printf("[%s] starting publish loop...", p.id)

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-subscriber:
			if !ok {
				return
			}

			// Marshal event to JSON
			payload, err := json.Marshal(event)
			if err != nil {
				log.Printf("[%s] error marshalling event to JSON: %v", p.id, err)
				continue
			}

			// Determine topic
			// Note: In a real-world scenario, gatewayId might come from config
			// For now, we assume ClientID represents the gatewayId.
			topic := fmt.Sprintf(topicTemplate, p.config.ClientID, event.DeviceID)

			// Publish to MQTT, ignore the token for async publishing
			_ = p.client.Publish(topic, p.config.QoS, false, payload)
		}
	}
}