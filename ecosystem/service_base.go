package ecosystem

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	dbhelper "github.com/CodeClarityCE/utility-dbhelper/helper"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type ServiceDatabases struct {
	CodeClarity *sql.DB
	Knowledge   *bun.DB
	Plugins     *bun.DB
}

type QueueConfig struct {
	Name     string
	Durable  bool
	Handler  func(d amqp.Delivery)
}

type ServiceBase struct {
	ConfigSvc *ConfigService
	DB        *ServiceDatabases
	conn      *amqp.Connection
	channels  map[string]*amqp.Channel
	queues    []QueueConfig
}

func NewServiceBase() (*ServiceBase, error) {
	configSvc, err := NewConfigService()
	if err != nil {
		return nil, fmt.Errorf("config service init failed: %w", err)
	}

	db, err := connectServiceDatabases(configSvc)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	conn, err := connectAMQP(configSvc)
	if err != nil {
		return nil, fmt.Errorf("AMQP connection failed: %w", err)
	}

	return &ServiceBase{
		ConfigSvc: configSvc,
		DB:        db,
		conn:      conn,
		channels:  make(map[string]*amqp.Channel),
		queues:    make([]QueueConfig, 0),
	}, nil
}

func connectServiceDatabases(configSvc *ConfigService) (*ServiceDatabases, error) {
	// CodeClarity Database (using sql.DB for compatibility)
	codeClarityDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		configSvc.Database.User, configSvc.Database.Password,
		configSvc.Database.Host, configSvc.Database.Port,
		dbhelper.Config.Database.Results)
	codeClaritySqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(codeClarityDSN), pgdriver.WithTimeout(50*time.Second)))
	if err := codeClaritySqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("codeclarity database ping failed: %w", err)
	}

	// Knowledge Database (using bun.DB)
	knowledgeDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		configSvc.Database.User, configSvc.Database.Password,
		configSvc.Database.Host, configSvc.Database.Port,
		dbhelper.Config.Database.Knowledge)
	knowledgeSqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(knowledgeDSN), pgdriver.WithTimeout(50*time.Second)))
	knowledgeDB := bun.NewDB(knowledgeSqlDB, pgdialect.New())

	// Plugins Database (using bun.DB)
	pluginsDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		configSvc.Database.User, configSvc.Database.Password,
		configSvc.Database.Host, configSvc.Database.Port,
		dbhelper.Config.Database.Plugins)
	pluginsSqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(pluginsDSN), pgdriver.WithTimeout(50*time.Second)))
	pluginsDB := bun.NewDB(pluginsSqlDB, pgdialect.New())

	return &ServiceDatabases{
		CodeClarity: codeClaritySqlDB,
		Knowledge:   knowledgeDB,
		Plugins:     pluginsDB,
	}, nil
}

func connectAMQP(configSvc *ConfigService) (*amqp.Connection, error) {
	// Use pre-constructed URL from config service
	conn, err := amqp.Dial(configSvc.AMQP.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return conn, nil
}

func (sb *ServiceBase) AddQueue(name string, durable bool, handler func(d amqp.Delivery)) {
	sb.queues = append(sb.queues, QueueConfig{
		Name:    name,
		Durable: durable,
		Handler: handler,
	})
}

func (sb *ServiceBase) StartListening() error {
	for _, queueConfig := range sb.queues {
		if err := sb.startQueueListener(queueConfig); err != nil {
			return fmt.Errorf("failed to start listener for queue %s: %w", queueConfig.Name, err)
		}
		log.Printf("Started listening on queue: %s", queueConfig.Name)
	}
	return nil
}

func (sb *ServiceBase) startQueueListener(config QueueConfig) error {
	ch, err := sb.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	sb.channels[config.Name] = ch

	q, err := ch.QueueDeclare(
		config.Name,    // name
		config.Durable, // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			start := time.Now()
			config.Handler(d)
			log.Printf("Queue %s processed message in %v", config.Name, time.Since(start))
		}
	}()

	return nil
}

func (sb *ServiceBase) SendMessage(queueName string, data []byte) error {
	ch, exists := sb.channels[queueName]
	if !exists {
		// Create new channel for sending if it doesn't exist
		var err error
		ch, err = sb.conn.Channel()
		if err != nil {
			return fmt.Errorf("failed to open channel for sending: %w", err)
		}
		sb.channels[queueName] = ch
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue for sending: %w", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (sb *ServiceBase) WaitForever() {
	forever := make(chan bool)
	<-forever
}

func (sb *ServiceBase) Close() error {
	for _, ch := range sb.channels {
		if err := ch.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}
	
	if sb.conn != nil {
		if err := sb.conn.Close(); err != nil {
			log.Printf("Error closing AMQP connection: %v", err)
		}
	}
	
	if sb.DB != nil {
		if sb.DB.CodeClarity != nil {
			if err := sb.DB.CodeClarity.Close(); err != nil {
				log.Printf("Error closing CodeClarity database: %v", err)
			}
		}
		if sb.DB.Knowledge != nil {
			if err := sb.DB.Knowledge.Close(); err != nil {
				log.Printf("Error closing Knowledge database: %v", err)
			}
		}
		if sb.DB.Plugins != nil {
			if err := sb.DB.Plugins.Close(); err != nil {
				log.Printf("Error closing Plugins database: %v", err)
			}
		}
	}
	
	return nil
}