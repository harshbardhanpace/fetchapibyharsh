package helpers

import (
	"encoding/json"
	"space/constants"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	ch            *amqp.Channel
	conn          *amqp.Connection
	connectionMux sync.Mutex
	isConnecting  bool
)

// InitializeRabbitMq establishes the initial connection to RabbitMQ
func InitializeRabbitMq() error {
	connectionMux.Lock()
	defer connectionMux.Unlock()

	if ch != nil && conn != nil && !conn.IsClosed() {
		logrus.Info("The channel is already initialized...")
		return nil
	}

	return connect()
}

// connect creates a new connection and channel to RabbitMQ
func connect() error {
	var err error
	heartBeatStr := strconv.Itoa(constants.RabbitMqHeartbeat)
	url := "amqp://" + constants.RabbitMqUser + ":" + constants.RabbitMqPassword + "@" + constants.RabbitMqAddress + "/?heartbeat=" + heartBeatStr

	if conn != nil {
		conn.Close()
	}
	if ch != nil {
		ch.Close()
	}

	conn, err = amqp.Dial(url)
	if err != nil {
		logrus.Error("connect, error in making connection: ", err)
		return err
	}

	connCloseChan := make(chan *amqp.Error)
	conn.NotifyClose(connCloseChan)
	go func() {
		closeErr := <-connCloseChan
		logrus.Warn("RabbitMQ connection closed, reason: ", closeErr)
		go reconnect()
	}()

	ch, err = conn.Channel()
	if err != nil {
		conn.Close()
		logrus.Error("connect, error in creating channel: ", err)
		return err
	}

	chanCloseChan := make(chan *amqp.Error)
	ch.NotifyClose(chanCloseChan)
	go func() {
		closeErr := <-chanCloseChan
		logrus.Warn("RabbitMQ channel closed, reason: ", closeErr)
		if conn != nil && !conn.IsClosed() {
			go recreateChannel()
		}
	}()

	err = ch.ExchangeDeclare(constants.TopicExchange, constants.Topic, true, false, false, false, nil)
	if err != nil {
		conn.Close()
		ch.Close()
		logrus.Error("connect, error in declaring exchange: ", err)
		return err
	}

	logrus.Info("RabbitMQ connection was initialized successfully")
	return nil
}

// reconnect attempts to reestablish connection with backoff
func reconnect() {
	connectionMux.Lock()
	if isConnecting {
		connectionMux.Unlock()
		return
	}
	isConnecting = true
	connectionMux.Unlock()

	defer func() {
		connectionMux.Lock()
		isConnecting = false
		connectionMux.Unlock()
	}()

	//exponential backoff for reconnection attempts
	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second
	for retries := 0; retries < 10; retries++ {
		logrus.Info("Attempting to reconnect to RabbitMQ, attempt: ", retries+1)

		connectionMux.Lock()
		err := connect()
		connectionMux.Unlock()

		if err == nil {
			logrus.Info("Successfully reconnected to RabbitMQ")
			return
		}

		logrus.Warn("Failed to reconnect to RabbitMQ: ", err)
		time.Sleep(backoff)

		//increase backoff time with a cap
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}

	logrus.Error("Alert Severity:P0-Critical, Failed to reconnect to RabbitMQ after multiple attempts")
}

// create a new channel if the existing one fails
func recreateChannel() {
	connectionMux.Lock()
	defer connectionMux.Unlock()

	if conn == nil || conn.IsClosed() {
		return
	}

	var err error
	ch, err = conn.Channel()
	if err != nil {
		logrus.Error("Failed to recreate channel: ", err)
		return
	}

	err = ch.ExchangeDeclare(constants.TopicExchange, constants.Topic, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		logrus.Error("Error in declaring exchange after recreation: ", err)
		return
	}

	logrus.Info("Successfully recreated RabbitMQ channel")
}

// verifies if the connection is healthy
func CheckConnection() bool {
	connectionMux.Lock()
	defer connectionMux.Unlock()

	if conn == nil || ch == nil || conn.IsClosed() {
		return false
	}

	//test channel to verify connection health
	testChannel, err := conn.Channel()
	if err != nil {
		return false
	}

	testChannel.Close()
	return true
}

// returns the current channel with connection verification
func GetChannel() *amqp.Channel {
	if !CheckConnection() {
		err := InitializeRabbitMq()
		if err != nil {
			logrus.Error("GetChannel, failed to initialize connection: ", err)
			return nil
		}
	}
	return ch
}

func PublishMessage(exchangeName, key string, data interface{}) error {
	channel := GetChannel()
	if channel == nil {
		err := InitializeRabbitMq()
		if err != nil {
			logrus.Error("PublishMessage, RabbitMQ connection unavailable: ", err)
			return err
		}
		channel = ch
	}

	body, err := json.Marshal(data)
	if err != nil {
		logrus.Error("PublishMessage, Failed to marshal struct to JSON: ", err)
		return err
	}

	err = channel.Publish(exchangeName, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	if err != nil {
		logrus.Error("PublishMessage, error in publishing message: ", err)

		//reconnect and publish again
		reconnErr := InitializeRabbitMq()
		if reconnErr != nil {
			logrus.Error("PublishMessage, failed to reconnect: ", reconnErr)
			return err
		}

		retryErr := ch.Publish(exchangeName, key, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

		if retryErr != nil {
			logrus.Error("PublishMessage, retry also failed: ", retryErr)
			return retryErr
		}
	}

	logrus.Info("PublishMessage, message was published successfully")
	return nil
}
